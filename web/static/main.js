function highlight_active() {
    // get current URL path and assign 'active' class
    var i;
    var link;
    var pathname = window.location.pathname;
    var navbar_links = document.getElementsByClassName("navbar_link");
    for (i = 0; i < navbar_links.length; i += 1) {
        link = navbar_links[i].getAttribute("href");
        if (link === pathname) {
            if (navbar_links[i].parentNode !== null) {
                navbar_links[i].parentNode.className = "active";
            }
        }
    }
}
window.highlight_active = highlight_active;
