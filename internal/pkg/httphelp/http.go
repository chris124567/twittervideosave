package httphelp

import (
	"bytes"
	"errors"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var EMPTY_HTTP_RESPONSE = &http.Response{}

func HttpOptionsHeaders(link string, headers map[string]string) (*http.Response, error) {
	client := &http.Client{}
	request, err := http.NewRequest("OPTIONS", link, nil)
	if err != nil {
		return EMPTY_HTTP_RESPONSE, err
	}

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	response, err := client.Do(request)
	if err != nil {
		return EMPTY_HTTP_RESPONSE, err
	}

	return response, nil
}

func HttpPostHeaders(link string, headers map[string]string, data url.Values) (*http.Response, error) {
	client := &http.Client{}
	request, err := http.NewRequest("POST", link, strings.NewReader(data.Encode()))
	if err != nil {
		return EMPTY_HTTP_RESPONSE, err
	}

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	response, err := client.Do(request)
	if err != nil {
		return EMPTY_HTTP_RESPONSE, err
	}

	return response, nil
}

func HttpGetHeaders(link string, headers map[string]string) (*http.Response, error) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return EMPTY_HTTP_RESPONSE, err
	}

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	response, err := client.Do(request)
	if err != nil {
		return EMPTY_HTTP_RESPONSE, err
	}

	if response.StatusCode != 200 {
		return EMPTY_HTTP_RESPONSE, errors.New("Non 200 response.")
	}

	return response, nil
}

func HttpGetHeadersCookies(link string, headers map[string]string, cookies map[string]string) (*http.Response, error) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return EMPTY_HTTP_RESPONSE, err
	}

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	for key, value := range cookies {
		request.AddCookie(&http.Cookie{Name: key, Value: value})
	}

	response, err := client.Do(request)
	if err != nil {
		return EMPTY_HTTP_RESPONSE, err
	}

	if response.StatusCode != 200 {
		return EMPTY_HTTP_RESPONSE, errors.New("Got non 200 response: " + strconv.Itoa(response.StatusCode))
	}

	return response, nil
}

func HttpPostHeadersCookies(link string, headers map[string]string, cookies map[string]string, data url.Values) (*http.Response, error) {
	client := &http.Client{}
	request, err := http.NewRequest("POST", link, strings.NewReader(data.Encode()))
	if err != nil {
		return EMPTY_HTTP_RESPONSE, err
	}

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	for key, value := range cookies {
		request.AddCookie(&http.Cookie{Name: key, Value: value})
	}

	response, err := client.Do(request)
	if err != nil {
		return EMPTY_HTTP_RESPONSE, err
	}

	if response.StatusCode != 200 {
		return EMPTY_HTTP_RESPONSE, errors.New("Got non 200 response: " + strconv.Itoa(response.StatusCode))
	}

	return response, nil
}

func HttpPostMultipartHeadersCookies(link string, headers map[string]string, cookies map[string]string, data map[string]string) (*http.Response, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	for key, value := range data {
		_ = writer.WriteField(key, value)
	}

	request, err := http.NewRequest("POST", link, body)
	if err != nil {
		return EMPTY_HTTP_RESPONSE, err
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	for key, value := range cookies {
		request.AddCookie(&http.Cookie{Name: key, Value: value})
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return EMPTY_HTTP_RESPONSE, err
	}

	if response.StatusCode != 200 {
		return EMPTY_HTTP_RESPONSE, errors.New("Got non 200 response: " + strconv.Itoa(response.StatusCode))
	}

	return response, nil
}

func GetBody(response *http.Response) string {
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return ""
	}
	return string(bodyBytes)
}
