package sources

import (
	"github.com/antchfx/htmlquery"
	htmlescape "html"
	"streamconvert/internal/pkg/htmlhelp"
	"streamconvert/internal/pkg/httphelp"
)

// TESTS := []string{"https://coub.com/view/2bshvx", "http://coub.com/view/5u5n1", "http://coub.com/view/237d5l5h"}

func CoubGetVideo(link string) (Video, error) {
	response, err := httphelp.HttpGetHeaders(link, httphelp.STANDARD_HEADERS)
	if err != nil {
		return EMPTY_VIDEO, err
	}
	defer response.Body.Close()

	doc, err := htmlquery.Parse(response.Body)
	if err != nil {
		return EMPTY_VIDEO, err
	}

	title := htmlescape.UnescapeString(htmlhelp.GetXpathValue(doc, "//meta[@property='og:title'][@content]/@content"))
	thumbnailUrl := htmlhelp.GetXpathValue(doc, "//meta[@property='og:image'][@content]/@content")

	jsonData := htmlhelp.GetXpathValue(doc, "//script[@id='coubPageCoubJson'][@type='text/json']")
	directUrl := htmlhelp.GetStringInBetween(jsonData, `"default":"`, `"`)
	return Video{Title: title, Links: []Link{{ThumbnailUrl: thumbnailUrl, Url: directUrl, Quality: "Default"}}}, nil
}
