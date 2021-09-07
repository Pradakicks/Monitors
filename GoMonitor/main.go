package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/bradhe/stopwatch"
	"github.com/elgs/gojq"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	helper "github.con/prada-monitors-go/helpers/mongo"
	Types "github.con/prada-monitors-go/helpers/types"
	TargetMonitor "github.con/prada-monitors-go/sites"
	BigLotsMonitor "github.con/prada-monitors-go/sites/BigLots"
	AcademyMonitor "github.con/prada-monitors-go/sites/academy"
	AmdMonitor "github.con/prada-monitors-go/sites/amd"
	BestBuyMonitor "github.con/prada-monitors-go/sites/bestBuy"
	FanaticsMonitor "github.con/prada-monitors-go/sites/fanatics"
	GameStopMonitor "github.con/prada-monitors-go/sites/gameStop"
	HomeDepot "github.con/prada-monitors-go/sites/homedepot"
	NewEggMonitor "github.con/prada-monitors-go/sites/newEgg"
	RestirMonitor "github.con/prada-monitors-go/sites/restir"
	Shopify "github.con/prada-monitors-go/sites/shopify"
	ShopifyProduct "github.con/prada-monitors-go/sites/shopifyProduct"
	SlickDealsMonitor "github.con/prada-monitors-go/sites/slickDeals"
	TargetNewTradingCards "github.con/prada-monitors-go/sites/targetNew"
	WalmartMonitor "github.con/prada-monitors-go/sites/walmart"
	WalmartNew "github.con/prada-monitors-go/sites/walmartNew"
	"go.mongodb.org/mongo-driver/bson"
)

type Monitor struct {
	Site          string `json:"site"`
	Sku           string `json:"sku"`
	PriceRangeMin int    `json:"priceRangeMin"`
	PriceRangeMax int    `json:"priceRangeMax"`
	SkuName       string `json:"skuName"`
}

type DB struct {
	Site string `json:"site"`
	Sku  string `json:"sku"`
}
type ItemInMonitorJson struct {
	Sku       string `json:"sku"`
	Site      string `json:"site"`
	Stop      bool   `json:"stop"`
	Name      string `json:"name"`
	Companies []Company
}
type Company struct {
	Company      string `json:"company,omitempty"`
	Webhook      string `json:"webhook,omitempty"`
	Color        string `json:"color,omitempty"`
	CompanyImage string `json:"companyImage,omitempty"`
}
type SiteInDB struct {
	Products []Product
}
type Product struct {
	Companies []Company
	Name      string `json:"name,omitempty"`
	Original  string `json:"original,omitempty"`
	Site      string `json:"site,omitempty"`
	Sku       string `json:"sku,omitempty"`
	Stop      bool   `json:"stop,omitempty"`
}
type ItemInCollection struct {
	Type string                 `json:"type"`
	Site map[string]interface{} `json:"sites"`
}
type DiscordIdsDB struct {
	Ids  []string `json:"ids,omitempty"`
	Type string   `json:"type,omitempty"`
}

var DBString string
var MainCollection = helper.ConnectDBMain()
var workerRunning bool = false

func getPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "Deals Monitor")
}
func getDB(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request Sent")

	watch := stopwatch.Start()

	defer func() {
		watch.Stop()
		fmt.Printf("Request Took : %v\n", watch.Milliseconds())
	}()

	parser, err := gojq.NewStringQuery(DBString)

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	var requestData DB
	_ = json.NewDecoder(r.Body).Decode(&requestData)

	requestedItems, err := parser.QueryToMap(fmt.Sprintf("%s.%s", requestData.Site, requestData.Sku))
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	fmt.Println(requestedItems)
	j, err := json.Marshal(requestedItems)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(j))
}
func getEntireDB(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request Sent")

	watch := stopwatch.Start()

	defer func() {
		watch.Stop()
		fmt.Printf("Request Took : %v\n", watch.Milliseconds())
	}()

	fmt.Fprintf(w, DBString)
}
func DBWorker() {
	if !workerRunning {
		// fmt.Println("Beginning Worker")
		workerRunning = true
		defer func() {
			// fmt.Println("Finished Worker")
			workerRunning = false
		}()
		cur, err := MainCollection.Find(context.TODO(), bson.M{"type": "sites"})
		if err != nil {
			fmt.Println(err)
		}
		defer cur.Close(context.TODO())
		for cur.Next(context.TODO()) {
			elements, err := cur.Current.Elements()
			if err != nil {
				fmt.Println(err)
			}
			for _, v := range elements {
				if strings.Contains(v.String(), "initialized") {
					parser, err := gojq.NewStringQuery(v.String())
					if err != nil {
						fmt.Println("DBWORKER")
						fmt.Println(err.Error())
						fmt.Println(err.Error())
						fmt.Println(err.Error())
						fmt.Println(err.Error())
						return
					}
					data, _ := parser.Query("site")
					jsonString, err := json.Marshal(data)
					if err != nil {
						fmt.Println("DBWORKER")
						fmt.Println(err.Error())
						fmt.Println(err.Error())
						fmt.Println(err.Error())
						fmt.Println(err.Error())
						return
					}
					// fmt.Println(string(jsonString))
					DBString = string(jsonString)

					continue
				}
			}
		}
	} else {
		fmt.Println("Worker Already Running")
		time.Sleep(time.Second * 3)
		DBWorker()
	}
}
func getProxies(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("Request For Proxies")
	watch := stopwatch.Start()
	defer func() {
		watch.Stop()
		fmt.Printf("Request For Proxies Took : %v\n", watch.Milliseconds())
	}()
	cur, err := MainCollection.Find(context.TODO(), bson.M{"type": "proxy"})
	if err != nil {
		fmt.Println(err)
	}

	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		elements, err := cur.Current.Elements()
		if err != nil {
			fmt.Println(err)
			fmt.Fprintf(w, err.Error())
		}
		for _, v := range elements {
			if strings.Contains(v.String(), "proxies") {
				w.Header().Set("Content-Type", "application/json")
				var jsonMap map[string]interface{}
				json.Unmarshal([]byte(v.String()), &jsonMap)
				for k, values := range jsonMap {
					if k == "proxies" {
						var proxyList = make([]string, 0)
						for _, proxy := range values.([]interface{}) {
							proxyList = append(proxyList, proxy.(string))
						}
						var currentResponse Types.ProxyResponseType = Types.ProxyResponseType{
							Proxies: proxyList,
						}
						re, err := json.Marshal(currentResponse)
						if err != nil {
							fmt.Println(err)
							fmt.Fprintf(w, err.Error())
						}
						w.Write(re)
					}
				}
			}
		}
	}
}
func getShopifyProxies(w http.ResponseWriter, r *http.Request) {
	watch := stopwatch.Start()
	defer func() {
		watch.Stop()
		fmt.Printf("Request For Proxies Took : %v\n", watch.Milliseconds())
	}()
	cur, err := MainCollection.Find(context.TODO(), bson.M{"type": "shopifyProxy"})
	if err != nil {
		fmt.Println(err)
	}
	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		elements, err := cur.Current.Elements()
		if err != nil {
			fmt.Println(err)
			fmt.Fprintf(w, err.Error())
		}
		for _, v := range elements {
			if strings.Contains(v.String(), "proxies") {
				w.Header().Set("Content-Type", "application/json")
				var jsonMap map[string]interface{}
				json.Unmarshal([]byte(v.String()), &jsonMap)
				for k, values := range jsonMap {
					if k == "proxies" {
						var proxyList = make([]string, 0)
						for _, proxies := range values.(map[string]interface{}) {
							proxyList = append(proxyList, proxies.(string))
						}
						var currentResponse Types.ProxyResponseType = Types.ProxyResponseType{
							Proxies: proxyList,
						}
						re, err := json.Marshal(currentResponse)
						if err != nil {
							fmt.Println(err)
							fmt.Fprintf(w, err.Error())
						}
						w.Write(re)
					}
				}
			}
		}
	}
}
func handleShopifyProducts() {
	fmt.Println("Handling Products")
	cur, err := Shopify.Collection.Find(context.TODO(), bson.M{})
	if err != nil {
		fmt.Println(err)
	}
	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		var product Types.ShopifyNewProduct
		err := cur.Decode(&product) // decode similar to deserialize process.
		if err != nil {
			fmt.Println(err)
		}
		// fmt.Println(product)
		go func() {
			url := "http://104.249.128.207:7243/SHOPIFYPRODUCT"

			var jsonData = []byte(fmt.Sprintf(`{
				"site": "%s",
				"skuName": "%s",
				"sku": "%s",
				"priceRangeMin": 1,
				"priceRangeMax": 100000
			  }`, product.Store, product.Handle, product.Store))

			req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
			req.Header.Add("Content-Type", "application/json")
			res, _ := http.DefaultClient.Do(req)
			defer res.Body.Close()
			body, _ := ioutil.ReadAll(res.Body)
			fmt.Println(string(body))
			// add item our array
		}()

	}
}
func updateSku(w http.ResponseWriter, r *http.Request) {
	watch := stopwatch.Start()
	var product Product
	defer func() {
		watch.Stop()
		fmt.Printf("Request For Updating Site %s Sku : %s Took : %v\n", product.Site, product.Sku, watch.Milliseconds())
	}()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewDecoder(r.Body).Decode(&product)
	// fmt.Println(product)
	testUpdate(product)
	fmt.Fprintf(w, "Updated Product")
}
func testUpdate(currentProduct Product) {
	if !workerRunning {
		// DBWorker()
		// fmt.Println("Beginning Worker")
		workerRunning = true
		defer func() {
			// fmt.Println("Finished Worker")
			workerRunning = false
			DBWorker()
		}()
		obj := make(map[string]interface{})
		err := json.Unmarshal([]byte(DBString), &obj)
		if err != nil {
			fmt.Println("TESTING UDPDATE", err)
			fmt.Println("TESTING UDPDATE", err)
			fmt.Println("TESTING UDPDATE", err)
			fmt.Println("TESTING UDPDATE", err)
			return
		}
		site := currentProduct.Site
		product := currentProduct.Sku
		newChange := currentProduct
		// fmt.Println(site, product)
		var sitePresent bool = false
		var productPresent bool = false
		for k, _ := range obj {
			if k == "_id" {
				delete(obj, k)
			}
			if k == site {
				sitePresent = true
				for keys, _ := range obj[k].(map[string]interface{}) {
					if keys == product {
						productPresent = true
						obj[k].(map[string]interface{})[product] = newChange
					}
				}

			}

		}

		if !sitePresent {
			fmt.Println("Adding New Site", site, product)
			obj[site] = map[string]interface{}{}
			obj[site].(map[string]interface{})[product] = newChange
		} else if !productPresent {
			fmt.Println("Adding New Product", site, product)
			obj[site].(map[string]interface{})[product] = newChange
		}
		filter := bson.M{"type": "sites"}
		replacedResult := MainCollection.FindOneAndDelete(context.TODO(), filter)
		if replacedResult.Err() != nil {
			fmt.Printf("remove fail %v\n", replacedResult)
			fmt.Printf("remove fail %v\n", replacedResult)
			fmt.Printf("remove fail %v\n", replacedResult)
			fmt.Println(errors.Cause(replacedResult.Err()))
		} else {
			time.Sleep(150 * time.Millisecond)
			results, err := MainCollection.InsertOne(context.TODO(), ItemInCollection{
				Type: "sites",
				Site: obj,
			})
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(results)
		}

	} else {
		fmt.Println("Worker Already Running")
		time.Sleep(3 * time.Second)
		testUpdate(currentProduct)
	}
}
func deleteSku(w http.ResponseWriter, r *http.Request) {
	watch := stopwatch.Start()
	var product Product
	defer func() {
		watch.Stop()
		fmt.Printf("Request For Deleting Site %s Sku : %s Took : %v\n", product.Site, product.Sku, watch.Milliseconds())
	}()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewDecoder(r.Body).Decode(&product)
	// fmt.Println(product)
	deleteSkuFunc(product.Site, product.Sku)
	fmt.Fprintf(w, "Deleted Product")
}
func deleteSkuFunc(site string, sku string) {
	if !workerRunning {
		// DBWorker()
		// fmt.Println("Beginning Worker")
		workerRunning = true
		defer func() {
			// fmt.Println("Finished Worker")
			workerRunning = false
			DBWorker()
		}()
		obj := make(map[string]interface{})
		err := json.Unmarshal([]byte(DBString), &obj)
		fmt.Println("Delete Sku ", err)
		fmt.Println("Deleting ", sku, " from ", site)
		for k, _ := range obj {
			if k == "_id" {
				delete(obj, k)
			}
			fmt.Println(k, site)
			if k == site {
				// fmt.Println(len(obj[k].(map[string]interface{})))
				if len(obj[k].(map[string]interface{})) == 1 {
					delete(obj, k)
					continue
				}
				for keys, _ := range obj[k].(map[string]interface{}) {
					if keys == sku {
						delete(obj[k].(map[string]interface{}), keys)
					}
				}

			}
		}

		filter := bson.M{"type": "sites"}
		replacedResult := MainCollection.FindOneAndDelete(context.TODO(), filter)
		if replacedResult.Err() != nil {
			fmt.Printf("remove fail %v\n", replacedResult)
			fmt.Printf("remove fail %v\n", replacedResult)
			fmt.Printf("remove fail %v\n", replacedResult)
			fmt.Println(errors.Cause(replacedResult.Err()))
		} else {
			time.Sleep(150 * time.Millisecond)
			results, err := MainCollection.InsertOne(context.TODO(), ItemInCollection{
				Type: "sites",
				Site: obj,
			})
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(results)
		}

	} else {
		fmt.Println("Worker Already Running")
		time.Sleep(3 * time.Second)
		deleteSkuFunc(site, sku)
	}
}
func getDiscordIds(w http.ResponseWriter, r *http.Request) {
	watch := stopwatch.Start()
	defer func() {
		watch.Stop()
		fmt.Printf("Request For Getting Discord IDS Took : %v\n", watch.Milliseconds())
	}()
	byteSlice, err := json.Marshal(discordIds())
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(w, err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(byteSlice)
}
func discordIds() []string {
	var discordIdsArr []string
	cur, err := MainCollection.Find(context.TODO(), bson.M{"type": "discordids"})
	if err != nil {
		fmt.Println(err)
	}
	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		elements, err := cur.Current.Elements()
		if err != nil {
			fmt.Println(err)
		}
		for k, _ := range elements {
			if elements[k].Key() == "ids" {
				arr := strings.Split(strings.Split(strings.Split(elements[k].Value().String(), "[")[1], "]")[0], ",")
				for _, v := range arr {
					discordIdsArr = append(discordIdsArr, strings.ReplaceAll(v, `"`, ""))
				}
				// fmt.Println(discordIdsArr)
			}
		}
	}
	return discordIdsArr
}
func addDiscordIds(w http.ResponseWriter, r *http.Request) {
	watch := stopwatch.Start()
	var ids DiscordIdsDB
	defer func() {
		watch.Stop()
		fmt.Printf("Request For Adding Discord Ids Took : %v\n", watch.Milliseconds())
	}()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewDecoder(r.Body).Decode(&ids)
	addDiscordId(ids.Ids)
	fmt.Fprintf(w, "Added IDS")
}
func addDiscordId(discordIds []string) {
	if !workerRunning {
		// fmt.Println("Beginning Worker")
		workerRunning = true
		defer func() {
			// fmt.Println("Finished Worker")
			workerRunning = false
			DBWorker()
		}()
		filter := bson.M{"type": "discordids"}
		replacedResult := MainCollection.FindOneAndDelete(context.TODO(), filter)
		if replacedResult.Err() != nil {
			fmt.Printf("remove fail %v\n", replacedResult)
			fmt.Printf("remove fail %v\n", replacedResult)
			fmt.Printf("remove fail %v\n", replacedResult)
			fmt.Println(errors.Cause(replacedResult.Err()))
		} else {
			time.Sleep(150 * time.Millisecond)
			results, err := MainCollection.InsertOne(context.TODO(), DiscordIdsDB{
				Type: "discordids",
				Ids:  discordIds,
			})
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(results)
		}

	} else {
		fmt.Println("Worker Already Running")
		time.Sleep(3 * time.Second)
		addDiscordId(discordIds)
	}
}

func handleRequests() {
	fmt.Println("Server Started")
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/TARGET", TargetMonitor.Target).Methods("POST")
	router.HandleFunc("/WALMART", WalmartMonitor.Walmart).Methods("POST")
	router.HandleFunc("/NEWEGG", NewEggMonitor.Newegg).Methods("POST")
	router.HandleFunc("/BIGLOTS", BigLotsMonitor.Biglots).Methods("POST")
	router.HandleFunc("/TARGETNEW", TargetNewTradingCards.TargetNew).Methods("POST")
	router.HandleFunc("/ACADEMY", AcademyMonitor.Academy).Methods("POST")
	router.HandleFunc("/BESTBUY", BestBuyMonitor.BestBuy).Methods("POST")
	router.HandleFunc("/AMD", AmdMonitor.Amd).Methods("POST")
	router.HandleFunc("/SLICK", SlickDealsMonitor.SlickDeals).Methods("POST")
	router.HandleFunc("/SLICKDEALS", SlickDealsMonitor.SlickDeals).Methods("POST")
	router.HandleFunc("/GAMESTOP", GameStopMonitor.GameStop).Methods("POST")
	router.HandleFunc("/WALMARTNEW", WalmartNew.WalmartNew).Methods("POST")
	router.HandleFunc("/SHOPIFY", Shopify.Shopify).Methods("POST")
	router.HandleFunc("/HOMEDEPOT", HomeDepot.HomeDepot).Methods("POST")
	router.HandleFunc("/SHOPIFYPRODUCT", ShopifyProduct.ShopifyProduct).Methods("POST")
	router.HandleFunc("/FANATICSNEWPRODUCTS", FanaticsMonitor.FanaticsNewProducts).Methods("POST")
	router.HandleFunc("/RESTIR", RestirMonitor.Restir).Methods("POST")

	// Helper Routes
	router.HandleFunc("/", getPage).Methods("GET")
	router.HandleFunc("/DB", getDB).Methods("POST")      // Post
	router.HandleFunc("/DB", getEntireDB).Methods("GET") // Get
	router.HandleFunc("/PROXY", getProxies).Methods("GET")
	router.HandleFunc("/SHOPIFYPROXY", getShopifyProxies).Methods("GET")
	router.HandleFunc("/UPDATESKU", updateSku).Methods("POST")
	router.HandleFunc("/DELETESKU", deleteSku).Methods("POST")
	router.HandleFunc("/DISCORDIDS", getDiscordIds).Methods("GET")
	router.HandleFunc("/DISCORDIDS", addDiscordIds).Methods("POST")

	log.Fatal(http.ListenAndServe(":7243", router))
}

func main() {
	fmt.Println("Initiating Server")
	DBWorker()
	go func() {
		time.Sleep(1000 * time.Millisecond)
		//	handleShopifyProducts()
	}()
	go func() {
		for {
			DBWorker()
			time.Sleep(5000 * (time.Millisecond))
		}
	}()
	handleRequests()

}
