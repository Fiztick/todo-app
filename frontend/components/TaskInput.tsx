"use client"

import { useState } from "react"

type Props = {
    onAdd: (title: string) => void
}

export default function TaskInput({ onAdd }: Props) {
    const [title, setTitle] = useState("")

    function handleSubmit() {
        if (title.trim() === "") return
        onAdd(title)
        setTitle("")
    }

    return (
        <div className="flex gap-2 mb-6">
            <input
                type="text"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                onKeyDown={(e) => e.key === "Enter" && handleSubmit()}
                placeholder="Add a new task..."
                className="flex-1 px-4 py-2 rounded-lg border border-gray-300 focus:outline-none focus:ring-2 focus:ring-blue-500 text-gray-800"
            />
            <button
                onClick={handleSubmit}
                className="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600"
            >
                Add
            </button>
        </div>
    )
}