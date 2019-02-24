package main

import (
	"errors"
	"fmt"
	"github.com/buger/jsonparser"
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

func video_stats(id string) ([]byte, error) {
	api_key := os.Getenv("YOUTUBE_API_KEY")
	if api_key == "" {
		return nil, errors.New("YOUTUBE_API_KEY environment variable not set!\n")
	}

	url := fmt.Sprintf("%s/videos?part=statistics&id=%s&key=%s", endpoint, id, api_key)

	json, err := fetch_url(url)
	if err != nil {
		return nil, fmt.Errorf("Error fetching YouTube URL: %v\n", err)
	}

	return json, nil
}

func get_stat(json []byte, stat string) (int64, error) {
	s, err := jsonparser.GetUnsafeString(json, "items", "[0]", "statistics", stat)
	if err != nil {
		return 0, fmt.Errorf("Error parsing JSON result for statistic \"%s\": %v\n%s\n", stat, err, json)
	}

	i, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return 0, fmt.Errorf("Error converting %s \"%s\" into int64: %v", stat, s, err)
	}

	return i, nil
}

func video_views(id string) (int64, error) {
	json, err := video_stats(video_id)
	if err != nil {
		return 0, err
	}

	views, err := get_stat(json, "viewCount")
	if err != nil {
		return 0, err
	}

	return views, nil
}

func main() {
	views, err := video_views(video_id)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%d\n", views)
}
