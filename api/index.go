package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/gocolly/colly"
	"github.com/mitchellh/mapstructure"
)

type person struct {
	Gender             int     `json:"gender"`
	Job                string  `json:"job"`
	Name               string  `json:"name"`
	ID 				int64    `json:"id"`
}

type Women struct {
	NumWomen int64 `json:"women"`
	Total int64 `json:"total"`
	Percentage float64 `json:"percentage"`
	NotFound []int64 `json:"notFound"`
}

var apiKey string = os.Getenv("TMDB_API_KEY")


func Handler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	query := r.URL.Query() //Get URL Params(type map)
	users, ok := query["users"]
	log.Println(len(users))
	if !ok || len(users) == 0 {
		http.Error(w, "no users", 400)
		return
	}
	wom := scrape(users[0])
	js, err := json.Marshal(wom)
	if err != nil {
		http.Error(w, "internal error", 500)
		return
	}
	w.Write(js)
}



func scrape(username string) Women {
	var total int64 
	var women int64
	var not_found []int64
	var not_found_lock = &sync.Mutex{}
	sitetovisit := "https://letterboxd.com/" + username + "/films"
	var wg sync.WaitGroup
	c := colly.NewCollector(
		colly.Async(true),
	)
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 100})
	c.OnHTML(".poster-container", func(e *colly.HTMLElement) { //primary scarer to get url of each film that contian full information
		e.ForEach("div.film-poster", func(i int, ein *colly.HTMLElement) {
			slug := ein.Attr("data-film-slug")
			url := ("https://letterboxd.com" + slug ) //start go routine to collect all film data
			wg.Add(1)
			go isWomen(url, &wg, not_found_lock, &not_found, &total, &women)
		})

	})
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if strings.Contains(link, "films/page") {
			e.Request.Visit(e.Request.AbsoluteURL(link))
		}
	})
	c.Visit(sitetovisit)
	c.Wait()
	wg.Wait()
	fmt.Printf("With Women: %d\n", women)
	fmt.Printf("Total Able To Get: %d\n", total)
	fmt.Printf("Total percentage %.3f\n", float64(women)/float64(total)*100)

	return Women{
		NumWomen: women,
		Total: total,
		Percentage: float64(women)/float64(total)*100,
		NotFound: not_found,
	}
}

func isWomen(url string, wg *sync.WaitGroup, not_found_lock *sync.Mutex, not_found *[]int64, total, women *int64) {

	isCountable := false // does the movie have gender data for a director
	hasWomen    := false // is at least one of those directors a women
	
	defer wg.Done()
	c := colly.NewCollector()
	var id string
	c.OnHTML("body", func(e *colly.HTMLElement) {
		id = e.Attr("data-tmdb-id")
	})
	c.Visit(url)

	tmdbURL := fmt.Sprintf(
		"%s%s%s?api_key=%s",
		"https://api.themoviedb.org/3/movie/",
		id,
		"/credits",
		apiKey,
	)
	res, err := http.Get(tmdbURL)
	if err != nil {
		log.Fatal(err)
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	var film_data map[string]interface{}
	jsonErr := json.Unmarshal(body, &film_data)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	if film_data["crew"] != nil {
		crew := film_data["crew"].([]interface{})
		for _, crew_member_json := range crew {
			var crew_member person
			mapstructure.Decode(crew_member_json, &crew_member)
			if crew_member.Job == "Director" {
				if crew_member.Gender == 0 {
					not_found_lock.Lock()
					*not_found = append(*not_found, crew_member.ID)
					not_found_lock.Unlock()
				} else if crew_member.Gender == 1 {
					isCountable = true
					hasWomen = true
				} else {
					isCountable = true
				}
			}
		}
	}
	if isCountable {
		atomic.AddInt64(total,1)
	}
	if hasWomen {
		atomic.AddInt64(women,1)
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
