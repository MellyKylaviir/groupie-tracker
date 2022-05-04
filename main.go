package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

type Band struct {
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locs         map[string][]string
}

type Relations struct {
	Index []struct {
		ID             int                 `json:"id"`
		DatesLocations map[string][]string `json:"datesLocations"`
	} `json:"index"`
}

type BandList []Band

func PageHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method not supported", http.StatusBadRequest)
		return
	}

	bandListData, relationData := loadData()

	for i, concertData := range relationData.Index {
		cleanData := improveData(&concertData.DatesLocations)
		bandListData[i].Locs = *cleanData
	}

	t, err := template.ParseFiles("templates/index.html")

	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, bandListData)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func loadData() (BandList, Relations) {

	var (
		bandListData BandList
		relationData Relations
	)

	url := "https://groupietrackers.herokuapp.com/api/"

	file := []string{"artists", "relation"}

	for _, link := range file {

		resp, err := http.Get(url + link)
		logError(err)
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)

		logError(err)

		switch link {
		case "artists":

			err := json.Unmarshal(body, &bandListData)

			logError(err)
		case "relation":

			err := json.Unmarshal(body, &relationData)
			logError(err)

		}
	}
	return bandListData, relationData
}

func improveData(concerts *map[string][]string) *map[string][]string {
	clean := make(map[string][]string, len(*concerts))

	for locations, dates := range *concerts {
		splitLocs := strings.Split(locations, "-")

		switch splitLocs[0] {
		case "st_louis", "st_gallen":
			splitLocs[0] = strings.Title(strings.ReplaceAll(splitLocs[0], "_", ". "))
		case "playa_del_carmen", "rio_de_janeiro":
			splitLocs[0] = strings.ReplaceAll(strings.Title(strings.ReplaceAll(splitLocs[0], "_", " ")), "D", "d")
		case "boulogne_billancourt", "freyming_merlebach":
			splitLocs[0] = strings.Title(strings.ReplaceAll(splitLocs[0], "_", "-"))
		case "pagney_derriere_barine", "westcliff_on_sea":
			splitLocs[0] = strings.ReplaceAll(strings.ReplaceAll(strings.Title(strings.ReplaceAll(splitLocs[0], "_", "-")), "D", "d"), "O", "o")
		default:
			splitLocs[0] = strings.Title(strings.ReplaceAll(splitLocs[0], "_", " "))
		}

		switch splitLocs[1] {
		case "usa", "uk":
			splitLocs[1] = strings.ToTitle(splitLocs[1])
		default:
			splitLocs[1] = strings.Title(strings.ReplaceAll(splitLocs[1], "_", " "))
		}
		locations = splitLocs[0] + ", " + splitLocs[1]
		clean[locations] = dates
	}

	return &clean
}

func logError(err error) {

	if err != nil {
		log.Fatal(err)
	}

}

func main() {

	http.HandleFunc("/", PageHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	fmt.Println("http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	logError(err)

}
