package ia

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"sync"
	"time"

	"github.com/wabarc/helper"
	"github.com/wabarc/logger"
)

type Archiver struct {
	Client *http.Client
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

// Wayback is the handle of saving webpages to archive.org
func (wbrc *Archiver) Wayback(links []string) (map[string]string, error) {
	collects, results := make(map[string]string), make(map[string]string)
	for _, link := range links {
		if !helper.IsURL(link) {
			logger.Info(link + " is invalid url.")
			continue
		}
		collects[link] = link
	}
	if wbrc.Client == nil {
		wbrc.Client = &http.Client{
			Timeout:       timeout,
			CheckRedirect: noRedirect,
		}
	}

	ch := make(chan string, len(collects))
	defer close(ch)

	var mu sync.Mutex
	var wg sync.WaitGroup
	for _, link := range collects {
		wg.Add(1)
		go func(link string) {
			mu.Lock()
			wbrc.archive(link, ch)
			results[link] = <-ch
			mu.Unlock()
			wg.Done()
		}(link)
	}
	wg.Wait()

	if len(results) == 0 {
		return results, fmt.Errorf("No results")
	}

	return results, nil
}

// Playback handle searching archived webpages from archive.is
func (wbrc *Archiver) Playback(links []string) (map[string]string, error) {
	collects, results := make(map[string]string), make(map[string]string)
	for _, link := range links {
		if !helper.IsURL(link) {
			logger.Info(link + " is invalid url.")
			continue
		}
		collects[link] = link
	}

	if wbrc.Client == nil {
		wbrc.Client = &http.Client{
			Timeout:       timeout,
			CheckRedirect: noRedirect,
		}
	}

	ch := make(chan string, len(collects))
	defer close(ch)

	var mu sync.Mutex
	var wg sync.WaitGroup
	for _, link := range collects {
		wg.Add(1)
		go func(link string) {
			mu.Lock()
			wbrc.search(link, ch)
			results[link] = <-ch
			mu.Unlock()
			wg.Done()
		}(link)
	}
	wg.Wait()

	if len(results) == 0 {
		return results, fmt.Errorf("No results")
	}

	return results, nil
}
func (wbrc *Archiver) archive(url string, ch chan<- string) {
	req, err := http.NewRequest("GET", base+url, nil)
	req.Header.Add("User-Agent", userAgent)
	resp, err := wbrc.Client.Do(req)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}
	defer resp.Body.Close()

	var loc string
	loc = resp.Header.Get("Content-Location")

	if len(loc) > 0 {
		ch <- loc
		return
	}

	loc = resp.Header.Get("Location")
	if len(loc) > 0 {
		ch <- loc
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

	got := wbrc.latest(url)

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

func (wbrc *Archiver) search(url string, ch chan<- string) {
	ch <- wbrc.latest(url)
}

func (wbrc *Archiver) latest(s string) string {
	// https://web.archive.org/*/https://example.org
	u := fmt.Sprintf("%s/*/%s", dest, s)

	if _, err := url.Parse(s); err != nil {
		return u
	}

	uri := endpoint + "?url=" + s
	resp, err := wbrc.Client.Get(uri)
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
			if closest["available"].(bool) && closest["status"] == "200" {
				return closest["url"].(string)
			}
		}
	}

	return u
}

func noRedirect(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}
