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
