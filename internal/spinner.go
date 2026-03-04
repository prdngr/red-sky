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
	StopSpinnerError("Operation failed")
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
	stopSpinner("")
}

func StopSpinnerError(message string) {
	stopSpinner("🛑" + " " + message + "\n")
}

func stopSpinner(message string) {
	if spin == nil || !spin.Active() {
		return
	}

	spin.FinalMSG = message
	spin.Stop()
}
