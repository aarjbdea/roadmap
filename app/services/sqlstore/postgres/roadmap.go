package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/gosimple/slug"
)

type dbRoadmapColumn struct {
	ID                int       `db:"id"`
	TenantID          int       `db:"tenant_id"`
	Name              string    `db:"name"`
	Slug              string    `db:"slug"`
	Position          int       `db:"position"`
	IsVisibleToPublic bool      `db:"is_visible_to_public"`
	CreatedAt         time.Time `db:"created_at"`
}

type dbRoadmapAssignment struct {
	ID           int       `db:"id"`
	PostID       int       `db:"post_id"`
	ColumnID     int       `db:"column_id"`
	TenantID     int       `db:"tenant_id"`
	Position     int       `db:"position"`
	AssignedAt   time.Time `db:"assigned_at"`
	AssignedByID int       `db:"assigned_by_id"`
}

func (r *dbRoadmapColumn) toModel() *entity.RoadmapColumn {
	return &entity.RoadmapColumn{
		ID:                r.ID,
		TenantID:          r.TenantID,
		Name:              r.Name,
		Slug:              r.Slug,
		Position:          r.Position,
		IsVisibleToPublic: r.IsVisibleToPublic,
		CreatedAt:         r.CreatedAt,
		Posts:             make([]*entity.Post, 0),
	}
}

func (r *dbRoadmapAssignment) toModel() *entity.RoadmapAssignment {
	return &entity.RoadmapAssignment{
		ID:           r.ID,
		PostID:       r.PostID,
		ColumnID:     r.ColumnID,
		TenantID:     r.TenantID,
		Position:     r.Position,
		AssignedAt:   r.AssignedAt,
		AssignedByID: r.AssignedByID,
	}
}

// GetRoadmapColumns returns all roadmap columns for a tenant
func GetRoadmapColumns(ctx context.Context, q *query.GetRoadmapColumns) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		dbColumns := make([]*dbRoadmapColumn, 0)
		query := `
			SELECT id, tenant_id, name, slug, position, is_visible_to_public, created_at
			FROM roadmap_columns
			WHERE tenant_id = $1
		`
		if !q.IncludePrivate {
			query += " AND is_visible_to_public = true"
		}
		query += " ORDER BY position ASC"
		
		err := trx.Select(&dbColumns, query, tenant.ID)
		if err != nil {
			return err
		}

		columns := make([]*entity.RoadmapColumn, len(dbColumns))
		for i, dbCol := range dbColumns {
			columns[i] = dbCol.toModel()
		}
		q.Result = columns
		return nil
	})
}

// GetRoadmapData returns roadmap columns with their assigned posts
func GetRoadmapData(ctx context.Context, q *query.GetRoadmapData) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		// Get all columns
		dbColumns := make([]*dbRoadmapColumn, 0)
		columnQuery := `
			SELECT id, tenant_id, name, slug, position, is_visible_to_public, created_at
			FROM roadmap_columns
			WHERE tenant_id = $1
		`
		if !q.IncludePrivate {
			columnQuery += " AND is_visible_to_public = true"
		}
		columnQuery += " ORDER BY position ASC"
		
		err := trx.Select(&dbColumns, columnQuery, tenant.ID)
		if err != nil {
			return err
		}

		columns := make([]*entity.RoadmapColumn, len(dbColumns))
		for i, dbCol := range dbColumns {
			columns[i] = dbCol.toModel()
			columns[i].Posts = make([]*entity.Post, 0)
		}

		// Get posts for each column
		for _, column := range columns {
			// Get post IDs for this column, ordered by position
			var postIDs []int
			err := trx.Select(&postIDs, `
				SELECT post_id
				FROM roadmap_post_assignments
				WHERE column_id = $1 AND tenant_id = $2
				ORDER BY position ASC
			`, column.ID, tenant.ID)
			if err != nil && err != sql.ErrNoRows {
				return err
			}

			// Get post details for each post ID
			for _, postID := range postIDs {
				post := &entity.Post{}
				err := trx.Get(post, `
					SELECT id, number, title, slug, description, created_at, status, votes_count, comments_count
					FROM posts
					WHERE id = $1
				`, postID)
				if err != nil {
					if err == sql.ErrNoRows {
						continue
					}
					return err
				}
				column.Posts = append(column.Posts, post)
			}
		}

		q.Result = columns
		return nil
	})
}

// GetPostRoadmapAssignment returns the roadmap assignment for a specific post
func GetPostRoadmapAssignment(ctx context.Context, q *query.GetPostRoadmapAssignment) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		assignment := &dbRoadmapAssignment{}
		err := trx.Get(assignment, `
			SELECT id, post_id, column_id, tenant_id, position, assigned_at, assigned_by_id
			FROM roadmap_post_assignments
			WHERE post_id = $1 AND tenant_id = $2
		`, q.PostID, tenant.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				q.Result = nil
				return nil
			}
			return err
		}
		q.Result = assignment.toModel()
		return nil
	})
}

// AssignPostToColumn assigns a post to a roadmap column
func AssignPostToColumn(ctx context.Context, c *cmd.AssignPostToColumn) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		// First remove any existing assignment
		_, err := trx.Execute(`
			DELETE FROM roadmap_post_assignments
			WHERE post_id = $1 AND tenant_id = $2
		`, c.PostID, tenant.ID)
		if err != nil {
			return err
		}

		// Insert new assignment
		assignment := &dbRoadmapAssignment{
			PostID:       c.PostID,
			ColumnID:     c.ColumnID,
			TenantID:     tenant.ID,
			Position:     c.Position,
			AssignedAt:   time.Now(),
			AssignedByID: c.AssignedByID,
		}

		err = trx.Get(assignment, `
			INSERT INTO roadmap_post_assignments (post_id, column_id, tenant_id, position, assigned_at, assigned_by_id)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id
		`, assignment.PostID, assignment.ColumnID, assignment.TenantID, assignment.Position, assignment.AssignedAt, assignment.AssignedByID)
		if err != nil {
			return err
		}

		c.Result = assignment.toModel()
		return nil
	})
}

// RemovePostFromRoadmap removes a post from the roadmap
func RemovePostFromRoadmap(ctx context.Context, c *cmd.RemovePostFromRoadmap) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		_, err := trx.Execute(`
			DELETE FROM roadmap_post_assignments
			WHERE post_id = $1 AND tenant_id = $2
		`, c.PostID, tenant.ID)
		return err
	})
}

// ReorderPostInColumn changes the position of a post within its column
func ReorderPostInColumn(ctx context.Context, c *cmd.ReorderPostInColumn) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		_, err := trx.Execute(`
			UPDATE roadmap_post_assignments
			SET position = $1
			WHERE post_id = $2
		`, c.NewPosition, c.PostID)
		return err
	})
}

// CreateRoadmapColumn creates a new roadmap column
func CreateRoadmapColumn(ctx context.Context, c *cmd.CreateRoadmapColumn) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		column := &dbRoadmapColumn{
			TenantID:          tenant.ID,
			Name:              c.Name,
			Slug:              c.Slug,
			Position:          c.Position,
			IsVisibleToPublic: c.IsVisibleToPublic,
			CreatedAt:         time.Now(),
		}

		err := trx.Get(column, `
			INSERT INTO roadmap_columns (tenant_id, name, slug, position, is_visible_to_public, created_at)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id
		`, column.TenantID, column.Name, column.Slug, column.Position, column.IsVisibleToPublic, column.CreatedAt)
		if err != nil {
			return err
		}

		c.Result = column.toModel()
		return nil
	})
}

// UpdateRoadmapColumn updates an existing roadmap column
func UpdateRoadmapColumn(ctx context.Context, c *cmd.UpdateRoadmapColumn) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		column := &dbRoadmapColumn{}
		err := trx.Get(column, `
			UPDATE roadmap_columns 
			SET name = $1, is_visible_to_public = $2
			WHERE id = $3 AND tenant_id = $4
			RETURNING id, tenant_id, name, slug, position, is_visible_to_public, created_at
		`, c.Name, c.IsVisibleToPublic, c.ColumnID, tenant.ID)
		if err != nil {
			return err
		}

		c.Result = column.toModel()
		return nil
	})
}

// DeleteRoadmapColumn deletes a roadmap column
func DeleteRoadmapColumn(ctx context.Context, c *cmd.DeleteRoadmapColumn) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		// First remove all assignments to this column
		_, err := trx.Execute(`
			DELETE FROM roadmap_post_assignments 
			WHERE column_id = $1 AND tenant_id = $2
		`, c.ColumnID, tenant.ID)
		if err != nil {
			return err
		}

		// Then delete the column
		_, err = trx.Execute(`
			DELETE FROM roadmap_columns 
			WHERE id = $1 AND tenant_id = $2
		`, c.ColumnID, tenant.ID)
		return err
	})
}

// ReorderRoadmapColumns changes the order of roadmap columns
func ReorderRoadmapColumns(ctx context.Context, c *cmd.ReorderRoadmapColumns) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		for i, columnID := range c.ColumnIDs {
			_, err := trx.Execute(`
				UPDATE roadmap_columns
				SET position = $1
				WHERE id = $2 AND tenant_id = $3
			`, i, columnID, tenant.ID)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// GetMaxRoadmapColumnPosition returns the maximum position for roadmap columns
func GetMaxRoadmapColumnPosition(ctx context.Context, q *query.GetMaxRoadmapColumnPosition) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		var maxPos int
		err := trx.Scalar(&maxPos, `
			SELECT COALESCE(MAX(position), 0)
			FROM roadmap_columns
			WHERE tenant_id = $1
		`, tenant.ID)
		if err != nil {
			return err
		}
		*q.Result = maxPos
		return nil
	})
}

// Action handlers that convert actions to commands

func handleCreateRoadmapColumnAction(ctx context.Context, action *actions.CreateRoadmapColumn) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		// Get the next position
		getMaxPos := &query.GetMaxRoadmapColumnPosition{
			TenantID: tenant.ID,
			Result:   new(int),
		}
		if err := GetMaxRoadmapColumnPosition(ctx, getMaxPos); err != nil {
			return err
		}

		// Generate slug from name
		slug := slug.Make(action.Name)

		// Create the command
		createCmd := &cmd.CreateRoadmapColumn{
			TenantID:          tenant.ID,
			Name:              action.Name,
			Slug:              slug,
			Position:          *getMaxPos.Result + 1,
			IsVisibleToPublic: action.IsVisibleToPublic,
			CreatedByID:       user.ID,
		}

		return CreateRoadmapColumn(ctx, createCmd)
	})
}

func handleUpdateRoadmapColumnAction(ctx context.Context, action *actions.UpdateRoadmapColumn) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		updateCmd := &cmd.UpdateRoadmapColumn{
			ColumnID:          action.ColumnID,
			Name:              action.Name,
			IsVisibleToPublic: action.IsVisibleToPublic,
			UpdatedByID:       user.ID,
		}

		return UpdateRoadmapColumn(ctx, updateCmd)
	})
}

func handleDeleteRoadmapColumnAction(ctx context.Context, action *actions.DeleteRoadmapColumn) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		deleteCmd := &cmd.DeleteRoadmapColumn{
			ColumnID:    action.ColumnID,
			TenantID:    tenant.ID,
			DeletedByID: user.ID,
		}

		return DeleteRoadmapColumn(ctx, deleteCmd)
	})
}

func handleReorderRoadmapColumnsAction(ctx context.Context, action *actions.ReorderRoadmapColumns) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		reorderCmd := &cmd.ReorderRoadmapColumns{
			TenantID:    tenant.ID,
			ColumnIDs:   action.ColumnIDs,
			UpdatedByID: user.ID,
		}

		return ReorderRoadmapColumns(ctx, reorderCmd)
	})
}

func handleAssignPostToRoadmapAction(ctx context.Context, action *actions.AssignPostToRoadmap) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		assignCmd := &cmd.AssignPostToColumn{
			PostID:       action.PostID,
			ColumnID:     action.ColumnID,
			Position:     action.Position,
			AssignedByID: user.ID,
		}

		return AssignPostToColumn(ctx, assignCmd)
	})
}

func handleReorderPostInRoadmapAction(ctx context.Context, action *actions.ReorderPostInRoadmap) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		reorderCmd := &cmd.ReorderPostInColumn{
			PostID:      action.PostID,
			NewPosition: action.NewPosition,
			UpdatedByID: user.ID,
		}

		return ReorderPostInColumn(ctx, reorderCmd)
	})
}

func init() {
	// Query handlers
	bus.AddHandler(GetRoadmapColumns)
	bus.AddHandler(GetRoadmapData)
	bus.AddHandler(GetPostRoadmapAssignment)
	bus.AddHandler(GetMaxRoadmapColumnPosition)
	
	// Command handlers
	bus.AddHandler(AssignPostToColumn)
	bus.AddHandler(RemovePostFromRoadmap)
	bus.AddHandler(ReorderPostInColumn)
	bus.AddHandler(CreateRoadmapColumn)
	bus.AddHandler(UpdateRoadmapColumn)
	bus.AddHandler(DeleteRoadmapColumn)
	bus.AddHandler(ReorderRoadmapColumns)
	
	// Action handlers
	bus.AddHandler(handleCreateRoadmapColumnAction)
	bus.AddHandler(handleUpdateRoadmapColumnAction)
	bus.AddHandler(handleDeleteRoadmapColumnAction)
	bus.AddHandler(handleReorderRoadmapColumnsAction)
	bus.AddHandler(handleAssignPostToRoadmapAction)
	bus.AddHandler(handleReorderPostInRoadmapAction)
}
