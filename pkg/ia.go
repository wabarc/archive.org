package ia

import (
	"log"
	"sync"

	"github.com/wabarc/helper"
)

// Wayback is the handle of saving webpages to archive.org
func (wbrc *Archiver) Wayback(links []string) (map[string]string, error) {
	collect := make(map[string]string)
	for _, link := range links {
		if !helper.IsURL(link) {
			log.Print(link + " is invalid url.")
			continue
		}
		collect[link] = link
	}

	ch := make(chan string, len(collect))
	defer close(ch)

	var wg sync.WaitGroup
	for link := range collect {
		wg.Add(1)
		go func(link string) {
			wbrc.fetch(link, ch)
			collect[link] = <-ch
			wg.Done()
		}(link)
	}
	wg.Wait()

	return collect, nil
}
