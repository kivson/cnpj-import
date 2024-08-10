package downloader

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/cavaliergopher/grab/v3"
	"github.com/schollz/progressbar/v3"
)

type Downloader struct {
	HttpClient *http.Client
	Graber     *grab.Client
	Dest       string
	BaseUrl    string
	filter     string
}

func (d Downloader) DownloadAll() {
	urls, err := d.retriveUrls()
	if err != nil {
		panic("unable to process url")
	}
	for _, url := range urls {
		status := fmt.Sprintf("Dowloading: %s", url)
		size, err := d.getFileSize(url)
		if err != nil {
			fmt.Printf("Error Processing url %s caused by %v\n", url, err)
			continue
		}
		bar := progressbar.DefaultBytes(int64(size), status)
		d.downloadFile(url, bar.Set)
	}
}

func (d Downloader) getFileSize(url string) (int, error) {
	resp, err := d.HttpClient.Head(url)
	if err != nil {
		return 0, err
	}
	return int(resp.ContentLength), nil
}

func (d Downloader) downloadFile(fileUrl string, notify func(num int) error) error {
	destination := path.Join(d.Dest, path.Base(fileUrl))

	req, err := grab.NewRequest(destination, fileUrl)
	if err != nil {
		return err
	}
	resp := d.Graber.Do(req)
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			notify(int(resp.BytesComplete()))

		case <-resp.Done:
			// download is complete
			break Loop
		}
	}
	notify(int(resp.BytesComplete()))
	return nil
}

func (d Downloader) retriveUrls() ([]string, error) {
	res, err := http.Get(d.BaseUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("url %s returned non 200", d.BaseUrl)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	urls := make([]string, 0, 20)
	doc.Find("a").Each(func(i int, elem *goquery.Selection) {
		ref, exist := elem.Attr("href")
		if exist && strings.HasSuffix(ref, "zip") {
			fileUrl, err := url.JoinPath(d.BaseUrl, ref)
			if err != nil {
				return
			}
			urls = append(urls, fileUrl)
		}
	})
	return urls, nil
}
