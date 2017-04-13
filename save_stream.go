package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ohtake/tokyosmart-tools/lib"
)

func main() {
	area := flag.String("area", "v-low-tokyo1", "Area (v-low-tokyo1, v-low-nagoya1, v-low-osaka1, v-low-fukuoka1)")
	service := flag.String("service", "2865", "Case-sensitive ServiceID (2865 for tokyo, 2C65 for nagoyoa, 3065 for osaka, 3865 for fukuoka)")
	help := flag.Bool("help", false, "Print help")
	flag.Parse()
	if *help {
		flag.Usage()
		return
	}

	newFileCh := make(chan string, 2*100)
	endpoint := lib.NewEndpoint(*area, *service)
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
				fmt.Println(result.Action, result.TsFile)
			} else {
				fmt.Fprintln(os.Stderr, result.Error)
				// TODO retry?
			}
		}
	}()
	select {}
}
