package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

type Data struct {
	A Artist
	R Relation
	L Location
	D Date
}
type Artist struct {
	Id           uint     `json:"id"`
	Name         string   `json:"name"`
	Image        string   `json:"image"`
	Members      []string `json:"members"`
	CreationDate uint     `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	Relations    string   `json:"relations"`
}
type Location struct {
	Locations []string `json:"locations"`
}
type Date struct {
	Dates []string `json:"dates"`
}
type Relation struct {
	DatesLocations map[string][]string `json:"datesLocations"`
}

var (
	artistInfo   []Artist
	locationMap  map[string]json.RawMessage
	locationInfo []Location
	datesMap     map[string]json.RawMessage
	datesInfo    []Date
	relationMap  map[string]json.RawMessage
	relationInfo []Relation
)

func ArtistData() []Artist {
	artist, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		log.Fatal()
	}
	artistData, err := io.ReadAll(artist.Body)
	if err != nil {
		log.Fatal()
	}
	json.Unmarshal(artistData, &artistInfo)
	return artistInfo
}

func LocationData() []Location {
	var bytes []byte
	location, err2 := http.Get("https://groupietrackers.herokuapp.com/api/locations")
	if err2 != nil {
		log.Fatal()
	}
	locationData, err3 := io.ReadAll(location.Body)
	if err3 != nil {
		log.Fatal()
	}
	err := json.Unmarshal(locationData, &locationMap)
	if err != nil {
		fmt.Println("error", err)
	}
	for _, m := range locationMap {
		for _, v := range m {
			bytes = append(bytes, v)
		}
	}
	err = json.Unmarshal(bytes, &locationInfo)
	if err != nil {
		fmt.Println("error :", err)
	}
	var cleanedBytes []byte
	for _, m := range locationMap {
		for _, v := range m {
			cleanedV := strings.ReplaceAll(string(v), "_", " ")
			cleanedV = strings.ReplaceAll(cleanedV, "-", ", ")
			cleanedBytes = append(cleanedBytes, []byte(cleanedV)...)
		}
	}
	err = json.Unmarshal(cleanedBytes, &locationInfo)
	return locationInfo
}

func DatesData() []Date {
	var bytes []byte
	dates, err2 := http.Get("https://groupietrackers.herokuapp.com/api/dates")
	if err2 != nil {
		log.Fatal()
	}
	datesData, err3 := io.ReadAll(dates.Body)
	if err3 != nil {
		log.Fatal()
	}
	err := json.Unmarshal(datesData, &datesMap)
	if err != nil {
		fmt.Println("error :", err)
	}
	for _, m := range datesMap {
		for _, v := range m {
			bytes = append(bytes, v)
		}
	}
	err = json.Unmarshal(bytes, &datesInfo)
	if err != nil {
		fmt.Println("error :", err)
	}
	return datesInfo
}

func RelationData() []Relation {
	var bytes []byte
	relation, err2 := http.Get("https://groupietrackers.herokuapp.com/api/relation")
	if err2 != nil {
		log.Fatal()
	}
	relationData, err3 := io.ReadAll(relation.Body)
	if err3 != nil {
		log.Fatal()
	}
	err := json.Unmarshal(relationData, &relationMap)
	if err != nil {
		fmt.Println("error :", err)
	}
	for _, m := range relationMap {
		for _, v := range m {
			bytes = append(bytes, v)
		}
	}
	err = json.Unmarshal(bytes, &relationInfo)
	if err != nil {
		fmt.Println("error :", err)
	}
	return relationInfo
}

func collectData() []Data {
	artistInfo := ArtistData()
	locationInfo := LocationData()
	datesInfo := DatesData()
	relationInfo := RelationData()
	dataData := make([]Data, len(artistInfo))
	for i := 0; i < len(artistInfo); i++ {
		dataData[i].A = artistInfo[i]
		dataData[i].R = relationInfo[i]
		dataData[i].L = locationInfo[i]
		dataData[i].D = datesInfo[i]
	}
	return dataData
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 Not Found", 404)
	}
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "text/html")
	tmpl, err := template.ParseFiles(filepath.Join(".", "templates", "index.html"))
	if err != nil {
		fmt.Println("Error parsing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	data := collectData()
	err = tmpl.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		fmt.Println("Error executing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func artistDataHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 || parts[2] == "" {
		http.Error(w, "Bad Request: Missing Artist ID", http.StatusBadRequest)
		return
	}
	if len(parts) > 3 {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}
	artistID := parts[2]
	var artistData Data
	for _, d := range collectData() {
		if strconv.FormatUint(uint64(d.A.Id), 10) == artistID {
			artistData = d
			break
		}
	}
	fmt.Println("Locations:", artistData.L)
	fmt.Println("Dates:", artistData.D)
	fmt.Println("Relations:", artistData.R)
	if artistData.A.Id == 0 {
		http.Error(w, "Not Found: Artist ID not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	tmpl, err := template.ParseFiles("templates/artistData.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	fmt.Println("Executing template with data:", artistData)
	err = tmpl.Execute(w, artistData)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(filepath.Join(".", "templates", "index.html"))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	q := r.URL.Query()
	search := q.Get("query")
	fmt.Print("THIS  IS SEARCH ", search)

	data := ArtistData()
	if len(search) != 0 {
		data = filter(data, search)
	}
	// fmt.Print(artistInfo)
	fmt.Println("THIS  IS DATA", data)
	err = tmpl.Execute(w, data) // Here, we should pass the filtered data, not the original data
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func filter(artists []Artist, search string) []Artist {
	var results []Artist
	lowerSearch := strings.ToLower(search)
	for _, a := range artists {
		if strings.Contains(strings.ToLower(a.Name), lowerSearch) ||
			containsInSlice(lowerSearch, a.Members) ||
			strings.Contains(strings.ToLower(a.Locations), lowerSearch) ||
			strings.Contains(strings.ToLower(a.FirstAlbum), lowerSearch) ||
			strings.Contains(strings.ToLower(strconv.Itoa(int(a.CreationDate))), lowerSearch) {
			results = append(results, a)
		}
	}
	return results
}

func containsInSlice(search string, slice []string) bool {
	for _, item := range slice {
		if strings.Contains(strings.ToLower(item), search) {
			return true
		}
	}
	return false
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/artistData/", artistDataHandler)
	http.HandleFunc("/search", SearchHandler)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	fmt.Println("Server is running on localhost:2020")
	http.ListenAndServe(":2020", nil)
}
