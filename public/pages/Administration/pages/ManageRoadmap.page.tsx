import "./ManageRoadmap.page.scss"

import React, { useState, useEffect } from "react"
import { RoadmapColumn } from "@fider/models"
import { Button, Input, Message, Modal, Toggle } from "@fider/components"
import { roadmap } from "@fider/services"
import { i18n } from "@lingui/core"
import { Trans } from "@lingui/react/macro"
import { AdminPageContainer } from "../components/AdminBasePage"

export interface ManageRoadmapPageState {
  columns: RoadmapColumn[]
  loading: boolean
  error?: string
  showCreateModal: boolean
  editingColumn?: RoadmapColumn
  newColumnName: string
  newColumnPublic: boolean
}

const ManageRoadmapPage = () => {
  const [state, setState] = useState<ManageRoadmapPageState>({
    columns: [],
    loading: true,
    showCreateModal: false,
    newColumnName: "",
    newColumnPublic: true,
  })

  useEffect(() => {
    loadColumns()
  }, [])

  const loadColumns = async () => {
    try {
      setState((prev) => ({ ...prev, loading: true, error: undefined }))
      const columns = await roadmap.getColumns()
      setState((prev) => ({ ...prev, columns, loading: false }))
    } catch (error) {
      setState((prev) => ({
        ...prev,
        loading: false,
        error: i18n._({ id: "admin.roadmap.error.loading", message: "Failed to load roadmap columns" }),
      }))
    }
  }

  const handleCreateColumn = async () => {
    if (!state.newColumnName.trim()) return

    try {
      await roadmap.createColumn(state.newColumnName.trim(), state.newColumnPublic)
      setState((prev) => ({
        ...prev,
        showCreateModal: false,
        newColumnName: "",
        newColumnPublic: true,
      }))
      await loadColumns()
    } catch (error) {
      setState((prev) => ({
        ...prev,
        error: i18n._({ id: "admin.roadmap.error.creating", message: "Failed to create roadmap column" }),
      }))
    }
  }

  const handleUpdateColumn = async (column: RoadmapColumn, name: string, isPublic: boolean) => {
    try {
      await roadmap.updateColumn(column.id, name, isPublic)
      await loadColumns()
    } catch (error) {
      setState((prev) => ({
        ...prev,
        error: i18n._({ id: "admin.roadmap.error.updating", message: "Failed to update roadmap column" }),
      }))
    }
  }

  const handleDeleteColumn = async (column: RoadmapColumn) => {
    if (
      !confirm(
        i18n._({
          id: "admin.roadmap.confirm.delete",
          message: `Are you sure you want to delete the "${column.name}" column? This will remove all posts from this column.`,
        })
      )
    ) {
      return
    }

    try {
      await roadmap.deleteColumn(column.id)
      await loadColumns()
    } catch (error) {
      setState((prev) => ({
        ...prev,
        error: i18n._({ id: "admin.roadmap.error.deleting", message: "Failed to delete roadmap column" }),
      }))
    }
  }

  if (state.loading) {
    return (
      <AdminPageContainer id="p-admin-roadmap" name="roadmap" title={i18n._({ id: "admin.roadmap.title", message: "Manage Roadmap" })} subtitle={i18n._({ id: "admin.roadmap.description", message: "Configure roadmap columns and their visibility settings." })}>
        <div className="text-center p-8">
          <Trans id="admin.roadmap.loading">Loading roadmap columns...</Trans>
        </div>
      </AdminPageContainer>
    )
  }

  return (
    <AdminPageContainer id="p-admin-roadmap" name="roadmap" title={i18n._({ id: "admin.roadmap.title", message: "Manage Roadmap" })} subtitle={i18n._({ id: "admin.roadmap.description", message: "Configure roadmap columns and their visibility settings." })}>
      <div className="p-admin-roadmap">
        <div className="p-admin-roadmap__header mb-6">
          <h1 className="text-2xl font-bold">
            <Trans id="admin.roadmap.title">Manage Roadmap</Trans>
          </h1>
          <p className="text-muted mt-2">
            <Trans id="admin.roadmap.description">Configure roadmap columns and their visibility settings.</Trans>
          </p>
        </div>

        {state.error && (
          <Message type="error" className="mb-4">
            {state.error}
          </Message>
        )}

        <div className="p-admin-roadmap__actions mb-6">
          <Button variant="primary" onClick={() => setState((prev) => ({ ...prev, showCreateModal: true }))}>
            <Trans id="admin.roadmap.create">Create New Column</Trans>
          </Button>
        </div>

        <div className="p-admin-roadmap__columns">
          {state.columns.length === 0 ? (
            <div className="text-center p-8 text-muted">
              <Trans id="admin.roadmap.empty">No roadmap columns configured yet.</Trans>
            </div>
          ) : (
            <div className="c-roadmap-columns-admin">
              {state.columns.map((column) => (
                <div key={column.id} className="c-roadmap-column-admin">
                  <div className="c-roadmap-column-admin__header">
                    <div className="c-roadmap-column-admin__info">
                      <h3 className="c-roadmap-column-admin__name">{column.name}</h3>
                      <div className="c-roadmap-column-admin__meta">
                        <span className="c-roadmap-column-admin__position">Position: {column.position}</span>
                        <span className={`c-roadmap-column-admin__visibility ${column.isVisibleToPublic ? "public" : "private"}`}>
                          {column.isVisibleToPublic ? "Public" : "Private"}
                        </span>
                      </div>
                    </div>
                    <div className="c-roadmap-column-admin__actions">
                      <Button
                        variant="secondary"
                        size="small"
                        onClick={() =>
                          setState((prev) => ({
                            ...prev,
                            editingColumn: column,
                          }))
                        }
                      >
                        <Trans id="admin.roadmap.edit">Edit</Trans>
                      </Button>
                      <Button variant="danger" size="small" onClick={() => handleDeleteColumn(column)}>
                        <Trans id="admin.roadmap.delete">Delete</Trans>
                      </Button>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Create Column Modal */}
        <Modal.Window isOpen={state.showCreateModal} onClose={() => setState((prev) => ({ ...prev, showCreateModal: false }))} size="small">
          <Modal.Header>
            <Trans id="admin.roadmap.create.title">Create New Column</Trans>
          </Modal.Header>
          <Modal.Content>
            <div className="mb-4">
              <Input
                field="name"
                label={i18n._({ id: "admin.roadmap.create.name.label", message: "Column Name" })}
                value={state.newColumnName}
                onChange={(value) => setState((prev) => ({ ...prev, newColumnName: value }))}
                placeholder={i18n._({ id: "admin.roadmap.create.name.placeholder", message: "e.g., In Progress" })}
              />
            </div>
            <div className="mb-4">
              <Toggle
                field="public"
                label={i18n._({ id: "admin.roadmap.create.public.label", message: "Visible to public" })}
                active={state.newColumnPublic}
                onToggle={(checked: boolean) => setState((prev) => ({ ...prev, newColumnPublic: checked }))}
              />
              <p className="text-sm text-muted mt-1">
                <Trans id="admin.roadmap.create.public.help">Public columns are visible to all users. Private columns are only visible to staff members.</Trans>
              </p>
            </div>
          </Modal.Content>
          <Modal.Footer>
            <Button variant="tertiary" onClick={() => setState((prev) => ({ ...prev, showCreateModal: false }))}>
              <Trans id="admin.roadmap.cancel">Cancel</Trans>
            </Button>
            <Button variant="primary" onClick={handleCreateColumn} disabled={!state.newColumnName.trim()}>
              <Trans id="admin.roadmap.create.button">Create Column</Trans>
            </Button>
          </Modal.Footer>
        </Modal.Window>

        {/* Edit Column Modal */}
        {state.editingColumn && (
          <EditColumnModal
            column={state.editingColumn}
            onClose={() => setState((prev) => ({ ...prev, editingColumn: undefined }))}
            onSave={handleUpdateColumn}
          />
        )}
      </div>
    </AdminPageContainer>
  )
}

interface EditColumnModalProps {
  column: RoadmapColumn
  onClose: () => void
  onSave: (column: RoadmapColumn, name: string, isPublic: boolean) => void
}

const EditColumnModal = (props: EditColumnModalProps) => {
  const [name, setName] = useState(props.column.name)
  const [isPublic, setIsPublic] = useState(props.column.isVisibleToPublic)

  const handleSave = () => {
    props.onSave(props.column, name, isPublic)
    props.onClose()
  }

  return (
    <Modal.Window isOpen={true} onClose={props.onClose} size="small">
      <Modal.Header>
        <Trans id="admin.roadmap.edit.title">Edit Column</Trans>
      </Modal.Header>
      <Modal.Content>
        <div className="mb-4">
          <Input
            field="name"
            label={i18n._({ id: "admin.roadmap.edit.name.label", message: "Column Name" })}
            value={name}
            onChange={setName}
            placeholder={i18n._({ id: "admin.roadmap.edit.name.placeholder", message: "e.g., In Progress" })}
          />
        </div>
        <div className="mb-4">
          <Toggle
            field="public"
            label={i18n._({ id: "admin.roadmap.edit.public.label", message: "Visible to public" })}
            active={isPublic}
            onToggle={(checked: boolean) => setIsPublic(checked)}
          />
        </div>
      </Modal.Content>
      <Modal.Footer>
        <Button variant="tertiary" onClick={props.onClose}>
          <Trans id="admin.roadmap.cancel">Cancel</Trans>
        </Button>
        <Button variant="primary" onClick={handleSave} disabled={!name.trim()}>
          <Trans id="admin.roadmap.save">Save Changes</Trans>
        </Button>
      </Modal.Footer>
    </Modal.Window>
  )
}

export default ManageRoadmapPage
