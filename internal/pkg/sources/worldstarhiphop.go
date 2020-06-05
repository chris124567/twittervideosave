package sources

import (
	"errors"
	"github.com/antchfx/htmlquery"
	htmlescape "html"
	"net/url"
	"streamconvert/internal/pkg/htmlhelp"
	"streamconvert/internal/pkg/httphelp"
	"strings"
)

// TESTS := []string{"http://www.worldstarhiphop.com/videos/video.php?v=wshh6a7q1ny0G34ZwuIO", "http://m.worldstarhiphop.com/android/video.php?v=wshh6a7q1ny0G34ZwuIO"}

const WSHH_DOMAIN = "worldstarhiphop.com"
const WSHH_VIDEOS_PATH = "/videos/video.php?v="

func WorldStarHipHopGetVideo(link string) (Video, error) {
	var links []Link

	link = wshhGetCanonicalLink(link)
	if link == "" {
		return EMPTY_VIDEO, errors.New("Could not find canonical URL")
	}

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
	desktopUrl := htmlhelp.GetXpathValue(doc, "//video[@id='video-player']//source[@src][@type='video/mp4']/@src")
	mobileUrl := strings.Replace(strings.Replace(desktopUrl, "hw-videos", "hw-mobile", 1), ".mp4", "_mobile.mp4", 1)
	links = append(links, Link{ThumbnailUrl: thumbnailUrl, Quality: "Desktop", Url: desktopUrl})
	links = append(links, Link{ThumbnailUrl: thumbnailUrl, Quality: "Mobile", Url: mobileUrl})

	return Video{Title: title, Links: links}, nil
}

func wshhGetCanonicalLink(link string) string {
	urlParsed, err := url.Parse(link)
	if err != nil {
		return ""
	}
	videoId := urlParsed.Query().Get("v")
	if videoId == "" {
		return ""
	}
	return "https://" + WSHH_DOMAIN + "/" + WSHH_VIDEOS_PATH + videoId
}
