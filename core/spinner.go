package core

import (
	"time"

	"github.com/briandowns/spinner"
)

var spin = spinner.New(spinner.CharSets[39], 100*time.Millisecond)

func StartSpinner() {
	spin.Suffix = " Working on stuff..."
	spin.Start()
}

func StopSpinner() {
	if spin != nil && spin.Active() {
		spin.Stop()
	}
}
