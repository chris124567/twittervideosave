package sources

import (
	"encoding/json"
	htmlescape "html"
	"streamconvert/internal/pkg/htmlhelp"
	"streamconvert/internal/pkg/httphelp"
)

// TESTS := []string{"https://rumble.com/v6809n-23abc-news-latest-headlines-august-14-10pm.html", "https://rumble.com/v98z8b-3-news-now-latest-headlines-april-24-11am.html", "https://rumble.com/v6729h-denver-7-latest-headlines-august-9-10pm.html", "https://rumble.com/v3008m-speed-drawing-olaf.html", "https://rumble.com/v30j7e-rainbow-rose-cake.html", "https://rumble.com/v99c47-who-survived-covid-19-theres-no-guarantee-you-wont-get-it-again.html"}

const RUMBLE_API_URL string = "https://rumble.com/embedJS/u3/?request=video&v="

type rumble struct {
	ThumbnailUrl string `json:"i"`
	Title        string `json:"title"`
	DirectLink   string `json:"u"`
}

func RumbleGetVideo(link string) (Video, error) {
	var rumbleData rumble

	response, err := httphelp.HttpGetHeaders(link, httphelp.STANDARD_HEADERS)
	if err != nil {
		return EMPTY_VIDEO, err
	}
	defer response.Body.Close()

	responseBody := httphelp.GetBody(response)

	rumbleId := htmlhelp.GetStringInBetween(responseBody, `"video":"`, `"`)
	apiUrlFormat := RUMBLE_API_URL + rumbleId

	jsonResponse, err := httphelp.HttpGetHeaders(apiUrlFormat, httphelp.STANDARD_HEADERS)
	if err != nil {
		return EMPTY_VIDEO, err
	}
	defer response.Body.Close()

	jsonResponseBody := httphelp.GetBody(jsonResponse)
	err = json.Unmarshal([]byte(jsonResponseBody), &rumbleData)
	if err != nil {
		return EMPTY_VIDEO, err
	}

	links := []Link{{ThumbnailUrl: rumbleData.ThumbnailUrl, Quality: "480", Url: rumbleData.DirectLink}}
	return Video{Title: htmlescape.UnescapeString(rumbleData.Title), Links: links}, nil
}
