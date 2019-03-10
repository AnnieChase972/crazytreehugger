package main

import (
	"fmt"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

const scrape_url = "https://www.youtube.com/watch?v="

var scrape_regex = regexp.MustCompile(`\\"viewCount\\":\\"(\d+)\\"`)

const endpoint = "https://www.googleapis.com/youtube/v3"
const api_key_variable string = "YOUTUBE_API_KEY"

type video struct {
	id    string
	title string
}

var videos = []video{
	video{
		id:    "IHNzOHi8sJs",
		title: "Ddu du-Ddu du",
	},
	video{
		id:    "dISNgvVpWlo",
		title: "Whistle",
	},
	video{
		id:    "Amq-qlqbjYA",
		title: "As If It's Your Last",
	},
	video{
		id:    "bwmSjveL3Lc",
		title: "Boombayah",
	},
	video{
		id:    "FzVR_fymZw4",
		title: "Stay",
	},
	video{
		id:    "9pdj4iJD08s",
		title: "Playing With Fire",
	},
	video{
		id:    "b73BI9eUkjM",
		title: "Solo",
	},
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

func video_stats(id string, api_key string) ([]byte, error) {
	url := fmt.Sprintf("%s/videos?part=statistics&id=%s&key=%s", endpoint, id, api_key)

	json, err := fetch_url(url)
	if err != nil {
		return nil, fmt.Errorf("Error fetching YouTube URL: %v\n", err)
	}

	return json, nil
}

func get_stat(json []byte, stat string) (string, error) {
	s, err := jsonparser.GetUnsafeString(json, "items", "[0]", "statistics", stat)
	if err != nil {
		return "", fmt.Errorf("Error parsing JSON result for statistic \"%s\": %v\n%s\n", stat, err, json)
	}

	return s, nil
}

func api_views(id string, api_key string) (string, error) {
	json, err := video_stats(id, api_key)
	if err != nil {
		return "", err
	}

	views, err := get_stat(json, "viewCount")
	if err != nil {
		return "", err
	}

	return views, nil
}

func scrape_views(id string) (string, error) {
	url := scrape_url + id

	data, err := fetch_url(url)
	if err != nil {
		return "", fmt.Errorf("Error fetching YouTube URL: %v\n", err)
	}

	match := scrape_regex.FindSubmatch(data)
	if match == nil {
		return "", fmt.Errorf("Couldn't scrape viewCount from URL: \"%s\"\n%s", url, data)
	}

	return string(match[1]), nil
}

func fetch_views(id string, api_key string) (string, error) {
	if api_key != "" {
		s, err := api_views(id, api_key)
		if err == nil {
			return s, nil
		}

		s, err2 := scrape_views(id)
		if err2 == nil {
			return s, nil
		}

		return "", fmt.Errorf("%v%v", err, err2)
	} else {
		return scrape_views(id)
	}
}

func video_views(id string, api_key string) (int64, error) {
	s, err := fetch_views(id, api_key)
	if err != nil {
		return 0, err
	}

	i, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return 0, fmt.Errorf("Error converting viewCount \"%s\" into int64: %v", s, err)
	}

	return i, nil
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
	api_key := os.Getenv("YOUTUBE_API_KEY")
	if api_key != "" {
		fmt.Printf("Using YouTube API key: \"%s\"\n", api_key)
	} else {
		fmt.Println("Scraping YouTube pages instead of using API calls.")
	}

	interval := 5 * time.Minute
	t := time.Now()
	for {
		// Sleep until the next interval.
		t = t.Truncate(interval).Add(interval)
		time.Sleep(t.Sub(time.Now()))

		for _, v := range videos {
			var data string

			views, err := video_views(v.id, api_key)
			if err != nil {
				data = fmt.Sprintf("%s\t%s\tERROR: %v\n", t.Format(time.RFC3339), v.id, err)
			} else {
				data = fmt.Sprintf("%s\t%s\t%d\t%s\n", t.Format(time.RFC3339), v.id, views, v.title)
			}

			fmt.Print(data)
			if err = append_to_file(v.id+".log", []byte(data)); err != nil {
				log.Fatal(err)
			}
		}
	}
}
