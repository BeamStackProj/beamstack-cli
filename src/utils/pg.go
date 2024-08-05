package utils

import (
	"fmt"

	"github.com/BeamStackProj/beamstack-cli/src/types"
	ansi "github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/viper"
)

func ProgressBar(taskDesc string, taskCount string, max int) *progressbar.ProgressBar {
	width := viper.GetInt("PROGRESS_BAR_WIDTH")
	bar := progressbar.NewOptions(max,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(int(width)),
		progressbar.OptionSetDescription(fmt.Sprintf("[cyan][%s][reset] %s\t", taskCount, taskDesc)),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
	return bar
}

func DisplayProgress(progChan *chan types.ProgCount, taskDesc string, taskCount string) {
	var bar *progressbar.ProgressBar
	for i := range *progChan {
		if i.OnInit {
			if i.Count == 1 {
				i.Count = 10
			} else {
				i.Count = i.Count * 10
			}
			bar = ProgressBar(taskDesc, taskCount, i.Count+10)
		} else {
			err := bar.Add(10)
			if err != nil {
				continue
			}
		}
	}
	bar.Add(10)
	fmt.Println()
}
