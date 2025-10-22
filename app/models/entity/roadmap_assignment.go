package entity

import "time"

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
