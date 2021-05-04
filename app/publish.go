package app

import (
	// "net/http"

	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
)

type fileResponse struct {
	Msg    string `json:"msg"`
	Status int    `json:"status"`
	Data   struct {
		FileId  int    `json:"fileId"`
		FileKey string `json:"fileKey"`
		Url     string `json:"url"`
	} `json:"data"`
}

type response byte

const (
	Init response = iota
	Success
	Fail
)

func downloadLang(lang string, gitPath string) <-chan response {
	r := make(chan response)
	go func() {
		defer close(r)
		client := resty.New()

		url := fmt.Sprintf("https://api.simplelocalize.io/api/v3/export?downloadFormat=single-language-json&languageKey=%s", lang)
		var fr = fileResponse{}
		_, err := client.R().SetHeader("x-simplelocalize-token", os.Getenv("API_KEY")).SetResult(&fr).Get(url)
		if err != nil {
			log.Print(err)
			r <- Fail
		}
		downloadUrl := fr.Data.Url
		filePath := fmt.Sprintf("%s/src/i18n/%s.json", gitPath, lang)
		_, err = client.R().SetOutput(filePath).Get(downloadUrl)

		if err != nil {
			log.Print(err)
			r <- Fail
		}
		r <- Success
	}()
	return r
}

func onPublish(c echo.Context) error {
	gitUrl := os.Getenv("GIT_URL")
	gitPath := getEnv("GIT_PATH", "./tmp")
	if gitPath == "" || gitPath == "." || gitPath == "./" || gitPath == "~" || gitPath == "/" {
		// ultimate fail-safe to protect my drive
		panic("GIT_PATH_ERROR. choose a proper path relative to the app")
	}
	err := os.RemoveAll(gitPath)

	checkError(err, "error while removing existing git path")
	log.Println("cloning...", gitUrl)
	gitClient, err := git.PlainClone(gitPath, false, &git.CloneOptions{URL: gitUrl})
	checkError(err, "")

	defaultLangs := strings.Split(os.Getenv("LANGS"), ",")

	for _, lang := range defaultLangs {
		downloadCh := downloadLang(lang, gitPath)
		<-downloadCh
	}

	branchName := plumbing.ReferenceName(fmt.Sprintf("refs/heads/lang-%d", time.Now().Unix()))

	headRef, err := gitClient.Head()
	checkError(err, "")
	ref := plumbing.NewHashReference(branchName, headRef.Hash())
	err = gitClient.Storer.SetReference((ref))
	checkError(err, "")

	w, err := gitClient.Worktree()
	checkError(err, "")
	for _, lang := range defaultLangs {
		langPath := fmt.Sprintf("%s/src/i18n/%s.json", gitPath, lang)
		_, err = w.Add(langPath)
		checkError(err, "")
	}

	status, err := w.Status()
	checkError(err, "")

	fmt.Println(status)

	return c.String(http.StatusOK, "publish webhook accepted!")
}
