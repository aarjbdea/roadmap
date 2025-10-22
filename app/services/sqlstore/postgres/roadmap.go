package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/dbx"
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
	return using(ctx).GetAll(&q.Result, `
		SELECT id, tenant_id, name, slug, position, is_visible_to_public, created_at
		FROM roadmap_columns 
		WHERE tenant_id = $1 
		`+dbx.If(q.IncludePrivate, "", "AND is_visible_to_public = true")+`
		ORDER BY position ASC
	`, q.TenantID)
}

// GetRoadmapData returns roadmap columns with their assigned posts
func GetRoadmapData(ctx context.Context, q *query.GetRoadmapData) error {
	// First get all columns
	dbColumns := make([]*dbRoadmapColumn, 0)
	err := using(ctx).GetAll(&dbColumns, `
		SELECT id, tenant_id, name, slug, position, is_visible_to_public, created_at
		FROM roadmap_columns 
		WHERE tenant_id = $1 
		`+dbx.If(q.IncludePrivate, "", "AND is_visible_to_public = true")+`
		ORDER BY position ASC
	`, q.TenantID)
	if err != nil {
		return err
	}

	columns := make([]*entity.RoadmapColumn, len(dbColumns))
	for i, dbCol := range dbColumns {
		columns[i] = dbCol.toModel()
	}

	// Get posts for each column
	for _, column := range columns {
		dbPosts := make([]*dbPost, 0)
		err := using(ctx).GetAll(&dbPosts, `
			SELECT p.id, p.number, p.title, p.slug, p.description, p.created_at,
				   p.votes_count, p.comments_count, p.status,
				   u.id as user_id, u.name as user_name, u.email as user_email, u.role as user_role,
				   u.avatar_type, u.avatar_blob_key,
				   CASE WHEN v.post_id IS NOT NULL THEN true ELSE false END as has_voted,
				   ARRAY_REMOVE(ARRAY_AGG(t.slug), NULL) as tags
			FROM roadmap_post_assignments rpa
			INNER JOIN posts p ON p.id = rpa.post_id
			INNER JOIN users u ON u.id = p.user_id
			LEFT JOIN post_votes v ON v.post_id = p.id AND v.user_id = $2
			LEFT JOIN post_tags pt ON pt.post_id = p.id
			LEFT JOIN tags t ON t.id = pt.tag_id
			WHERE rpa.column_id = $3 AND rpa.tenant_id = $1
			GROUP BY p.id, p.number, p.title, p.slug, p.description, p.created_at,
					 p.votes_count, p.comments_count, p.status,
					 u.id, u.name, u.email, u.role, u.avatar_type, u.avatar_blob_key, v.post_id
			ORDER BY rpa.position ASC
		`, q.TenantID, getUserID(ctx), column.ID)
		if err != nil {
			return err
		}
		
		posts := make([]*entity.Post, len(dbPosts))
		for j, dbPost := range dbPosts {
			posts[j] = dbPost.toModel(ctx)
		}
		column.Posts = posts
	}

	q.Result = columns
	return nil
}

// GetPostRoadmapAssignment returns the roadmap assignment for a specific post
func GetPostRoadmapAssignment(ctx context.Context, q *query.GetPostRoadmapAssignment) error {
	assignment := &dbRoadmapAssignment{}
	err := using(ctx).GetFirst(assignment, `
		SELECT id, post_id, column_id, tenant_id, position, assigned_at, assigned_by_id
		FROM roadmap_post_assignments 
		WHERE post_id = $1 AND tenant_id = $2
	`, q.PostID, q.TenantID)
	if err != nil {
		if err == sql.ErrNoRows {
			q.Result = nil
			return nil
		}
		return err
	}
	q.Result = assignment.toModel()
	return nil
}

// AssignPostToColumn assigns a post to a roadmap column
func AssignPostToColumn(ctx context.Context, c *cmd.AssignPostToColumn) error {
	// First remove any existing assignment
	_, err := using(ctx).Exec(`
		DELETE FROM roadmap_post_assignments 
		WHERE post_id = $1 AND tenant_id = (SELECT tenant_id FROM posts WHERE id = $1)
	`, c.PostID)
	if err != nil {
		return err
	}

	// Get tenant_id from post
	var tenantID int
	err = using(ctx).GetFirst(&tenantID, `SELECT tenant_id FROM posts WHERE id = $1`, c.PostID)
	if err != nil {
		return err
	}

	// Insert new assignment
	assignment := &dbRoadmapAssignment{
		PostID:       c.PostID,
		ColumnID:     c.ColumnID,
		TenantID:     tenantID,
		Position:     c.Position,
		AssignedAt:   time.Now(),
		AssignedByID: c.AssignedByID,
	}

	err = using(ctx).GetFirst(assignment, `
		INSERT INTO roadmap_post_assignments (post_id, column_id, tenant_id, position, assigned_at, assigned_by_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, assignment.PostID, assignment.ColumnID, assignment.TenantID, assignment.Position, assignment.AssignedAt, assignment.AssignedByID)
	if err != nil {
		return err
	}

	c.Result = assignment.toModel()
	return nil
}

// RemovePostFromRoadmap removes a post from the roadmap
func RemovePostFromRoadmap(ctx context.Context, c *cmd.RemovePostFromRoadmap) error {
	_, err := using(ctx).Exec(`
		DELETE FROM roadmap_post_assignments 
		WHERE post_id = $1 AND tenant_id = $2
	`, c.PostID, c.TenantID)
	return err
}

// ReorderPostInColumn changes the position of a post within its column
func ReorderPostInColumn(ctx context.Context, c *cmd.ReorderPostInColumn) error {
	_, err := using(ctx).Exec(`
		UPDATE roadmap_post_assignments 
		SET position = $1
		WHERE post_id = $2
	`, c.NewPosition, c.PostID)
	return err
}

// CreateRoadmapColumn creates a new roadmap column
func CreateRoadmapColumn(ctx context.Context, c *cmd.CreateRoadmapColumn) error {
	column := &dbRoadmapColumn{
		TenantID:          c.TenantID,
		Name:              c.Name,
		Slug:              c.Slug,
		Position:          c.Position,
		IsVisibleToPublic: c.IsVisibleToPublic,
		CreatedAt:         time.Now(),
	}

	err := using(ctx).GetFirst(column, `
		INSERT INTO roadmap_columns (tenant_id, name, slug, position, is_visible_to_public, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, column.TenantID, column.Name, column.Slug, column.Position, column.IsVisibleToPublic, column.CreatedAt)
	if err != nil {
		return err
	}

	c.Result = column.toModel()
	return nil
}

// UpdateRoadmapColumn updates an existing roadmap column
func UpdateRoadmapColumn(ctx context.Context, c *cmd.UpdateRoadmapColumn) error {
	column := &dbRoadmapColumn{}
	err := using(ctx).GetFirst(column, `
		UPDATE roadmap_columns 
		SET name = $1, is_visible_to_public = $2
		WHERE id = $3
		RETURNING id, tenant_id, name, slug, position, is_visible_to_public, created_at
	`, c.Name, c.IsVisibleToPublic, c.ColumnID)
	if err != nil {
		return err
	}

	c.Result = column.toModel()
	return nil
}

// DeleteRoadmapColumn deletes a roadmap column
func DeleteRoadmapColumn(ctx context.Context, c *cmd.DeleteRoadmapColumn) error {
	// First remove all assignments to this column
	_, err := using(ctx).Exec(`
		DELETE FROM roadmap_post_assignments 
		WHERE column_id = $1 AND tenant_id = $2
	`, c.ColumnID, c.TenantID)
	if err != nil {
		return err
	}

	// Then delete the column
	_, err = using(ctx).Exec(`
		DELETE FROM roadmap_columns 
		WHERE id = $1 AND tenant_id = $2
	`, c.ColumnID, c.TenantID)
	return err
}

// ReorderRoadmapColumns changes the order of roadmap columns
func ReorderRoadmapColumns(ctx context.Context, c *cmd.ReorderRoadmapColumns) error {
	for i, columnID := range c.ColumnIDs {
		_, err := using(ctx).Exec(`
			UPDATE roadmap_columns 
			SET position = $1
			WHERE id = $2 AND tenant_id = $3
		`, i, columnID, c.TenantID)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetMaxRoadmapColumnPosition returns the maximum position for roadmap columns
func GetMaxRoadmapColumnPosition(ctx context.Context, q *query.GetMaxRoadmapColumnPosition) error {
	var maxPos int
	err := using(ctx).GetFirst(&maxPos, `
		SELECT COALESCE(MAX(position), 0) 
		FROM roadmap_columns 
		WHERE tenant_id = $1
	`, q.TenantID)
	if err != nil {
		return err
	}
	*q.Result = maxPos
	return nil
}

func init() {
	bus.AddHandler(GetRoadmapColumns)
	bus.AddHandler(GetRoadmapData)
	bus.AddHandler(GetPostRoadmapAssignment)
	bus.AddHandler(GetMaxRoadmapColumnPosition)
	bus.AddHandler(AssignPostToColumn)
	bus.AddHandler(RemovePostFromRoadmap)
	bus.AddHandler(ReorderPostInColumn)
	bus.AddHandler(CreateRoadmapColumn)
	bus.AddHandler(UpdateRoadmapColumn)
	bus.AddHandler(DeleteRoadmapColumn)
	bus.AddHandler(ReorderRoadmapColumns)
}
