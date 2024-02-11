package ui

import "io"

type Renderer interface {
	Render(w io.Writer, name string, data any) error
}
