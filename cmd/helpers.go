package cmd

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Hit Toggl's api for a given url
func queryTogglApi(url string) []byte {
	togglUsername := os.Getenv("TOGGL_USERNAME")
	togglPassword := os.Getenv("TOGGL_PASSWORD")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.SetBasicAuth(togglUsername, togglPassword)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return body
}
