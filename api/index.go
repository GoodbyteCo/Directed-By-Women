package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gocolly/colly"
	"github.com/mitchellh/mapstructure"
)

type person struct {
	Gender             int     `json:"gender"`
	Job                string  `json:"job"`
	Name               string  `json:"name"`
}

type Women struct {
	NumWomen int `json:"women"`
	Total int `json:"total"`
	Percentage float64 `json:"percentage"`
}

var apiKey string = os.Getenv("TMDB_API_KEY")

func main() {
	getWomenHandler := http.HandlerFunc(Handler)
	http.Handle("/api", getWomenHandler)
	log.Println("serving at :8080")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	http.ListenAndServe(":"+port, nil)
}

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
	var total int 
	var women int
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
			go isWomen(url, &wg, &total, &women)
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
	}
}

func isWomen(url string, wg *sync.WaitGroup, total, women *int) {
	// fmt.Println("running")
	hasWomen := false
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
	var cd map[string]interface{}
	jsonErr := json.Unmarshal(body, &cd)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	if cd["crew"] != nil {
		crew := cd["crew"].([]interface{})
		*total++
		for _, v := range crew {
			var crew person
			mapstructure.Decode(v, &crew)
			if crew.Job == "Director" {
				if crew.Gender == 1 {
					hasWomen = true
				}
			}
		}
	}
	if hasWomen {
		*women++
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
