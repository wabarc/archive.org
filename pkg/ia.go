package ia

import (
	"fmt"
	"log"
	"sync"

	"github.com/wabarc/helper"
)

// Wayback is the handle of saving webpages to archive.org
func (wbrc *Archiver) Wayback(links []string) (map[string]string, error) {
	collect, results := make(map[string]string), make(map[string]string)
	for _, link := range links {
		if !helper.IsURL(link) {
			log.Print(link + " is invalid url.")
			continue
		}
		collect[link] = link
	}

	ch := make(chan string, len(collect))
	defer close(ch)

	var mu sync.Mutex
	var wg sync.WaitGroup
	for link := range collect {
		wg.Add(1)
		go func(link string) {
			wbrc.fetch(link, ch)
			mu.Lock()
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
