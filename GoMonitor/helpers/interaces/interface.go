package Interfaces

import (
	Types "github.con/prada-monitors-go/helpers/types"
)

type Monitor struct {
	Site          string `json:"site"`
	Sku           string `json:"sku"`
	PriceRangeMin int    `json:"priceRangeMin"`
	PriceRangeMax int    `json:"priceRangeMax"`
	SkuName       string `json:"skuName"`
}

type Site interface {
	getSiteName() string
	getSitePackage() Types.Monitor
}

// func test(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println(r.URL)
// 	w.Header().Set("Content-Type", "application/json")
// 	fmt.Fprintf(w, "Target Monitor")
// 	var currentMonitor Monitor
// 	_ = json.NewDecoder(r.Body).Decode(&currentMonitor)
// 	fmt.Println(currentMonitor)
// 	go getSitePackage.NewMonitor(currentMonitor.Sku, 1, 100000)
// 	json.NewEncoder(w).Encode(currentMonitor)
// }
