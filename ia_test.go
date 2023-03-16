package ia

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/wabarc/helper"
)

const available = `{
  "archived_snapshots":{
    "closest":{
      "available":true,
      "url":"http://web.archive.org/web/20060101064348/http://www.example.com:80/",
      "timestamp":"20060101064348",
      "status":"200"
    }
  }
}`

func TestWayback(t *testing.T) {
	uri := "https://example.com"
	u, err := url.Parse(uri)
	if err != nil {
		t.Fatal(err)
	}
	wbrc := &Archiver{}
	got, err := wbrc.Wayback(context.Background(), u)
	if err != nil {
		t.Log(got)
		t.Fatal(err)
	}
}

func TestPlayback(t *testing.T) {
	httpClient, mux, server := helper.MockServer()
	defer server.Close()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(available))
	})

	uri := "https://example.com"
	u, err := url.Parse(uri)
	if err != nil {
		t.Fatal(err)
	}
	wbrc := &Archiver{Client: httpClient}
	got, err := wbrc.Playback(context.Background(), u)
	if err != nil {
		t.Log(got)
		t.Fatal(err)
	}
}
