package sources

var MP3_HOSTS = []string{"coub.com", "streamable.com", "mobile.twitter.com", "twitter.com", "vimeo.com", "m.facebook.com", "facebook.com", "instagram.com", "m.worldstarhiphop.com", "worldstarhiphop.com"}

type Link struct {
	ThumbnailUrl string
	Quality      string
	Url          string
}

type Video struct {
	Title       string
	OriginalUrl string
	Source      string
	Links       []Link
	Mp3Support  bool
}

var EMPTY_VIDEO = Video{}
var EMPTY_LINK = Link{}
var EMPTY_VIDEO_ARRAY = []Video{}
var EMPTY_LINK_ARRAY = []Link{}
