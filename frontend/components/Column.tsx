"use client"

import { useState } from "react"
import { useDroppable } from "@dnd-kit/core"
import { SortableContext, verticalListSortingStrategy } from "@dnd-kit/sortable"
import { Column as ColumnType, Task } from "@/lib/api"
import TaskCard from "./TaskCard"

type Props = {
    column: ColumnType
    tasks: Task[]
    onComplete: (id: number) => void
    onDelete: (id: number) => void
    onAddTask: (title: string, columnId: number) => void
    onDeleteColumn: (id: number) => void
    onRenameColumn: (id: number, newTitle: string) => void
}

export default function Column({
    column,
    tasks,
    onComplete,
    onDelete,
    onAddTask,
    onDeleteColumn,
    onRenameColumn,
}: Props) {
    const { setNodeRef, isOver } = useDroppable({ id: column.id })
    const [input, setInput] = useState("")
    const [isRenaming, setIsRenaming] = useState(false)
    const [renameValue, setRenameValue] = useState(column.title)

    function handleAddTask() {
        if (input.trim() === "") return
        onAddTask(input, column.id)
        setInput("")
    }

    function handleRename() {
        if (renameValue.trim() === "") return
        onRenameColumn(column.id, renameValue)
        setIsRenaming(false)
    }

    return (
        <div className="bg-gray-100 rounded-x1 p-3 w-72 flex-shrink-0">
            <div className="flex items-center justify-between mb-3">
                {isRenaming ? (
                    <input
                        autoFocus
                        value={renameValue}
                        onChange={(e) => setRenameValue(e.target.value)}
                        onBlur={handleRename}
                        onKeyDown={(e) => e.key === "Enter" && handleRename()}
                        className="text-sm font-semibold bg-white border border-blue-400 rounded px-2 py-1 w-full text-gray-800"
                    />
                ) : (
                    <h3
                        onClick={() => setIsRenaming(true)}
                        className="text-sm font-semibold text-gray-700 cursor-pointer hover:text-blue-500"
                    >
                        {column.title}
                        <span className="ml-1 text-gray-400 text-xs">({tasks.length})</span>
                    </h3>
                )}
                <button
                    className="text-gray-400 hover:text-red-500 text-xs ml-2"
                    onClick={() => onDeleteColumn(column.id)}
                >
                    x
                </button>
            </div>

            <div
                ref={setNodeRef}
                className={`min-h-20 rounded-lg transition-colors ${isOver ? "bg-blue-50" : ""}`}
            >
                <SortableContext
                    items={tasks.map((t) => t.id)}
                    strategy={verticalListSortingStrategy}
                >
                    {tasks.map((task) => (
                        <TaskCard
                            key={task.id}
                            task={task}
                            onComplete={onComplete}
                            onDelete={onDelete}
                        />
                    ))}
                </SortableContext>
            </div>

            <div className="flex gap-1 mt-3">
                <input
                    type="text"
                    value={input}
                    onChange={(e) => setInput(e.target.value)}
                    onKeyDown={(e) => e.key === "Enter" && handleAddTask()}
                    placeholder="Add a task..."
                    className="flex-1 text-sm px-2 py-1 rounded border border-gray-300 focus:outline-none focus:ring-1 focus:ring-blue-400 text-gray-800"
                />
                <button
                    onClick={handleAddTask}
                    className="text-sm px-2 py-1 bg-blue-500 text-white rounded hover:bg-blue-600"
                >
                    +
                </button>
            </div>
        </div>
    )
}