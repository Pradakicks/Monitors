package WalmartNew

import (
	"bytes"
	"encoding/gob"
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
	skuName          string // Only for new Egg
	startDelay       int
	discord          string
	site             string
	priceRangeMax    int
	priceRangeMin    int
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
	products            []string
	keywords            []string
	CurrentCompanies    []Company
	Query               string
}
type Product struct {
	name        string
	stockNumber int
	offerId     string
	price       int
	image       string
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

func NewMonitor(query string, sku string) *Monitor {
	// fmt.Println("TESTING", sku)
	m := Monitor{}
	m.Availability = false
	// var err error
	m.Client = http.Client{Timeout: 10 * time.Second}
	m.Config.site = "WalmartNew"
	m.Config.startDelay = 3000
	m.Query = query
	m.Config.sku = sku
	m.Config.skuName = sku
	// 	m.file, err = os.Create("./testing.txt")
	m.Client = http.Client{
		Timeout: 60 * time.Second,
	}
	m.Config.discord = "https://discord.com/api/webhooks/833902775196450818/wCEYgKpT7eJaOtNERfwe5AlietWcFomGp10zTP3JyDEc8Kk3f2ujqY-BpdXPmpQYANiT"
	m.monitorProduct.name = "Testing Product"
	m.monitorProduct.stockNumber = 10

	proxyList := FetchProxies.Get()
	// fmt.Println(timeout)
	//m.Availability = "OUT_OF_STOCK"
	//fmt.Println(m)
	i := true
	time.Sleep(15000 * (time.Millisecond))
	go m.checkStop()
	time.Sleep(3000 * (time.Millisecond))
	for i {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Recovering from panic in printAllOperations error is: %v \n", r)
			}
		}()

		if !m.stop {
			currentProxy := m.getProxy(proxyList)
			splittedProxy := strings.Split(currentProxy, ":")
			proxy := Proxy{splittedProxy[0], splittedProxy[1], splittedProxy[2], splittedProxy[3]}
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
			m.monitor()
		} else {
			fmt.Println(m.Config.sku, "STOPPED STOPPED STOPPED")
			i = false
		}

	}
	return &m
}

func (m *Monitor) getCookies() {
	jar, _ := cookiejar.New(nil)
				m.Client = http.Client{
					Jar: jar,
				}

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

func (m *Monitor) monitor() error {
	watch := stopwatch.Start()
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovering from panic in printAllOperations error is: %v \n", r)
		}
	}()
	

	url := fmt.Sprintf("https://www.walmart.com/search/api/preso?%s", m.Query)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("authority", "www.walmart.com")
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
		fmt.Printf("Walmart New - Status Code : %d Milliseconds elapsed: %v\n", res.StatusCode, watch.Milliseconds())
	}()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		go MonitorLogger.LogError(m.Config.site, m.Config.sku, err)
		return nil
	}
	if res.StatusCode != 200 {
		if res.StatusCode == 412 || res.StatusCode == 444 {
			fmt.Println("Blocked by PX 12341")
			fmt.Println("Blocked by PX1234")
			fmt.Println("Blocked by PX 1234")
			fmt.Println("Blocked by PX 1234")
			fmt.Println("Blocked by PX 1234")
			m.getCookies()
		}
		return nil
	}

	var newList []string
	parser, err := gojq.NewStringQuery(string(body))
	if err != nil {
		fmt.Println(err)
	}
	products, err := parser.Query("items")
	if err != nil {
		fmt.Println(err)
	}
	for key := range products.([]interface{}) {
		var isPresent bool
		pid, err := parser.Query(fmt.Sprintf("items.[%d].usItemId", key))
		if err != nil {
			fmt.Println(err)
			fmt.Println(parser.Query(fmt.Sprintf("items.[%d]", key)))
		}
		sellerName, err := parser.Query(fmt.Sprintf("items.[%d].sellerName", key))
		if err != nil {
			fmt.Println(err)
		}
		title, err := parser.Query(fmt.Sprintf("items.[%d].title", key))
		if err != nil {
			fmt.Println(err)
		}
		image, err := parser.Query(fmt.Sprintf("items.[%d].imageUrl", key))
		if err != nil {
			fmt.Println(err)
		}
		stockNum, err := parser.Query(fmt.Sprintf("items.[%d].quantity", key))
		if err != nil {
			fmt.Println(err)
		}
		offerId, err := parser.Query(fmt.Sprintf("items.[%d].primaryOffer.offerId", key))
		if err != nil {
			fmt.Println(err)
		}
		price, err := parser.Query(fmt.Sprintf("items.[%d].primaryOffer.offerPrice", key))
		if err != nil {
			fmt.Println(err)
		}
		newList = append(newList, pid.(string))
		fmt.Println(len(m.products))
		for _, v := range m.products {
			if v == pid {
				isPresent = true
			}
		}
		fmt.Println(isPresent)
		if !isPresent {
			m.products = append(m.products, pid.(string))
			go m.sendWebhook(pid.(string), offerId.(string), price.(float64), title.(string), image.(string), stockNum.(float64), sellerName.(string))

		}
	}

	m.products = newList
	return nil
}
func (m *Monitor) getProxy(proxyList []string) string {

	//fmt.Scanln()
	// rand.Seed(time.Now().UnixNano())
	// randomPosition := rand.Intn(len(proxyList)-0) + 0
	if m.Config.proxyCount+1 == len(proxyList) {
		m.Config.proxyCount = 0
	}
	m.Config.proxyCount++
	//fmt.Println(proxyList[m.Config.proxyCount])
	return proxyList[m.Config.proxyCount]
}
func (m *Monitor) sendWebhook(sku string, offerId string, price float64, productName string, image string, stockNum float64, sellerName string) error {
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
		go m.webHookSend(comp, m.Config.site, sku, productName, int(price), offerId, t, image, int(stockNum), sellerName)
	}
	return nil
}
func (m *Monitor) webHookSend(c Company, site string, sku string, name string, price int, offerId string, time string, image string, stockNum int, sellerName string) {
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
				"name": "Stock Quantity",
				"value": "%d",
				"inline": true
			  },
			  {
				"name": "Seller Name",
				"value": "%s",
				"inline": true
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
	  }`, site, sku, c.Color, name, price, sku, offerId, stockNum, sellerName, sku, sku, time, image, c.CompanyImage))
	req, err := http.NewRequest("POST", c.Webhook, payload)
	if err != nil {
		fmt.Println(err)
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
	}
	defer res.Body.Close()
	fmt.Println(res)
	fmt.Println(payload)
}
func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
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
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println(err)

		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
		}
		var currentObject ItemInMonitorJson
		err = json.Unmarshal(body, &currentObject)
		if err != nil {
			fmt.Println(err)
		}
		m.stop = currentObject.Stop
		m.CurrentCompanies = currentObject.Companies
		fmt.Println(m.CurrentCompanies)
		time.Sleep(5000 * (time.Millisecond))
	}
	return nil
}
