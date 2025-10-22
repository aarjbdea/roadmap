import "./RoadmapColumn.scss"

import React from "react"
import { RoadmapColumn as RoadmapColumnModel } from "@fider/models"
import { RoadmapPostCard } from "./RoadmapPostCard"

interface RoadmapColumnProps {
  column: RoadmapColumnModel
  isStaff: boolean
  onPostMoved: (postNumber: number, fromColumnId: number, toColumnId: number, newPosition: number) => void
  onPostRemoved: (postNumber: number) => void
}

export const RoadmapColumn = (props: RoadmapColumnProps) => {
  const { column, isStaff, onPostMoved, onPostRemoved } = props

  return (
    <div className="c-roadmap-column">
      <div className="c-roadmap-column__header">
        <h3 className="c-roadmap-column__title">{column.name}</h3>
      </div>
      
      <div className="c-roadmap-column__posts">
        {column.posts && column.posts.length > 0 ? (
          column.posts.map((post, index) => (
            <RoadmapPostCard
              key={post.id}
              post={post}
              columnId={column.id}
              position={index}
              isStaff={isStaff}
              onMoved={(postNumber, toColumnId, newPosition) => onPostMoved(postNumber, column.id, toColumnId, newPosition)}
              onRemoved={onPostRemoved}
            />
          ))
        ) : (
          <div className="c-roadmap-column__empty">
            <p className="text-muted text-sm">No posts in this column</p>
          </div>
        )}
      </div>
    </div>
  )
}