package handlers

import (
	"github.com/getfider/fider/app/pkg/web"
)

// RoadmapPage renders the roadmap page
func RoadmapPage() web.HandlerFunc {
	return func(c *web.Context) error {
		return c.Page(web.Props{
			Title:       "Roadmap",
			Description: "View our product roadmap",
			Data: web.Map{
				"roadmap": true,
			},
		})
	}
}
