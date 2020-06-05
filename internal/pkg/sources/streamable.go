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

// TESTS := []string{"https://streamable.com/uu486", "https://streamable.com/m/rendon-s-huge-2-run-homer", "https://streamable.com/dnd1", "https://streamable.com/moo", "https://streamable.com/e/dnd1", "https://streamable.com/s/okkqk/drxjds"}

var getStreamableVideoObjectRegex = regexp.MustCompile(`var videoObject = (.*);`)

type streamableVideo struct {
	Height int    `json:"height"`
	URL    string `json:"url"`
}

type streamableFiles struct {
	Mp4       streamableVideo `json:"mp4"`
	Mp4Mobile streamableVideo `json:"mp4-mobile"`
}

type streamable struct {
	Files streamableFiles `json:"files"`
	// Title   string   `json:"original_name"`
	// ThumbnailUrl   string   `json:"dynamic_thumbnail_url"`
}

func StreamableGetVideo(link string) (Video, error) {
	var links []Link
	var streamableData streamable

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

	title := htmlescape.UnescapeString(htmlhelp.GetXpathValue(doc, "//title/text()"))
	thumbnailUrl := httphelp.ResolveRelativeScheme(htmlhelp.GetXpathValue(doc, "//video[@poster]/@poster"))
	results := getStreamableVideoObjectRegex.FindStringSubmatch(responseBody)
	if len(results) < 2 {
		return EMPTY_VIDEO, errors.New("Could not find video JSON data")
	}
	jsonData := results[1]

	err = json.Unmarshal([]byte(jsonData), &streamableData)
	if err != nil {
		return EMPTY_VIDEO, err
	}

	if streamableData.Files.Mp4.URL != "" {
		streamableData.Files.Mp4.URL = httphelp.ResolveRelativeScheme(streamableData.Files.Mp4.URL)
		streamableData.Files.Mp4.URL = htmlescape.UnescapeString(streamableData.Files.Mp4.URL)
		links = append(links, Link{ThumbnailUrl: thumbnailUrl, Quality: strconv.Itoa(streamableData.Files.Mp4.Height), Url: streamableData.Files.Mp4.URL})
	}
	if streamableData.Files.Mp4Mobile.URL != "" {
		streamableData.Files.Mp4Mobile.URL = httphelp.ResolveRelativeScheme(streamableData.Files.Mp4Mobile.URL)
		streamableData.Files.Mp4Mobile.URL = htmlescape.UnescapeString(streamableData.Files.Mp4Mobile.URL)
		links = append(links, Link{ThumbnailUrl: thumbnailUrl, Quality: "Mobile", Url: streamableData.Files.Mp4Mobile.URL})
	}

	return Video{Title: title, Links: links}, nil
}
