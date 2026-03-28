"use client"

import { useState } from "react"
import { useSortable } from "@dnd-kit/sortable"
import { CSS } from "@dnd-kit/utilities"
import { Task } from "@/lib/api"

type Props = {
    task: Task
    onCompleteTask: (id: number) => void
    onDeleteTask: (id: number) => void
    onEditTask: (id: number, title: string) => void
}

export default function TaskCard({ task, onCompleteTask, onDeleteTask, onEditTask }: Props) {
    const {
        attributes,
        listeners,
        setNodeRef,
        transform,
        transition,
        isDragging,
    } = useSortable({ id: task.id })

    const [isEditing, setIsEditing] = useState(false)
    const [editValue, setEditValue] = useState(task.title)

    const style = {
        transform: CSS.Transform.toString(transform),
        transition,
        opacity: isDragging ? 0.4 : 1,
        width: isDragging ? "288px" : undefined,
    }

    function handleEdit() {
        if (editValue.trim() === "") return
        if (editValue === task.title) {
            setIsEditing(false)
            return
        }
        onEditTask(task.id, editValue)
        setIsEditing(false)
    }

    return (
        <div ref={setNodeRef} style={style}  {...attributes} {...listeners}
            className="bg-white rounded-lg shadow-sm border-gray-100 p-3 mb-2">
            <div className="flex items-center justify-between">
                <div className="flex items-center gap-2 flex-1">
                    <div
                        className="cursor-grab active:cursor-grabbing text-gray-300 hover:text-gray-500 select-none">
                        ⠿
                    </div>
                    <input
                        type="checkbox"
                        checked={task.completed}
                        onChange={() => onCompleteTask(task.id)}
                        className="w-4 h-4 accent-blue-500"
                    />
                    {isEditing ? (
                        <input
                            autoFocus
                            type="text"
                            value={editValue}
                            onChange={(e) => setEditValue(e.target.value)}
                            onBlur={handleEdit}
                            onKeyDown={(e) => {
                                if (e.key === "Enter") handleEdit()
                                if (e.key === "Escape") {
                                    setEditValue(task.title)
                                    setIsEditing(false)
                                }
                            }}
                            className="flex-1 text-sm text-gray-800 border border-blue-400 rounded px-1 focus:outline-none"
                        />
                    ) : (
                        <span
                            onClick={() => setIsEditing(true)}
                            className="{`flex-1 text-gray-800 text-sm cursor-pointer hover:text-blue-500 ${task.completed ? 'line-through text-gray-400' : ''`}"
                        >
                            {task.title}
                        </span>
                    )}
                </div>
                <button 
                    onClick={() => onDeleteTask(task.id)}    
                    className="text-red-400 hover:text-red-600 text-sm ml-2 cursor-pointer"
                >
                    x
                </button>
            </div>
        </div>
    )
}