package ia

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

const (
	userAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Safari/537.36"
	timeout   = time.Duration(30) * time.Second
)

var (
	host = "archive.org"
	dest = "https://web." + host
	base = "https://web.archive.org/save/"
)

func fetch(url string, ch chan<- string) {
	start := time.Now()

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

	if resp.StatusCode != http.StatusOK {
		ch <- fmt.Sprintf("status code error: %d %s", resp.StatusCode, resp.Status)
		return
	}

	loc := resp.Header.Get("Content-Location")

	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}

	secs := time.Since(start).Seconds()
	fmt.Printf("%.2fs %7d %s\n", secs, nbytes, url)

	ch <- fmt.Sprintf("%v%v", dest, loc)
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
