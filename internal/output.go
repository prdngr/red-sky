package internal

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

const redSkyBanner = `
        ____           _______ __
       / __ \___  ____/ / ___// /____  __
      / /_/ / _ \/ __  /\__ \/ //_/ / / /
     / _, _/  __/ /_/ /___/ / ,< / /_/ /
    /_/ |_|\___/\__,_//____/_/|_|\__, /
                                /____/
	`

func PrintBanner() {
	fmt.Fprintln(os.Stderr, redSkyBanner)
}

func PrintHeader(header string) {
	color.Yellow("\n" + header)
	color.Yellow(strings.Repeat("-", len(header)) + "\n\n")
}
