package downloader

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

type ContentPageData struct {
	Data Data `json:"data"`
}

type Data struct {
	MainContent MainContent `json:"mainContent"`
	SeasonList  SeasonList  `json:"seasonList"`
}

type MainContent struct {
	Medias  []Media `json:"medias"`
	Title   string  `json:"fancyUrl"`
	Season  int     `json:"season"`
	Episode int     `json:"episode"`
}

type SeasonList struct {
	Type    string   `json:"type"`
	Seasons []Season `json:"items"`
}

type Season struct {
	Name     string    `json:"name"`
	Contents []Content `json:"contents"`
}

type Content struct {
	Url string `json:"url"`
}

type Media struct {
	Src       Source      `json:"src"`
	Subtitles []Subtitles `json:"subtitles"`
}

type Subtitles struct {
	Src      string `json:"src"`
	FileName string `json:"filename"`
	SrcLang  string `json:"srclang"`
}

type Source struct {
	File string `json:"file"`
}

func DownloadSingle(url string, subtitleLang string, parentDirectory string) {
	id := ExtractContentId(url)
	data := GetContentPageData(id)
	downloadUrl := GetDownloadUrl(data)

	name := fmt.Sprintf("%s_%d_%d", data.Data.MainContent.Title, data.Data.MainContent.Season, data.Data.MainContent.Episode)

	if subtitleLang != "" {
		subtitles := data.Data.MainContent.Medias[0].Subtitles

		if len(subtitles) == 0 {
			fmt.Println("No subtitles for this media")
			return
		}

		fmt.Println("Fetching subtitles")

		for _, subtitle := range subtitles {
			if subtitle.SrcLang == subtitleLang {
				parentDirectory = filepath.Join(parentDirectory, name)
				_ = os.Mkdir(parentDirectory, os.ModePerm) // dont care if directory fails to create

				subitleFileName := fmt.Sprintf("%s_%s", name, subtitle.FileName)
				subitleFilePath := filepath.Join(parentDirectory, subitleFileName)

				downloadFile(subtitle.Src, subitleFilePath)
				fmt.Println("Subtitles downloaded successfully")
				break
			}
		}
	}

	filenameExt := fmt.Sprintf("%s.mp4", name)

	path := filepath.Join(parentDirectory, filenameExt)

	fmt.Printf("Downloading %s to %s\n", url, path)
	downloadFile(downloadUrl, path)
	fmt.Println("Finished downloading")
	fmt.Println()
}

func DownloadSeason(url string, seasonName string, subtitleLang string, maxConcurrent int) {
	id := ExtractContentId(url)
	data := GetContentPageData(id)
	var wg sync.WaitGroup

	currentConcurrent := 1

	title := data.Data.MainContent.Title
	seasonList := data.Data.SeasonList
	if seasonList.Type != "seasonal" {
		panic("Invalid seasontype")
	}

	for _, season := range seasonList.Seasons {
		if season.Name == seasonName {
			for _, seasonContent := range season.Contents {
				parentDirName := title
				_ = os.Mkdir(parentDirName, os.ModePerm) // dont care if directory fails to create

				if currentConcurrent < maxConcurrent {
					wg.Add(1)
					currentConcurrent++
					go func() {
						defer wg.Done()
						fmt.Println("Running in parallel")
						DownloadSingle(seasonContent.Url, subtitleLang, parentDirName)
					}()
				} else {
					DownloadSingle(seasonContent.Url, subtitleLang, parentDirName)
				}
			}
		}
	}
	wg.Wait()
}

func GetContentPageData(contentId string) *ContentPageData {
	var data ContentPageData

	url := fmt.Sprintf("https://services.err.ee/api/v2/vodContent/getContentPageData?contentId=%s&rootId=3905", contentId)

	resp, err := http.Get(url)
	if err != nil {
		panic("Invalid url")
	}

	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic("Failing to read the body")
	}

	err = json.Unmarshal(bytes, &data)
	if err != nil {
		panic(err)
	}

	return &data
}

func ExtractContentId(url string) string {
	r := regexp.MustCompile(`\d+`)
	id := r.FindString(url)
	return id
}

func GetDownloadUrl(data *ContentPageData) string {
	return fmt.Sprintf("https:%s", data.Data.MainContent.Medias[0].Src.File)
}

func downloadFile(url string, filepath string) {
	out, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	resp, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		panic("Bad status")
	}

	_, err = io.Copy(out, resp.Body)

	if err != nil {
		panic(err)
	}
}
