"use client"

import { useState, useEffect } from "react"
import { DndContext, DragEndEvent, PointerSensor, useSensor, useSensors, pointerWithin, rectIntersection, CollisionDetection } from "@dnd-kit/core"
import { Column as ColumnType, Task, getColumns, getTasks, createColumn, updateColumn, deleteColumn, createTask, completeTask, moveTask, deleteTask } from "@/lib/api"
import Column from "./Column"

export default function Board() {
    const [columns, setColumns] = useState<ColumnType[]>([])
    const [tasks, setTasks] = useState<Task[]>([])
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState("")
    const [newColumnTitle, setNewColumnTitle] = useState("")

    const sensors = useSensors(
        useSensor(PointerSensor, {
            activationConstraint: { distance: 5 },
        })
    )

    useEffect(() => {
        fetchAll()
    }, [])

    async function fetchAll() {
        try {
            const [cols, tsks] = await Promise.all([getColumns(), getTasks()])
            setColumns(cols)
            setTasks(tsks)
        } catch {
            setError("Failed to load board")
        } finally {
            setLoading(false)
        }
    }

    function getTasksForColumn(columnId: number) {
        return tasks
            .filter((t) => t.column_id === columnId)
            .sort((a, b) => a.position - b.position)
    }

    async function handleAddTask(title: string, columnId: number) {
        try {
            const newTask = await createTask(title, columnId)
            setTasks((prev) => [...prev, newTask])
        } catch {
            setError("Failed to add task")
        }
    }

    async function handleCompleteTask(id: number) {
        try {
            const updated = await completeTask(id)
            setTasks((prev) => prev.map((t) => (t.id === id ? updated : t)))
        } catch {
            setError("Failed to complete task")
        }
    }

    async function handleDeleteTask(id: number) {
        try {
            await deleteTask(id)
            setTasks((prev) => prev.filter((t) => t.id !== id))
        } catch {
            setError("Failed to delete task")
        }
    }

    async function handleAddColumn() {
        if (newColumnTitle.trim() === "") return
        try {
            const col = await createColumn(newColumnTitle)
            setColumns((prev) => [...prev, col])
            setNewColumnTitle("")
        } catch {
            setError("Failed to add column")
        }
    }

    async function handleDeleteColumn(id: number) {
        try {
            await deleteColumn(id)
            setColumns((prev) => prev.filter((c) => c.id !== id))
            setTasks((prev) => prev.filter((t) => t.column_id !== id))
        } catch {
            setError("Failed to delete column")
        }
    }

    async function handleRenameColumn(id: number, title: string) {
        try {
            const updated = await updateColumn(id, title)
            setColumns((prev) => prev.map((c) => c.id === id ? updated : c))
        } catch {
            setError("Failed to rename column")
        }
    }

    const collisionDetection: CollisionDetection = (args) => {
        const pointerCollisions = pointerWithin(args)
        if (pointerCollisions.length > 0) {
            return pointerCollisions
        }
        return rectIntersection(args)
    }

    async function onDragEnd(event: DragEndEvent) {
        const { active, over } = event
        if (!over) return

        const activeTask = tasks.find((t) => t.id === active.id)
        if (!activeTask) return

        const overId = over.id as number
        const overTask = tasks.find((t) => t.id === overId)
        const isOverColumn = columns.some((c) => c.id === overId)
        const targetColumnId = isOverColumn ? overId : overTask ? overTask.column_id : overId

        if (activeTask.column_id === targetColumnId && !isOverColumn) return

        const newPosition = tasks
            .filter((t) => t.column_id === targetColumnId && t.id !== activeTask.id)
            .length

        try {
            await moveTask(activeTask.id, targetColumnId, newPosition)
            setTasks((prev) => prev.map((t) =>
                t.id === activeTask.id
                    ? { ...t, column_id: targetColumnId, position: newPosition }
                    : t
            ))
        } catch {
            setError("Failed to move task")
        }
    }

    if (loading) {
        return (
            <div className="flex items-center justify-center min-h-screen">
                <p className="text-gray-400">Loading board...</p>
            </div>
        )
    }

    return (
        <div className="min-h-screen bg-gray-50 p-6">
            <div className="mb-6">
                <h1 className="text-3xl font-bold text-gray-800 mb-1">My Board</h1>
                <p className="text-gray-400 text-sm">Drag cards between columns</p>
            </div>

            {error && (
                <p className="text-red-400 mb-4 text-sm">{error}</p>
            )}

            <DndContext
                sensors={sensors}
                collisionDetection={collisionDetection}
                onDragEnd={onDragEnd}
            >
                <div className="flex gap-4 overflow-x-auto pb-4">
                    {columns.map((col) => (
                        <Column
                            key={col.id}
                            column={col}
                            tasks={getTasksForColumn(col.id)}
                            onCompleteTask={handleCompleteTask}
                            onDeleteTask={handleDeleteTask}
                            onAddTask={handleAddTask}
                            onDeleteColumn={handleDeleteColumn}
                            onRenameColumn={handleRenameColumn}
                        />
                    ))}
                    <div className="flex flex-col gap-2 w-72 flex-shrink-0">
                        <input
                            type="text"
                            value={newColumnTitle}
                            onChange={(e) => setNewColumnTitle(e.target.value)}
                            onKeyDown={(e) => e.key === "Enter" && handleAddColumn()}
                            placeholder="New column name..."
                            className="text-sm px-3 py-2 rounded-lg border border-gray-300 focus:outline-none focus:ring-2 focus:ring-blue-400 text-gray-800"
                        />
                        <button
                            onClick={handleAddColumn}
                            className="text-sm px-3 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600"
                        >
                            + Add Column
                        </button>
                    </div>
                </div>
            </DndContext>
        </div>
    )
}