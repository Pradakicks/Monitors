package WalmartMonitor

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
)

type Config struct {
	sku           string
	startDelay    int
	discord       string
	site          string
	priceRangeMax int
	priceRangeMin int
	image         string
	proxyCount    int
	indexMonitorJson int
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
type ItemInMonitorJson struct {
	Sku  string `json:"sku"`
	Site string `json:"site"`
	Stop bool `json:"stop"`
	Name string `json:"name"`
}

var file os.File

// func walmartMonitor(sku string) {
// 	go NewMonitor(sku, 1, 1000)
// 	fmt.Scanln()
// }

func NewMonitor(sku string, priceRangeMin int, priceRangeMax int) *Monitor {
defer func() {
	     if r := recover(); r != nil {
	        fmt.Printf("Recovering from panic in printAllOperations error is: %v \n", r)
	    }
	  }()	
	m := Monitor{}
	m.Availability = false
	var err error
	m.Client = http.Client{Timeout: 5 * time.Second}
	m.Config.site = "Walmart"
	m.Config.startDelay = 3000
	m.Config.sku = sku
	m.file, err = os.Create("./testing.txt")
	m.Client = http.Client{Timeout: 60 * time.Second}
	m.Config.discord = "https://discord.com/api/v8/webhooks/826289643455643658/tRuYU2WQGSoyD5gH2QL8dKecI59F8IyH_wds5_pio7pOst79cBWs6wEe0jdkGI1qeYMC"
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

	data, err := ioutil.ReadFile("GoMonitors.json")
	if err != nil {
		fmt.Print(err)
	}
		// fmt.Println(string(data))
	var monitorCheckJson []interface{}
	err = json.Unmarshal(data, &monitorCheckJson)
	fmt.Println(monitorCheckJson)
	for key, value := range monitorCheckJson {
		var currentObject ItemInMonitorJson
		currentObject.Site = value.(map[string]interface{})["site"].(string)
		currentObject.Stop = value.(map[string]interface{})["stop"].(bool)
		currentObject.Name = value.(map[string]interface{})["name"].(string)
		currentObject.Sku = value.(map[string]interface{})["sku"].(string)
		if currentObject.Sku == m.Config.sku {
			m.Config.indexMonitorJson = key
			fmt.Println(currentObject, key)
		}
	}

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

	// fmt.Println(timeout)
	//m.Availability = "OUT_OF_STOCK"
	//fmt.Println(m)
	i := true
	for i == true {
			defer func() {
	     if r := recover(); r != nil {
	        fmt.Printf("Recovering from panic in printAllOperations error is: %v \n", r)
	    }
	  }()
		data, err := ioutil.ReadFile("GoMonitors.json")
		if err != nil {
			fmt.Print(err)
		}
		// fmt.Println(string(data))
		var monitorCheckJson []interface{}
		err = json.Unmarshal(data, &monitorCheckJson)
		var currentObject ItemInMonitorJson
		currentObject.Site = monitorCheckJson[m.Config.indexMonitorJson].(map[string]interface{})["site"].(string)
		currentObject.Stop = monitorCheckJson[m.Config.indexMonitorJson].(map[string]interface{})["stop"].(bool)
		currentObject.Name = monitorCheckJson[m.Config.indexMonitorJson].(map[string]interface{})["name"].(string)
		currentObject.Sku = monitorCheckJson[m.Config.indexMonitorJson].(map[string]interface{})["sku"].(string)
		if !currentObject.Stop {
		currentProxy := m.getProxy(proxyList)
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
		// fmt.Println(m.Availability)
		} else {
			fmt.Println(currentObject.Sku , "STOPPED STOPPED STOPPED")
			i = false
		}
	
	}
	return &m
}

func (m *Monitor) monitor() error {
//a	fmt.Println("Monitoring")
		defer func() {
	     if r := recover(); r != nil {
	        fmt.Printf(" Product : %s Recovering from panic in printAllOperations error is: %v \n", m.Config.site, r)
	    }
	  }()
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
//	req.Header.Add("authority", "discord.com")
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
		if res.StatusCode == 412 {
			fmt.Println("Blocked by PX")
		}
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
	// fmt.Println(m)
	// fmt.Println(realBody["payload"].(map[string]interface{})["selected"].(map[string]interface{}))
	selectedProduct := realBody["payload"].(map[string]interface{})["selected"].(map[string]interface{})["product"].(string)
	m.monitorProduct.name = realBody["payload"].(map[string]interface{})["products"].(map[string]interface{})[selectedProduct].(map[string]interface{})["productAttributes"].(map[string]interface{})["productName"].(string)
	// fmt.Println(selectedProduct, m.monitorProduct.name)
	productImageName := realBody["payload"].(map[string]interface{})["selected"].(map[string]interface{})["defaultImage"].(string)
	m.Config.image = realBody["payload"].(map[string]interface{})["images"].(map[string]interface{})[productImageName].(map[string]interface{})["assetSizeUrls"].(map[string]interface{})["DEFAULT"].(string)
	var offerList []string
	offerArray := realBody["payload"].(map[string]interface{})["products"].(map[string]interface{})[selectedProduct].(map[string]interface{})["offers"].([]interface{})
	for _, value := range offerArray {
		offerList = append(offerList, value.(string))
	}

	offers := realBody["payload"].(map[string]interface{})["offers"].(map[string]interface{})
	for key, value := range offers {
		for _, v := range offerList {
			if key == v {
				// fmt.Printf("%s === %s\n", key, v)
				var currentAvailability interface{}
				var currentPrice float64
				var currentPrice1 int
				if value.(map[string]interface{}) != nil {
					currentOffer := value.(map[string]interface{})
					if currentOffer["productAvailability"].(map[string]interface{})["availabilityStatus"] != nil {
						currentAvailability = currentOffer["productAvailability"].(map[string]interface{})["availabilityStatus"]
						if currentOffer["pricesInfo"].(map[string]interface{})["priceMap"].(map[string]interface{})["CURRENT"].(map[string]interface{})["price"] != nil {
							currentPrice = currentOffer["pricesInfo"].(map[string]interface{})["priceMap"].(map[string]interface{})["CURRENT"].(map[string]interface{})["price"].(float64)
							currentPrice1 = int(currentPrice)
						}

						if err != nil {
							fmt.Println(err)
							m.file.WriteString(err.Error() + "\n")
							return nil
						}
					}
				}

				if currentAvailability == "IN_STOCK" && m.Config.priceRangeMin < currentPrice1 && currentPrice1 < m.Config.priceRangeMax {
					fmt.Println(currentAvailability, m.Config.priceRangeMin, currentPrice1, m.Config.priceRangeMax)
					monitorAvailability = true
					m.monitorProduct.offerId = key
					m.monitorProduct.price = currentPrice1
				}
			}
		}
	}

	// To Do - Rewrite for loop so it only loops through specific offerIds
	fmt.Println("Walmart : ", monitorAvailability, m.monitorProduct.offerId, m.monitorProduct.price, m.Config.sku)
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

func (m *Monitor) sendWebhook() error {
	defer func() {
	     if r := recover(); r != nil {
	        fmt.Printf("Recovering from panic in printAllOperations error is: %v \n", r)
	    }
	  }()
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
