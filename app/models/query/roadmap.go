package query

import (
	"github.com/getfider/fider/app/models/entity"
)

// GetRoadmapColumns returns all roadmap columns for a tenant
type GetRoadmapColumns struct {
	TenantID           int
	IncludePrivate     bool
	Result             []*entity.RoadmapColumn
}

// GetRoadmapData returns roadmap columns with their assigned posts
type GetRoadmapData struct {
	TenantID           int
	IncludePrivate     bool
	Result             []*entity.RoadmapColumn
}

// GetPostRoadmapAssignment returns the roadmap assignment for a specific post
type GetPostRoadmapAssignment struct {
	PostID   int
	TenantID int
	Result   *entity.RoadmapAssignment
}

// GetMaxRoadmapColumnPosition returns the maximum position for roadmap columns
type GetMaxRoadmapColumnPosition struct {
	TenantID int
	Result   *int
}
