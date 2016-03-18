package main

import (
	"encoding/json"
	"github.com/thisissoon/fm-deepmind/fm"
	"github.com/thisissoon/fm-deepmind/rnd"
	"io/ioutil"
	"log"
	"net/http"
)

func totalSum(t []fm.TotalGenres) int {
	sum := 0
	for _, e := range t {
		sum += e.Total
	}
	return sum
}

type Dataset struct {
	Name      string
	Prob      float64
	TotalHits int
}

func main() {
	url := "https://api.thisissoon.fm/player/stats"
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	// read body
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != 200 {
		log.Fatal("Unexpected status code", res.StatusCode)
	}

	user := fm.Stats{}

	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Fatal(err)
	}

	crippled := map[string]float64{
		"indie r&b": .001,
		"pop":       2,
	}

	gengers := user.MostPlayedGenres
	for index, element := range gengers {
		total := element.Total
		if val, exists := crippled[element.Item]; exists {
			total = int(float64(element.Total) * val)
		}
		gengers[index].Total = total
	}

	totalSum := totalSum(gengers)
	datasets := []Dataset{}
	weights := []float64{}
	for _, element := range gengers {
		prob := float64(element.Total) / float64(totalSum)
		datasets = append(datasets, Dataset{Name: element.Item, Prob: prob, TotalHits: 0})
		weights = append(weights, prob)
	}

	runs := 100000
	for i := 0; i < runs; i++ {
		datasetIndex := rnd.Weight(weights)
		datasets[datasetIndex].TotalHits += 1
	}

	for index, value := range datasets {
		log.Printf("%2d %-20s %20f %20f", index, value.Name, value.Prob, float64(value.TotalHits)/float64(runs))
	}
}
