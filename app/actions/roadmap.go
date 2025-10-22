package actions

import (
	"context"

	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/pkg/validate"
)

// AssignPostToRoadmap is the action to assign a post to a roadmap column
type AssignPostToRoadmap struct {
	PostID   int `json:"postId"`
	ColumnID int `json:"columnId"`
	Position int `json:"position"`
}

// IsAuthorized returns true if current user is authorized to perform this action
func (a *AssignPostToRoadmap) IsAuthorized(ctx context.Context, user *entity.User) bool {
	return user != nil && user.IsCollaborator()
}

// Validate if current model is valid
func (a *AssignPostToRoadmap) Validate(ctx context.Context, user *entity.User) *validate.Result {
	result := validate.Success()

	if a.PostID <= 0 {
		result.AddFieldFailure("postId", "Post ID is required")
	}

	if a.ColumnID <= 0 {
		result.AddFieldFailure("columnId", "Column ID is required")
	}

	if a.Position < 0 {
		result.AddFieldFailure("position", "Position must be non-negative")
	}

	return result
}

// RemovePostFromRoadmap is the action to remove a post from the roadmap
type RemovePostFromRoadmap struct {
	PostID int `json:"postId"`
}

// IsAuthorized returns true if current user is authorized to perform this action
func (a *RemovePostFromRoadmap) IsAuthorized(ctx context.Context, user *entity.User) bool {
	return user != nil && user.IsCollaborator()
}

// Validate if current model is valid
func (a *RemovePostFromRoadmap) Validate(ctx context.Context, user *entity.User) *validate.Result {
	result := validate.Success()

	if a.PostID <= 0 {
		result.AddFieldFailure("postId", "Post ID is required")
	}

	return result
}

// ReorderPostInRoadmap is the action to reorder a post within a roadmap column
type ReorderPostInRoadmap struct {
	PostID      int `json:"postId"`
	NewPosition int `json:"newPosition"`
}

// IsAuthorized returns true if current user is authorized to perform this action
func (a *ReorderPostInRoadmap) IsAuthorized(ctx context.Context, user *entity.User) bool {
	return user != nil && user.IsCollaborator()
}

// Validate if current model is valid
func (a *ReorderPostInRoadmap) Validate(ctx context.Context, user *entity.User) *validate.Result {
	result := validate.Success()

	if a.PostID <= 0 {
		result.AddFieldFailure("postId", "Post ID is required")
	}

	if a.NewPosition < 0 {
		result.AddFieldFailure("newPosition", "Position must be non-negative")
	}

	return result
}

// CreateRoadmapColumn is the action to create a new roadmap column
type CreateRoadmapColumn struct {
	Name              string `json:"name"`
	IsVisibleToPublic bool   `json:"isVisibleToPublic"`
}

// IsAuthorized returns true if current user is authorized to perform this action
func (a *CreateRoadmapColumn) IsAuthorized(ctx context.Context, user *entity.User) bool {
	return user != nil && user.IsAdministrator()
}

// Validate if current model is valid
func (a *CreateRoadmapColumn) Validate(ctx context.Context, user *entity.User) *validate.Result {
	result := validate.Success()

	if a.Name == "" {
		result.AddFieldFailure("name", "Name is required")
	} else if len(a.Name) > 100 {
		result.AddFieldFailure("name", "Name must be less than 100 characters")
	}

	return result
}

// UpdateRoadmapColumn is the action to update an existing roadmap column
type UpdateRoadmapColumn struct {
	ColumnID          int    `json:"columnId"`
	Name              string `json:"name"`
	IsVisibleToPublic bool   `json:"isVisibleToPublic"`
}

// IsAuthorized returns true if current user is authorized to perform this action
func (a *UpdateRoadmapColumn) IsAuthorized(ctx context.Context, user *entity.User) bool {
	return user != nil && user.IsAdministrator()
}

// Validate if current model is valid
func (a *UpdateRoadmapColumn) Validate(ctx context.Context, user *entity.User) *validate.Result {
	result := validate.Success()

	if a.ColumnID <= 0 {
		result.AddFieldFailure("columnId", "Column ID is required")
	}

	if a.Name == "" {
		result.AddFieldFailure("name", "Name is required")
	} else if len(a.Name) > 100 {
		result.AddFieldFailure("name", "Name must be less than 100 characters")
	}

	return result
}

// DeleteRoadmapColumn is the action to delete a roadmap column
type DeleteRoadmapColumn struct {
	ColumnID int `json:"columnId"`
}

// IsAuthorized returns true if current user is authorized to perform this action
func (a *DeleteRoadmapColumn) IsAuthorized(ctx context.Context, user *entity.User) bool {
	return user != nil && user.IsAdministrator()
}

// Validate if current model is valid
func (a *DeleteRoadmapColumn) Validate(ctx context.Context, user *entity.User) *validate.Result {
	result := validate.Success()

	if a.ColumnID <= 0 {
		result.AddFieldFailure("columnId", "Column ID is required")
	}

	return result
}

// ReorderRoadmapColumns is the action to reorder roadmap columns
type ReorderRoadmapColumns struct {
	ColumnIDs []int `json:"columnIds"`
}

// IsAuthorized returns true if current user is authorized to perform this action
func (a *ReorderRoadmapColumns) IsAuthorized(ctx context.Context, user *entity.User) bool {
	return user != nil && user.IsAdministrator()
}

// Validate if current model is valid
func (a *ReorderRoadmapColumns) Validate(ctx context.Context, user *entity.User) *validate.Result {
	result := validate.Success()

	if len(a.ColumnIDs) == 0 {
		result.AddFieldFailure("columnIds", "At least one column ID is required")
	}

	return result
}
