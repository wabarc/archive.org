package ia

import (
	"fmt"
	"net/http"
	"regexp"
	"time"
)

type Archiver struct {
}

const (
	userAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Safari/537.36"
	timeout   = time.Duration(30) * time.Second
)

var (
	host = "archive.org"
	dest = "https://web." + host
	base = "https://web.archive.org/save/"
)

func (wbrc *Archiver) fetch(url string, ch chan<- string) {
	client := &http.Client{
		Timeout: timeout,
	}
	req, err := http.NewRequest("GET", base+url, nil)
	req.Header.Add("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}
	defer resp.Body.Close()

	var loc string
	loc = resp.Header.Get("Content-Location")

	if len(loc) > 0 {
		ch <- fmt.Sprintf("%v%v", dest, loc)
		return
	}

	loc = resp.Header.Get("Location")
	if len(loc) > 0 {
		ch <- fmt.Sprintf("%v%v", dest, loc)
		return
	}

	links := resp.Header.Get("Link")
	re := regexp.MustCompile(`(?m)http[s]?:\/\/web\.archive\.org/web/[-a-zA-Z0-9@:%_\+.~#?&//=]*`)
	if match := re.FindAllString(links, -1); len(match) > 0 {
		loc = match[len(match)-1]
		ch <- fmt.Sprintf("%v", loc)
		return
	}

	ch <- fmt.Sprintf("The Internet Archive: %v %v for url: %v", resp.StatusCode, http.StatusText(resp.StatusCode), base+url)
}

func isURL(str string) bool {
	re := regexp.MustCompile(`(http(s)?:\/\/.)?(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)`)
	match := re.FindAllString(str, -1)
	for _, el := range match {
		if len(el) > 2 {
			return true
		}
	}
	return false
}
