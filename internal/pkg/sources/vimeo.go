package sources

import (
	"encoding/json"
	"errors"
	"github.com/antchfx/htmlquery"
	htmlescape "html"
	"regexp"
	"strconv"
	"streamconvert/internal/pkg/htmlhelp"
	"streamconvert/internal/pkg/httphelp"
	"strings"
)

// TESTS := []string{"https://vimeo.com/20", "http://vimeo.com/56015672#at=0", "http://player.vimeo.com/video/54469442", "http://vimeo.com/channels/keypeele/75629013", "http://vimeo.com/76979871", "https://player.vimeo.com/video/98044508", "https://vimeo.com/33951933", "https://vimeo.com/channels/tributes/6213729 ***", "https://vimeo.com/groups/travelhd/videos/22439234", "https://vimeo.com/album/2632481/video/79010983", "https://vimeo.com/showcase/3373663/video/126543769", "https://vimeo.com/7809605", "https://vimeo.com/160743502/abd0e13fb4"}

var getIdRegex = regexp.MustCompile(`\d{1,}`)

const VIMEO_DOMAIN string = "vimeo.com"

// structs to deserialize player config json
type vimeoProgressive struct {
	Profile int    `json:"profile"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	Quality string `json:"quality"`
	URL     string `json:"url"`
}

type vimeoFiles struct {
	Progressive []vimeoProgressive `json:"progressive"`
}

type vimeoRequest struct {
	Files vimeoFiles `json:"files"`
}

type vimeoVideo struct {
	Title string `json:"title"`
}

type vimeo struct {
	Request vimeoRequest `json:"request"`
	Video   vimeoVideo   `json:"video"`
}

func VimeoGetVideo(link string) (Video, error) {
	var links []Link
	var vimeoData vimeo

	response, err := httphelp.HttpGetHeaders(link, httphelp.STANDARD_HEADERS)
	if err != nil {
		return EMPTY_VIDEO, err
	}

	responseBody := httphelp.GetBody(response)
	response.Body.Close()

	doc, err := htmlquery.Parse(strings.NewReader(responseBody))
	if err != nil {
		return EMPTY_VIDEO, err
	}

	canonicalUrl := htmlhelp.GetXpathValue(doc, "//link[@rel='canonical'][@href]/@href")
	if canonicalUrl == "" {
		return EMPTY_VIDEO, errors.New("Could not get canonical URL")
	}
	if canonicalUrl != response.Request.URL.String() {
		return VimeoGetVideo(canonicalUrl)
	}

	thumbnailUrl := htmlhelp.JsonUnescape(htmlhelp.GetXpathValue(doc, "//meta[@property='og:image'][@content]/@content"))
	title := htmlescape.UnescapeString(htmlhelp.GetXpathValue(doc, "//meta[@property='og:title'][@content]/@content"))

	configUrl := htmlhelp.JsonUnescape(htmlhelp.GetStringInBetween(responseBody, `{"config_url":"`, `"`))
	configResponse, err := httphelp.HttpGetHeaders(configUrl, httphelp.STANDARD_HEADERS)
	if err != nil {
		return EMPTY_VIDEO, err
	}
	configBody := httphelp.GetBody(configResponse)
	configResponse.Body.Close()

	err = json.Unmarshal([]byte(configBody), &vimeoData)
	if err != nil {
		return EMPTY_VIDEO, err
	}

	for _, video := range vimeoData.Request.Files.Progressive {
		links = append(links, Link{Quality: strconv.Itoa(video.Height), Url: video.URL, ThumbnailUrl: thumbnailUrl})
	}

	return Video{Title: title, Links: links}, nil
}
