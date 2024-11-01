package downloader

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
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

func DownloadSingle(url string, filename string, subtitleLang string) {
	id := ExtractContentId(url)
	data := GetContentPageData(id)
	downloadUrl := GetDownloadUrl(data)
	dirCreated := false

	if filename == "" {
		filename = fmt.Sprintf("%s_%d_%d", data.Data.MainContent.Title, data.Data.MainContent.Season, data.Data.MainContent.Episode)
	}

	if subtitleLang != "" {
		subtitles := data.Data.MainContent.Medias[0].Subtitles

		if len(subtitles) == 0 {
			fmt.Println("No subtitles for this media")
			return
		}

		fmt.Println("Fetching subtitles")

		for _, subtitle := range subtitles {
			if subtitle.SrcLang == subtitleLang {
				_ = os.Mkdir(filename, os.ModePerm) // dont care if directory fails to create
				subitleFileName := fmt.Sprintf("%s/%s_%s", filename, filename, subtitle.FileName)
				downloadFile(subtitle.Src, subitleFileName)
				fmt.Println("Subtitles downloaded successfully")
				dirCreated = true
				break
			}
		}
	}
	filepath := fmt.Sprintf("%s.mp4", filename)

	if dirCreated {
		filepath = fmt.Sprintf("%s/%s.mp4", filename, filename)
	}

	fmt.Printf("Downloading %s to %s\n", url, filename)
	downloadFile(downloadUrl, filepath)
	fmt.Printf("Finished Downloading %s to %s\n", url, filename)
}

func DownloadSeason(url string, seasonName string, subtitleLang string) {
	id := ExtractContentId(url)
	// fmt.Println(id)
	data := GetContentPageData(id)

	title := data.Data.MainContent.Title
	seasonList := data.Data.SeasonList
	if seasonList.Type != "seasonal" {
		panic("Invalid seasontype")
	}

	for _, season := range seasonList.Seasons {
		if season.Name == seasonName {
			for _, seasonContent := range season.Contents {
				_ = os.Mkdir(title, os.ModePerm) // dont care if directory fails to create
				// fileName := fmt.Sprintf("%s/episood_%s", title, strconv.Itoa(i+1))
				DownloadSingle(seasonContent.Url, "", subtitleLang)
			}
		}
	}
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
