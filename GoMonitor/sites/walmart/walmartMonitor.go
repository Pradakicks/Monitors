package WalmartMonitor

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Config struct {
	sku           string
	startDelay    int
	discord       string
	site          string
	priceRangeMax int
	priceRangeMin int
	image         string
}
type Monitor struct {
	Config              Config
	monitorProduct      Product
	Availability        bool
	currentAvailability string
	Client              http.Client
	file                *os.File
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

var file os.File

func walmartMonitor(sku string) {
	go NewMonitor(sku, 1, 1000)
	fmt.Scanln()
}

func NewMonitor(sku string, priceRangeMin int, priceRangeMax int) *Monitor {
	m := Monitor{}
	m.Availability = false
	var err error
	m.Client = http.Client{Timeout: 5 * time.Second}
	m.Config.site = "Walmart"
	m.Config.startDelay = 3000
	m.Config.sku = sku
	m.file, err = os.Create("./testing.txt")
	m.Client = http.Client{Timeout: 60 * time.Second}
	m.Config.discord = "https://discord.com/api/v8/webhooks/801223334703661068/TuQyXOOa0McqUFH2FohthDH7nSKfol_lQtt460DNXh5qqfwrkGYEM7N2aexGWN9NDtEY?wait=true"
	m.monitorProduct.name = "Testing Product"
	m.monitorProduct.stockNumber = 10
	m.Config.priceRangeMax = priceRangeMax
	m.Config.priceRangeMin = priceRangeMin

	if err != nil {
		fmt.Println(err)
		m.file.WriteString(err.Error() + "\n")
		return nil
	}
	defer file.Close()
	// fmt.Println(timeout)
	//m.Availability = "OUT_OF_STOCK"
	//fmt.Println(m)
	i := true
	for i == true {
		currentProxy := getProxy()
		splittedProxy := strings.Split(currentProxy, ":")
		proxy := Proxy{splittedProxy[0], splittedProxy[1], splittedProxy[2], splittedProxy[3]}
		//	fmt.Println(proxy, proxy.ip)
		prox1y := fmt.Sprintf("http://%s:%s@%s:%s", proxy.userAuth, proxy.userPass, proxy.ip, proxy.port)
		proxyUrl, err := url.Parse(prox1y)
		if err != nil {
			fmt.Println(err)
			m.file.WriteString(err.Error() + "\n")
			return nil
		}
		defaultTransport := &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
		m.Client.Transport = defaultTransport
		go m.monitor()
		time.Sleep(500 * (time.Millisecond))
		fmt.Println(m.Availability)
	}
	return &m
}

func (m *Monitor) monitor() error {
	fmt.Println("Monitoring")

	// url := "https://httpbin.org/ip"

	// req, _ := http.NewRequest("GET", url, nil)

	// res, _ := m.Client.Do(req)

	// defer res.Body.Close()
	// body, _ := ioutil.ReadAll(res.Body)

	// fmt.Println(res)
	// fmt.Println(string(body))

	url := fmt.Sprintf("https://www.walmart.com/terra-firma/item/%s", m.Config.sku)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		m.file.WriteString(err.Error() + "\n")
		return nil
	}
	req.Header.Add("authority", "discord.com")
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
		m.file.WriteString(err.Error() + "\n")
		return nil
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		m.file.WriteString(err.Error() + "\n")
		return nil
	}
	//	fmt.Println(res)
	fmt.Println(res.StatusCode)
	if res.StatusCode != 200 {
		return nil
	}
	var realBody map[string]interface{}
	err = json.Unmarshal([]byte(body), &realBody)
	if err != nil {
		fmt.Println(err)
		m.file.WriteString(err.Error() + "\n")
		return nil
	}
	var monitorAvailability bool
	monitorAvailability = false
	nameTheory := realBody["payload"].(map[string]interface{})["products"].(map[string]interface{})
	for _, value := range nameTheory {
		if value.(map[string]interface{})["productAttributes"] != nil {
			m.monitorProduct.name = value.(map[string]interface{})["productAttributes"].(map[string]interface{})["productName"].(string)
			//	fmt.Println(m.monitorProduct.name, key)
		}

	}

	image := realBody["payload"].(map[string]interface{})["images"].(map[string]interface{})
	for _, value := range image {
		m.Config.image = value.(map[string]interface{})["assetSizeUrls"].(map[string]interface{})["DEFAULT"].(string)
	}
	offers := realBody["payload"].(map[string]interface{})["offers"].(map[string]interface{})
	for key, value := range offers {
		currentOffer := value.(map[string]interface{})
		currentAvailability := currentOffer["productAvailability"].(map[string]interface{})["availabilityStatus"]
		currentPrice := currentOffer["pricesInfo"].(map[string]interface{})["priceMap"].(map[string]interface{})["CURRENT"].(map[string]interface{})["price"].(float64)
		currentPrice1 := int(currentPrice)
		if err != nil {
			fmt.Println(err)
			m.file.WriteString(err.Error() + "\n")
			return nil
		}
		if currentAvailability == "IN_STOCK" && m.Config.priceRangeMin < currentPrice1 && currentPrice1 < m.Config.priceRangeMax {
			monitorAvailability = true
			m.monitorProduct.offerId = key
			m.monitorProduct.price = currentPrice1
		}
		//	fmt.Println(key, currentAvailability, currentPrice)

	}
	fmt.Println(monitorAvailability)
	// // log.Printf("%+v", m.Availability)
	if m.Availability == false && monitorAvailability == true {
		fmt.Println("Item in Stock")
		m.sendWebhook()

	}
	if m.Availability == true && monitorAvailability == false {
		fmt.Println("Item Out Of Stock")
	}
	m.Availability = monitorAvailability
	return nil
}

func getProxy() string {
	path := "test.txt"
	var proxyList = make([]string, 0)
	buf, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = buf.Close(); err != nil {
			log.Fatal(err)
		}
	}()

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
	err = snl.Err()
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Scanln()
	rand.Seed(time.Now().UnixNano())
	randomPosition := rand.Intn(len(proxyList)-0) + 0
	return proxyList[randomPosition]
}

func (m *Monitor) sendWebhook() error {
	for _, letter := range m.monitorProduct.name {
		if string(letter) == `"` {
			m.monitorProduct.name = strings.Replace(m.monitorProduct.name, `"`, "", -1)
		}
	}
	// payload := strings.NewReader("{\"content\":null,\"embeds\":[{\"title\":\"Target Monitor\",\"url\":\"https://discord.com/developers/docs/resources/channel#create-message\",\"color\":507758,\"fields\":[{\"name\":\"Product Name\",\"value\":\"%s\"},{\"name\":\"Product Availability\",\"value\":\"In Stock\\u0021\",\"inline\":true},{\"name\":\"Stock Number\",\"value\":\"%s\",\"inline\":true},{\"name\":\"Links\",\"value\":\"[Product](https://www.walmart.com/ip/prada/%s)\"}],\"footer\":{\"text\":\"Prada#4873\"},\"timestamp\":\"2021-04-01T18:40:00.000Z\",\"thumbnail\":{\"url\":\"https://cdn.discordapp.com/attachments/815507198394105867/816741454922776576/pfp.png\"}}],\"avatar_url\":\"https://cdn.discordapp.com/attachments/815507198394105867/816741454922776576/pfp.png\"}")
	payload := strings.NewReader(fmt.Sprintf(`{
  "content": null,
  "embeds": [
    {
      "title": "%s Monitor",
      "url": "https://www.walmart.com/ip/prada/%s",
      "color": 507758,
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
      "timestamp": "2021-04-01T18:40:00.000Z",
      "thumbnail": {
        "url": "%s"
      }
    }
  ],
  "avatar_url": "https://cdn.discordapp.com/attachments/815507198394105867/816741454922776576/pfp.png"
}`, m.Config.site, m.Config.sku, m.monitorProduct.name, m.monitorProduct.price, m.Config.sku, m.monitorProduct.offerId, m.Config.sku, m.Config.sku, m.Config.image))
	req, err := http.NewRequest("POST", m.Config.discord, payload)
	if err != nil {
		fmt.Println(err)
		fmt.Println(payload)
		m.file.WriteString(err.Error() + "\n")
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
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		fmt.Println(payload)
		m.file.WriteString(err.Error() + "\n")
		return nil
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		fmt.Println(payload)
		m.file.WriteString(err.Error() + "\n")
		return nil
	}
	fmt.Println(res)
	fmt.Println(string(body))
	fmt.Println(payload)
	return nil
}
