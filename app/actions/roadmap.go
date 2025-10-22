package actions

import (
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/validate"
	"github.com/gosimple/slug"
)

// AssignPostToRoadmapColumn assigns a post to a roadmap column
type AssignPostToRoadmapColumn struct {
	PostNumber int `json:"postNumber"`
	ColumnID   int `json:"columnId"`
	Position   int `json:"position"`
}

// IsAuthorized returns true if current user is authorized to perform this action
func (a *AssignPostToRoadmapColumn) IsAuthorized(user *entity.User, tenant *entity.Tenant) bool {
	return user != nil && user.IsCollaborator()
}

// Validate if current model is valid
func (a *AssignPostToRoadmapColumn) Validate(user *entity.User, tenant *entity.Tenant) *validate.Result {
	result := validate.Success()

	if a.PostNumber <= 0 {
		result.AddFieldFailure("postNumber", "Post number is required")
	}

	if a.ColumnID <= 0 {
		result.AddFieldFailure("columnId", "Column ID is required")
	}

	if a.Position < 0 {
		result.AddFieldFailure("position", "Position must be non-negative")
	}

	return result
}

// Execute performs the action
func (a *AssignPostToRoadmapColumn) Execute(ctx bus.Context) error {
	assignCmd := &cmd.AssignPostToColumn{
		PostID:       a.PostNumber, // Assuming postNumber maps to post ID
		ColumnID:     a.ColumnID,
		Position:     a.Position,
		AssignedByID: ctx.User().ID,
	}

	return bus.Dispatch(ctx, assignCmd)
}

// ReorderRoadmapPost reorders a post within its roadmap column
type ReorderRoadmapPost struct {
	PostNumber  int `json:"postNumber"`
	NewPosition int `json:"newPosition"`
}

// IsAuthorized returns true if current user is authorized to perform this action
func (a *ReorderRoadmapPost) IsAuthorized(user *entity.User, tenant *entity.Tenant) bool {
	return user != nil && user.IsCollaborator()
}

// Validate if current model is valid
func (a *ReorderRoadmapPost) Validate(user *entity.User, tenant *entity.Tenant) *validate.Result {
	result := validate.Success()

	if a.PostNumber <= 0 {
		result.AddFieldFailure("postNumber", "Post number is required")
	}

	if a.NewPosition < 0 {
		result.AddFieldFailure("newPosition", "Position must be non-negative")
	}

	return result
}

// Execute performs the action
func (a *ReorderRoadmapPost) Execute(ctx bus.Context) error {
	reorderCmd := &cmd.ReorderPostInColumn{
		PostID:       a.PostNumber,
		NewPosition: a.NewPosition,
		UpdatedByID:  ctx.User().ID,
	}

	return bus.Dispatch(ctx, reorderCmd)
}

// ManageRoadmapColumn manages roadmap column creation/updates
type ManageRoadmapColumn struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	IsVisibleToPublic bool   `json:"isVisibleToPublic"`
}

// IsAuthorized returns true if current user is authorized to perform this action
func (a *ManageRoadmapColumn) IsAuthorized(user *entity.User, tenant *entity.Tenant) bool {
	return user != nil && user.IsAdministrator()
}

// Validate if current model is valid
func (a *ManageRoadmapColumn) Validate(user *entity.User, tenant *entity.Tenant) *validate.Result {
	result := validate.Success()

	if a.Name == "" {
		result.AddFieldFailure("name", "Name is required")
	} else if len(a.Name) > 50 {
		result.AddFieldFailure("name", "Name must be less than 50 characters")
	}

	return result
}

// Execute performs the action
func (a *ManageRoadmapColumn) Execute(ctx bus.Context) error {
	if a.ID > 0 {
		// Update existing column
		updateCmd := &cmd.UpdateRoadmapColumn{
			ColumnID:          a.ID,
			Name:              a.Name,
			IsVisibleToPublic: a.IsVisibleToPublic,
			UpdatedByID:       ctx.User().ID,
		}
		return bus.Dispatch(ctx, updateCmd)
	} else {
		// Create new column
		// Get next position
		var maxPosition int
		err := bus.Dispatch(ctx, &cmd.GetMaxRoadmapColumnPosition{TenantID: ctx.Tenant().ID, Result: &maxPosition})
		if err != nil {
			return err
		}

		createCmd := &cmd.CreateRoadmapColumn{
			TenantID:          ctx.Tenant().ID,
			Name:              a.Name,
			Slug:              slug.Make(a.Name),
			Position:          maxPosition + 1,
			IsVisibleToPublic: a.IsVisibleToPublic,
			CreatedByID:       ctx.User().ID,
		}
		return bus.Dispatch(ctx, createCmd)
	}
}

// GetMaxRoadmapColumnPosition gets the maximum position for roadmap columns
type GetMaxRoadmapColumnPosition struct {
	TenantID int
	Result   *int
}

// IsAuthorized returns true if current user is authorized to perform this action
func (a *GetMaxRoadmapColumnPosition) IsAuthorized(user *entity.User, tenant *entity.Tenant) bool {
	return user != nil && user.IsAdministrator()
}

// Validate if current model is valid
func (a *GetMaxRoadmapColumnPosition) Validate(user *entity.User, tenant *entity.Tenant) *validate.Result {
	return validate.Success()
}

// Execute performs the action
func (a *GetMaxRoadmapColumnPosition) Execute(ctx bus.Context) error {
	var maxPos int
	err := bus.Dispatch(ctx, &query.GetMaxRoadmapColumnPosition{TenantID: a.TenantID, Result: &maxPos})
	if err != nil {
		return err
	}
	*a.Result = maxPos
	return nil
}
