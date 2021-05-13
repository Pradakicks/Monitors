package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	TargetMonitor "github.con/prada-monitors-go/sites"
	NewEggMonitor "github.con/prada-monitors-go/sites/newEgg"
	WalmartMonitor "github.con/prada-monitors-go/sites/walmart"
	BigLotsMonitor "github.con/prada-monitors-go/sites/BigLots"
	TargetNewTradingCards "github.con/prada-monitors-go/sites/targetNew"
	AcademyMonitor"github.con/prada-monitors-go/sites/academy"
	BestBuyMonitor "github.con/prada-monitors-go/sites/bestBuy"
	AmdMonitor "github.con/prada-monitors-go/sites/amd"
	SlickDealsMonitor "github.con/prada-monitors-go/sites/slickDeals"
)

type Monitor struct {
	Site          string `json:"site"`
	Sku           string `json:"sku"`
	PriceRangeMin int    `json:"priceRangeMin"`
	PriceRangeMax int    `json:"priceRangeMax"`
	SkuName       string `json:"skuName"`
}

type KeyWordMonitor struct {
	Endpoint string `json:"endpoint"`
	Keywords []string `json:"keywords"`
}

func target(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "Target Monitor")
	var currentMonitor Monitor
	_ = json.NewDecoder(r.Body).Decode(&currentMonitor)
	fmt.Println(currentMonitor)
	go TargetMonitor.NewMonitor(currentMonitor.Sku, 1, 100000)
	json.NewEncoder(w).Encode(currentMonitor)
}
func walmart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "Walmart Monitor")
	fmt.Println("Walmart")
	var currentMonitor Monitor
	_ = json.NewDecoder(r.Body).Decode(&currentMonitor)
	fmt.Println(currentMonitor)
	go WalmartMonitor.NewMonitor(currentMonitor.Sku, currentMonitor.PriceRangeMin, currentMonitor.PriceRangeMax)
	json.NewEncoder(w).Encode(currentMonitor)
}
func newegg(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "NewEgg Monitor")
	fmt.Println("New Egg")
	var currentMonitor Monitor
	_ = json.NewDecoder(r.Body).Decode(&currentMonitor)
	fmt.Println(currentMonitor)
	go NewEggMonitor.NewMonitor(currentMonitor.Sku, currentMonitor.SkuName, currentMonitor.PriceRangeMin, currentMonitor.PriceRangeMax)
	json.NewEncoder(w).Encode(currentMonitor)
}
func bigLots(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "Big Lots Monitor")
	fmt.Println("Big Lots")
	var currentMonitor Monitor
	_ = json.NewDecoder(r.Body).Decode(&currentMonitor)
	fmt.Println(currentMonitor)
	go BigLotsMonitor.NewMonitor(currentMonitor.Sku, currentMonitor.PriceRangeMin, currentMonitor.PriceRangeMax)
	json.NewEncoder(w).Encode(currentMonitor)
}
func targetNew(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "Target New Products")
	fmt.Println("Big Lots")
	var currentMonitor KeyWordMonitor
	_ = json.NewDecoder(r.Body).Decode(&currentMonitor)
	fmt.Println(currentMonitor)
	go TargetNewTradingCards.NewMonitor(currentMonitor.Endpoint, currentMonitor.Keywords)
	json.NewEncoder(w).Encode(currentMonitor)
}
func academy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "Academy Monitor")
	fmt.Println("Academy")
	var currentMonitor Monitor
	_ = json.NewDecoder(r.Body).Decode(&currentMonitor)
	fmt.Println(currentMonitor)
	go AcademyMonitor.NewMonitor(currentMonitor.Sku)
	json.NewEncoder(w).Encode(currentMonitor)
}
func bestBuy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "Best Buy Monitor")
	fmt.Println("Best Buy")
	var currentMonitor Monitor
	_ = json.NewDecoder(r.Body).Decode(&currentMonitor)
	fmt.Println(currentMonitor)
	go BestBuyMonitor.NewMonitor(currentMonitor.Sku)
	json.NewEncoder(w).Encode(currentMonitor)
}
func amd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "Amd Monitor")
	fmt.Println("Amd")
	var currentMonitor Monitor
	_ = json.NewDecoder(r.Body).Decode(&currentMonitor)
	fmt.Println(currentMonitor)
	go AmdMonitor.NewMonitor(currentMonitor.Sku)
	json.NewEncoder(w).Encode(currentMonitor)
}
func slickDeals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "Slick Deals Monitor")
	fmt.Println("Slick Deals")
	var currentMonitor Monitor
	_ = json.NewDecoder(r.Body).Decode(&currentMonitor)
	fmt.Println(currentMonitor)
	go SlickDealsMonitor.NewMonitor()
	json.NewEncoder(w).Encode(currentMonitor)
}
func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/target", target).Methods("POST")
	router.HandleFunc("/walmart", walmart).Methods("POST")
	router.HandleFunc("/newEgg", newegg).Methods("POST")
	router.HandleFunc("/bigLots", bigLots).Methods("POST")
	router.HandleFunc("/targetNew", targetNew).Methods("POST")
	router.HandleFunc("/academy", academy).Methods("POST")
	router.HandleFunc("/bestBuy", bestBuy).Methods("POST")
	router.HandleFunc("/amd", amd).Methods("POST")
	router.HandleFunc("/slick", slickDeals).Methods("POST")
	log.Fatal(http.ListenAndServe(":7243", router))

}

func main() {
	handleRequests()
}
