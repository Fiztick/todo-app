"use client"

import { useEffect, useState } from "react"
import { getTasks, createTask, completeTask, deleteTask, Task } from "@/lib/api"
import TaskInput from "@/components/TaskInput"
import TaskList from "@/components/TaskList"

export default function Home() {
  const [tasks, setTasks] = useState<Task[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState("")

  useEffect(() => {
    fetchTasks()
  }, [])

  async function fetchTasks() {
    try {
      const data = await getTasks()
      setTasks(data)
    } catch {
      setError("Failed to fetch tasks")
    } finally {
      setLoading(false)
    }
  }

  async function handleAdd(title: string) {
    try {
      const newTask = await createTask(title)
      setTasks((prev) => [...prev, newTask])
    } catch {
      setError("Failed to create task")
    }
  }

  async function handleComplete(id: number) {
    try {
      const updated = await completeTask(id)
      setTasks((prev) =>
        prev.map((task) => task.id === id ? updated : task))
    } catch {
      setError("Failed to complete task")
    }
  }

  async function handleDelete(id: number) {
    try {
      await deleteTask(id)
      setTasks((prev) => prev.filter((task) => task.id !== id))
    } catch {
      setError("Failed to delete task")
    }
  }

  if (loading) {
    return (
      <main className="min-h-screen bg-gray-50 flex items-center justify-center">
        <p className="text-gray-400">Loading tasks...</p>
      </main>
    )
  }

  return (
    <main className="min-h-screen bg-gray-50 py-10">
      <div className="max-w-xl mx-auto px-4">
        <h1 className="text-3xl font-bold text-gray-800 mb-2">
          My Todo List
        </h1>
        <p className="text-gray-400 mb-6">
          Tasks are queued in order of creation
        </p>
        {error && (
          <p className="text-red-400 mb-4">{error}</p>
        )}
        <TaskInput onAdd={handleAdd} />
        <TaskList
          tasks={tasks}
          onComplete={handleComplete}
          onDelete={handleDelete}
        />
      </div>
    </main>
  )
}