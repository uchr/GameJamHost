package templates

import (
	"html/template"

	"github.com/gomarkdown/markdown"
)

func renderContent(content string) template.HTML {
	return template.HTML(markdown.ToHTML(markdown.NormalizeNewlines([]byte(content)), nil, nil))
}
