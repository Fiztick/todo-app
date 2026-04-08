package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
	"todo-app/backend/db"
	"todo-app/backend/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if body.Username == "" || body.Password == "" {
		respondError(w, http.StatusBadRequest, "Username and password are required")
		return
	}

	if len(body.Password) < 6 {
		respondError(w, http.StatusBadRequest, "Password must be at least 6 characters")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	var user models.User
	err = db.DB.QueryRow(
		`INSERT INTO users (username, password)
		VALUES ($1, $2)
		RETURNING id, username, created_at`,
		body.Username, string(hashedPassword),
	).Scan(&user.ID, &user.Username, &user.CreatedAt)
	if err != nil {
		respondError(w, http.StatusConflict, "Username already exists")
		return
	}

	err = createDefaultBoard(user.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create default board")
		return
	}

	respondJson(w, http.StatusCreated, map[string]string{
		"message": "User registered successfully",
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if body.Username == "" || body.Password == "" {
		respondError(w, http.StatusBadRequest, "Username and password are required")
		return
	}

	var user models.User
	err = db.DB.QueryRow(
		`SELECT id, username, password, created_at FROM users WHERE username = $1`,
		body.Username,
	).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	token, err := generateToken(user.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	respondJson(w, http.StatusOK, map[string]interface{}{
		"token":    token,
		"user_id":  user.ID,
		"username": user.Username,
	})
}

func generateToken(userID int) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func createDefaultBoard(userID int) error {
	var boardID int
	err := db.DB.QueryRow(
		`INSERT INTO boards (title, user_id)
			VALUES ($1, $2)
			RETURNING id`,
		"My Board", userID,
	).Scan(&boardID)
	if err != nil {
		return err
	}

	defaultColumns := []string{"To Do", "In Progress", "Done"}
	for i, title := range defaultColumns {
		_, err = db.DB.Exec(
			`INSERT INTO columns (title, position, board_id)
			VALUES ($1, $2, $3)`,
			title, i, boardID,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
