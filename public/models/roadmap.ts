import { Post } from "./post"

export interface RoadmapColumn {
  id: number
  name: string
  slug: string
  position: number
  isVisibleToPublic: boolean
  posts: Post[]
}

export interface RoadmapData {
  columns: RoadmapColumn[]
}

export interface RoadmapAssignment {
  id: number
  postId: number
  columnId: number
  tenantId: number
  position: number
  assignedAt: string
  assignedById: number
}
