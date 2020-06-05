package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"streamconvert/web/app"
	"strings"
	"time"
)

// make sure to check base_variables.html.tmpl too!
const REDIRECT_DOMAIN = "fftwitter.com"
const SUPPORTED_SITE string = "Twitter"
const SUPPORTED_SITE_LOWER string = "twitter"
const SUPPORTED_SITE_DOMAIN string = "twitter.com"
const SITE_NAME string = "TwitterVideoSave"
const SITE_DOMAIN string = "TwitterVideoSave.com"
const SITE_DOMAIN_LOWER string = "twittervideosave.com"
const SITE_DOMAIN_HOST = "https://" + SITE_DOMAIN_LOWER

func main() {
	log.Print("Starting streamconvert")

	// urls := []string{"https://coub.com/view/2bshvx", "http://coub.com/view/5u5n1", "http://coub.com/view/237d5l5h"}
	// urls := []string{"https://rumble.com/v6809n-23abc-news-latest-headlines-august-14-10pm.html", "https://rumble.com/v98z8b-3-news-now-latest-headlines-april-24-11am.html", "https://rumble.com/v6729h-denver-7-latest-headlines-august-9-10pm.html", "https://rumble.com/v3008m-speed-drawing-olaf.html", "https://rumble.com/v30j7e-rainbow-rose-cake.html", "https://rumble.com/v99c47-who-survived-covid-19-theres-no-guarantee-you-wont-get-it-again.html"}
	// urls := []string{"https://streamable.com/uu486", "https://streamable.com/m/rendon-s-huge-2-run-homer", "https://streamable.com/dnd1", "https://streamable.com/moo", "https://streamable.com/e/dnd1", "https://streamable.com/s/okkqk/drxjds"}
	// urls := []string{"https://twitter.com/RenwaX23/status/1254827144603189252", "https://twitter.com/CollegeBoard/status/1254792622113263617", "https://twitter.com/starwars/status/665052190608723968", "https://twitter.com/freethenipple/status/643211948184596480", "https://twitter.com/jaydingeer/status/700207533655363584", "https://twitter.com/Filmdrunk/status/713801302971588609", "https://twitter.com/captainamerica/status/719944021058060289", "https://twitter.com/OPP_HSD/status/779210622571536384", "https://twitter.com/news_al3alm/status/852138619213144067", "https://twitter.com/i/web/status/910031516746514432", "https://twitter.com/LisPower1/status/1001551623938805763", "https://twitter.com/Twitter/status/1087791357756956680", "https://twitter.com/ViviEducation/status/1136534865145286656"}
	// urls := []string{"https://vidmax.com/video/195124-the-covid-world-order-arrives-at-the-mcdonald-s-drive-thru", "https://vidmax.com/video/195125-brazilian-soccer-matches-be-like", "https://vidmax.com/video/195118-when-this-kid-is-naughty-his-mother-pretends-to-call-trump-kid-pisses-himself", "https://vidmax.com/video/195087-baby-delivered-during-car-crash-gets-lost-cops-find-it", "https://www.vidmax.com/video/65911-And-now-the-dumbest-car-modification-EVER", "https://vidmax.com/video/56601-", "https://www.vidmax.com/video/58838-Middle-Eastern-Woman-Held-Down-and-Publicly-Spanked"}
	// urls := []string{"https://vimeo.com/20", "http://vimeo.com/56015672#at=0", "http://player.vimeo.com/video/54469442", "http://vimeo.com/channels/keypeele/75629013", "http://vimeo.com/76979871", "https://player.vimeo.com/video/98044508", "https://vimeo.com/33951933", "https://vimeo.com/channels/tributes/6213729 ***", "https://vimeo.com/groups/travelhd/videos/22439234", "https://vimeo.com/album/2632481/video/79010983", "https://vimeo.com/showcase/3373663/video/126543769", "https://vimeo.com/7809605", "https://vimeo.com/160743502/abd0e13fb4"}
	// urls := []string{"https://vk.com/videos-77521?z=video-77521_162222515%2Fclub77521", "http://vk.com/video205387401_165548505", "https://vk.com/video_ext.php?oid=-77521&id=162222515&hash=87b046504ccd8bfa", "https://vk.com/video-140332_456239111", "https://vk.com/video205387401_164765225", "http://new.vk.com/video205387401_165548505"}
	// urls := []string{"https://www.facebook.com/video.php?v=274175099429670", "https://www.facebook.com/cnn/videos/10155529876156509", "https://www.facebook.com/yaroslav.korpan/videos/1417995061575415/", "https://www.facebook.com/groups/1024490957622648/permalink/1396382447100162/", "https://zh-hk.facebook.com/peoplespower/videos/1135894589806027/", "https://www.facebook.com/WatchESLOne/videos/359649331226507/"}
	// urls := []string{"https://www.instagram.com/p/BlIka1ZFCNr", "https://instagram.com/p/aye83DjauH/?foo=bar#abc", "https://www.instagram.com/p/BA-pQFBG8HZ/?taken-by=britneyspears", "https://www.instagram.com/p/BQ0eAlwhDrw/", "https://www.instagram.com/p/-Cmh1cukG2/", "https://www.instagram.com/p/9o6LshA7zy/embed/", "https://www.instagram.com/tv/aye83DjauH/"}
	// urls := []string{"https://www.worldstarhiphop.com/videos/video.php?v=wshh6a7q1ny0G34ZwuIO" , "http://m.worldstarhiphop.com/android/video.php?v=wshh6a7q1ny0G34ZwuIO"}

	// for _, url := range urls {
	// 	video, err := sources.DynamicGetVideo(url)
	// 	log.Printf("%v\n%v\n", video, err)
	// }

	httpServer := &http.Server{
		Addr:         ":443",
		ReadTimeout:  90 * time.Second,  // getting URLs may take a while especially under heave load
		WriteTimeout: 25 * time.Second,  // should not take this long, just have to send a URL (~2000 bytes max)
		IdleTimeout:  240 * time.Second, // keep alive idle timeout

		Handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Set("Connection", "close") // cloudflare

			host, _, _ := net.SplitHostPort(request.Host)
			if host == "" {
				host = request.Host
			}
			host = strings.ToLower(host)

			switch host {
			case SITE_DOMAIN_LOWER, "www." + SITE_DOMAIN_LOWER:
				http.DefaultServeMux.ServeHTTP(writer, request)
				return
			case REDIRECT_DOMAIN, "www." + REDIRECT_DOMAIN:
				if request.URL.Path == "" || request.URL.Path == "/" || request.URL.Path == "/sitemap.xml" {
					http.Error(writer, http.StatusText(http.StatusNotFound), 404)
					return
				} else if request.URL.Path == "/robots.txt" {
					fmt.Fprintf(writer, "User-agent: *\nDisallow: *\n")
					return
				} else {
					urlRedirect := SITE_DOMAIN_HOST + "/video?url=" + SUPPORTED_SITE_DOMAIN + request.URL.Path
					http.Redirect(writer, request, urlRedirect, 302)
					return
				}
			default:
				http.Error(writer, http.StatusText(http.StatusNotFound), 404)
				return
			}

		}),
	}

	http.HandleFunc("/", web.HomeHandler)
	http.HandleFunc("/static/", web.StaticHandler)
	http.HandleFunc("/robots.txt", web.RobotsHandler)
	http.HandleFunc("/sitemap.xml", web.SitemapHandler)
	http.HandleFunc("/video", web.VideoHandler)
	http.HandleFunc("/terms-of-service", web.TermsOfServiceHandler)
	http.HandleFunc("/privacy-policy", web.PrivacyPolicyHandler)
	http.HandleFunc("/contact", web.ContactHandler)
	http.HandleFunc("/about", web.AboutHandler)
	http.HandleFunc("/frequently-asked-questions", web.FAQHandler)
	http.HandleFunc("/ios-shortcut", web.IosShortcutHandler)
	http.HandleFunc("/getmp3", web.Mp3Handler)

	// how tos
	http.HandleFunc("/how-to", web.HowToHandler)

	// production DISABLE
	// http.HandleFunc("/blank", web.BlankHandler)

	// go http.ListenAndServe(":80", certManager.HTTPHandler(nil))
	httpServer.ListenAndServeTLS("server.crt", "server.key")
}
