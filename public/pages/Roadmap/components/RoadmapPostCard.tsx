import "./RoadmapPostCard.scss"

import React from "react"
import { Post, PostStatus } from "@fider/models"
import { VoteCounter, ShowTag, ShowPostStatus } from "@fider/components"
import { navigator } from "@fider/services"

interface RoadmapPostCardProps {
  post: Post
  columnId: number
  position: number
  isStaff: boolean
  onMoved: (postNumber: number, toColumnId: number, newPosition: number) => void
  onRemoved: (postNumber: number) => void
}

export const RoadmapPostCard = (props: RoadmapPostCardProps) => {
  const { post, isStaff, onRemoved } = props

  const handleClick = () => {
    navigator.goTo(`/posts/${post.number}/${post.slug}`)
  }

  const handleRemove = (e: React.MouseEvent) => {
    e.stopPropagation()
    onRemoved(post.number)
  }

  return (
    <div
      className={`c-roadmap-post-card ${isStaff ? "c-roadmap-post-card--draggable" : ""}`}
      onClick={handleClick}
    >
      <div className="c-roadmap-post-card__header">
        <h4 className="c-roadmap-post-card__title">{post.title}</h4>
        {isStaff && (
        <button
          className="c-roadmap-post-card__remove"
          onClick={handleRemove}
          title="Remove from roadmap"
        >
            ×
          </button>
        )}
      </div>
      
      <div className="c-roadmap-post-card__content">
        <p className="c-roadmap-post-card__description">
          {post.description.length > 100 
            ? `${post.description.substring(0, 100)}...` 
            : post.description
          }
        </p>
      </div>
      
      <div className="c-roadmap-post-card__footer">
        <div className="c-roadmap-post-card__meta">
          <VoteCounter post={post} />
          <ShowPostStatus status={PostStatus.Get(post.status)} />
        </div>
        
      {post.tags && post.tags.length > 0 && (
        <div className="c-roadmap-post-card__tags">
          {post.tags.slice(0, 2).map((tagSlug: string) => (
            <ShowTag key={tagSlug} tag={{ name: tagSlug, slug: tagSlug, color: "" }} />
          ))}
          {post.tags.length > 2 && <span className="c-roadmap-post-card__more-tags">+{post.tags.length - 2}</span>}
        </div>
      )}
      </div>
    </div>
  )
}