package cmd

import (
	"fmt"
	"os"
	"time"

	"strconv"
	"subcli/download"
	"subcli/provider"
	"subcli/style"
	"subcli/ui"

	"github.com/charmbracelet/bubbles/table"
)

func errMsg(str string) {
	fmt.Printf("%v\n", style.ErrColor.Render(str))
}

func someFunc() {
	time.Sleep(time.Second * 2)
}

func Main() {
	query := ui.Input()
	if query == "" {
		errMsg("No input given")
		os.Exit(1)
	}

	// Gets INDB id
	movieId := provider.Search(query)
	subs, subErr := provider.GetSubtitles(movieId)

	if subErr != nil {
		fmt.Printf("\n%v\n", style.ErrColor.Render(subErr.Error()))
		os.Exit(1)
	}

	chosenLanguage := ui.Select(subs.AvailableLangs())
	if chosenLanguage == "" {
		errMsg("You didn't choose any language")
		os.Exit(1)
	}

	subsByLang := subs.ByLang()

	columns := []table.Column{
		{Title: "ID", Width: 3},
		{Title: "Title", Width: 70},
		{Title: "Uploader", Width: 14},
		{Title: "Rating", Width: 6},
	}

	rows := []table.Row{}

	for _, sub := range subsByLang[chosenLanguage] {
		rows = append(rows, []string{strconv.Itoa(sub.Id), sub.Title, sub.Uploader, sub.Rating})
	}

	var chosenSubId = ui.Table(columns, rows)
	id, err := strconv.Atoi(chosenSubId)

	if err != nil {
		errMsg("No file selected, exiting program...")
		os.Exit(1)
	}

	download.Download(subs.FindById(id).Url)
}
