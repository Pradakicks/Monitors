package BigLotsMonitor

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

	"github.com/bradhe/stopwatch"
)

type Config struct {
	sku string
	//	skuName string // Only for new Egg
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
}
type Product struct {
	name        string
	stockNumber int
	offerId     string
	price       string
	image       string
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

func NewMonitor(sku string, priceRangeMin int, priceRangeMax int) *Monitor {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Site : %s, Product : %s Recovering from panic in printAllOperations error is: %v \n", sku, sku, r)
		}
	}()
	fmt.Println("TESTING", sku, priceRangeMin, priceRangeMax)
	m := Monitor{}
	m.Availability = false
	var err error
	m.Client = http.Client{Timeout: 5 * time.Second}
	m.Config.site = "Big Lots"
	m.Config.startDelay = 3000
	m.Config.sku = sku
	// 	m.file, err = os.Create("./testing.txt")
	m.Client = http.Client{Timeout: 5 * time.Second}
	m.Config.discord = "https://discord.com/api/webhooks/833531825951080478/DZcTzNJbZmfcq8KpRJFNJVunFnQj48QdGg6EIecHvmUkucldj-0q6UZdhZv7H75OWdqj"
	m.monitorProduct.name = "Testing Product"
	m.monitorProduct.stockNumber = 10
	m.Config.priceRangeMax = priceRangeMax
	m.Config.priceRangeMin = priceRangeMin

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
	go m.checkStop()
	i := true
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
			watch := stopwatch.Start()
			m.monitor()
			watch.Stop()
			fmt.Printf("Milliseconds elapsed: %v\n", watch.Milliseconds())
			// time.Sleep(500 * (time.Millisecond))
			//fmt.Println(m.Availability)
		} else {
			fmt.Println(m.Config.sku, "STOPPED STOPPED STOPPED")
			i = false
		}

	}
	return &m
}

func (m *Monitor) monitor() error {
	//	fmt.Println("Monitoring")
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Site : %s, Product : %s Recovering from panic in printAllOperations error is: %v \n", m.Config.site, m.Config.sku, r)
		}
	}()
	// url := "https://httpbin.org/ip"

	// req, _ := http.NewRequest("GET", url, nil)

	// res, _ := m.Client.Do(req)

	// defer res.Body.Close()
	// body, _ := ioutil.ReadAll(res.Body)

	// fmt.Println(res)
	// fmt.Println(string(body))

	//	url := fmt.Sprintf("%s", m.Config.sku)
	req, err := http.NewRequest("GET", m.Config.sku, nil)
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
	// fmt.Println(l)
	fmt.Println(res.StatusCode)
	if res.StatusCode != 200 {
		return nil
	}
	var monitorAvailability bool
	// g, err := os.Create("testing.html")
	// 	if err != nil {
	// 		fmt.Println(err)
	//
	// 		return nil
	// 	}
	// 	l, err := g.WriteString(string(body))
	// 	if err != nil {
	// 		fmt.Println(err)
	//
	// 		return nil
	// 	}
	// 	fmt.Println(l)
	itemOutOfStockCheck := strings.Contains(string(body), "This item is currently out of stock")
	inStockOnlineCheck := strings.Contains(string(body), "inStockOnline")
	fmt.Println(itemOutOfStockCheck, inStockOnlineCheck)
	if itemOutOfStockCheck == true && inStockOnlineCheck == false {
		monitorAvailability = false
	} else if itemOutOfStockCheck == false && inStockOnlineCheck == true {

		isItReallyInStock := strings.Split(string(body), "isInStock :	")[1]
		isItReallyInStock = isItReallyInStock[0:3]

		if isItReallyInStock == "true" {
			monitorAvailability = true
		} else if isItReallyInStock == "fals" || isItReallyInStock == "fal" {
			fmt.Println("Out of Stock Online")
			monitorAvailability = true
		}

		fmt.Println("Big Lots", isItReallyInStock)

		m.monitorProduct.image = strings.Split((strings.Split(string(body), `data-resolvechain="set=`)[1]), `"`)[0]
		m.monitorProduct.name = strings.Split(strings.Split(string(body), `<div class="product-name">
					<h1>`)[1], "</h1>")[0]
		m.monitorProduct.price = (strings.Split(strings.Split(strings.Split(string(body), `<div class="regular-price">`)[1], "$")[1], "</div>")[0])

	}

	fmt.Println("Big Lots : ", res.StatusCode, itemOutOfStockCheck, inStockOnlineCheck, monitorAvailability, m.monitorProduct.name, m.monitorProduct.image)
	//	l, err := g.WriteString(string(body))
	// htmlSplit := (strings.Split(string(body), "inStockOnline : ")[1])
	// finalSplit := (strings.Split(htmlSplit, "<br")[0])
	// fmt.Println(finalSplit)
	// switch finalSplit {
	// case "true":
	// 	monitorAvailability = true
	// 	break
	// case "false":
	// 	monitorAvailability = false
	// default:
	// 	fmt.Println(finalSplit, "Unknown")
	// 	fmt.Println(finalSplit, "Unknown")
	// 	fmt.Println(finalSplit, "Unknown")
	// 	g, err := os.Create("testing.html")
	// 	if err != nil {
	// 		fmt.Println(err)
	//
	// 		return nil
	// 	}
	// 	l, err := g.WriteString(string(body))
	// 	if err != nil {
	// 		fmt.Println(err)
	//
	// 		return nil
	// 	}
	// 	fmt.Println(l)
	// }
	//fmt.Println("Other Check ", m.Availability, monitorAvailability)
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
	now := time.Now()
	currentTime := strings.Split(now.String(), "-0400")[0]
	// payload := strings.NewReader("{\"content\":null,\"embeds\":[{\"title\":\"Target Monitor\",\"url\":\"https://discord.com/developers/docs/resources/channel#create-message\",\"color\":507758,\"fields\":[{\"name\":\"Product Name\",\"value\":\"%s\"},{\"name\":\"Product Availability\",\"value\":\"In Stock\\u0021\",\"inline\":true},{\"name\":\"Stock Number\",\"value\":\"%s\",\"inline\":true},{\"name\":\"Links\",\"value\":\"[Product](https://www.walmart.com/ip/prada/%s)\"}],\"footer\":{\"text\":\"Prada#4873\"},\"timestamp\":\"2021-04-01T18:40:00.000Z\",\"thumbnail\":{\"url\":\"https://cdn.discordapp.com/attachments/815507198394105867/816741454922776576/pfp.png\"}}],\"avatar_url\":\"https://cdn.discordapp.com/attachments/815507198394105867/816741454922776576/pfp.png\"}")
	payload := strings.NewReader(fmt.Sprintf(`{
  "content": null,
  "embeds": [
    {
      "title": "%s Monitor",
      "url": "%s",
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
          "value": "$%s",
          "inline": true
        },
        {
          "name": "Links",
          "value": "[Product](%s) | [Cart](https://www.biglots.com/cart/cart.jsp)"
        }
      ],
      "footer": {
        "text": "Prada#4873"
      },
      "timestamp": "%s",
      "thumbnail": {
        "url": "https://images.biglots.com/images?set=key[resolve.pixelRatio],value[1]&set=key[resolve.width],value[600]&set=key[resolve.height],value[10000]&set=key[resolve.imageFit],value[containerwidth]&set=key[resolve.allowImageUpscaling],value[0]&set=env[prod],%s"
      }
    }
  ],
  "avatar_url": "https://cdn.discordapp.com/attachments/815507198394105867/816741454922776576/pfp.png"
}`, m.Config.site, m.Config.sku, m.monitorProduct.name, m.monitorProduct.price, m.Config.sku, currentTime, m.monitorProduct.image))
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
