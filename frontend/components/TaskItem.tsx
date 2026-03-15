"use client"

import { Task } from "@/lib/api"

type Props = {
    task: Task
    onComplete: (id: number) => void
    onDelete: (id: number) => void
}

export default function TaskItem({ task, onComplete, onDelete }: Props) {
    return (
        <div className="flex items-center justify-between p-4 mb-2 bg-white rounded-lg shadow-sm border border-gray-100">
            <div className="flex items-center gap-3">
                <input
                    type="checkbox"
                    checked={task.completed}
                    onChange={() => onComplete(task.id)}
                    className="w-4 h-4 accent-blue-500"
                />
                <span className={`text-gray-800 ${task.completed ? "line-through text-gray-400" : ""}`}>
                    {task.title}
                </span>
            </div>
            <button
                onClick={() => onDelete(task.id)}
                className="text-red-400 hover:text-red-600 text-sm"
            >
                Delete
            </button>
        </div>
    )
}