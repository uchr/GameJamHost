package validationerr

import "fmt"

type ValidationErrors struct {
	fieldErrors map[string]string
}

func (v ValidationErrors) Error() string {
	result := ""
	for key, message := range v.fieldErrors {
		result += fmt.Sprintf("%s: %s\n", key, message)
	}
	return result
}

func (v *ValidationErrors) Errors() map[string]string {
	if v == nil {
		return nil
	}

	return v.fieldErrors
}

func (v *ValidationErrors) Add(key, message string) {
	v.fieldErrors[key] = message
}

func (v ValidationErrors) HasErrors() bool {
	return len(v.fieldErrors) > 0
}

func New() ValidationErrors {
	return ValidationErrors{fieldErrors: make(map[string]string)}
}
