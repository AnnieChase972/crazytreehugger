package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

const endpoint = "https://www.googleapis.com/youtube/v3"
const api_key_variable string = "YOUTUBE_API_KEY"
const video_id string = "IHNzOHi8sJs"

func fetch_url(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func video_views(id string) (uint64, error) {
	api_key := os.Getenv("YOUTUBE_API_KEY")

	if api_key == "" {
		return 0, errors.New("YOUTUBE_API_KEY environment variable not set!")
	}

	url := fmt.Sprintf("%s/videos?part=statistics&id=%s&key=%s", endpoint, id, api_key)

	body, err := fetch_url(url)
	if err != nil {
		return 0, fmt.Errorf("Error fetching YouTube URL: %v", err)
	}

	type Result struct {
		ViewCount string
	}

	var result Result

	err = json.Unmarshal(body, &result)
	if err != nil {
		return 0, fmt.Errorf("Error parsing JSON result: %v", err)
	}

	count, err := strconv.ParseUint(result.ViewCount, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Error converting string \"%s\" to uint64: %v", result.ViewCount, err)
	}

	return count, nil
}

func main() {
	views, err := video_views(video_id)
	if err != nil {
		fmt.Printf("ERROR: %v", err)
		os.Exit(1)
	}

	fmt.Printf("%d\n", views)
}
