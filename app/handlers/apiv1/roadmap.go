package apiv1

import (
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/web"
)

// GetRoadmap returns roadmap data for the current tenant
func GetRoadmap() web.HandlerFunc {
	return func(c *web.Context) error {
		includePrivate := c.User() != nil && c.User().IsCollaborator()
		
		getRoadmap := &query.GetRoadmapData{
			TenantID:       c.Tenant().ID,
			IncludePrivate: includePrivate,
		}

		if err := bus.Dispatch(c, getRoadmap); err != nil {
			return c.Failure(err)
		}

		return c.Ok(getRoadmap.Result)
	}
}

// AssignPostToColumn assigns a post to a roadmap column
func AssignPostToColumn() web.HandlerFunc {
	return func(c *web.Context) error {
		action := new(actions.AssignPostToRoadmapColumn)
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		if err := bus.Dispatch(c, action); err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{
			"success": true,
		})
	}
}

// ReorderPostInColumn reorders a post within its column
func ReorderPostInColumn() web.HandlerFunc {
	return func(c *web.Context) error {
		action := new(actions.ReorderRoadmapPost)
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		if err := bus.Dispatch(c, action); err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{
			"success": true,
		})
	}
}

// RemovePostFromRoadmap removes a post from the roadmap
func RemovePostFromRoadmap() web.HandlerFunc {
	return func(c *web.Context) error {
		postNumber := c.ParamAsInt("number")
		if postNumber <= 0 {
			return c.BadRequest(web.Map{
				"error": "Invalid post number",
			})
		}

		removeCmd := &cmd.RemovePostFromRoadmap{
			PostID:       postNumber,
			TenantID:     c.Tenant().ID,
			RemovedByID:  c.User().ID,
		}

		if err := bus.Dispatch(c, removeCmd); err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{
			"success": true,
		})
	}
}

// GetRoadmapColumns returns all roadmap columns for admin management
func GetRoadmapColumns() web.HandlerFunc {
	return func(c *web.Context) error {
		getColumns := &query.GetRoadmapColumns{
			TenantID:       c.Tenant().ID,
			IncludePrivate: true,
		}

		if err := bus.Dispatch(c, getColumns); err != nil {
			return c.Failure(err)
		}

		return c.Ok(getColumns.Result)
	}
}

// CreateRoadmapColumn creates a new roadmap column
func CreateRoadmapColumn() web.HandlerFunc {
	return func(c *web.Context) error {
		action := new(actions.ManageRoadmapColumn)
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		if err := bus.Dispatch(c, action); err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{
			"success": true,
		})
	}
}

// UpdateRoadmapColumn updates an existing roadmap column
func UpdateRoadmapColumn() web.HandlerFunc {
	return func(c *web.Context) error {
		columnID := c.ParamAsInt("id")
		if columnID <= 0 {
			return c.BadRequest(web.Map{
				"error": "Invalid column ID",
			})
		}

		action := new(actions.ManageRoadmapColumn)
		action.ID = columnID
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		if err := bus.Dispatch(c, action); err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{
			"success": true,
		})
	}
}

// DeleteRoadmapColumn deletes a roadmap column
func DeleteRoadmapColumn() web.HandlerFunc {
	return func(c *web.Context) error {
		columnID := c.ParamAsInt("id")
		if columnID <= 0 {
			return c.BadRequest(web.Map{
				"error": "Invalid column ID",
			})
		}

		deleteCmd := &cmd.DeleteRoadmapColumn{
			ColumnID:    columnID,
			TenantID:    c.Tenant().ID,
			DeletedByID: c.User().ID,
		}

		if err := bus.Dispatch(c, deleteCmd); err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{
			"success": true,
		})
	}
}

// ReorderColumns reorders roadmap columns
func ReorderColumns() web.HandlerFunc {
	return func(c *web.Context) error {
		var columnIDs []int
		if err := c.BindJSON(&columnIDs); err != nil {
			return c.BadRequest(web.Map{
				"error": "Invalid column IDs",
			})
		}

		reorderCmd := &cmd.ReorderRoadmapColumns{
			TenantID:    c.Tenant().ID,
			ColumnIDs:   columnIDs,
			UpdatedByID: c.User().ID,
		}

		if err := bus.Dispatch(c, reorderCmd); err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{
			"success": true,
		})
	}
}
