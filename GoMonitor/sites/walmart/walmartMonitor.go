package WalmartMonitor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"net/http/cookiejar"
	"github.com/bradhe/stopwatch"
	"github.com/elgs/gojq"
	MonitorLogger "github.con/prada-monitors-go/helpers/logging"
	FetchProxies "github.con/prada-monitors-go/helpers/proxy"
)

type Config struct {
	sku              string
	startDelay       int
	discord          string
	site             string
	priceRangeMax    int
	priceRangeMin    int
	image            string
	proxyCount       int
	indexMonitorJson int
}
type Monitor struct {
	Config              Config
	monitorProduct      Product
	Availability        bool
	currentAvailability string
	Client              http.Client
	file                *os.File
	stop                bool
	CurrentCompanies    []Company
	useProxy            bool
}

type Product struct {
	name        string
	stockNumber int
	offerId     string
	price       int
}
type Proxy struct {
	ip       string
	port     string
	userAuth string
	userPass string
}
type ItemInMonitorJson struct {
	Sku       string `json:"sku"`
	Site      string `json:"site"`
	Stop      bool   `json:"stop"`
	Name      string `json:"name"`
	Companies []Company
}
type Company struct {
	Company      string `json:"company"`
	Webhook      string `json:"webhook"`
	Color        string `json:"color"`
	CompanyImage string `json:"companyImage"`
}

func NewMonitor(sku string, priceRangeMin int, priceRangeMax int) *Monitor {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Site : %s, Product : %s Recovering from panic in printAllOperations error is: %v \n", sku, sku, r)
		}
	}()
	m := Monitor{}
	m.Availability = false
	m.Client = http.Client{Timeout: 10 * time.Second}
	m.Config.site = "Walmart"
	m.Config.startDelay = 3000
	m.Config.sku = sku
	// 	m.file, err = os.Create("./testing.txt")
	m.Client = http.Client{Timeout: 60 * time.Second}
	m.Config.discord = "https://discord.com/api/v8/webhooks/826289643455643658/tRuYU2WQGSoyD5gH2QL8dKecI59F8IyH_wds5_pio7pOst79cBWs6wEe0jdkGI1qeYMC"
	m.monitorProduct.name = "Testing Product"
	m.monitorProduct.stockNumber = 10
	m.Config.priceRangeMax = priceRangeMax
	m.Config.priceRangeMin = priceRangeMin
	m.useProxy = true
	proxyList := FetchProxies.Get()

	// fmt.Println(timeout)
	//m.Availability = "OUT_OF_STOCK"
	//fmt.Println(m)
	time.Sleep(15000 * (time.Millisecond))
	go m.checkStop()
	time.Sleep(3000 * (time.Millisecond))

	i := true
	for i {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Site : %s, Product : %s Recovering from panic in printAllOperations error is: %v \n", m.Config.site, m.Config.sku, r)
			}
		}()

		if !m.stop {
			if m.useProxy {
				currentProxy := m.getProxy(proxyList)
				splittedProxy := strings.Split(currentProxy, ":")
				proxy := Proxy{splittedProxy[0], splittedProxy[1], splittedProxy[2], splittedProxy[3]}
				//	fmt.Println(proxy, proxy.ip)
				prox1y := fmt.Sprintf("http://%s:%s@%s:%s", proxy.userAuth, proxy.userPass, proxy.ip, proxy.port)
				proxyUrl, err := url.Parse(prox1y)
				if err != nil {
					fmt.Println(err)
					return nil
				}
				defaultTransport := &http.Transport{
					Proxy: http.ProxyURL(proxyUrl),
				}
				m.Client.Transport = defaultTransport
			} else {
				fmt.Println("No Proxy")
				m.Client.Transport = http.DefaultTransport
				jar, _ := cookiejar.New(nil)
				m.Client = http.Client{
					Jar: jar,
				}
				m.getCookies()
			}

			m.monitor()
			// time.Sleep(500 * (time.Millisecond))
			// fmt.Println(m.Availability)
		} else {
			fmt.Println(m.stop, "STOPPED STOPPED STOPPED")
			i = false
		}
	}
	return &m
}
func (m *Monitor) monitor() error {
	watch := stopwatch.Start()
	var monitorAvailability bool
	var walmartOffer string
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Site : %s Product : %s Recovering from panic in printAllOperations error is: %v \n", m.Config.site, m.Config.sku, r)
		}
	}()

	url := fmt.Sprintf("https://www.walmart.com/terra-firma/item/%s", m.Config.sku)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		go MonitorLogger.LogError(m.Config.site, m.Config.sku, err)
		return nil
	}
	// req.Header.Add("cookie", `vtc=Vfo2bC9cKIoO5EMKBoV7Cs; TS01b0be75=01538efd7cce98dab74ecc0161aba8abd8d31116a1500fbccd7e68e83c0b6c1f4d3fa8338c433194cae117a5aabe34e5a1e0780d89; TS013ed49a=01538efd7c0deff8b24f3a182d57ce08b06f181c86da1f7161cf403d14a43fa7b56bab4b52bd3cc8df3a1065d376fa4278a8489b5e; TS011baee6=01c5a4e2f993fd028656da8dedad4caf75a0670e3d75c5aa42a7f0d3fd4223f91b69f79b2ec688a8f3eaf669e6d9f5e079288f4f61; akavpau_p1=1624067190~id%3Dbb8d8093dad3913c14ff9dcad2b2f5b5`)
	req.Header.Add("authority", "www.walmart.com")
	req.Header.Add("pragma", "no-cache")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("sec-ch-ua", `" Not;A Brand";v="99", "Google Chrome";v="91", "Chromium";v="91"`)
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("dnt", "1")
	req.Header.Add("upgrade-insecure-requests", "1")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36")
	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Add("service-worker-navigation-preload", "true")
	req.Header.Add("sec-fetch-site", "none")
	req.Header.Add("sec-fetch-mode", "navigate")
	req.Header.Add("sec-fetch-user", "?1")
	req.Header.Add("sec-fetch-dest", "document")
	req.Header.Add("accept-language", "en-US,en;q=0.9")
	req.Header.Set("Connection", "close")
	req.Close = true
	res, err := m.Client.Do(req)
	if err != nil {
		fmt.Println(err)
		go MonitorLogger.LogError(m.Config.site, m.Config.sku, err)
		return nil
	}
	defer res.Body.Close()
	defer func() {
		watch.Stop()
		fmt.Printf("Walmart - Status Code : %d : %s %s %d %s:  Milliseconds elapsed: %v\n", res.StatusCode, m.Config.sku, m.monitorProduct.offerId, m.monitorProduct.price, walmartOffer, watch.Milliseconds())
	}()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		go MonitorLogger.LogError(m.Config.site, m.Config.sku, err)
		return nil
	}
	if res.StatusCode != 200 {
		if res.StatusCode == 412 || res.StatusCode == 444 {
			fmt.Println("Blocked by PX")
			fmt.Println("Blocked by PX")
			fmt.Println("Blocked by PX")
			fmt.Println("Blocked by PX")
			fmt.Println("Blocked by PX")
			m.useProxy = !m.useProxy
		}

		return nil
	} else {
		m.useProxy = true
	}

	monitorAvailability = false
	parser, err := gojq.NewStringQuery(string(body))
	if err != nil {
		go MonitorLogger.LogError(m.Config.site, m.Config.sku, err)
		fmt.Println(err)
	}
	walmartOffers, err := parser.QueryToMap("payload.sellers")
	for key, _ := range walmartOffers {
		if err != nil {
			fmt.Println(err)
		}
		sell := fmt.Sprintf("payload.sellers.%s.sellerName", key)
		sellerName, err := parser.QueryToString(sell)
		if err != nil {
			fmt.Println(err)
		}
		if sellerName == "Walmart.com" {
			walmartOffer = key
		}
	}
	selectedProduct, err := parser.Query("payload.selected.product")
	par := fmt.Sprintf("payload.products.%s.productAttributes.productName", selectedProduct)
	name, err := parser.Query(par)
	if err != nil {
		fmt.Println(err)
	}
	m.monitorProduct.name = name.(string)
	im, err := parser.Query("payload.selected.defaultImage")
	productImageName := im.(string)
	ima, err := parser.Query(fmt.Sprintf("payload.images.%s.assetSizeUrls.DEFAULT", productImageName))
	if err != nil {
		fmt.Println(err)
	} else {
		m.Config.image = ima.(string)

	}
	arr, err := parser.Query(fmt.Sprintf("payload.products.%s.offers", selectedProduct))
	if err != nil {
		fmt.Println(err)
	}
	var offerList []string
	for _, value := range arr.([]interface{}) {
		offerList = append(offerList, value.(string))
	}
	off, err := parser.Query("payload.offers")
	for key, _ := range off.(map[string]interface{}) {
		if err != nil {
			fmt.Println(err)
			return nil
		}
		for _, v := range offerList {
			if key == v {
				//	var currentAvailability interface{}
				var currentPrice1 int
				ca, err := parser.Query((fmt.Sprintf("payload.offers.%s.productAvailability.availabilityStatus", key)))
				if err != nil {
					fmt.Println(err)
					break
				}
				currentAvailability := ca.(string)
				CP, err := parser.Query(fmt.Sprintf("payload.offers.%s.pricesInfo.priceMap.CURRENT.price", key))
				if err != nil {
					fmt.Println(err)
					break
				}

				currentSellerId, err := parser.Query(fmt.Sprintf("payload.offers.%s.sellerId", key))
				currentPrice1 = int(CP.(float64))
				if err != nil {
					fmt.Println(err)

					break
				}
				// fmt.Println(key, walmartOffer)
				if currentAvailability == "IN_STOCK" && currentSellerId == walmartOffer {
					fmt.Println(currentAvailability, key, walmartOffer)
					monitorAvailability = true
					m.monitorProduct.offerId = key
					m.monitorProduct.price = currentPrice1
				}
			}
		}
	}

	if m.Availability == false && monitorAvailability == true {
		fmt.Println("Item in Stock")
		go m.sendWebhook()
	}
	if m.Availability == true && monitorAvailability == false {
		fmt.Println("Item Out Of Stock")
	}
	m.Availability = monitorAvailability
	return nil
}
func (m *Monitor) getProxy(proxyList []string) string {
	if m.Config.proxyCount+1 == len(proxyList) {
		m.Config.proxyCount = 0
	}
	m.Config.proxyCount++
	return proxyList[m.Config.proxyCount]
}
func (m *Monitor) sendWebhook() error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Site : %s, Product : %s Recovering from panic in printAllOperations error is: %v \n", m.Config.site, m.Config.sku, r)
		}
	}()
	for _, letter := range m.monitorProduct.name {
		if string(letter) == `"` {
			m.monitorProduct.name = strings.Replace(m.monitorProduct.name, `"`, "", -1)
		}
	}
	t := time.Now().UTC().Format("2006-01-02T15:04:05Z")

	for _, comp := range m.CurrentCompanies {
		fmt.Println(comp.Company)
		go webHookSend(comp, m.Config.site, m.Config.sku, m.monitorProduct.name, m.monitorProduct.price, m.monitorProduct.offerId, t, m.Config.image)
	}
	return nil
}
func (m *Monitor) getCookies() {
	req1, _ := http.NewRequest("GET", "https://www.walmart.com/checkout", nil)
	fmt.Println(m.Client.Jar)
	req1.Header.Add("authority", "www.walmart.com")
	req1.Header.Add("sec-ch-ua", `" Not;A Brand";v="99", "Google Chrome";v="91", "Chromium";v="91"`)
	req1.Header.Add("sec-ch-ua-mobile", "?0")
	req1.Header.Add("upgrade-insecure-requests", "1")
	req1.Header.Add("dnt", "1")
	req1.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36")
	req1.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req1.Header.Add("service-worker-navigation-preload", "true")
	req1.Header.Add("sec-fetch-site", "same-origin")
	req1.Header.Add("sec-fetch-mode", "navigate")
	req1.Header.Add("sec-fetch-user", "?1")
	req1.Header.Add("sec-fetch-dest", "document")
	req1.Header.Add("referer", "https://www.walmart.com/pac?id=818d2a49-28d8-4eb2-9c90-b4f52b0fd0d6&quantity=1&cv=137")
	req1.Header.Add("accept-language", "en-US,en;q=0.9")

	res1, _ := m.Client.Do(req1)
	fmt.Println(res1.StatusCode)
	fmt.Println(m.Client.Jar)
	defer res1.Body.Close()
}
func webHookSend(c Company, site string, sku string, name string, price int, offerId string, time string, image string) {
	payload := strings.NewReader(fmt.Sprintf(`{
		"content": null,
		"embeds": [
		  {
			"title": "%s Monitor",
			"url": "https://www.walmart.com/ip/prada/%s",
			"color": %s,
			"fields": [
			  {
				"name": "Product Name",
				"value": "%s"
			  },
			  {
				"name": "Product Availability",
				"value": "In Stock",
				"inline": true
			  },
			  {
				"name": "Price",
				"value": "%d",
				"inline": true
			  },
			  {
				"name": "Sku",
				"value": "%s",
				"inline": true
			  },
			  {
				"name": "OfferId",
				"value": "%s"
			  },
			  {
				"name": "Links",
				"value": "[Product](https://www.walmart.com/ip/prada/%s) | [ATC](https://affil.walmart.com/cart/buynow?items=%s) | [Checkout](https://www.walmart.com/checkout/) | [Cart](https://www.walmart.com/cart)"
			  }
			],
			"footer": {
			  "text": "Prada#4873"
			},
			"timestamp": "%s",
			"thumbnail": {
			  "url": "%s"
			}
		  }
		],
		"avatar_url": "%s"
	  }`, site, sku, c.Color, name, price, sku, offerId, sku, sku, time, image, c.CompanyImage))
	req, err := http.NewRequest("POST", c.Webhook, payload)
	if err != nil {
		fmt.Println(err)
		go MonitorLogger.LogError(site, sku, err)
	}
	req.Header.Add("pragma", "no-cache")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("accept", "application/json")
	req.Header.Add("dnt", "1")
	req.Header.Add("accept-language", "en")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("sec-fetch-site", "cross-site")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-dest", "empty")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		go MonitorLogger.LogError(site, sku, err)
	}
	defer res.Body.Close()
	fmt.Println(res)
	fmt.Println(payload)
	return
}
func (m *Monitor) checkStop() error {
	for !m.stop {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Site : %s, Product : %s  CHECK STOP Recovering from panic in printAllOperations error is: %v \n", m.Config.site, m.Config.sku, r)
			}
		}()
		getDBPayload := strings.NewReader(fmt.Sprintf(`{
			"site" : "%s",
			"sku" : "%s"
		  }`, strings.ToUpper(m.Config.site), m.Config.sku))
		url := "http://172.93.100.112:7243/DB"
		req, err := http.NewRequest("POST", url, getDBPayload)
		if err != nil {
			fmt.Println(err)
			go MonitorLogger.LogError(m.Config.site, m.Config.sku, err)
			return nil
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println(err)
			res.Body.Close()
			go MonitorLogger.LogError(m.Config.site, m.Config.sku, err)
			return nil

		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			res.Body.Close()
			go MonitorLogger.LogError(m.Config.site, m.Config.sku, err)
			return nil
		}
		var currentObject ItemInMonitorJson
		err = json.Unmarshal(body, &currentObject)
		if err != nil {
			fmt.Println(err)
			res.Body.Close()
			go MonitorLogger.LogError(m.Config.site, m.Config.sku, err)
			return nil
		}
		m.stop = currentObject.Stop
		m.CurrentCompanies = currentObject.Companies
		fmt.Println(m.CurrentCompanies)
		res.Body.Close()
		time.Sleep(5000 * (time.Millisecond))
	}
	return nil
}
