package sources

import (
	"github.com/antchfx/htmlquery"
	htmlescape "html"
	"streamconvert/internal/pkg/htmlhelp"
	"streamconvert/internal/pkg/httphelp"
	"strings"
)

// TESTS := []string{"https://www.facebook.com/video.php?v=274175099429670", "https://www.facebook.com/cnn/videos/10155529876156509", "https://www.facebook.com/yaroslav.korpan/videos/1417995061575415/", "https://www.facebook.com/groups/1024490957622648/permalink/1396382447100162/", "https://zh-hk.facebook.com/peoplespower/videos/1135894589806027/", "https://www.facebook.com/WatchESLOne/videos/359649331226507/"}

var FB_MP4_TYPES = map[string]string{"contentUrl": QUALITY_UNKNOWN, "sd_src": "SD", "sd_src_no_ratelimit": "SD", "hd_src": "HD", "hd_src_no_ratelimit": "HD"}

func FacebookGetVideo(link string) (Video, error) {
	var links []Link
	var title string
	var thumbnailUrl string
	var directVideoLink string

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

	thumbnailUrl = htmlescape.UnescapeString(htmlhelp.GetXpathValue(doc, "//meta[@property='og:image'][@content]/@content"))
	if thumbnailUrl == "" {
		thumbnailUrl = htmlhelp.GetXpathValue(doc, "//div[@class='_3fnx _1jto _4x6d _8yzm _3htz']//img[@class='_1p6f _3fnw img'][@src][@alt]/@src")
	}
	title = htmlescape.UnescapeString(htmlhelp.GetXpathValue(doc, "//meta[@property='og:title'][@content]/@content"))
	if title == "" {
		title = htmlhelp.GetXpathValue(doc, "//div[@class='_8z4_']//div[@class='_8z50']//a[@target][@id][@href]/text()")
		if title == "" {
			title = htmlhelp.GetXpathValue(doc, "//title/text()")
		}
	}

	for linkType, quality := range FB_MP4_TYPES {
		// keep updating the link, so we get the highest quality link
		// sometimes the variables are quoted, sometimes they aren't
		directVideoLink = htmlhelp.GetStringInBetween(responseBody, linkType+":\"", "\"")
		if directVideoLink != "" {
			links = append(links, Link{Quality: quality, Url: htmlhelp.JsonUnescape(directVideoLink)}) // get rid of JSON escaping in URL
		} else {
			directVideoLink = htmlhelp.GetStringInBetween(responseBody, "\""+linkType+"\":\"", "\"")
			if directVideoLink != "" {
				links = append(links, Link{Quality: quality, Url: htmlhelp.JsonUnescape(directVideoLink), ThumbnailUrl: thumbnailUrl}) // get rid of JSON escaping in URL
			}
		}
	}

	return Video{Title: title, Links: links}, nil
}
