package main

import (
	"errors"
	"fmt"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const endpoint = "https://www.googleapis.com/youtube/v3"
const api_key_variable string = "YOUTUBE_API_KEY"

var video_ids = []string{
	"IHNzOHi8sJs",
	"dISNgvVpWlo",
	"Amq-qlqbjYA",
	"bwmSjveL3Lc",
	"FzVR_fymZw4",
	"9pdj4iJD08s",
	"b73BI9eUkjM",
}

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
	json, err := video_stats(id)
	if err != nil {
		return 0, err
	}

	views, err := get_stat(json, "viewCount")
	if err != nil {
		return 0, err
	}

	return views, nil
}

func append_to_file(file string, data []byte) error {
	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	if _, err := f.Write(data); err != nil {
		f.Close()
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	return nil
}

func main() {
	interval := 5 * time.Minute
	t := time.Now()
	for {
		// Sleep until the next interval.
		t = t.Truncate(interval).Add(interval)
		time.Sleep(t.Sub(time.Now()))

		for _, video_id := range video_ids {
			var data string

			views, err := video_views(video_id)
			if err != nil {
				data = fmt.Sprintf("%s\t%s\tERROR: %v\n", t.Format(time.RFC3339), video_id, err)
			} else {
				data = fmt.Sprintf("%s\t%s\t%d\n", t.Format(time.RFC3339), video_id, views)
			}

			fmt.Print(data)
			if err = append_to_file(video_id+".log", []byte(data)); err != nil {
				log.Fatal(err)
			}
		}
	}
}
