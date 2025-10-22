import { http } from "@fider/services"
import { RoadmapData, RoadmapColumn } from "@fider/models"

export const roadmap = {
  async getRoadmap(): Promise<RoadmapData> {
    const response = await http.get<RoadmapData>("/api/v1/roadmap")
    return response.data
  },

  async assignPostToColumn(postNumber: number, columnId: number, position: number): Promise<void> {
    await http.post(`/api/v1/roadmap/posts/${postNumber}/assign`, {
      postNumber,
      columnId,
      position,
    })
  },

  async removePostFromRoadmap(postNumber: number): Promise<void> {
    await http.delete(`/api/v1/roadmap/posts/${postNumber}/assign`)
  },

  async reorderPostInColumn(postNumber: number, newPosition: number): Promise<void> {
    await http.put(`/api/v1/roadmap/posts/${postNumber}/position`, {
      postNumber,
      newPosition,
    })
  },

  async getColumns(): Promise<RoadmapColumn[]> {
    const response = await http.get<RoadmapColumn[]>("/api/v1/admin/roadmap/columns")
    return response.data
  },

  async createColumn(name: string, isVisibleToPublic: boolean): Promise<RoadmapColumn> {
    const response = await http.post<RoadmapColumn>("/api/v1/admin/roadmap/columns", {
      name,
      isVisibleToPublic,
    })
    return response.data
  },

  async updateColumn(id: number, name: string, isVisibleToPublic: boolean): Promise<RoadmapColumn> {
    const response = await http.put<RoadmapColumn>(`/api/v1/admin/roadmap/columns/${id}`, {
      name,
      isVisibleToPublic,
    })
    return response.data
  },

  async deleteColumn(id: number): Promise<void> {
    await http.delete(`/api/v1/admin/roadmap/columns/${id}`)
  },

  async reorderColumns(columnIds: number[]): Promise<void> {
    await http.put("/api/v1/admin/roadmap/reorder-columns", columnIds)
  },
}
