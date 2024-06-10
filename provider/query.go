package provider

import (
	"fmt"
	"strconv"
	"strings"

	"subcli/ui"
	"subcli/util"

	"github.com/gocolly/colly"
)

type QueryResult struct {
	Id         int
	Title, Url string
}

type queryModel struct {
	list []QueryResult
}

func (m queryModel) GetTitles() []string {
	var titles []string
	for index, i := range m.list {
		titles = append(titles, fmt.Sprintf("%v. %s", index+1, i.Title))

	}
	return titles
}

func (m queryModel) getId(index int) string {
	url := m.list[index].Url
	if url != "" {
		id := strings.Split(url, "/")[4]
		return id
	}
	return ""
}

func Search(query string) string {
	spinner := ui.NewSpinner("Searching...")

	q := strings.Join(strings.Split(query, " "), "+")
	url := fmt.Sprintf("https://www.startpage.com/do/search?cmd=process_search&query=%s+imdb:language=english_en", q)

	results := queryModel{}

	c := colly.NewCollector(
		colly.AllowedDomains("www.startpage.com"),
		colly.IgnoreRobotsTxt(),
	)

	c.UserAgent = util.UserAgent

	c.OnHTML("div.result", func(e *colly.HTMLElement) {
		e.ForEach("a.result-title", func(index int, el *colly.HTMLElement) {
			title := el.ChildText("h2")
			url := el.Attr("href")
			urlParams := strings.Split(strings.Split(url, "//")[1], "/")
			if urlParams[0] == "www.imdb.com" && urlParams[len(urlParams)-1] == "" {
				results.list = append(results.list, QueryResult{Id: index, Title: title, Url: url})
			}

		})

	})

	c.OnRequest(func(r *colly.Request) {
		spinner.Start()
	})

	c.OnScraped(func(r *colly.Response) {
		spinner.Stop()
	})

	c.Visit(url)

	selected := ui.Select(results.GetTitles())
	chosen, _ := strconv.Atoi(strings.Split(selected, ".")[0])

	return results.getId(chosen - 1)
}
