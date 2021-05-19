package GameStopMonitor

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

	"github.com/elgs/gojq"
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
	m.Availability = "Not Available"
	var err error
	//	m.Client = http.Client{Timeout: 5 * time.Second}
	m.Config.site = "GameStop"
	m.Config.startDelay = 3000
	m.Config.sku = sku
	// 	m.file, err = os.Create("./testing.txt")
	m.Client = http.Client{Timeout: 5 * time.Second}
	//	m.Config.discord = "https://discord.com/api/webhooks/838637042119213067/V7WQ7z-9u32rNh5SO4YyxS5kibcHadXW4FxjVJTosO5cSGRoSqv4CY5g3GrAcIcwZhkF"
	m.monitorProduct.name = "Testing Product"
	m.monitorProduct.stockNumber = ""

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
	time.Sleep(5000 * (time.Millisecond))
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
			fmt.Println("Gamestop : ", m.Availability, m.Config.sku)
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

	url := fmt.Sprintf("https://www.gamestop.com/on/demandware.store/Sites-gamestop-us-Site/default/Product-Variation?dwvar_%s_condition=New&pid=%s&quantity=1&redesignFlag=true&rt=productDetailsRedesign", "11107421", "11107421")

	req, _ := http.NewRequest("GET", url, nil)

	//	req.Header.Add("cookie", "dwac_a78eb5a11975e4c9cbc84f55ad=lDPRds9spMpfapQx2a5IJvTxNZRHi6FaIZ4%253D%7Cdw-only%7C%7C%7CUSD%7Cfalse%7CUS%252FCentral%7Ctrue; cqcid=ab5UPFZX9yyAfmyTsKUgaaHFzv; cquid=%7C%7C; customergroups=Everyone; sid=lDPRds9spMpfapQx2a5IJvTxNZRHi6FaIZ4; dwanonymous_420142ceefb9f0c103b3815e84e9fcef=ab5UPFZX9yyAfmyTsKUgaaHFzv; userInfo=S%3DN%7CR%3DN%7CC%3D0; __cq_dnt=0; dw_dnt=0; dwsid=JfGqofdGhc81RmDriGdLF0UzTkN_oYBfS7KwKQPnygMp8iGqHSsVqNEFDYL5a9uJD-pAPTI6JwAdz4I0JZZG4g%3D%3D; ak_bmsc=26E883DAE75F5DAB598B6520FAB3C1AF172F3B79004F0000A208A560F792450A~plYFxZ4LLyXtlAD7MSqCNyGfcpA%2BWSiKCoR95Nfat1xmFMP3elINjIeE0GfbQqkAkePHOfzTQql%2B3i0Iu3OdnxOByFIDCKtv41uY5DMVseM7zfvEUjSmK37AwXRksT27pC%2F%2Fr3K2xc0HE8Vh1MsTvqGeRGByEft4Bs8MvOXiIMuW5KahPxceboP%2FPzuvr2BhYFKnZQGK44XUoTYq8%2F0BoLlxhrNxaeZR3FD27sP3uQ09k%3D; akaas_ChatThrottling=2147483647~rv%3D96~id%3D63cbc0fee3058ac39b466275a70905af~rn%3D; bm_sz=153F0FCB0D702CB72287C3144158AC33~YAAQeTsvFx4uE4J5AQAAzLmphAu7dPhH9S0t7AQUZDlvUSl3nExOfSNMp18%2Fsm3I%2FefgUtZV4UEJsEFLuRR5j0MFhJrLnU8i8AlGpMe%2BEpZpKhfpEu3axLjzD1V%2B%2FdbtQRr66TW8SCau9i3N6x2ZOUWepDFw9nptQSynheLmogi0vKi%2FJuryNnsaFVW8cBMDl5c%3D; _abck=87E16D9492B4FF25B48E128780C8B242~-1~YAAQeTsvFx8uE4J5AQAAzLmphAUAxZb3DtCi23urHp7Lq1grwAu7WIP2xYHzVxEQlJE%2Bx9M%2FOJf8PcyWd0Bv0f3ijc4uSvE7op5LCyGNhMH3uWru3VTaoG%2B7Mt5S%2By4TULji2CwzJud6bs%2B9XVBCkXPaCTP8guUwNLWPwSue%2F4cGt%2BjPeFZ42%2FlFlORmfTWZN1v9HY58Jro6%2BEDTANpmpxRvyR26Uv5%2B5FT7Lrs7Szb5nW4gmkYj8TSZ1x5wBOAusRvn2ab1GowWf0jCYfK1uAmsVEBdy5JzqUrMCQZekwTJjLPRHytwPMMuTpjYjrzRyQ%2FeUo1kgNqC3mMPez4arryvm7jwTUOcRSXEsNKgJvyHR3qngSiZnWPHlrAMwA%3D%3D~-1~-1~-1")
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
	req.Header.Add("referer", "https://www.gamestop.com/video-games/pc-gaming/virtual-reality/products/oculus-quest-2-64gb/11107421.html")
	req.Header.Add("accept-language", "en-US,en;q=0.9")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	parser, err := gojq.NewStringQuery(string(body))
	monitorAvailability, err := parser.QueryToString("gtmData.productInfo.availability")
	m.monitorProduct.name, err = parser.QueryToString("gtmData.productInfo.name")
	m.monitorProduct.sku, err = parser.QueryToString("gtmData.productInfo.sku")
	m.monitorProduct.productId, err = parser.QueryToString("gtmData.productInfo.productID")
	m.monitorProduct.price, err = parser.QueryToString("gtmData.price.sellingPrice")
	//m.monitorProduct.link, err = parser.QueryToString("__mccEvents.[0].[1].[0].url")
	m.monitorProduct.image, err = parser.QueryToString("product.images.large.[0].url")
	fmt.Println(m.monitorProduct.link, m.monitorProduct.image)
	if err != nil {
		fmt.Println(err)
	}
	if m.Availability == "Not Available" && monitorAvailability == "Available" {
		fmt.Println("Item in Stock")
		m.sendWebhook()
	} else if m.Availability == "Available" && monitorAvailability == "Not Available" {
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
	// now := time.Now()
	// currentTime := strings.Split(now.String(), "-0400")[0]
	for _, comp := range m.CurrentCompanies {
		fmt.Println(comp.Company)
		go webHookSend(comp, m.Config.site, m.Config.sku, m.monitorProduct.sku, m.monitorProduct.name, m.monitorProduct.price, "test", m.monitorProduct.image)
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
			"timestamp": "2021-05-13 13:57:26.5157268",
			"thumbnail": {
			  "url": "%s"
			}
		  }
		],
		"avatar_url": "%s"
	  }`, site, pid, c.Color, name, price, pid, pid, image, c.CompanyImage))
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
