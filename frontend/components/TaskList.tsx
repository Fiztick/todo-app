"use client"

import { Task } from "@/lib/api"
import TaskItem from "./TaskItem"

type Props = {
    tasks: Task[]
    onComplete: (id: number) => void
    onDelete: (id: number) => void
}

export default function TaskList({ tasks, onComplete, onDelete }: Props) {
    if (tasks.length === 0) {
        return (
            <div className="text-center text-gray-400 py-10">
                No tasks yet. Add one above!
            </div>
        )
    }

    return (
        <div>
            {tasks.map((task) => (
                <TaskItem
                    key={task.id}
                    task={task}
                    onComplete={onComplete}
                    onDelete={onDelete}
                />
            ))}
        </div>
    )
}