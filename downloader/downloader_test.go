package downloader_test

import (
	"jupiter_downloader/downloader"
	"net/http"
	"testing"
)

func TestIdExtract(t *testing.T) {
	testUrl := "https://jupiter.err.ee/1038278/aktuaalne-kaamera"
	expectedId := "1038278"

	actual := downloader.ExtractContentId(testUrl)

	if expectedId != actual {
		t.Fatalf("Expected %s, got %s", expectedId, actual)
	}
}

func TestGetContentPageData(t *testing.T) {
	testUrl := "https://jupiter.err.ee/1038278/aktuaalne-kaamera"
	id := downloader.ExtractContentId(testUrl)
	data := downloader.GetContentPageData(id)
	medias := data.Data.MainContent.Medias

	if len(medias) == 0 {
		t.Fatalf("Expected medias to contain items, got 0")
	}

	for _, media := range medias {
		file := media.Src.File

		if len(file) == 0 {
			t.Fatalf("Invalid file name gotten")
		}

		t.Logf("Media is %s", media.Src.File)
	}
}

func TestDownloadUrl(t *testing.T) {
	testUrl := "https://jupiter.err.ee/1038278/aktuaalne-kaamera"
	id := downloader.ExtractContentId(testUrl)
	data := downloader.GetContentPageData(id)
	url := downloader.GetDownloadUrl(data)

	resp, err := http.Get(url)

	if err != nil {
		t.Fatal(err)
	}
    
	if resp.StatusCode != http.StatusOK {
	    t.Fatalf("bad status: %s", resp.Status)
	}
}
