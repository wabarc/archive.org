package ia

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

type Archiver struct {
}

const (
	userAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Safari/537.36"
	timeout   = 120 * time.Second
)

var (
	host = "archive.org"
	dest = "https://web." + host
	base = "https://web.archive.org/save/"

	endpoint = "https://archive.org/wayback/available"
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

	loc = resp.Request.URL.String()
	if match := re.FindAllString(loc, -1); len(match) > 0 {
		ch <- fmt.Sprintf("%v", loc)
		return
	}

	got := latest(url)

	// HTTP 509 Bandwidth Limit Exceeded
	if resp.StatusCode == 509 {
		ch <- fmt.Sprint(got)
		return
	}

	if resp.StatusCode != 200 {
		ch <- fmt.Sprint(got)
		return
	}

	ch <- fmt.Sprintf("The Internet Archive: %v %v for url: %v", resp.StatusCode, http.StatusText(resp.StatusCode), base+url)
}

func isURL(str string) bool {
	re := regexp.MustCompile(`https?://?[-a-zA-Z0-9@:%._\+~#=]{1,255}\.[a-z]{0,63}\b(?:[-a-zA-Z0-9@:%_\+.~#?&//=]*)`)
	match := re.FindAllString(str, -1)

	return len(match) >= 1
}

func latest(s string) string {
	// https://web.archive.org/*/https://example.org
	u := fmt.Sprintf("%s/*/%s", dest, s)

	if _, err := url.Parse(s); err != nil {
		return u
	}

	endpoint += "?url=" + s
	resp, err := http.Get(endpoint)
	if err != nil {
		return u
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return u
	}

	var dat map[string]interface{}
	if err := json.Unmarshal(data, &dat); err != nil {
		return u
	}

	if archived, ok := dat["archived_snapshots"].(map[string]interface{}); ok {
		if closest, ok := archived["closest"].(map[string]interface{}); ok {
			if closest["available"].(bool) {
				return closest["url"].(string)
			}
		}
	}

	return u
}
