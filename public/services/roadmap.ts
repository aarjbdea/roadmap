import { httpClient } from "@fider/services"
import { RoadmapData, RoadmapColumn } from "@fider/models"

export const roadmap = {
  async getRoadmap(): Promise<RoadmapData> {
    const response = await httpClient.get("/api/v1/roadmap")
    return response.data
  },

  async assignPostToColumn(postNumber: number, columnId: number, position: number): Promise<void> {
    await httpClient.post(`/api/v1/roadmap/posts/${postNumber}/assign`, {
      postNumber,
      columnId,
      position,
    })
  },

  async removePostFromRoadmap(postNumber: number): Promise<void> {
    await httpClient.delete(`/api/v1/roadmap/posts/${postNumber}/assign`)
  },

  async reorderPostInColumn(postNumber: number, newPosition: number): Promise<void> {
    await httpClient.put(`/api/v1/roadmap/posts/${postNumber}/position`, {
      postNumber,
      newPosition,
    })
  },

  async getColumns(): Promise<RoadmapColumn[]> {
    const response = await httpClient.get("/api/v1/admin/roadmap/columns")
    return response.data
  },

  async createColumn(name: string, isVisibleToPublic: boolean): Promise<RoadmapColumn> {
    const response = await httpClient.post("/api/v1/admin/roadmap/columns", {
      name,
      isVisibleToPublic,
    })
    return response.data
  },

  async updateColumn(id: number, name: string, isVisibleToPublic: boolean): Promise<RoadmapColumn> {
    const response = await httpClient.put(`/api/v1/admin/roadmap/columns/${id}`, {
      name,
      isVisibleToPublic,
    })
    return response.data
  },

  async deleteColumn(id: number): Promise<void> {
    await httpClient.delete(`/api/v1/admin/roadmap/columns/${id}`)
  },

  async reorderColumns(columnIds: number[]): Promise<void> {
    await httpClient.put("/api/v1/admin/roadmap/columns/reorder", columnIds)
  },
}
