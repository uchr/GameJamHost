package templatemanager

import (
	"fmt"
	"html/template"
	"io"
	"path/filepath"

	"GameJamPlatform/internal/log"
	"GameJamPlatform/internal/web/pagedata"
)

type templateManager struct {
	templates map[string]*template.Template

	templateFolder string
}

var _ = (TemplateManager)(&templateManager{})

func NewManager(templateFolder string) (*templateManager, error) {
	tm := templateManager{
		templates:      make(map[string]*template.Template),
		templateFolder: templateFolder,
	}

	if err := tm.loadTemplate("index"); err != nil {
		return nil, err
	}
	if err := tm.loadTemplate("error"); err != nil {
		return nil, err
	}
	if err := tm.loadTemplate("jam_list"); err != nil {
		return nil, err
	}
	if err := tm.loadTemplate("jam_overview"); err != nil {
		return nil, err
	}
	if err := tm.loadTemplate("jam_edit_form"); err != nil {
		return nil, err
	}
	if err := tm.loadTemplate("jam_entries"); err != nil {
		return nil, err
	}
	if err := tm.loadTemplate("game_edit_form"); err != nil {
		return nil, err
	}
	if err := tm.loadTemplate("game_overview"); err != nil {
		return nil, err
	}
	if err := tm.loadTemplate("user_registration"); err != nil {
		return nil, err
	}
	if err := tm.loadTemplate("user_login"); err != nil {
		return nil, err
	}
	if err := tm.loadTemplate("user_profile"); err != nil {
		return nil, err
	}
	if err := tm.loadTemplate("user_edit_form"); err != nil {
		return nil, err
	}

	return &tm, nil
}

func (tm *templateManager) Render(wr io.Writer, name string, data interface{}) {
	t, ok := tm.templates[name]
	if !ok {
		log.Error(fmt.Errorf("template %s not found", name))
	}

	err := t.ExecuteTemplate(wr, "base", data)
	if err != nil {
		log.Error(err)
	}
}

func (tm *templateManager) RenderError(wr io.Writer, errorCode int, err error) {
	if err != nil {
		log.Error(err)
	}
	tm.Render(wr, "error", pagedata.NewErrorPageData(errorCode, err))
}

func (tm *templateManager) loadTemplate(templateName string) error {
	templateBase := filepath.Join(tm.templateFolder, "base.gohtml")
	templateFile := filepath.Join(tm.templateFolder, templateName+".gohtml")
	t, err := template.ParseFiles(templateFile, templateBase)
	if err != nil {
		return err
	}

	tm.templates[templateName] = t
	return nil
}
