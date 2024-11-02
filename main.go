package main

import (
	// "fmt"
	"jupiter_downloader/downloader"
	"flag"
)

func main() {
	// downloader.DownloadSingle("https://jupiter.err.ee/1038278/aktuaalne-kaamera", "", "")
	urlParam := flag.String("url", "", "URL of the Jupiter show or movie you want to download")

	seasonNameParam := flag.String("seasonName", "", "Season of the show you want to download")
    maxConcurrentParam := flag.Int("maxConcurrent", 4, "Parameter to toggle how many episodes to download at the same time")

	subtitleLanguageParam := flag.String("subtitleLanguage", "ET", "Parameter to toggle what subtitles you want to download. (ET, EN). NB! Jupiter may not have subtitles in your language of choice.")
    
    flag.Parse()
    
    if *seasonNameParam != "" {
        downloader.DownloadSeason(*urlParam, *seasonNameParam, *subtitleLanguageParam, *maxConcurrentParam)
    } else {
        downloader.DownloadSingle(*urlParam, *subtitleLanguageParam, "")
    }
    // fmt.Println(*urlParam, *seasonNameParam, *maxConcurrentParam, *subtitleLanguageParam)
}
