package templatemanager

import (
	"io"
)

type TemplateManager interface {
	Render(wr io.Writer, name string, data interface{})
	RenderError(wr io.Writer, errorCode int, err error)
}
