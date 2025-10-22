package handlers

import (
	"github.com/getfider/fider/app/pkg/web"
	"net/http"
)

// RoadmapPage renders the roadmap page
func RoadmapPage() web.HandlerFunc {
	return func(c *web.Context) error {
		return c.Page(http.StatusOK, web.Props{
			Page:  "Roadmap.page",
			Title: "Roadmap",
		})
	}
}
