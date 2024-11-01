package main

import (
    "jupiter_downloader/downloader"
)

func main() {
    // downloader.Download("1038278")
    downloader.DownloadSingle("https://jupiter.err.ee/1038278/aktuaalne-kaamera", "", "")
    // downloader.DownloadSingle("https://jupiter.err.ee/1609406782/babulon-berliin", "ET", "")
    // downloader.DownloadSeason("https://jupiter.err.ee/1235599/babulon-berliin", "4", "ET") 
    // fs := http.FileServer(http.Dir("static/"))
    //
    // http.Handle("/", fs)
    //
    // fmt.Println("Listening on 8080")
    // http.ListenAndServe(":8080", nil)
}
