package download

import (
	"fmt"
	"path/filepath"
	"subcli/ui"
	"subcli/util"

	"github.com/gocolly/colly"
	"github.com/imroc/req"
)

const domain = "yifysubtitles.ch"

func Download(url string) {

	var (
		downloadLink string
		fileName     string
	)

	spinner := ui.NewSpinner(fmt.Sprintf("Downloading %s", fileName))

	c := colly.NewCollector(
		colly.AllowedDomains(domain),
	)

	c.UserAgent = util.UserAgent

	c.OnHTML("div.movie-main-info", func(e *colly.HTMLElement) {
		url := e.ChildAttr("a.download-subtitle", "href")
		fileName = filepath.Base(url)
		downloadLink = fmt.Sprintf("https://%s%s", domain, url)
	})

	c.OnRequest(func(r *colly.Request) {
		spinner.Start()
	})

	c.OnScraped(func(r *colly.Response) {
		spinner.Stop()
	})

	c.Visit(fmt.Sprintf("https://%s%s", domain, url))

	if downloadLink != "" {
		r, err := req.Get(downloadLink)
		if err == nil {
			r.ToFile(fileName)
		}
	}
}
