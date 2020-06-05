package web

import (
	"github.com/antchfx/htmlquery"
	"net/http"
	"net/url"
	"streamconvert/internal/pkg/htmlhelp"
	"streamconvert/internal/pkg/httphelp"
	"streamconvert/internal/pkg/miscutil"
)

const FREE_FILE_CONVERT_BASE string = "https://www.freefileconvert.com"

func Mp3Handler(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		GenericErrorHandler(writer, request, 500)
		return
	}
	videoUrl := request.Form.Get("url")

	urlParsed, err := url.Parse(videoUrl)
	if err != nil {
		GenericErrorHandler(writer, request, 500)
		return
	}

	urlParsed.Fragment = ""
	convertId := FreeFileConvert(urlParsed.String())

	if convertId == "" {
		GenericErrorHandler(writer, request, 500)
		return
	}

	Redirect(FREE_FILE_CONVERT_BASE+"/file/"+convertId, writer, request)
}

func FreeFileConvert(link string) string {
	var xsrfToken string
	var freefileconvertSession string

	headers := miscutil.CopyMap(httphelp.STANDARD_HEADERS)
	headers["Referer"] = FREE_FILE_CONVERT_BASE + "/mp4-mp3"

	response, err := httphelp.HttpGetHeaders(FREE_FILE_CONVERT_BASE, headers)
	if err != nil {
		return ""
	}

	doc, err := htmlquery.Parse(response.Body)
	if err != nil {
		return ""
	}
	defer response.Body.Close()

	phpSessionUploadProgress := htmlhelp.GetXpathValue(doc, "//input[@id='progress_key'][@name='PHP_SESSION_UPLOAD_PROGRESS'][@type='hidden'][@value]/@value")
	csrfToken := htmlhelp.GetXpathValue(doc, "//meta[@name='csrf-token'][@content]/@content")
	if csrfToken == "" || phpSessionUploadProgress == "" {
		return ""
	}

	for _, cookie := range response.Cookies() {
		if cookie.Name == "XSRF-TOKEN" {
			xsrfToken = cookie.Value
		}
		if cookie.Name == "freefileconvert_session" {
			freefileconvertSession = cookie.Value
		}
	}

	if xsrfToken == "" || freefileconvertSession == "" {
		return ""
	}

	var cookies = map[string]string{
		"cookieconsent_dismissed": "yes",
		"XSRF-TOKEN":              xsrfToken,
		"freefileconvert_session": freefileconvertSession,
	}

	headers["X-CSRF-TOKEN"] = csrfToken
	headers["X-Requested-With"] = "XMLHttpRequest"
	urlData := map[string]string{
		"_token":        csrfToken,
		"url":           link,
		"output_format": "mp3",
		"progress_key":  phpSessionUploadProgress,
	}

	urlResponse, err := httphelp.HttpPostMultipartHeadersCookies(FREE_FILE_CONVERT_BASE+"/file/url", headers, cookies, urlData)
	if err != nil {
		return ""
	}

	urlResponseBody := httphelp.GetBody(urlResponse)
	convertId := htmlhelp.GetStringInBetween(urlResponseBody, `"id":"`, `"`)

	if convertId == "" {
		return ""
	}

	return convertId
}
