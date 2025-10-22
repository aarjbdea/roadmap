package cmd

import (
	"github.com/getfider/fider/app/models/entity"
)

// AssignPostToColumn assigns a post to a roadmap column
type AssignPostToColumn struct {
	PostID       int
	ColumnID     int
	Position     int
	AssignedByID int
	Result       *entity.RoadmapAssignment
}

// RemovePostFromRoadmap removes a post from the roadmap
type RemovePostFromRoadmap struct {
	PostID       int
	TenantID     int
	RemovedByID  int
}

// ReorderPostInColumn changes the position of a post within its column
type ReorderPostInColumn struct {
	PostID       int
	NewPosition  int
	UpdatedByID  int
}

// CreateRoadmapColumn creates a new roadmap column
type CreateRoadmapColumn struct {
	TenantID          int
	Name              string
	Slug              string
	Position          int
	IsVisibleToPublic bool
	CreatedByID       int
	Result            *entity.RoadmapColumn
}

// UpdateRoadmapColumn updates an existing roadmap column
type UpdateRoadmapColumn struct {
	ColumnID          int
	Name              string
	IsVisibleToPublic bool
	UpdatedByID       int
	Result            *entity.RoadmapColumn
}

// DeleteRoadmapColumn deletes a roadmap column
type DeleteRoadmapColumn struct {
	ColumnID     int
	TenantID     int
	DeletedByID  int
}

// ReorderRoadmapColumns changes the order of roadmap columns
type ReorderRoadmapColumns struct {
	TenantID     int
	ColumnIDs    []int
	UpdatedByID  int
}

// GetMaxRoadmapColumnPosition gets the maximum position for roadmap columns
type GetMaxRoadmapColumnPosition struct {
	TenantID int
	Result   *int
}
