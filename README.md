# Source code for [twittervideosave.com](https://twittervideosave.com) and [fftwitter.com](https://fftwitter.com)

## Overview

### ```scripts```
Contains scripts for running/building the application and cleaning the binary.  These are all called in the makefile

### ```cmd```
Contains the code that calls all the web server route handlers.  I'm using Cloudflare for my sites, so I use a self-signed certificate and allow Cloudflare to handle the rest.  L59-61 of ```cmd/streamconvert/main.go``` is just calling the HTTP default mux to handle twittervideosave.com.  L62-73 handle ```fftwitter.com```, which is the redirect site (you can just put "ff" in front of any Tweet URL to be redirected to the video download page).  fftwitter.com blocks everything in robots.txt and has no sitemap.xml to prevent search engines from getting confused.

### ```internal/pkg```

#### ```internal/pkg/htmlhelp```
This is mostly just miscellaneous utilities that I found useful when parsing HTML with Xpath and just plain stringss.

#### ```internal/pkg/httphelp```
This just simplifies making requests and makes it so I don't have to set the headers manually on every request.  We impersonate a Chrome browser on Windows 10 to make the requests seem more legitimate.

#### ```internal/pkg/miscutil```
This just contains a function for copying a map.

#### ```internal/pkg/sources```
This contains all the scrapers for getting direct video links from a bunch of video sharing and social media platforms.  The site is styled after Twitter but really supports many more sources, including Coub, Rumble, Streamable, Twitter, Vidmax, Facebook, Instagram, and Worldstarhiphop.  You can change the targeting of the site by modifying the variables in ```web/templates/base_variables.html.tmpl``` 

### ```web```

#### ```web/app```
Basically just calls the functions in ```internal/pkg/sources``` and outputs the pretty templates in ```web/template```.  Also, there is ```web/app/mp3handler.go``` which leverages the (unofficial) API of freefileconvert.com to get free MP4 conversions to MP3 for basically no CPU cost.

#### ```web/static```
Static resources like bootstrap, css, and the favicon.

#### ```web/template```
Contains the HTML templates for the website.  ```base.html.tmpl``` and ```base_variables.html.tmpl``` are used in all the templates.