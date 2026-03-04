package internal

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/briandowns/spinner"
)

type spinnerAwareWriter struct {
	base io.Writer
}

var spin = spinner.New(spinner.CharSets[14], 100*time.Millisecond)

func (writer spinnerAwareWriter) Write(data []byte) (int, error) {
	StopSpinner()
	return writer.base.Write(data)
}

func ConfigureLogger() {
	log.SetOutput(spinnerAwareWriter{base: os.Stderr})
}

func StartSpinner(message string) {
	if spin == nil || spin.Active() {
		return
	}

	spin.Suffix = " " + message
	spin.Start()
}

func StopSpinner() {
	if spin == nil || !spin.Active() {
		return
	}

	spin.FinalMSG = ""
	spin.Stop()
}
