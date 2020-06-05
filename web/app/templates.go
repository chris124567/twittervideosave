package web

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"streamconvert/internal/pkg/sources"
	texttemplate "text/template"
)

var templateFuncMap = template.FuncMap{
	// The name "inc" is what the function will be called in the template text.
	"inc": func(i int) int {
		return i + 1
	},
	"noescape": func(str string) template.HTML {
		return template.HTML(str)
	},
}

func WriteTemplateData(writer http.ResponseWriter, request *http.Request, path string, data sources.Video) error {
	// Initialize a slice containing the paths to the two files. Note that
	// home.page.tmpl file must be the *first* file in the slice.
	files := []string{
		path,
		"./web/template/base.html.tmpl",
		"./web/template/base_variables.html.tmpl",
	}

	// Use the template.ParseFiles() function to read the files and store the
	// templates in a template set. Notice that we can pass the slice of file paths
	// as a variadic parameter?
	templateName := filepath.Base(path)
	template, err := template.New(templateName).Funcs(templateFuncMap).ParseFiles(files...)
	if err != nil {
		log.Fatal(err)
		return err
	}

	// write to template output to client
	err = template.ExecuteTemplate(writer, templateName, &data)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func WriteTemplate(writer http.ResponseWriter, request *http.Request, path string) error {
	// Initialize a slice containing the paths to the two files. Note that
	// home.page.tmpl file must be the *first* file in the slice.
	files := []string{
		path,
		"./web/template/base.html.tmpl",
		"./web/template/base_variables.html.tmpl",
	}

	// Use the template.ParseFiles() function to read the files and store the
	// templates in a template set. Notice that we can pass the slice of file paths
	// as a variadic parameter?
	templates, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatal(err)
		return err
	}
	templates = templates.Funcs(templateFuncMap)

	// write to template output to client
	err = templates.Execute(writer, nil)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func WriteTemplateNoBase(writer http.ResponseWriter, request *http.Request, path string) error {
	// Initialize a slice containing the paths to the two files. Note that
	// home.page.tmpl file must be the *first* file in the slice.
	files := []string{
		path,
		"./web/template/base_variables.html.tmpl",
	}

	// Use the template.ParseFiles() function to read the files and store the
	// templates in a template set. Notice that we can pass the slice of file paths
	// as a variadic parameter?
	templates, err := texttemplate.ParseFiles(files...)
	if err != nil {
		log.Fatal(err)
		return err
	}
	// templates = templates.Funcs(templateFuncMap)

	// write to template output to client
	err = templates.Execute(writer, nil)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
