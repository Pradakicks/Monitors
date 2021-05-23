package BestBuyMonitor

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
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
	Availability        string
	currentAvailability string
	Client              http.Client
	file                *os.File
	stop                bool
	CurrentCompanies    []Company
}
type Product struct {
	name        string
	stockNumber string
	productId   string
	price       int
	image       string
	link        string
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
type bestBuyResponse []struct {
	Sku struct {
		ButtonState struct {
			ButtonState string `json:"buttonState"`
			DisplayText string `json:"displayText"`
			SkuID       string `json:"skuId"`
		} `json:"buttonState"`
		Condition        string        `json:"condition"`
		InkSubscriptions []interface{} `json:"inkSubscriptions"`
		Names            struct {
			Short string `json:"short"`
		} `json:"names"`
		Price struct {
			CurrentPrice float64 `json:"currentPrice"`
			PriceDomain  struct {
				CurrentAsOfDate string  `json:"currentAsOfDate"`
				CurrentPrice    float64 `json:"currentPrice"`
				IsMAP           bool    `json:"isMAP"`
				PriceEventType  string  `json:"priceEventType"`
				RegularPrice    float64 `json:"regularPrice"`
				SkuID           string  `json:"skuId"`
			} `json:"priceDomain"`
			PricingType        string `json:"pricingType"`
			SmartPricerEnabled bool   `json:"smartPricerEnabled"`
		} `json:"price"`
		ProductType string `json:"productType"`
		SkuID       string `json:"skuId"`
		Subclass    struct {
			DisplayName string `json:"displayName"`
			ID          string `json:"id"`
		} `json:"subclass"`
		URL string `json:"url"`
	} `json:"sku"`
}

var file os.File

// func walmartMonitor(sku string) {
// 	go NewMonitor(sku, 1, 1000)
// 	fmt.Scanln()
// }

func NewMonitor(sku string) *Monitor {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Site : %s, Product : %s Recovering from panic in printAllOperations error is: %v \n", sku, sku, r)
		}
	}()
	fmt.Println("TESTING")
	m := Monitor{}
	m.Availability = "SOLD_OUT"
	var err error
	//	m.Client = http.Client{Timeout: 5 * time.Second}
	m.Config.site = "BestBuy"
	m.Config.startDelay = 3000
	m.Config.sku = sku
	// 	m.file, err = os.Create("./testing.txt")
	m.Client = http.Client{Timeout: 5 * time.Second}
	m.Config.discord = "https://discord.com/api/webhooks/827263591114997790/chAZK84Gnad7rjHDlh4BnF7dz5KQ7-0l4atsFzJGgcTkAaeZno6ePYB_A-WiiClS3FpY"
	m.monitorProduct.name = "Testing Product"
	m.monitorProduct.stockNumber = ""
	m.getProductDetails()

	path := "cloud.txt"
	var proxyList = make([]string, 0)
	buf, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	// defer func() {
	// 	if err = buf.Close(); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }()

	snl := bufio.NewScanner(buf)
	for snl.Scan() {
		proxy := snl.Text()
		proxyList = append(proxyList, proxy)
		splitProxy := strings.Split(string(proxy), ":")
		newProxy := Proxy{}
		newProxy.userAuth = splitProxy[2]
		newProxy.userPass = splitProxy[3]
		newProxy.ip = splitProxy[0]
		newProxy.port = splitProxy[1]
		//	go NewMonitor(newProxy)
		//	time.Sleep(5 * time.Second)
	}
	buf.Close()
	err = snl.Err()
	if err != nil {
		fmt.Println(err)
	}

	// fmt.Println(timeout)
	//m.Availability = "OUT_OF_STOCK"
	//fmt.Println(m)
	i := true
	go m.checkStop()
	time.Sleep(3000 * (time.Millisecond))
	for i == true {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Site : %s, Product : %s Recovering from panic in printAllOperations error is: %v \n", m.Config.site, m.Config.sku, r)
			}
		}()

		if !m.stop {
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
			m.monitor()
			//	time.Sleep(500 * (time.Millisecond))
			fmt.Println("Best Buy : ", m.Availability, m.Config.sku)
		} else {
			fmt.Println(m.Config.sku, "STOPPED STOPPED STOPPED")
			i = false
		}

	}
	return &m
}

func (m *Monitor) monitor() error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Site : %s, Product : %s Recovering from panic in printAllOperations error is: %v \n", m.Config.site, m.Config.sku, r)
		}
	}()

	url := fmt.Sprintf("https://www.bestbuy.com/api/3.0/priceBlocks?skus=%s", m.Config.sku)
	fmt.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)

		return nil
	}
	// req.Header.Add("authority", "discord.com")
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
	// req.Header.Add("cookie", "TealeafAkaSid=r5S-XRsuxWbk94tkqVB3CruTmaJKz32Z")
	res, err := m.Client.Do(req)
	if err != nil {
		fmt.Println(err)

		return nil
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)

		return nil
	}
	//	fmt.Println(res)
	fmt.Println("Best Buy", res.StatusCode)
	if res.StatusCode != 200 {
		return nil
	}
	var realBody bestBuyResponse
	err = json.Unmarshal([]byte(body), &realBody)
	if err != nil {
		fmt.Println(err)

		return nil
	}
	// fmt.Println(realBody)
	var monitorAvailability string
	monitorAvailability = realBody[0].Sku.ButtonState.ButtonState
	m.monitorProduct.name = realBody[0].Sku.Names.Short
	m.monitorProduct.price = int(realBody[0].Sku.Price.CurrentPrice)
	if m.Availability == "SOLD_OUT" && monitorAvailability == "ADD_TO_CART" {
		fmt.Println("Item in Stock")
		go m.sendWebhook()
	}
	if m.Availability == "ADD_TO_CART" && monitorAvailability == "SOLD_OUT" {
		fmt.Println("Item Out Of Stock")
	}
	m.Availability = monitorAvailability
	return nil
}

func (m *Monitor) getProxy(proxyList []string) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Site : %s, Product : %s Recovering from panic in printAllOperations error is: %v \n", m.Config.site, m.Config.sku, r)
		}
	}()
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

func (m *Monitor) sendWebhook() error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Site : %s, Product : %s Recovering from panic in printAllOperations error is: %v \n", m.Config.site, m.Config.sku, r)
		}
	}()
	// now := time.Now()
	// currentTime := strings.Split(now.String(), "-0400")[0]
	for _, letter := range m.monitorProduct.name {
		if string(letter) == `"` {
			m.monitorProduct.name = strings.Replace(m.monitorProduct.name, `"`, "", -1)
		}
	}
	fmt.Println("Testing Here : ", m.monitorProduct.name, "Here")
	if strings.HasSuffix(m.monitorProduct.name, "                       ") {
		m.monitorProduct.name = strings.Replace(m.monitorProduct.name, "                       ", "", -1)
	}
	fmt.Println("Testing Here : ", m.monitorProduct.name, "Here")
	for _, comp := range m.CurrentCompanies {
		fmt.Println(comp.Company)
		go webHookSend(comp, m.Config.site, m.Config.sku, m.monitorProduct.name, m.monitorProduct.price, "test", m.monitorProduct.image)
	}
	// payload := strings.NewReader("{\"content\":null,\"embeds\":[{\"title\":\"Target Monitor\",\"url\":\"https://discord.com/developers/docs/resources/channel#create-message\",\"color\":507758,\"fields\":[{\"name\":\"Product Name\",\"value\":\"%s\"},{\"name\":\"Product Availability\",\"value\":\"In Stock\\u0021\",\"inline\":true},{\"name\":\"Stock Number\",\"value\":\"%s\",\"inline\":true},{\"name\":\"Links\",\"value\":\"[Product](https://www.walmart.com/ip/prada/%s)\"}],\"footer\":{\"text\":\"Prada#4873\"},\"timestamp\":\"2021-04-01T18:40:00.000Z\",\"thumbnail\":{\"url\":\"https://cdn.discordapp.com/attachments/815507198394105867/816741454922776576/pfp.png\"}}],\"avatar_url\":\"https://cdn.discordapp.com/attachments/815507198394105867/816741454922776576/pfp.png\"}")
	return nil
}
func (m *Monitor) getProductDetails() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Site : %s, Product : %s Recovering from panic in printAllOperations error is: %v \n", m.Config.site, m.Config.sku, r)
		}
	}()

	url := fmt.Sprintf("https://www.bestbuy.com/site/prada/%s.p", m.Config.sku)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)

		return
	}
	req.Header.Add("pragma", "no-cache")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("accept", "application/json, text/plain, */*")
	req.Header.Add("dnt", "1")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("accept-language", "en-US,en;q=0.9")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)

		return
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	m.monitorProduct.image = doc.Find(".primary-image").AttrOr("src", "https://cdn.discordapp.com/attachments/815507198394105867/816741454922776576/pfp.png")
	fmt.Println(m.monitorProduct.image)
}
func webHookSend(c Company, site string, sku string, name string, price int, time string, image string) {
	payload := strings.NewReader(fmt.Sprintf(`{
		"content": null,
		"embeds": [
		  {
			"title": "%s Monitor",
			"url": "https://www.bestbuy.com/site/prada/%s.p?skuId=%s",
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
				"name": "Product Sku",
				"value": "%s",
				"inline": true
			  },
			  {
				"name": "Links",
				"value": "[ATC](https://api.bestbuy.com/click/tempo/%s/cart) | [Cart](https://www.bestbuy.com/cart)"
			  }
			],
			"footer": {
			  "text": "Prada#4873"
			},
			"timestamp": "2021-05-13 13:57:26.5157268",
			"thumbnail": {
			  "url": "%s"
			}
		  }
		],
		"avatar_url": "%s"
	  }`, site, sku, sku, c.Color, name, price, sku, sku, image, c.CompanyImage))
	req, err := http.NewRequest("POST", c.Webhook, payload)
	if err != nil {
		fmt.Println(err)
		fmt.Println(payload)
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
		fmt.Println(payload)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		fmt.Println(payload)
	}
	fmt.Println(res)
	fmt.Println(string(body))
	fmt.Println(payload)
	return
}
func (m *Monitor) checkStop() error {
	for !m.stop {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Site : %s, Product : %s Recovering from panic in printAllOperations error is: %v \n", m.Config.site, m.Config.sku, r)
			}
		}()
		url := fmt.Sprintf("https://monitors-9ad2c-default-rtdb.firebaseio.com/monitor/%s/%s.json", strings.ToUpper(m.Config.site), m.Config.sku)
		req, _ := http.NewRequest("GET", url, nil)
		res, _ := http.DefaultClient.Do(req)

		body, _ := ioutil.ReadAll(res.Body)
		var currentObject ItemInMonitorJson
		err := json.Unmarshal(body, &currentObject)
		if err != nil {
			fmt.Println(err)

		}
		m.stop = currentObject.Stop
		m.CurrentCompanies = currentObject.Companies
		fmt.Println(m.CurrentCompanies)
		res.Body.Close()
		//	fmt.Println(currentObject)
		time.Sleep(5000 * (time.Millisecond))
	}
	return nil
}
