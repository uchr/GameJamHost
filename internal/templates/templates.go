package templates

import "html/template"

const TimeLayout = "2006-01-02T15:04"

type Templates struct {
	JamListTemplate     *template.Template
	JamNewTemplate      *template.Template
	JamOverviewTemplate *template.Template
	JamEditTemplate     *template.Template
	JamEntriesTemplate  *template.Template

	GameNewTemplate      *template.Template
	GameOverviewTemplate *template.Template
	GameEditTemplate     *template.Template

	ErrorTemplate *template.Template
}

func NewTemplates(templateFolder string) (Templates, error) {
	templates := Templates{}

	var err error
	templates.JamListTemplate, err = template.ParseFiles(templateFolder+"/jam_list.html", templateFolder+"/base.html")
	if err != nil {
		return templates, err
	}

	templates.JamNewTemplate, err = template.ParseFiles(templateFolder+"/jam_new.html", templateFolder+"/base.html")
	if err != nil {
		return templates, err
	}

	templates.JamOverviewTemplate, err = template.ParseFiles(templateFolder+"/jam_overview.html", templateFolder+"/base.html")
	if err != nil {
		return templates, err
	}

	templates.JamEditTemplate, err = template.ParseFiles(templateFolder+"/jam_edit.html", templateFolder+"/base.html")
	if err != nil {
		return templates, err
	}

	templates.JamEntriesTemplate, err = template.ParseFiles(templateFolder+"/jam_entries.html", templateFolder+"/base.html")
	if err != nil {
		return templates, err
	}

	templates.GameNewTemplate, err = template.ParseFiles(templateFolder+"/game_new.html", templateFolder+"/base.html")
	if err != nil {
		return templates, err
	}

	templates.GameEditTemplate, err = template.ParseFiles(templateFolder+"/game_edit.html", templateFolder+"/base.html")
	if err != nil {
		return templates, err
	}

	templates.GameOverviewTemplate, err = template.ParseFiles(templateFolder+"/game_overview.html", templateFolder+"/base.html")
	if err != nil {
		return templates, err
	}

	templates.ErrorTemplate, err = template.ParseFiles(templateFolder+"/error.html", templateFolder+"/base.html")
	if err != nil {
		return templates, err
	}

	return templates, nil
}
