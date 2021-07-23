package GameStopMonitor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/bradhe/stopwatch"
	"github.com/elgs/gojq"
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
	Availability        string
	AvailabilityBool     bool
	currentAvailability string
	Client              http.Client
	file                *os.File
	stop                bool
	CurrentCompanies    []Company
	statusCode          int
}
type Product struct {
	name        string
	stockNumber string
	productId   string
	price       string
	image       string
	link        string
	sku         string
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

func NewMonitor(sku string) *Monitor {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Site : %s, Product : %s Recovering from panic in printAllOperations error is: %v \n", sku, sku, r)
		}
	}()
	fmt.Println("TESTING")
	m := Monitor{}
	m.Availability = "Not Available"
	// var err error
	//	m.Client = http.Client{Timeout: 10 * time.Second}
	m.Config.site = "GameStop"
	m.Config.startDelay = 3000
	m.Config.sku = sku
	// 	m.file, err = os.Create("./testing.txt")
	m.Client = http.Client{Timeout: 10 * time.Second}
	//	m.Config.discord = "https://discord.com/api/webhooks/838637042119213067/V7WQ7z-9u32rNh5SO4YyxS5kibcHadXW4FxjVJTosO5cSGRoSqv4CY5g3GrAcIcwZhkF"
	m.monitorProduct.name = "Testing Product"
	m.monitorProduct.stockNumber = ""

	proxyList := FetchProxies.Get()
	time.Sleep(15000 * (time.Millisecond))
	go m.checkStop()
	time.Sleep(3000 * (time.Millisecond))
	i := true
	for i {
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
			time.Sleep(100 * (time.Millisecond))
		} else {
			fmt.Println(m.Config.sku, "STOPPED STOPPED STOPPED")
			i = false
		}

	}
	return &m
}

func (m *Monitor) monitor() error {
	watch := stopwatch.Start()
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Site : %s, Product : %s Recovering from panic in printAllOperations error is: %v \n", m.Config.site, m.Config.sku, r)
		}
	}()
	url := fmt.Sprintf("https://www.gamestop.com/on/demandware.store/Sites-gamestop-us-Site/default/Product-Variation?dwvar_%s_condition=New&pid=%s&quantity=1&redesignFlag=true&rt=productDetailsRedesign", m.Config.sku, m.Config.sku)
	req, _ := http.NewRequest("GET", url, nil)
	gmeRef := fmt.Sprintf("https://www.gamestop.com/products/prada/%s.html", m.Config.sku)
	req.Header.Add("authority", "www.gamestop.com")
	req.Header.Add("pragma", "no-cache")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("accept", "*/*")
	req.Header.Add("dnt", "1")
	req.Header.Add("x-requested-with", "XMLHttpRequest")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("referer", gmeRef)
	req.Header.Add("accept-language", "en-US,en;q=0.9")
	req.Header.Set("Connection", "close")
	req.Close = true
	res, err := m.Client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer res.Body.Close()
	defer func() {
		watch.Stop()
		fmt.Printf("Game Stop - Status Code : %d, %s, %s, Milliseconds elapsed: %v \n", m.statusCode, m.Availability, m.Config.sku, watch.Milliseconds())
	}()
	m.statusCode = res.StatusCode
	if res.StatusCode != 200 {
		return nil
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	parser, err := gojq.NewStringQuery(string(body))
	if err != nil {
		fmt.Println(err)
		return nil
	}

	monitorAvailability, err := parser.QueryToString("product.available")
	if err != nil {
		fmt.Println(err)
	}
	m.monitorProduct.name, err = parser.QueryToString("gtmData.productInfo.name")
	if err != nil {
		fmt.Println(err)
	}
	m.monitorProduct.sku, err = parser.QueryToString("gtmData.productInfo.sku")
	if err != nil {
		fmt.Println(err)
	}
	m.monitorProduct.productId, err = parser.QueryToString("gtmData.productInfo.productID")
	if err != nil {
		fmt.Println(err)
	}
	m.monitorProduct.price, err = parser.QueryToString("gtmData.price.sellingPrice")
	//m.monitorProduct.link, err = parser.QueryToString("__mccEvents.[0].[1].[0].url")
	if err != nil {
		fmt.Println(err)
	}
	m.monitorProduct.image, err = parser.QueryToString("product.images.large.[0].url")

	//	fmt.Println(m.monitorProduct.link, m.monitorProduct.image)
	if err != nil {
		fmt.Println(err)
	}
	var currentAvailabilityBool bool = false
	if monitorAvailability == "true" {
		currentAvailabilityBool = true
	} else {
		currentAvailabilityBool = false
	}

	if !m.AvailabilityBool && currentAvailabilityBool {
		fmt.Println("Item in Stock")
		m.sendWebhook()
	} else if m.AvailabilityBool && !currentAvailabilityBool {
		fmt.Println("Item Out Of Stock", m.Config.sku)
	}
	m.AvailabilityBool = currentAvailabilityBool
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
	fmt.Println("Testing Here : ", m.monitorProduct.name, "Here")
	if strings.HasSuffix(m.monitorProduct.name, "                       ") {
		m.monitorProduct.name = strings.Replace(m.monitorProduct.name, "                       ", "", -1)
	}
	fmt.Println("Testing Here : ", m.monitorProduct.name, "Here")
	// now := time.Now()
	// currentTime := strings.Split(now.String(), "-0400")[0]
	t := time.Now().UTC().Format("2006-01-02T15:04:05Z")

	for _, comp := range m.CurrentCompanies {
		fmt.Println(comp.Company)
		go webHookSend(comp, m.Config.site, m.Config.sku, m.monitorProduct.sku, m.monitorProduct.name, m.monitorProduct.price, t, m.monitorProduct.image)
	}
	// payload := strings.NewReader("{\"content\":null,\"embeds\":[{\"title\":\"Target Monitor\",\"url\":\"https://discord.com/developers/docs/resources/channel#create-message\",\"color\":507758,\"fields\":[{\"name\":\"Product Name\",\"value\":\"%s\"},{\"name\":\"Product Availability\",\"value\":\"In Stock\\u0021\",\"inline\":true},{\"name\":\"Stock Number\",\"value\":\"%s\",\"inline\":true},{\"name\":\"Links\",\"value\":\"[Product](https://www.walmart.com/ip/prada/%s)\"}],\"footer\":{\"text\":\"Prada#4873\"},\"timestamp\":\"2021-04-01T18:40:00.000Z\",\"thumbnail\":{\"url\":\"https://cdn.discordapp.com/attachments/815507198394105867/816741454922776576/pfp.png\"}}],\"avatar_url\":\"https://cdn.discordapp.com/attachments/815507198394105867/816741454922776576/pfp.png\"}")
	return nil
}
func webHookSend(c Company, site string, pid string, sku string, name string, price string, time string, image string) {
	payload := strings.NewReader(fmt.Sprintf(`{
		"content": null,
		"embeds": [
		  {
			"title": "%s Monitor",
			"url": "https://www.gamestop.com/products/prada/%s.html",
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
				"value": "%s",
				"inline": true
			  },
			  {
				"name": "Product ID",
				"value": "%s",
				"inline": true
			  },
			  {
				"name": "Links",
				"value": "[Product](https://www.gamestop.com/products/prada/%s.html) | [Cart](https://www.gamestop.com/cart/) | [Checkout](https://www.gamestop.com/checkout/)"
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
	  }`, site, pid, c.Color, name, price, pid, pid, time, image, c.CompanyImage))
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
