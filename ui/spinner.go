package ui

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
)

func NewSpinner(suffix string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Color("cyan")
	s.Suffix = fmt.Sprintf(" %s", suffix)
	return s
}
