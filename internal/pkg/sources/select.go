package sources

import (
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"streamconvert/internal/pkg/htmlhelp"
	"strings"
)

var getNumberRegexp = regexp.MustCompile(`\d+`)

func AddVideoInformation(host string, link string, video Video) Video {
	video.Mp3Support = htmlhelp.StringInSlice(MP3_HOSTS, strings.Replace(host, "www.", "", 1))
	video.OriginalUrl = link
	video.Source = host
	sort.Slice(video.Links, func(i, j int) bool {
		firstQualityString := getNumberRegexp.FindString(video.Links[i].Quality)
		firstQualityInt, err := strconv.Atoi(firstQualityString)
		if err != nil {
			return false
		}
		secondQualityString := getNumberRegexp.FindString(video.Links[j].Quality)
		secondQualityInt, err := strconv.Atoi(secondQualityString)
		if err != nil {
			return false
		}
		return firstQualityInt < secondQualityInt // make lowest quality first
	})
	return video
}

func DynamicGetVideo(link string) (Video, error) {
	var err error
	var host string
	var video Video

	if !strings.HasPrefix(link, "http") {
		link = "https://" + link
	}
	urlParsed, err := url.Parse(link)
	if err != nil { // allow us to use twitter fallback
		host = link
	} else {
		urlParsed.Fragment = ""
		link = urlParsed.String()
		host = urlParsed.Host		
	}

	switch {
	case strings.Contains(host, "coub.com"):
		video, err = CoubGetVideo(link)
	case strings.Contains(host, "rumble.com"):
		video, err = RumbleGetVideo(link)
	case strings.Contains(host, "streamable.com"):
		video, err = StreamableGetVideo(link)
	case strings.Contains(host, "twitter.com"):
		video, err = TwitterGetVideo(link)
	case strings.Contains(host, "vidmax.com"):
		video, err = VidmaxGetVideo(link)
	// case strings.Contains(host, "vimeo.com"): // Vimeo blocks Vultr IPs
	// video, err = VimeoGetVideo(link)
	// case strings.Contains(host, "vk.com"): // VK urls are IP specific
	// video, err = VkontakteGetVideo(link)
	case strings.Contains(host, "facebook.com"):
		video, err = FacebookGetVideo(link)
	case strings.Contains(host, "instagram.com"):
		video, err = InstagramGetVideo(link)
	case strings.Contains(host, "worldstarhiphop.com"):
		video, err = WorldStarHipHopGetVideo(link)
	default:
		video, err = TwitterGetVideo(link)
	}
	if err != nil {
		return EMPTY_VIDEO, err
	}
	video = AddVideoInformation(host, link, video)

	return video, nil
}
