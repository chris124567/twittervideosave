package sources

import (
	"github.com/antchfx/htmlquery"
	htmlescape "html"
	"streamconvert/internal/pkg/htmlhelp"
	"streamconvert/internal/pkg/httphelp"
)

// TESTS := []string{"https://vidmax.com/video/195124-the-covid-world-order-arrives-at-the-mcdonald-s-drive-thru", "https://vidmax.com/video/195125-brazilian-soccer-matches-be-like", "https://vidmax.com/video/195118-when-this-kid-is-naughty-his-mother-pretends-to-call-trump-kid-pisses-himself", "https://vidmax.com/video/195087-baby-delivered-during-car-crash-gets-lost-cops-find-it", "https://www.vidmax.com/video/65911-And-now-the-dumbest-car-modification-EVER", "https://vidmax.com/video/56601-", "https://www.vidmax.com/video/58838-Middle-Eastern-Woman-Held-Down-and-Publicly-Spanked"}

func VidmaxGetVideo(link string) (Video, error) {
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
	directUrl := htmlhelp.GetXpathValue(doc, "//video[@id='thisPlayer']//source[@src][@type='video/mp4']/@src")
	links := []Link{{ThumbnailUrl: thumbnailUrl, Quality: QUALITY_UNKNOWN, Url: directUrl}}

	return Video{Title: title, Links: links}, nil
}
