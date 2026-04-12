const TOKEN_KEY = "todo_token"
const USER_KEY = "todo_user"

export type AuthUser = {
    user_id: number
    username: string
    token: string
}

export function saveAuth(data: AuthUser) {
    localStorage.setItem(TOKEN_KEY, data.token)
    localStorage.setItem(USER_KEY, JSON.stringify(data))
}

export function getToken(): string | null {
    return localStorage.getItem(TOKEN_KEY)
}

export function getUser(): AuthUser | null {
    const raw = localStorage.getItem(USER_KEY)
    if (!raw) return null
    return JSON.parse(raw)
}

export function clearAuth() {
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(USER_KEY)
}

export function isAuthenticated(): boolean {
    return getToken() !== null
}
