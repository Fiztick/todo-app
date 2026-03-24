"use client"

import { useSortable } from "@dnd-kit/sortable"
import { CSS } from "@dnd-kit/utilities"
import { Task } from "@/lib/api"

type Props = {
    task: Task
    onCompleteTask: (id: number) => void
    onDeleteTask: (id: number) => void
}

export default function TaskCard({ task, onCompleteTask, onDeleteTask }: Props) {
    const {
        attributes,
        listeners,
        setNodeRef,
        transform,
        transition,
        isDragging,
    } = useSortable({ id: task.id })

    const style = {
        transform: CSS.Transform.toString(transform),
        transition,
        opacity: isDragging ? 0.4 : 1,
    }

    return (
        <div ref={setNodeRef} style={style} {...attributes} {...listeners}
            className="bg-white rounded-lg shadow-sm border-gray-100 p-mb-2">
            <div className="flex items-center justify-between">
                <div className="flex items-center gap-2 flex-1">
                    <div {...attributes} {...listeners}
                        className="cursor-grab active:cursor-grabbing text-gray-300 hover:text-gray-500 select-none">
                        ⠿
                    </div>
                    <input
                        type="text"
                        checked={task.completed}
                        onChange={() => onCompleteTask(task.id)}
                        className="w-4 h-4 accent-blue-500"
                    />
                    <span className={`text-gray-800 text-sm ${task.completed ? "line-through text-gray-400" : ""}`}>
                        {task.title}
                    </span>
                </div>
                <button 
                    className="text-red-400 hover:text-red-600 text-sm ml-2 cursor-pointer"
                    onClick={() => onDeleteTask(task.id)}    
                >
                    x
                </button>
            </div>
        </div>
    )
}