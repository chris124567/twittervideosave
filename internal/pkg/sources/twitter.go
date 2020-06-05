package sources

import (
	"encoding/json"
	"errors"
	"github.com/antchfx/xmlquery"
	htmlescape "html"
	"net/url"
	"regexp"
	"streamconvert/internal/pkg/htmlhelp"
	"streamconvert/internal/pkg/httphelp"
	"streamconvert/internal/pkg/miscutil"
)

// TESTS := []string{"https://twitter.com/starwars/status/665052190608723968", "https://twitter.com/freethenipple/status/643211948184596480", "https://twitter.com/jaydingeer/status/700207533655363584", "https://twitter.com/Filmdrunk/status/713801302971588609", "https://twitter.com/captainamerica/status/719944021058060289", "https://twitter.com/OPP_HSD/status/779210622571536384", "https://twitter.com/news_al3alm/status/852138619213144067", "https://twitter.com/i/web/status/910031516746514432", "https://twitter.com/LisPower1/status/1001551623938805763", "https://twitter.com/Twitter/status/1087791357756956680", "https://twitter.com/ViviEducation/status/1136534865145286656"}

var getTwitterIdRegex = regexp.MustCompile(`(\d{15,})`) // get a 15 digit number in the URL. usernames are limited to 15 characters so this excludes the regex from catching usernames
var getResolutionRegex = regexp.MustCompile(`(([\d ]{2,5}[x][\d ]{2,5}))`)

const TWITTER_API_BASE string = "https://api.twitter.com/1.1/"
const TWITTER_AUTHORIZATION string = "Bearer AAAAAAAAAAAAAAAAAAAAAPYXBAAAAAAACLXUNDekMxqa8h%2F40K4moUkGsoc%3DTYfbDKbT3jJPCEVnMYqilB28NHfOPqkca3qaAxGfsyKCs0wRbw"
const TWITTER_GUEST_URL string = TWITTER_API_BASE + "guest/activate.json"

type guestTokenStruct struct {
	Guest_Token string `json:"guest_token"`
}

type twitterVideoConfig struct {
	Track struct {
		VmapURL string `json:"vmapUrl"`
	} `json:"track"`
}

// URLEntity represents a URL which has been parsed from text.
type URLEntity struct {
	// Indices     Indices `json:"indices"`
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	URL         string `json:"url"`
}

// MediaEntity represents media elements associated with a Tweet.
type MediaEntity struct {
	URLEntity
	ID                int64     `json:"id"`
	IDStr             string    `json:"id_str"`
	MediaURL          string    `json:"media_url"`
	MediaURLHttps     string    `json:"media_url_https"`
	SourceStatusID    int64     `json:"source_status_id"`
	SourceStatusIDStr string    `json:"source_status_id_str"`
	Type              string    `json:"type"`
	VideoInfo         VideoInfo `json:"video_info"`
}

// ExtendedEntity contains media information.
// https://dev.twitter.com/overview/api/entities-in-twitter-objects#extended_entities
type ExtendedEntity struct {
	Media []MediaEntity `json:"media"`
}

// VideoInfo is available on video media objects.
type VideoInfo struct {
	// AspectRatio    [2]int         `json:"aspect_ratio"`
	// DurationMillis int            `json:"duration_millis"`
	Variants []VideoVariant `json:"variants"`
}

// VideoVariant describes one of the available video formats.
type VideoVariant struct {
	ContentType string `json:"content_type"`
	// Bitrate     int    `json:"bitrate"`
	URL string `json:"url"`
}

type twitterShow struct {
	Text             string         `json:"full_text"`
	ExtendedEntities ExtendedEntity `json:"extended_entities"`
}

func TwitterGetVideo(link string) (Video, error) {
	var links []Link

	guestToken, personalizationId, guestId := getTwitterGuestToken()
	if guestToken == "" || personalizationId == "" || guestId == "" {
		return EMPTY_VIDEO, errors.New("Failed to get guest token and ID")
	}

	var cookies = map[string]string{
		"guest_id":           guestId,
		"personalization_id": personalizationId,
	}

	idMatches := getTwitterIdRegex.FindStringSubmatch(link)
	if len(idMatches) < 2 {
		return EMPTY_VIDEO, errors.New("Failed to get ID from URL")
	}
	videoId := idMatches[1]

	vmapUrl := getVmapUrl(guestToken, cookies, videoId)
	vmapUrls := getVmapMp4s(guestToken, cookies, vmapUrl)

	links = append(links, vmapUrls...)
	title, apiLinks := twitterGetTitle(guestToken, cookies, videoId)
	links = append(links, apiLinks...)

	return Video{Title: title, Links: links}, nil
}

func getVmapMp4s(guestToken string, cookies map[string]string, vmapUrl string) []Link {
	var links []Link

	headers := miscutil.CopyMap(httphelp.STANDARD_HEADERS)
	headers["Authorization"] = TWITTER_AUTHORIZATION
	headers["x-guest-token"] = guestToken
	response, err := httphelp.HttpGetHeadersCookies(vmapUrl, headers, cookies)
	if err != nil {
		return links
	}
	defer response.Body.Close()

	doc, err := xmlquery.Parse(response.Body)
	if err != nil {
		return links
	}

	xmlLinks := xmlquery.Find(doc, "//tw:videoVariant[@url][@content_type='video/mp4']")
	for _, xmlLink := range xmlLinks {
		link := htmlhelp.GetXpathXmlValue(xmlLink, "/@url")
		link, err = url.PathUnescape(link)
		if err != nil {
			continue
		}
		linkParsed, err := url.Parse(link)
		if err != nil {
			continue
		}
		linkParsed.RawQuery = ""
		link = linkParsed.String()

		links = append(links, Link{Url: link, Quality: getResolutionRegex.FindString(link)})
	}

	return links
}

func getVmapUrl(guestToken string, cookies map[string]string, videoId string) string {
	var videoConfigData twitterVideoConfig

	headers := miscutil.CopyMap(httphelp.STANDARD_HEADERS)
	headers["Authorization"] = TWITTER_AUTHORIZATION
	headers["x-guest-token"] = guestToken

	jsonUrl := TWITTER_API_BASE + "videos/tweet/config/" + videoId + ".json"
	response, err := httphelp.HttpGetHeadersCookies(jsonUrl, headers, cookies)
	if err != nil {
		return ""
	}

	responseBody := httphelp.GetBody(response)
	response.Body.Close()

	err = json.Unmarshal([]byte(responseBody), &videoConfigData)
	if err != nil {
		return ""
	}
	vmapUrl := videoConfigData.Track.VmapURL

	return vmapUrl
}

func twitterGetTitle(guestToken string, cookies map[string]string, videoId string) (string, []Link) {
	var tweetData twitterShow
	var links []Link

	headers := miscutil.CopyMap(httphelp.STANDARD_HEADERS)
	headers["Authorization"] = TWITTER_AUTHORIZATION
	headers["x-guest-token"] = guestToken
	response, err := httphelp.HttpGetHeadersCookies("https://api.twitter.com/1.1/statuses/show.json?tweet_mode=extended&id="+videoId, headers, cookies)
	if err != nil {
		return "", EMPTY_LINK_ARRAY
	}
	responseBody := httphelp.GetBody(response)
	response.Body.Close()

	err = json.Unmarshal([]byte(responseBody), &tweetData)
	if err != nil {
		return "", EMPTY_LINK_ARRAY
	}

	for _, media := range tweetData.ExtendedEntities.Media {
		for _, video := range media.VideoInfo.Variants {
			if video.ContentType == "video/mp4" {
				links = append(links, Link{Url: video.URL, Quality: getResolutionRegex.FindString(video.URL)})
			}
		}
	}
	return htmlescape.UnescapeString(tweetData.Text), links
}

func getTwitterGuestToken() (string, string, string) {
	var guestToken guestTokenStruct
	var personalization_id, guest_id string
	headers := miscutil.CopyMap(httphelp.STANDARD_HEADERS)
	headers["Authorization"] = TWITTER_AUTHORIZATION

	httphelp.HttpOptionsHeaders(TWITTER_GUEST_URL, headers)

	response, err := httphelp.HttpPostHeaders(TWITTER_GUEST_URL, headers, nil)
	if err != nil {
		return "", "", ""
	}
	responseBody := httphelp.GetBody(response)
	response.Body.Close()

	err = json.Unmarshal([]byte(responseBody), &guestToken)
	if err != nil {
		return "", "", ""
	}

	for _, cookie := range response.Cookies() {
		if cookie.Name == "personalization_id" {
			personalization_id = cookie.Value
		}
		if cookie.Name == "guest_id" {
			guest_id = cookie.Value
		}
	}

	return guestToken.Guest_Token, personalization_id, guest_id
}
