import "./AssignToRoadmapModal.scss"

import React, { useState, useEffect } from "react"
import { RoadmapColumn, Post } from "@fider/models"
import { Modal, Button, Select, SelectOption, Message } from "@fider/components"
import { roadmap } from "@fider/services"
import { i18n } from "@lingui/core"
import { Trans } from "@lingui/react/macro"

interface AssignToRoadmapModalProps {
  post: Post
  isOpen: boolean
  onClose: () => void
  onAssigned?: () => void
}

export const AssignToRoadmapModal = (props: AssignToRoadmapModalProps) => {
  const { post, isOpen, onClose, onAssigned } = props
  const [columns, setColumns] = useState<RoadmapColumn[]>([])
  const [selectedColumnId, setSelectedColumnId] = useState<number>(0)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string>("")

  useEffect(() => {
    if (isOpen) {
      loadColumns()
    }
  }, [isOpen])

  const loadColumns = async () => {
    try {
      const roadmapData = await roadmap.getRoadmap()
      setColumns(roadmapData.columns)
    } catch (err) {
      setError(i18n._({ id: "roadmap.modal.error.loading", message: "Failed to load roadmap columns" }))
    }
  }

  const handleAssign = async () => {
    if (!selectedColumnId) return
    try {
      setLoading(true)
      setError("")
      await roadmap.assignPostToColumn(post.number, selectedColumnId, 0)
      onAssigned?.()
      onClose()
    } catch (err) {
      setError(i18n._({ id: "roadmap.modal.error.assigning", message: "Failed to assign post to roadmap" }))
    } finally {
      setLoading(false)
    }
  }

  const handleRemove = async () => {
    try {
      setLoading(true)
      setError("")
      await roadmap.removePostFromRoadmap(post.number)
      onAssigned?.()
      onClose()
    } catch (err) {
      setError(i18n._({ id: "roadmap.modal.error.removing", message: "Failed to remove post from roadmap" }))
    } finally {
      setLoading(false)
    }
  }

  const columnOptions = columns.map((col) => ({
    value: col.id.toString(),
    label: col.name,
  }))

  return (
    <Modal.Window isOpen={isOpen} onClose={onClose} size="small">
      <Modal.Header>
        <Trans id="roadmap.modal.title">Assign to Roadmap</Trans>
      </Modal.Header>
      <Modal.Content>
        {error && (
          <Message type="error" className="mb-4">
            {error}
          </Message>
        )}
        <div className="mb-4">
          <p className="text-sm text-muted mb-2">
            <Trans id="roadmap.modal.description">
              Choose which roadmap column this post should be assigned to.
            </Trans>
          </p>
          <Select
            field="column"
            label={i18n._({ id: "roadmap.modal.column.label", message: "Roadmap Column" })}
            options={columnOptions}
            defaultValue={selectedColumnId > 0 ? selectedColumnId.toString() : undefined}
            onChange={(option: SelectOption | undefined) => setSelectedColumnId(option ? parseInt(option.value) : 0)}
          />
        </div>
        <div className="text-sm text-muted">
          <Trans id="roadmap.modal.post.info">Post: <strong>{post.title}</strong></Trans>
        </div>
      </Modal.Content>
      <Modal.Footer>
        <Button variant="tertiary" onClick={onClose}>
          <Trans id="roadmap.modal.cancel">Cancel</Trans>
        </Button>
        <Button variant="danger" onClick={handleRemove} disabled={loading} className="mr-2">
          <Trans id="roadmap.modal.remove">Remove from Roadmap</Trans>
        </Button>
        <Button variant="primary" onClick={handleAssign} disabled={loading || !selectedColumnId}>
          <Trans id="roadmap.modal.assign">Assign to Roadmap</Trans>
        </Button>
      </Modal.Footer>
    </Modal.Window>
  )
}
