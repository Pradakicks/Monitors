package AcademyMonitor 

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
	skuName string // Only for new Egg
	startDelay    int
	discord       string
	site          string
	priceRangeMax int
	priceRangeMin int
	proxyCount    int
	indexMonitorJson int
}
type Monitor struct {
	Config              Config
	monitorProduct      Product
	Availability        string
	currentAvailability string
	Client              http.Client
	file                *os.File
}
type Product struct {
	name        string
	stockNumber 	string
	productId     string
	price       	int
	image         string
	link string
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
type AcademyResponse struct {
	Online []struct {
		Specialorder      string `json:"SPECIALORDER"`
		AddToCart         string `json:"addToCart"`
		AvailableQuantity string `json:"availableQuantity"`
		DeliveryMessage   struct {
			OnlineDeliveryMessage struct {
				AdditionalInfoKey   string `json:"additionalInfoKey"`
				AdditionalInfoValue string `json:"additionalInfoValue"`
				Key                 string `json:"key"`
				ShowTick            string `json:"showTick"`
				Value               string `json:"value"`
			} `json:"onlineDeliveryMessage"`
			StoreDeliveryMessage struct {
				Key          string `json:"key"`
				ShowTick     string `json:"showTick"`
				StoreInvType string `json:"storeInvType"`
				Value        string `json:"value"`
			} `json:"storeDeliveryMessage"`
		} `json:"deliveryMessage"`
		InventoryStatus             string `json:"inventoryStatus"`
		SkuID                       string `json:"skuId"`
	} `json:"online"`
	ProductID string `json:"productId"`
}

type ProductDetails struct {
	Productinfo []struct {
		Attributes []struct {
			Name     string `json:"name"`
			UniqueID string `json:"uniqueID"`
			Usage    string `json:"usage"`
			Values   string `json:"values"`
		} `json:"attributes"`
		DefaultSkuID    string   `json:"defaultSkuId"`
		FullImage              string        `json:"fullImage"`
		ID                     string        `json:"id"`
		ImageURL               string        `json:"imageURL"`
		Name                   string        `json:"name"`
		PriceUSD               float64       `json:"priceUSD"`
		ProductPrice struct {
			ListPrice string `json:"listPrice"`
		} `json:"productPrice"`
		UpdatedSeoURL     string   `json:"updatedSeoURL"`
	} `json:"productinfo"`
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
	m.Availability = "OUT_OF_STOCK_ONLINE"
	var err error
//	m.Client = http.Client{Timeout: 5 * time.Second}
	m.Config.site = "Academy"
	m.Config.startDelay = 3000
	m.Config.sku = sku
	m.file, err = os.Create("./testing.txt")
	m.Client = http.Client{Timeout: 5 * time.Second}
	m.Config.discord = "https://discord.com/api/webhooks/838637042119213067/V7WQ7z-9u32rNh5SO4YyxS5kibcHadXW4FxjVJTosO5cSGRoSqv4CY5g3GrAcIcwZhkF"
	m.monitorProduct.name = "Testing Product"
	m.monitorProduct.stockNumber = ""
	m.getProductDetails()
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

	path := "cloud.txt"
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
	        	        fmt.Printf("Site : %s, Product : %s Recovering from panic in printAllOperations error is: %v \n", m.Config.site, m.Config.sku, r)
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
		m.monitor()
	//	time.Sleep(500 * (time.Millisecond))
		fmt.Println("Academy : ", m.Availability, m.Config.sku)
		} else {
			fmt.Println(currentObject.Sku , "STOPPED STOPPED STOPPED")
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
//	fmt.Println("Monitoring")
	// 	defer func() {
	//      if r := recover(); r != nil {
	//         	        fmt.Printf("Site : %s, Product : %s Recovering from panic in printAllOperations error is: %v \n", m.Config.site, m.Config.sku, r)
	//     }
	//   }()
	// url := "https://httpbin.org/ip"

	// req, _ := http.NewRequest("GET", url, nil)

	// res, _ := m.Client.Do(req)

	// defer res.Body.Close()
	// body, _ := ioutil.ReadAll(res.Body)

	// fmt.Println(res)
	// fmt.Println(string(body))

	url := fmt.Sprintf("https://www.academy.com/api/inventory/v2?productId=%s&bopisEnabled=true&isSTSEnabled=true", m.Config.sku)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		m.file.WriteString(err.Error() + "\n")
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
	var realBody AcademyResponse
	err = json.Unmarshal([]byte(body), &realBody)
	if err != nil {
		fmt.Println(err)
		m.file.WriteString(err.Error() + "\n")
		return nil
	}
	// fmt.Println(realBody)
	var monitorAvailability string
	monitorAvailability = realBody.Online[0].DeliveryMessage.OnlineDeliveryMessage.Key
	m.monitorProduct.stockNumber = realBody.Online[0].AvailableQuantity
	m.monitorProduct.productId = realBody.Online[0].SkuID
	if m.Availability == "OUT_OF_STOCK_ONLINE" && monitorAvailability == "IN_STOCK_ONLINE" {
		fmt.Println("Item in Stock")
		m.sendWebhook()
	}
	if m.Availability == "IN_STOCK_ONLINE" && monitorAvailability == "OUT_OF_STOCK_ONLINE" {
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
	// payload := strings.NewReader("{\"content\":null,\"embeds\":[{\"title\":\"Target Monitor\",\"url\":\"https://discord.com/developers/docs/resources/channel#create-message\",\"color\":507758,\"fields\":[{\"name\":\"Product Name\",\"value\":\"%s\"},{\"name\":\"Product Availability\",\"value\":\"In Stock\\u0021\",\"inline\":true},{\"name\":\"Stock Number\",\"value\":\"%s\",\"inline\":true},{\"name\":\"Links\",\"value\":\"[Product](https://www.walmart.com/ip/prada/%s)\"}],\"footer\":{\"text\":\"Prada#4873\"},\"timestamp\":\"2021-04-01T18:40:00.000Z\",\"thumbnail\":{\"url\":\"https://cdn.discordapp.com/attachments/815507198394105867/816741454922776576/pfp.png\"}}],\"avatar_url\":\"https://cdn.discordapp.com/attachments/815507198394105867/816741454922776576/pfp.png\"}")
	payload := strings.NewReader(fmt.Sprintf(`{
  "content": null,
  "embeds": [
    {
      "title": "%s Monitor",
      "url": "https://www.academy.com/%s",
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
          "name": "Product ID",
          "value": "%s",
          "inline": true
        },
        {
          "name": "Stock Number",
          "value": "%s"
        },
        {
          "name": "Links",
          "value": "[Product](https://www.academy.com/%s) | [Cart](https://www.academy.com/shop/cart)"
        }
      ],
      "footer": {
        "text": "Prada#4873"
      },
      "timestamp": "2021-04-01T18:40:00.000Z",
      "thumbnail": {
        "url": "https:%s"
      }
    }
  ],
  "avatar_url": "https://cdn.discordapp.com/attachments/815507198394105867/816741454922776576/pfp.png"
}`, m.Config.site, m.monitorProduct.link, m.monitorProduct.name, m.monitorProduct.price, m.Config.sku, m.monitorProduct.stockNumber, m.monitorProduct.link, m.monitorProduct.image))
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

func (m *Monitor) getProductDetails() {
	defer func() {
	     if r := recover(); r != nil {
	        	        fmt.Printf("Site : %s, Product : %s Recovering from panic in printAllOperations error is: %v \n", m.Config.site, m.Config.sku, r)
	    }
	  }()

	url := fmt.Sprintf("https://www.academy.com/api/productinfo/v2?productIds=%s&summary=true&dbFallback=true", m.Config.sku)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		m.file.WriteString(err.Error() + "\n")
		return
	}
	req.Header.Add("authority", "www.academy.com")
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
			m.file.WriteString(err.Error() + "\n")
			return 
		}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
			fmt.Println(err)
			m.file.WriteString(err.Error() + "\n")
			return
		}
		var productBody ProductDetails
			err = json.Unmarshal([]byte(body), &productBody)
			if err != nil {
				fmt.Println(err)
				m.file.WriteString(err.Error() + "\n")
				return
			}
	m.monitorProduct.price = int(productBody.Productinfo[0].PriceUSD)
	m.monitorProduct.link = productBody.Productinfo[0].UpdatedSeoURL
	m.monitorProduct.name = productBody.Productinfo[0].Name
	m.monitorProduct.image = productBody.Productinfo[0].FullImage
}
