package main

import (
	"flag"
	"fmt"
	"jupiter_downloader/downloader"
)

func main() {
	urlParam := flag.String("url", "", "URL of the Jupiter show or movie you want to download")
	seasonNameParam := flag.String("seasonName", "", "Season of the show you want to download")
	maxConcurrentParam := flag.Int("maxConcurrent", 1, "Parameter to toggle how many episodes to download at the same time")
	subtitleLanguageParam := flag.String("subtitleLanguage", "ET", "Parameter to toggle what subtitles you want to download. (ET, EN). NB! Jupiter may not have subtitles in your language of choice.")

	flag.Parse()

	if *urlParam == "" {
		fmt.Print("Enter the url: ")
		_, err := fmt.Scanln(urlParam)
		if err != nil {
			panic(err)
		}
	}

	if *seasonNameParam != "" {
		downloader.DownloadSeason(*urlParam, *seasonNameParam, *subtitleLanguageParam, *maxConcurrentParam)
	} else {
		downloader.DownloadSingle(*urlParam, *subtitleLanguageParam, "")
	}
}
