const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export type Task = {
    id: number
    title: string
    completed: boolean
    column_id: number
    position: number
    created_at: string
}

export type Column = {
    id: number
    title:string
    position: number
    created_at: string
}

export async function getColumns(): Promise<Column[]> {
    const res = await fetch(`${API_URL}/columns`)
    const data = await res.json()
    return data ?? []
}

export async function createColumn(title: string): Promise<Column> {
    const res = await fetch(`${API_URL}/columns`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ title })
    })
    return res.json()
}

export async function updateColumn(id: number, title: string): Promise<Column> {
    const res = await fetch(`${API_URL}/columns/${id}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ title })
    })
    return res.json()
}

export async function deleteColumn(id: number): Promise<void> {
    await fetch(`${API_URL}/columns/${id}`, {
        method: "DELETE",
    })
}

export async function getTasks(): Promise<Task[]> {
    const res = await fetch(`${API_URL}/tasks`)
    const data = await res.json()
    return data ?? []
}

export async function createTask(title: string, column_Id: number): Promise<Task> {
    const res = await fetch(`${API_URL}/tasks`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ title, column_id: column_Id })
    })
    return res.json()
}

export async function completeTask(id: number): Promise<Task> {
    const res = await fetch(`${API_URL}/tasks/${id}/complete`, {
        method: "PATCH",
    })
    return res.json()
}

export async function moveTask(id: number, column_Id: number, position: number): Promise<Task> {
    const res = await fetch(`${API_URL}/tasks/${id}/move`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({column_id: column_Id, position})
    })
    return res.json()
}

export async function deleteTask(id: number): Promise<void> {
    await fetch(`${API_URL}/tasks/${id}`, {
        method: "DELETE",
    })
}

export async function editTask(id: number, title: string): Promise<Task> {
    const res = await fetch(`${API_URL}/TASKS/${id}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({title})
    })
    return res.json()
}