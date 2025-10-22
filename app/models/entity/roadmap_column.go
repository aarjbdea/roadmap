package entity

import "time"

// RoadmapColumn represents a roadmap column
type RoadmapColumn struct {
	ID                int       `json:"id"`
	TenantID          int       `json:"tenantId"`
	Name              string    `json:"name"`
	Slug              string    `json:"slug"`
	Position          int       `json:"position"`
	IsVisibleToPublic bool      `json:"isVisibleToPublic"`
	CreatedAt         time.Time `json:"createdAt"`
	Posts             []*Post   `json:"posts"`
}

// RoadmapAssignment represents a post assignment to a roadmap column
type RoadmapAssignment struct {
	ID           int       `json:"id"`
	PostID       int       `json:"postId"`
	ColumnID     int       `json:"columnId"`
	TenantID     int       `json:"tenantId"`
	Position     int       `json:"position"`
	AssignedAt   time.Time `json:"assignedAt"`
	AssignedByID int       `json:"assignedById"`
}
