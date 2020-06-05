package sources

import (
	// "log"
	"encoding/json"
	"errors"
	"github.com/antchfx/htmlquery"
	htmlescape "html"
	"streamconvert/internal/pkg/htmlhelp"
	"streamconvert/internal/pkg/httphelp"
	"strings"
)

// TESTS := []string{"https://www.instagram.com/p/BlIka1ZFCNr", "https://instagram.com/p/aye83DjauH/?foo=bar#abc", "https://www.instagram.com/p/BA-pQFBG8HZ/?taken-by=britneyspears", "https://www.instagram.com/p/BQ0eAlwhDrw/", "https://www.instagram.com/p/-Cmh1cukG2/", "https://www.instagram.com/p/9o6LshA7zy/embed/", "https://www.instagram.com/tv/aye83DjauH/"}

type instagram struct {
	EntryData struct {
		PostPage []struct {
			Graphql struct {
				ShortcodeMedia struct {
					DisplayURL  string `json:"display_url"`
					VideoURL    string `json:"video_url"`
					EdgeSidecar struct {
						Edges []struct {
							Node struct {
								TypeName   string `json:"__typename"`
								DisplayURL string `json:"display_url"`
								VideoURL   string `json:"video_url"`
							} `json:"node"`
						} `json:"edges"`
					} `json:"edge_sidecar_to_children"`
				} `json:"shortcode_media"`
			} `json:"graphql"`
		} `json:"PostPage"`
	} `json:"entry_data"`
}

func InstagramGetVideo(link string) (Video, error) {
	var links []Link
	var instagramData instagram

	link = instagramGetCanonicalLink(link)
	response, err := httphelp.HttpGetHeaders(link, httphelp.STANDARD_HEADERS)
	if err != nil {
		return EMPTY_VIDEO, err
	}
	defer response.Body.Close()

	responseBody := httphelp.GetBody(response)

	doc, err := htmlquery.Parse(strings.NewReader(responseBody))
	if err != nil {
		return EMPTY_VIDEO, err
	}

	title := htmlescape.UnescapeString(htmlhelp.GetXpathValue(doc, "//title/text()"))

	dataStrings := htmlhelp.MatchOneOf(responseBody, `window\._sharedData\s*=\s*(.*);`)
	if dataStrings == nil || len(dataStrings) < 2 {
		return EMPTY_VIDEO, errors.New("Failed to get JSON data")
	}
	dataString := dataStrings[1]

	err = json.Unmarshal([]byte(dataString), &instagramData)
	if err != nil {
		return EMPTY_VIDEO, errors.New("Failed to parse JSON data")
	}

	if len(instagramData.EntryData.PostPage) < 1 {
		return EMPTY_VIDEO, errors.New("Invalid JSON data")
	}

	if instagramData.EntryData.PostPage[0].Graphql.ShortcodeMedia.VideoURL != "" {
		links = append(links, Link{Quality: QUALITY_UNKNOWN, Url: instagramData.EntryData.PostPage[0].Graphql.ShortcodeMedia.VideoURL, ThumbnailUrl: instagramData.EntryData.PostPage[0].Graphql.ShortcodeMedia.DisplayURL})
	}

	for _, edge := range instagramData.EntryData.PostPage[0].Graphql.ShortcodeMedia.EdgeSidecar.Edges {
		if edge.Node.TypeName == "GraphVideo" {
			links = append(links, Link{Quality: QUALITY_UNKNOWN, Url: edge.Node.VideoURL, ThumbnailUrl: edge.Node.DisplayURL})
		}
	}

	return Video{Title: title, Links: links}, nil
}

func instagramGetCanonicalLink(link string) string {
	link = strings.Replace(link, "/embed/", "", 1)
	link = strings.Replace(link, "/embed", "", 1)
	return link
}
