package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	TargetMonitor "github.con/prada-monitors-go/sites"
	WalmartMonitor "github.con/prada-monitors-go/sites/walmart"
)

type Monitor struct {
	Site string `json:"site"`
	Sku  string `json:"sku"`
	PriceRangeMin int `json:"priceRangeMin"`
	PriceRangeMax int `json:"priceRangeMax"`
	
}

func target(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "Target Monitor")
	var currentMonitor Monitor
	_ = json.NewDecoder(r.Body).Decode(&currentMonitor)
	fmt.Println(currentMonitor)
	go TargetMonitor.TargetMonitor(currentMonitor.Sku)
	json.NewEncoder(w).Encode(currentMonitor)
}

func walmart(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "Walmart Monitor")
	var currentMonitor Monitor
	_ = json.NewDecoder(r.Body).Decode(&currentMonitor)
	fmt.Println(currentMonitor)
	go WalmartMonitor.NewMonitor(currentMonitor.Sku, currentMonitor.PriceRangeMin, currentMonitor.PriceRangeMax)
	json.NewEncoder(w).Encode(currentMonitor)
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/target", target).Methods("GET")
	router.HandleFunc("/walmart", walmart).Methods("POST")
	log.Fatal(http.ListenAndServe(":7243", router))

}

func main() {
	handleRequests()
}
