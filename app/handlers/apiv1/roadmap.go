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
		action := new(actions.AssignPostToRoadmap)
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
		action := new(actions.ReorderPostInRoadmap)
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
		postNumber, err := c.ParamAsInt("number")
		if err != nil || postNumber <= 0 {
			return c.BadRequest(web.Map{
				"error": "Invalid post number",
			})
		}

		removeCmd := &cmd.RemovePostFromRoadmap{
			PostID: postNumber,
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
		action := new(actions.CreateRoadmapColumn)
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
		columnID, err := c.ParamAsInt("id")
		if err != nil || columnID <= 0 {
			return c.BadRequest(web.Map{
				"error": "Invalid column ID",
			})
		}

		action := new(actions.UpdateRoadmapColumn)
		action.ColumnID = columnID
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
		columnID, err := c.ParamAsInt("id")
		if err != nil || columnID <= 0 {
			return c.BadRequest(web.Map{
				"error": "Invalid column ID",
			})
		}

		deleteCmd := &cmd.DeleteRoadmapColumn{
			ColumnID: columnID,
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
		action := new(actions.ReorderRoadmapColumns)
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

// ManageRoadmapColumn handles both PUT (update) and DELETE operations on a roadmap column
func ManageRoadmapColumn() web.HandlerFunc {
	return func(c *web.Context) error {
		columnID, err := c.ParamAsInt("id")
		if err != nil || columnID <= 0 {
			return c.BadRequest(web.Map{
				"error": "Invalid column ID",
			})
		}

		if c.Request.Method == "DELETE" {
			deleteCmd := &cmd.DeleteRoadmapColumn{
				ColumnID: columnID,
			}

			if err := bus.Dispatch(c, deleteCmd); err != nil {
				return c.Failure(err)
			}

			return c.Ok(web.Map{
				"success": true,
			})
		}

		// Handle PUT (update)
		action := new(actions.UpdateRoadmapColumn)
		action.ColumnID = columnID
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