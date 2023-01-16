package templates

import "html/template"

type Templates struct {
	IndexTemplate *template.Template

	JamListTemplate     *template.Template
	JamOverviewTemplate *template.Template
	JamEditFormTemplate *template.Template
	JamEntriesTemplate  *template.Template

	GameOverviewTemplate *template.Template
	GameEditFormTemplate *template.Template

	UserRegistrationTemplate *template.Template
	UserLoginTemplate        *template.Template
	UserEditFormTemplate     *template.Template
	UserProfileTemplate      *template.Template

	ErrorTemplate *template.Template
}

func NewTemplates(templateFolder string) (Templates, error) {
	templates := Templates{}

	var err error
	templates.IndexTemplate, err = template.ParseFiles(templateFolder+"/index.gohtml", templateFolder+"/base.gohtml")
	if err != nil {
		return templates, err
	}

	templates.JamListTemplate, err = template.ParseFiles(templateFolder+"/jam_list.gohtml", templateFolder+"/base.gohtml")
	if err != nil {
		return templates, err
	}

	templates.JamOverviewTemplate, err = template.ParseFiles(templateFolder+"/jam_overview.gohtml", templateFolder+"/base.gohtml")
	if err != nil {
		return templates, err
	}

	templates.JamEditFormTemplate, err = template.ParseFiles(templateFolder+"/jam_edit_form.gohtml", templateFolder+"/base.gohtml")
	if err != nil {
		return templates, err
	}

	templates.JamEntriesTemplate, err = template.ParseFiles(templateFolder+"/jam_entries.gohtml", templateFolder+"/base.gohtml")
	if err != nil {
		return templates, err
	}

	templates.GameEditFormTemplate, err = template.ParseFiles(templateFolder+"/game_edit_form.gohtml", templateFolder+"/base.gohtml")
	if err != nil {
		return templates, err
	}

	templates.GameOverviewTemplate, err = template.ParseFiles(templateFolder+"/game_overview.gohtml", templateFolder+"/base.gohtml")
	if err != nil {
		return templates, err
	}

	templates.UserRegistrationTemplate, err = template.ParseFiles(templateFolder+"/user_registration.gohtml", templateFolder+"/base.gohtml")
	if err != nil {
		return templates, err
	}

	templates.UserLoginTemplate, err = template.ParseFiles(templateFolder+"/user_login.gohtml", templateFolder+"/base.gohtml")
	if err != nil {
		return templates, err
	}

	templates.ErrorTemplate, err = template.ParseFiles(templateFolder+"/error.gohtml", templateFolder+"/base.gohtml")
	if err != nil {
		return templates, err
	}

	return templates, nil
}
