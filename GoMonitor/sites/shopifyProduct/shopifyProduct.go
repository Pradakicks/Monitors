package ShopifyProduct

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/bradhe/stopwatch"

	FetchProxies "github.con/prada-monitors-go/helpers/proxy"
	Types "github.con/prada-monitors-go/helpers/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	products            []int64
	stop                bool
	CurrentCompanies    []Company
	collection          *mongo.Collection
	zerosArray          []string
	prod                Types.ShopifyNewProduct
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

func NewMonitor(sku string, skuName string, collection *mongo.Collection) *Monitor {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Site : %s, Product : %s Recovering from panic in printAllOperations error is: %v \n", sku, sku, r)
		}
	}()
	fmt.Println("Shopify Product", sku, skuName)
	m := Monitor{}
	// m.Availability = "OUT_OF_STOCK_ONLINE"
	m.collection = collection
	m.Config.skuName = skuName
	// var err error
	//	m.Client = http.Client{Timeout: 10 * time.Second}
	m.Config.site = "ShopifyProduct"
	m.Config.startDelay = 3000
	m.Config.sku = sku
	// 	m.file, err = os.Create("./testing.txt")
	m.Client = http.Client{Timeout: 10 * time.Second}
	m.monitorProduct.name = skuName
	m.monitorProduct.stockNumber = ""
	proxyList := FetchProxies.Get()

	// fmt.Println(timeout)
	//m.Availability = "OUT_OF_STOCK"
	//fmt.Println(m)
	time.Sleep(15000 * (time.Millisecond))
	go m.checkStop()
	time.Sleep(3000 * (time.Millisecond))

	var product Types.ShopifyNewProduct
	filter := bson.M{"handle": skuName}
	err := collection.FindOne(context.TODO(), filter).Decode(&product)
	if err != nil {
		fmt.Println(err)
	} else {
		m.prod = product
		fmt.Println(m.prod.Handle)
	}
	i := true
	for i {
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
		} else {
			fmt.Println(m.Config.sku, "STOPPED STOPPED STOPPED")
			i = false
		}

	}
	return &m
}

func (m *Monitor) monitor() error {
	watch := stopwatch.Start()
	t := time.Now().UTC().UnixNano()
	var url string
	switch m.Config.sku {
	case "ShopNiceKicks":
		url = fmt.Sprintf("https://%s.com/products/%s.js?limit=%d", m.Config.sku, m.Config.skuName, t)
		break
	default:
		url = fmt.Sprintf("https://www.%s.com/products/%s.js?limit=%d", m.Config.sku, m.Config.skuName, t)

	}
	// fmt.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return nil
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
	req.Header.Set("Connection", "close")
	req.Close = true
	res, err := m.Client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer res.Body.Close()
	var currentAvailability bool

	defer func() {
		watch.Stop()
		fmt.Printf("Shopify Product - Status Code : %d Cache: %s ,Availability : %t Current : %t Milliseconds elapsed: %v\n", res.StatusCode, res.Header["X-Cache"], m.Availability, currentAvailability, watch.Milliseconds())
	}()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	if res.StatusCode != 200 {
		return nil
	}

	var jsonResponse Types.ShopifyProductJS

	err = json.Unmarshal([]byte(body), &jsonResponse)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	currentAvailability = jsonResponse.Available

	if !m.Availability && currentAvailability {
		fmt.Println("Item In Stock")
		link := fmt.Sprintf("https://www.%s.com/products/%s", m.Config.sku, jsonResponse.Handle)
		var price int64
		var image string
		if len(jsonResponse.Variants) < 1 || len(jsonResponse.Images) == 0 {
			price = 0000
			image = "https://cdn.discordapp.com/attachments/866714782554914857/866806869845737472/Prada_Solutions_transparent2x.png"
		} else {
			price = jsonResponse.Variants[0].Price
			image = jsonResponse.FeaturedImage
		}
		m.Availability = currentAvailability
		go m.sendWebhook(m.Config.sku, jsonResponse.Title, price, link,image)
	}
	if m.Availability && !currentAvailability {
		fmt.Println("Out of Stock")
	}
	m.Availability = currentAvailability
	return nil
}

func (m *Monitor) getProxy(proxyList []string) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Site : %s, Product : %s Recovering from panic in printAllOperations error is: %v \n", m.Config.site, m.Config.sku, r)
		}
	}()
	if m.Config.proxyCount+1 == len(proxyList) {
		m.Config.proxyCount = 0
	}
	m.Config.proxyCount++
	return proxyList[m.Config.proxyCount]
}

func (m *Monitor) sendWebhook(site string, name string, price int64, link string, image string) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Site : %s, Product : %s Recovering from panic in printAllOperations error is: %v \n", m.Config.site, m.Config.sku, r)
		}
	}()
	for _, letter := range name {
		if string(letter) == `"` {
			name = strings.Replace(name, `"`, "", -1)
		}
	}
	fmt.Println("Testing Here : ", name, "Here")
	if strings.HasSuffix(name, "                       ") {
		name = strings.Replace(name, "                       ", "", -1)
	}
	fmt.Println("Testing Here : ", name, "Here")
	// now := time.Now()
	// currentTime := strings.Split(now.String(), "-0400")[0]
	t := time.Now().UTC().Format("2006-01-02T15:04:05Z")

	for _, comp := range m.CurrentCompanies {
		fmt.Println(comp.Company)
		go m.webHookSend(comp, site, name, price, link, t, image)
	}
	// payload := strings.NewReader("{\"content\":null,\"embeds\":[{\"title\":\"Target Monitor\",\"url\":\"https://discord.com/developers/docs/resources/channel#create-message\",\"color\":507758,\"fields\":[{\"name\":\"Product Name\",\"value\":\"%s\"},{\"name\":\"Product Availability\",\"value\":\"In Stock\\u0021\",\"inline\":true},{\"name\":\"Stock Number\",\"value\":\"%s\",\"inline\":true},{\"name\":\"Links\",\"value\":\"[Product](https://www.walmart.com/ip/prada/%s)\"}],\"footer\":{\"text\":\"Prada#4873\"},\"timestamp\":\"2021-04-01T18:40:00.000Z\",\"thumbnail\":{\"url\":\"https://cdn.discordapp.com/attachments/815507198394105867/816741454922776576/pfp.png\"}}],\"avatar_url\":\"https://cdn.discordapp.com/attachments/815507198394105867/816741454922776576/pfp.png\"}")
	return nil
}

func (m *Monitor) webHookSend(c Company, site string, name string, price int64, link string, time string, image string) {
	payload := strings.NewReader(fmt.Sprintf(`{
		"content": null,
		"embeds": [
		  {
			"title": "%s Monitor",
			"url": "%s",
			"color": %s,
			"fields": [
			  {
				"name": "Product Name",
				"value": "%s"
			  },
			  {
				"name": "Price",
				"value": "%d",
				"inline": true
			  },
			  {
				"name": "Links",
				"value": "[In Development](%s)"
			  }
			],
			"footer": {
			  "text": "Prada#4873"
			},
			"timestamp": "%s",
			"thumbnail": {
			  "url": "https:%s"
			}
		  }
		],
		"avatar_url": "%s"
	  }`, site, link, c.Color, name, price, link, time, image, c.CompanyImage))
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
		  }`, strings.ToUpper(m.Config.site), m.Config.skuName))
		  fmt.Println("Contract", strings.ToUpper(m.Config.site), m.Config.skuName)
		url := "http://104.249.128.207:7243/DB"
		req, err := http.NewRequest("POST", url, getDBPayload)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println(err)
			res.Body.Close()
			return nil

		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			res.Body.Close()
			return nil
		}
		var currentObject ItemInMonitorJson
		err = json.Unmarshal(body, &currentObject)
		if err != nil {
			fmt.Println(err)
			res.Body.Close()
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
