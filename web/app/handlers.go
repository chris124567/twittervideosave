package web

import (
	"net/http"
	"streamconvert/internal/pkg/sources"
)

func VideoHandler(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		GenericErrorHandler(writer, request, 500)
	}
	videoUrl := request.Form.Get("url")

	video, err := sources.DynamicGetVideo(videoUrl)
	if err != nil {
		GenericErrorHandler(writer, request, templateError)
		return
	}

	err = WriteTemplateData(writer, request, "./web/template/video.html.tmpl", video)
	if err != nil {
		GenericErrorHandler(writer, request, templateError)
		return
	}
}

func HomeHandler(writer http.ResponseWriter, request *http.Request) {
	// if we don't have a root URL, return.  We have to do this because golang by default puts everything under "/" thats not handled separately under this function.
	if request.URL.Path != "/" {
		GenericErrorHandler(writer, request, 404)
		return
	}

	err := WriteTemplate(writer, request, "./web/template/home.html.tmpl")
	if err != nil {
		GenericErrorHandler(writer, request, templateError)
		return
	}
}

func TermsOfServiceHandler(writer http.ResponseWriter, request *http.Request) {
	err := WriteTemplate(writer, request, "./web/template/tos.html.tmpl")
	if err != nil {
		GenericErrorHandler(writer, request, templateError)
		return
	}
}

func PrivacyPolicyHandler(writer http.ResponseWriter, request *http.Request) {
	err := WriteTemplate(writer, request, "./web/template/privacypolicy.html.tmpl")
	if err != nil {
		GenericErrorHandler(writer, request, templateError)
		return
	}
}

func ContactHandler(writer http.ResponseWriter, request *http.Request) {
	err := WriteTemplate(writer, request, "./web/template/contact.html.tmpl")
	if err != nil {
		GenericErrorHandler(writer, request, templateError)
		return
	}
}

func AboutHandler(writer http.ResponseWriter, request *http.Request) {
	err := WriteTemplate(writer, request, "./web/template/about.html.tmpl")
	if err != nil {
		GenericErrorHandler(writer, request, templateError)
		return
	}
}

func FAQHandler(writer http.ResponseWriter, request *http.Request) {
	err := WriteTemplate(writer, request, "./web/template/faq.html.tmpl")
	if err != nil {
		GenericErrorHandler(writer, request, templateError)
		return
	}
}

func HowToHandler(writer http.ResponseWriter, request *http.Request) {
	err := WriteTemplate(writer, request, "./web/template/howto.html.tmpl")
	if err != nil {
		GenericErrorHandler(writer, request, templateError)
		return
	}
}

func IosShortcutHandler(writer http.ResponseWriter, request *http.Request) {
	err := WriteTemplate(writer, request, "./web/template/ios.html.tmpl")
	if err != nil {
		GenericErrorHandler(writer, request, templateError)
		return
	}
}

func RobotsHandler(writer http.ResponseWriter, request *http.Request) {
	err := WriteTemplateNoBase(writer, request, "./web/template/robots.txt.tmpl")
	if err != nil {
		GenericErrorHandler(writer, request, templateError)
		return
	}
}

func SitemapHandler(writer http.ResponseWriter, request *http.Request) {
	err := WriteTemplateNoBase(writer, request, "./web/template/sitemap.xml.tmpl")
	if err != nil {
		GenericErrorHandler(writer, request, templateError)
		return
	}
}
