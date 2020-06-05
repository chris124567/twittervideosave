package sources

import (
	"errors"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"streamconvert/internal/pkg/htmlhelp"
	"streamconvert/internal/pkg/httphelp"
	"strings"
)

// TESTS := []string{"http://vk.com/videos-77521?z=video-77521_162222515%2Fclub77521", "http://vk.com/video205387401_165548505", "https://vk.com/video_ext.php?oid=-77521&id=162222515&hash=87b046504ccd8bfa", "https://vk.com/video-140332_456239111", "https://vk.com/video205387401_164765225", "http://new.vk.com/video205387401_165548505"}

const VK_DOMAIN string = "vk.com"
const VK_MOBILE_DOMAIN string = "m.vk.com"
const VK_VIDEO_PATH string = "/video"

// var VK_MP4_TYPES = [...]string{"mp4_240", "url240", "mp4_360", "url360", "mp4_480", "url480", "mp4_720", "url720", "mp4_1080", "url1080", "postlive_mp4"}
var VK_MP4_TYPES = map[string]string{"url240": "240", "url360": "360", "url480": "480", "url720": "720", "url1080": "1080", "postlive_mp4": "Unknown"}

func VkontakteGetVideo(link string) (Video, error) {
	var links []Link
	var tempDirectVideoLink string

	response, err := httphelp.HttpGetHeaders(link, httphelp.STANDARD_HEADERS)
	if err != nil {
		return EMPTY_VIDEO, err
	}
	defer response.Body.Close()

	doc, err := htmlquery.Parse(response.Body)
	if err != nil {
		return EMPTY_VIDEO, err
	}

	canonicalLink := vkGetCanonicalLink(link, doc)
	if canonicalLink == "" {
		return EMPTY_VIDEO, errors.New("Could not get canonical URL")
	}

	canonicalResponse, err := httphelp.HttpGetHeaders(canonicalLink, httphelp.STANDARD_HEADERS)
	if err != nil {
		return EMPTY_VIDEO, err
	}
	canonicalResponseBody := httphelp.GetBody(canonicalResponse)
	canonicalResponse.Body.Close()

	canonicalDoc, err := htmlquery.Parse(strings.NewReader(canonicalResponseBody))
	if err != nil {
		return EMPTY_VIDEO, err
	}
	videoTitle := htmlhelp.GetXpathValue(canonicalDoc, "//meta[@property='og:title'][@content]/@content")
	thumbnailUrl := htmlhelp.GetXpathValue(canonicalDoc, "//span[@itemprop='thumbnail'][@itemscope][@itemtype='http://schema.org/ImageObject']//link[@itemprop='url'][@href]/@href")

	for linkType, quality := range VK_MP4_TYPES {
		// keep updating the link, so we get the highest quality link
		tempDirectVideoLink = htmlhelp.GetStringInBetween(canonicalResponseBody, "\""+linkType+"\":\"", "\"")
		if tempDirectVideoLink != "" {
			links = append(links, Link{Quality: quality, Url: htmlhelp.JsonUnescape(tempDirectVideoLink), ThumbnailUrl: thumbnailUrl}) // get rid of JSON escaping in URL
		}
	}

	return Video{Title: videoTitle, Links: links}, nil
}

func vkGetCanonicalLink(link string, node *html.Node) string {
	var videoId string
	var mobileLink string
	var actualLink string

	if !strings.Contains(link, "video_ext.php") {
		mobileLink = htmlhelp.GetXpathValue(node, "//link[@rel='alternate'][@media='only screen and (max-width: 640px)'][@href]/@href")
		if !strings.Contains(mobileLink, VK_MOBILE_DOMAIN) {
			actualLink = htmlhelp.GetXpathValue(node, "//link[@rel='canonical'][@href]/@href")
		} else {
			actualLink = strings.Replace(mobileLink, VK_MOBILE_DOMAIN, VK_DOMAIN, 1)
		}
	} else {
		videoId = strings.Replace(htmlhelp.GetXpathValue(node, "//div[@id='page_wrap']//div[@id][@class='video_box_wrap']/@id"), "video_box_wrap", "", 1)
		if videoId != "" {
			actualLink = "https://" + VK_DOMAIN + VK_VIDEO_PATH + videoId
		}
	}

	if !strings.Contains(actualLink, VK_VIDEO_PATH) { // if we don't have a video link
		return ""
	}
	return actualLink
}
