import "./Roadmap.page.scss"

import React, { useEffect, useState } from "react"
import { RoadmapData } from "@fider/models"
import { Loader, Message } from "@fider/components"
import { roadmap } from "@fider/services"
import { useFider } from "@fider/hooks"
import { RoadmapColumn as RoadmapColumnComponent } from "./components/RoadmapColumn"
import { i18n } from "@lingui/core"
import { Trans } from "@lingui/react/macro"

export interface RoadmapPageState {
  loading: boolean
  roadmapData?: RoadmapData
  error?: string
}

const RoadmapPage = () => {
  const fider = useFider()
  const [state, setState] = useState<RoadmapPageState>({
    loading: true,
  })

  useEffect(() => {
    loadRoadmap()
  }, [])

  const loadRoadmap = async () => {
    try {
      setState({ loading: true })
      const roadmapData = await roadmap.getRoadmap()
      setState({ loading: false, roadmapData })
    } catch (error) {
      setState({
        loading: false,
        error: i18n._({ id: "roadmap.error.loading", message: "Failed to load roadmap" }),
      })
    }
  }

  const handlePostMoved = async (
    postNumber: number,
    fromColumnId: number,
    toColumnId: number,
    newPosition: number
  ) => {
    try {
      await roadmap.assignPostToColumn(postNumber, toColumnId, newPosition)
      // Reload roadmap to get updated data
      await loadRoadmap()
    } catch (error) {
      console.error("Failed to move post:", error)
      // Reload to revert any optimistic updates
      await loadRoadmap()
    }
  }

  const handlePostRemoved = async (postNumber: number) => {
    try {
      await roadmap.removePostFromRoadmap(postNumber)
      await loadRoadmap()
    } catch (error) {
      console.error("Failed to remove post:", error)
      await loadRoadmap()
    }
  }

  if (state.loading) {
    return <Loader />
  }

  if (state.error) {
    return <Message type="error">{state.error}</Message>
  }

  if (!state.roadmapData || state.roadmapData.columns.length === 0) {
    return (
      <div className="text-center p-8">
        <Message type="warning">
          <Trans id="roadmap.empty">No roadmap columns have been configured yet.</Trans>
        </Message>
      </div>
    )
  }

  const isStaff = fider.session.isAuthenticated && fider.session.user.isCollaborator

  return (
    <div id="p-roadmap" className="page">
      <div className="container">
        <div className="p-roadmap__header mb-6">
          <h1 className="text-2xl font-bold">
            <Trans id="roadmap.title">Roadmap</Trans>
          </h1>
          <p className="text-muted mt-2">
            <Trans id="roadmap.description">
              Track the progress of feature requests and see what&apos;s coming next.
            </Trans>
          </p>
        </div>

        <div className="p-roadmap__columns">
          <div className="c-roadmap-columns">
            {state.roadmapData.columns.map((column) => (
              <RoadmapColumnComponent
                key={column.id}
                column={column}
                isStaff={isStaff}
                onPostMoved={handlePostMoved}
                onPostRemoved={handlePostRemoved}
              />
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}

export default RoadmapPage
