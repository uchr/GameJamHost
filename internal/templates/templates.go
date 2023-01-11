package templates

import "html/template"

type Templates struct {
	JamListTemplate     *template.Template
	JamOverviewTemplate *template.Template
	JamEditFormTemplate *template.Template
	JamEntriesTemplate  *template.Template

	GameOverviewTemplate *template.Template
	GameEditFormTemplate *template.Template

	ErrorTemplate *template.Template
}

func NewTemplates(templateFolder string) (Templates, error) {
	templates := Templates{}

	var err error
	templates.JamListTemplate, err = template.ParseFiles(templateFolder+"/jam_list.html", templateFolder+"/base.html")
	if err != nil {
		return templates, err
	}

	templates.JamOverviewTemplate, err = template.ParseFiles(templateFolder+"/jam_overview.html", templateFolder+"/base.html")
	if err != nil {
		return templates, err
	}

	templates.JamEditFormTemplate, err = template.ParseFiles(templateFolder+"/jam_edit_form.html", templateFolder+"/base.html")
	if err != nil {
		return templates, err
	}

	templates.JamEntriesTemplate, err = template.ParseFiles(templateFolder+"/jam_entries.html", templateFolder+"/base.html")
	if err != nil {
		return templates, err
	}

	templates.GameEditFormTemplate, err = template.ParseFiles(templateFolder+"/game_edit_form.html", templateFolder+"/base.html")
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
