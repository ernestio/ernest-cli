package view

import (
	"fmt"

	"github.com/fatih/color"
)

// EnvDry : Pretty print for env Dry
func EnvDry(lines []string) {
	if len(lines) == 0 {
		fmt.Println("")
		color.Green("This definition is up to date with latest changes. Nothing will be applied")
		fmt.Println("")
		return
	}

	color.Green("Applying this definition will:")
	fmt.Println("")
	for i := range lines {
		fmt.Println(" - " + lines[i])
	}
	fmt.Println("")
	fmt.Println("If you're agree with these changes please rerun apply without --dry option")
}
