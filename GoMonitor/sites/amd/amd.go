package AmdMonitor

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
	Sku  string `json:"sku"`
	Site string `json:"site"`
	Stop bool   `json:"stop"`
	Name string `json:"name"`
}

var file os.File

// func walmartMonitor(sku string) {
// 	go NewMonitor(sku, 1, 1000)
// 	fmt.Scanln()
// }

func NewMonitor(sku string) *Monitor {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovering from panic in printAllOperations error is: %v \n", r)
		}
	}()
	fmt.Println("TESTING")
	m := Monitor{}
	m.Availability = "false"
	var err error
	//	m.Client = http.Client{Timeout: 5 * time.Second}
	m.Config.site = "Amd"
	m.Config.startDelay = 3000
	m.Config.sku = sku
	// 	m.file, err = os.Create("./testing.txt")
	m.Client = http.Client{Timeout: 5 * time.Second}
	m.Config.discord = "https://discord.com/api/webhooks/827263591114997790/chAZK84Gnad7rjHDlh4BnF7dz5KQ7-0l4atsFzJGgcTkAaeZno6ePYB_A-WiiClS3FpY"
	m.monitorProduct.name = "Testing Product"
	m.monitorProduct.stockNumber = ""

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
	buf.Close()
	err = snl.Err()
	if err != nil {
		fmt.Println(err)
	}

	// fmt.Println(timeout)
	//m.Availability = "OUT_OF_STOCK"
	//fmt.Println(m)
	go m.checkStop()
	i := true
	for i == true {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Recovering from panic in printAllOperations error is: %v \n", r)
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
			fmt.Println("AMD : ", m.Availability, m.Config.sku)
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
			fmt.Printf("Recovering from panic in printAllOperations error is: %v \n", r)
		}
	}()
	//	fmt.Println("Monitoring")
	// 	defer func() {
	//      if r := recover(); r != nil {
	//         fmt.Printf("Recovering from panic in printAllOperations error is: %v \n", r)
	//     }
	//   }()
	// url := "https://httpbin.org/ip"

	// req, _ := http.NewRequest("GET", url, nil)

	// res, _ := m.Client.Do(req)

	// defer res.Body.Close()
	// body, _ := ioutil.ReadAll(res.Body)

	// fmt.Println(res)
	// fmt.Println(string(body))

	url := fmt.Sprintf("https://www.amd.com/en/direct-buy/add-to-cart/%s", m.Config.sku)
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
	fmt.Println("AMD", res.StatusCode)
	if res.StatusCode != 200 {
		return nil
	}
	jsonBody := []byte(strings.Split(strings.Split(string(body), "<textarea>")[1], "</textarea>")[0])
	fmt.Println(string(jsonBody))
	// var realBody bestBuyResponse
	// err = json.Unmarshal(jsonBody, &realBody)
	// if err != nil {
	// 	fmt.Println(err)
	//
	// 	return nil
	// }
	// // fmt.Println(realBody)
	// var monitorAvailability string
	// monitorAvailability = realBody[0].Sku.ButtonState.ButtonState
	// m.monitorProduct.name = realBody[0].Sku.Names.Short
	// m.monitorProduct.price = int(realBody[0].Sku.Price.CurrentPrice)
	// if m.Availability == "SOLD_OUT" && monitorAvailability == "ADD_TO_CART" {
	// 	fmt.Println("Item in Stock")
	// 	m.sendWebhook()
	// }
	// if m.Availability == "ADD_TO_CART" && monitorAvailability == "SOLD_OUT" {
	// 	fmt.Println("Item Out Of Stock")
	// }
	// m.Availability = monitorAvailability
	return nil
}

func (m *Monitor) getProxy(proxyList []string) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovering from panic in printAllOperations error is: %v \n", r)
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
			fmt.Printf("Recovering from panic in printAllOperations error is: %v \n", r)
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
      "url": "https://www.bestbuy.com/site/prada/%s.p?skuId=%s",
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
          "name": "Links",
          "value": "[ATC](https://api.bestbuy.com/click/tempo/%s/cart) | [Cart](https://www.bestbuy.com/cart)"
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
}`, m.Config.site, m.Config.sku, m.Config.sku, m.monitorProduct.name, m.monitorProduct.price, m.Config.sku, m.Config.sku, m.monitorProduct.image))
	req, err := http.NewRequest("POST", m.Config.discord, payload)
	if err != nil {
		fmt.Println(err)
		fmt.Println(payload)

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

		return nil
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		fmt.Println(payload)

		return nil
	}
	fmt.Println(res)
	fmt.Println(string(body))
	fmt.Println(payload)
	return nil
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
		res.Body.Close()
		fmt.Println(currentObject)
		time.Sleep(3500 * (time.Millisecond))
	}
	return nil
}
