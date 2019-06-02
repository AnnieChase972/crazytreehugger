package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/buger/jsonparser"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"time"
)

const scrape_url = "https://www.youtube.com/watch?v="

var scrape_regex = regexp.MustCompile(`\\"viewCount\\":\\"(\d+)\\"`)
var parse_regex = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}[+-]\d{2}:\d{2})\t[^\t]*\t(\d+)($|\t)`)

const endpoint = "https://www.googleapis.com/youtube/v3"
const api_key_variable string = "YOUTUBE_API_KEY"

type video struct {
	id    string
	title string
	views int64
	when  time.Time
}

var videos = []video{
	video{
		id:    "IHNzOHi8sJs",
		title: "Ddu du-Ddu du",
	},
	video{
		id:    "2S24-y0Ij3Y",
		title: "Kill This Love",
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
		return "", fmt.Errorf("Error parsing JSON result for statistic %q: %v\n%s\n", stat, err, json)
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
		return "", fmt.Errorf("Couldn't scrape viewCount from URL: %q\n%s", url, data)
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
		return 0, fmt.Errorf("Error converting viewCount %q into int64: %v", s, err)
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

func gather() error {
	api_key := os.Getenv("YOUTUBE_API_KEY")
	if api_key != "" {
		fmt.Printf("Using YouTube API key: %q\n", api_key)
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
				return err
			}
		}
	}
}

func max() error {
	api_key := os.Getenv("YOUTUBE_API_KEY")
	for i := range videos {
		var err error

		v := &videos[i]
		v.when = time.Now()
		v.views, err = video_views(v.id, api_key)
		if err != nil {
			return fmt.Errorf("Error on video %q (%s): %v\n", v.id, v.title, err)
		}
	}

	sorted := videos
	sort.Slice(sorted, func(i, j int) bool { return videos[j].views < videos[i].views })

	for _, v := range sorted {
		fmt.Printf("%s\t%s\t%d\t%s\n", v.when.Format(time.RFC3339), v.id, v.views, v.title)
	}
	return nil
}

func billion(file string) error {
	var line string
	var err, file_err error

	f, file_err := os.Open(file)
	if file_err != nil {
		return fmt.Errorf("Error opening file %q: %v\n", file, file_err)
	}
	defer f.Close()

	reader := bufio.NewReader(f)

	var start_time, start_count, end_time, end_count string
	for file_err == nil {
		line, file_err = reader.ReadString('\n')

		if match := parse_regex.FindStringSubmatch(line); match != nil {
			if start_time == "" {
				start_time, start_count = match[1], match[2]
			} else {
				end_time, end_count = match[1], match[2]
			}
		}
	}
	if file_err != io.EOF {
		return fmt.Errorf("Error reading file %q: %v\n", file, file_err)
	}

	var t1, t2 time.Time
	if t1, err = time.Parse(time.RFC3339, start_time); err != nil {
		return fmt.Errorf("Error parsing starting timestamp %q: %v\n", start_time, err)
	}
	if t2, err = time.Parse(time.RFC3339, end_time); err != nil {
		return fmt.Errorf("Error parsing ending timestamp %q: %v\n", end_time, err)
	}

	var c1, c2 int64
	if c1, err = strconv.ParseInt(start_count, 0, 64); err != nil {
		return fmt.Errorf("Error converting starting view count %q into int64: %v", start_count, err)
	}
	if c2, err = strconv.ParseInt(end_count, 0, 64); err != nil {
		return fmt.Errorf("Error converting ending view count %q into int64: %v", end_count, err)
	}

	fmt.Printf("Starting time: %v\n", t1)
	fmt.Printf("Starting count: %v\n", c1)
	fmt.Printf("Ending time: %v\n", t2)
	fmt.Printf("Ending count: %v\n", c2)

	return nil
}

func main() {
	gatherPtr := flag.Bool("gather", false, "gather data")
	maxPtr := flag.Bool("max", false, "sort by max views")
	billionPtr := flag.Bool("billion", false, "when ddu-du ddu-du hits a billion views")
	flag.Parse()
	switch {
	case *gatherPtr:
		if err := gather(); err != nil {
			log.Fatal(err)
		}
	case *maxPtr:
		if err := max(); err != nil {
			log.Fatal(err)
		}
	case *billionPtr:
		if err := billion(videos[0].id + ".log"); err != nil {
			log.Fatal(err)
		}
	default:
		flag.Usage()
		os.Exit(1)
	}
}
