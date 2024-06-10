package provider

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"subcli/ui"
	"subcli/util"

	"github.com/gocolly/colly"
)

type Subtitle struct {
	Id                                 int
	Lang, Title, Rating, Uploader, Url string
}

type SubtitleByLang = map[string][]Subtitle

type err error

type subtitleModel struct {
	subtitles []Subtitle
}

func (s subtitleModel) ByLang() SubtitleByLang {
	subtitlesByLang := make(SubtitleByLang)
	for _, sub := range s.subtitles {
		subtitles, ok := subtitlesByLang[sub.Lang]
		if !ok {
			subtitles = []Subtitle{}
		}
		subtitlesByLang[sub.Lang] = append(subtitles, sub)

	}
	return subtitlesByLang

}

func (s subtitleModel) AvailableLangs() []string {
	uniqueMap := make(map[string]bool)
	var langs []string
	for _, i := range s.subtitles {
		if _, ok := uniqueMap[i.Lang]; !ok {
			uniqueMap[i.Lang] = true
			langs = append(langs, i.Lang)
		}
	}
	return langs
}

func (s subtitleModel) FindById(id int) Subtitle {
	for _, sub := range s.subtitles {
		if sub.Id == id {
			return sub
		}
	}
	return Subtitle{}
}

func GetSubtitles(id string) (subtitleModel, err) {
	spinner := ui.NewSpinner("Searching for subtitle...")
	url := fmt.Sprintf("https://yifysubtitles.ch/movie-imdb/%s", id)

	var subtitles []Subtitle

	c := colly.NewCollector(
		colly.AllowedDomains("yifysubtitles.ch"),
		colly.CacheDir("./subtitles"),
	)

	c.UserAgent = util.UserAgent

	c.OnHTML("table.other-subs > tbody", func(e *colly.HTMLElement) {

		subtitle := Subtitle{}

		e.ForEach("tr", func(index int, el *colly.HTMLElement) {
			subtitle.Id = index + 1
			subtitle.Rating = el.ChildText("td.rating-cell")
			subtitle.Title = strings.Fields(el.ChildText("a"))[1]
			subtitle.Lang = el.ChildText("td.flag-cell > span.sub-lang")
			subtitle.Url = el.ChildAttr("a", "href")
			subtitle.Uploader = el.ChildText("td.uploader-cell")
			subtitles = append(subtitles, subtitle)

		})
	})

	c.OnRequest(func(r *colly.Request) {
		spinner.Start()
	})

	c.OnScraped(func(r *colly.Response) {
		spinner.Stop()
	})

	c.Visit(url)

	// Soring by rating
	sort.Slice(subtitles, func(i, j int) bool {
		return subtitles[i].Rating > subtitles[j].Rating
	})

	if len(subtitles) != 0 {

		return subtitleModel{subtitles: subtitles}, nil
	}
	return subtitleModel{subtitles: nil}, errors.New("No subtitles found!")
}
