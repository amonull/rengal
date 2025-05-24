package integration

import (
	"github.com/amonull/rengal/integration/anilist"
	"github.com/amonull/rengal/source"
)

// Integrator is the interface that wraps the basic integration methods.
type Integrator interface {
	// MarkRead marks a chapter as read
	MarkRead(chapter *source.Chapter) error
}

var (
	Anilist Integrator = anilist.New()
)
