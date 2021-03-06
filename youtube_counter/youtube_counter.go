package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/buger/jsonparser"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
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

const scrapeURL = "https://www.youtube.com/watch?v="

var scrapeRegex = regexp.MustCompile(`\\"viewCount\\":\\"(\d+)\\"`)
var parseRegex = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}[+-]\d{2}:\d{2})\t[^\t]*\t(\d+)($|\t)`)

const endpoint = "https://www.googleapis.com/youtube/v3"
const apiKeyVariable string = "YOUTUBE_API_KEY"

const plotFile = "/mnt/c/Users/Cassie/Desktop/points.png"

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

func fetchURL(url string) ([]byte, error) {
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

func videoStats(id string, apiKey string) ([]byte, error) {
	url := fmt.Sprintf("%s/videos?part=statistics&id=%s&key=%s", endpoint, id, apiKey)

	json, err := fetchURL(url)
	if err != nil {
		return nil, fmt.Errorf("Error fetching YouTube URL: %v\n", err)
	}

	return json, nil
}

func getStat(json []byte, stat string) (string, error) {
	s, err := jsonparser.GetUnsafeString(json, "items", "[0]", "statistics", stat)
	if err != nil {
		return "", fmt.Errorf("Error parsing JSON result for statistic %q: %v\n%s\n", stat, err, json)
	}

	return s, nil
}

func apiViews(id string, apiKey string) (string, error) {
	json, err := videoStats(id, apiKey)
	if err != nil {
		return "", err
	}

	views, err := getStat(json, "viewCount")
	if err != nil {
		return "", err
	}

	return views, nil
}

func scrapeViews(id string) (string, error) {
	url := scrapeURL + id

	data, err := fetchURL(url)
	if err != nil {
		return "", fmt.Errorf("Error fetching YouTube URL: %v\n", err)
	}

	match := scrapeRegex.FindSubmatch(data)
	if match == nil {
		return "", fmt.Errorf("Couldn't scrape viewCount from URL: %q\n%s", url, data)
	}

	return string(match[1]), nil
}

func fetchViews(id string, apiKey string) (string, error) {
	if apiKey != "" {
		s, err := apiViews(id, apiKey)
		if err == nil {
			return s, nil
		}

		s, err2 := scrapeViews(id)
		if err2 == nil {
			return s, nil
		}

		return "", fmt.Errorf("%v%v", err, err2)
	} else {
		return scrapeViews(id)
	}
}

func videoViews(id string, apiKey string) (int64, error) {
	s, err := fetchViews(id, apiKey)
	if err != nil {
		return 0, err
	}

	i, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return 0, fmt.Errorf("Error converting viewCount %q into int64: %v", s, err)
	}

	return i, nil
}

func appendToFile(file string, data []byte) error {
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
	apiKey := os.Getenv(apiKeyVariable)
	if apiKey != "" {
		fmt.Printf("Using YouTube API key: %q\n", apiKey)
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

			views, err := videoViews(v.id, apiKey)
			if err != nil {
				data = fmt.Sprintf("%s\t%s\tERROR: %v\n", t.Format(time.RFC3339), v.id, err)
			} else {
				data = fmt.Sprintf("%s\t%s\t%d\t%s\n", t.Format(time.RFC3339), v.id, views, v.title)
			}

			fmt.Print(data)
			if err = appendToFile(v.id+".log", []byte(data)); err != nil {
				return err
			}
		}
	}
}

func max() error {
	apiKey := os.Getenv(apiKeyVariable)
	for i := range videos {
		var err error

		v := &videos[i]
		v.when = time.Now()
		v.views, err = videoViews(v.id, apiKey)
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

func videoData(file string) ([]video, error) {
	var line string
	var fileErr error
	var data []video

	f, fileErr := os.Open(file)
	if fileErr != nil {
		return nil, fmt.Errorf("Error opening file %q: %v\n", file, fileErr)
	}
	defer f.Close()

	reader := bufio.NewReader(f)

	for fileErr == nil {
		line, fileErr = reader.ReadString('\n')

		if match := parseRegex.FindStringSubmatch(line); match != nil {
			timeString, countString := match[1], match[2]

			timeValue, err := time.Parse(time.RFC3339, timeString)
			if err != nil {
				return nil, fmt.Errorf("Error parsing timestamp %q: %v\n", timeString, err)
			}

			countValue, err := strconv.ParseInt(countString, 0, 64)
			if err != nil {
				return nil, fmt.Errorf("Error converting view count %q into int64: %v", countString, err)
			}

			data = append(data, video{when: timeValue, views: countValue})
		}
	}
	if fileErr != io.EOF {
		return nil, fmt.Errorf("Error reading file %q: %v\n", file, fileErr)
	}

	return data, nil
}

func billion(file string) error {
	data, err := videoData(file)
	if err != nil {
		return err
	}

	t1 := data[0].when
	c1 := data[0].views

	t2 := data[len(data)-1].when
	c2 := data[len(data)-1].views

	fmt.Printf("Starting time: %v\n", t1)
	fmt.Printf("Starting count: %v\n", c1)
	fmt.Printf("Ending time: %v\n", t2)
	fmt.Printf("Ending count: %v\n", c2)

	difference := t2.Sub(t1)
	fmt.Printf("duration = %v\n", difference)

	days := difference.Hours() / 24
	fmt.Printf("number of days = %v\n", days)

	increase := c2 - c1
	fmt.Printf("view difference = %v\n", increase)

	average := float64(increase) / days
	fmt.Printf("average number of views per day = %f\n", average)

	remainingDays := float64(1000000000-c2) / average
	fmt.Printf("days it takes to hit a billion views = %f\n", remainingDays)

	targetDate := t2.Add(time.Duration(remainingDays * 24 * float64(time.Hour)))
	fmt.Printf("target date to hit a billion = %v\n", targetDate)

	return nil
}

func viewPoints(file string) (plotter.XYs, error) {
	data, err := videoData(file)
	if err != nil {
		return nil, err
	}

	//timeConv := plot.UnixTimeIn(time.Local)

	pts := make(plotter.XYs, len(data))
	for i := range pts {
		pts[i].X = float64(data[i].when.Unix())
		pts[i].Y = float64(data[i].views)
	}

	return pts, nil
}

func plotViews(file, outFile string) error {
	pts, err := viewPoints(file)
	if err != nil {
		return err
	}

	p, err := plot.New()
	if err != nil {
		return err
	}

	p.Title.Text = "Blackpink music video views"
	p.X.Label.Text = "Date"
	p.Y.Label.Text = "Views"

	err = plotutil.AddLinePoints(p, file, pts)
	if err != nil {
		return err
	}

	// Save the plot to a PNG file.
	if err = p.Save(9*vg.Inch, 9*vg.Inch, outFile); err != nil {
		return err
	}

	return nil
}

func main() {
	gatherPtr := flag.Bool("gather", false, "gather data")
	maxPtr := flag.Bool("max", false, "sort by max views")
	billionPtr := flag.Bool("billion", false, "when ddu-du ddu-du hits a billion views")
	plotPtr := flag.Bool("plot", false, "plot views of ddu-du ddu-du")
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
	case *plotPtr:
		if err := plotViews(videos[0].id+".log", plotFile); err != nil {
			log.Fatal(err)
		}
	default:
		flag.Usage()
		os.Exit(1)
	}
}
