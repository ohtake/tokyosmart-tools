package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ohtake/tokyosmart-tools/lib"
)

func main() {
	newFileCh := make(chan string, 2*100)
	endpoint := lib.NewEndpoint("v-low-tokyo1", "2865")
	go func() {
		streamWather := lib.NewStreamWatcher(endpoint, newFileCh)
		for {
			err := streamWather.FetchList()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			time.Sleep(5 * time.Second)
		}
	}()
	go func() {
		downloader := lib.NewTSDownloader("output", endpoint, newFileCh)
		downloader.Prepare()
		for {
			result := downloader.Next()
			if result.Success {
				fmt.Printf("%s %s\n", result.Action, result.TsFile)
			} else {
				fmt.Fprintln(os.Stderr, result.Error)
				// TODO retry?
			}
		}
	}()
	select {}
}
