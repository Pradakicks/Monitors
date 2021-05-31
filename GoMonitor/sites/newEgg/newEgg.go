package NewEggMonitor

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
	CurrentCompanies    []Company
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
type NewEggResponse struct {
	MainItem struct {
		Description struct {
			ShortTitle string `json:"ShortTitle"`
			Title      string `json:"Title"`
		} `json:"Description"`
		FinalPrice float64 `json:"FinalPrice"`
		Image      struct {
			ImagePathPattern []struct {
				PathPattern string `json:"PathPattern"`
				Size        int64  `json:"Size"`
			} `json:"ImagePathPattern"`
			ImageType         int64  `json:"ImageType"`
			IsDFISImage       bool   `json:"IsDFISImage"`
			ItemCellImageName string `json:"ItemCellImageName"`
			Normal            struct {
				ImageName     string `json:"ImageName"`
				ImageNameList string `json:"ImageNameList"`
			} `json:"Normal"`
		} `json:"Image"`
		Instock  bool `json:"Instock"`
		NewImage struct {
			DFIS360ImgFlag interface{} `json:"DFIS360ImgFlag"`
			ImageName      string      `json:"ImageName"`
			ImageNameList  string      `json:"ImageNameList"`
		} `json:"NewImage"`
		Stock         int64 `json:"Stock"`
		StockForCombo int64 `json:"StockForCombo"`
	} `json:"MainItem"`
}

func NewMonitor(sku string, skuName string, priceRangeMin int, priceRangeMax int) *Monitor {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Site : %s, Product : %s Recovering from panic in printAllOperations error is: %v \n", sku, sku, r)
		}
	}()
	fmt.Println("TESTING", sku, skuName, priceRangeMin, priceRangeMax)
	m := Monitor{}
	m.Availability = false
	// var err error
	m.Client = http.Client{Timeout: 10 * time.Second}
	m.Config.site = "NewEgg"
	m.Config.startDelay = 3000
	m.Config.sku = sku
	m.Config.skuName = skuName
	// 	m.file, err = os.Create("./testing.txt")
	m.Client = http.Client{Timeout: 60 * time.Second}
	m.Config.discord = "https://discord.com/api/webhooks/816740348222767155/2APr1EdhzNO4hRWznexhMRlO0g7qOiCkI7HFtmuU7_r48PCWnGYmSTGJmRVX0LPCNN_t"
	m.monitorProduct.name = "Testing Product"
	m.monitorProduct.stockNumber = 10
	m.Config.priceRangeMax = priceRangeMax
	m.Config.priceRangeMin = priceRangeMin

	proxyList := FetchProxies.Get()

	// fmt.Println(timeout)
	//m.Availability = "OUT_OF_STOCK"
	//fmt.Println(m)
	go m.checkStop()
	time.Sleep(3000 * (time.Millisecond))
	i := true
	for i == true {
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
			go m.monitor()
			time.Sleep(250 * (time.Millisecond))
		} else {
			fmt.Println(m.Config.skuName, "STOPPED STOPPED STOPPED")
			i = false
		}

	}
	return &m
}

func (m *Monitor) monitor() error {
	watch := stopwatch.Start()
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Site : %s, Product : %s Recovering from panic in printAllOperations error is: %v \n", m.Config.site, m.Config.skuName, r)
		}
	}()
	url := fmt.Sprintf("https://www.newegg.com/product/api/ProductRealtime?ItemNumber=%s", m.Config.skuName)
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
	var monitorAvailability bool
	defer res.Body.Close()
	defer func() {
		watch.Stop()
		fmt.Printf("New Egg Monitor : %s, %d, %t, %t, Status Code : %d, %d Milliseconds elapsed: %v \n", m.monitorProduct.name, m.monitorProduct.price, m.Availability, monitorAvailability, m.monitorProduct.stockNumber, res.StatusCode, watch.Milliseconds())
	}()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if res.StatusCode != 200 {
		return nil
	}
	var realBody NewEggResponse
	err = json.Unmarshal(body, &realBody)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	m.monitorProduct.name = realBody.MainItem.Description.Title
	m.monitorProduct.price = int(realBody.MainItem.FinalPrice)
	monitorAvailability = realBody.MainItem.Instock
	m.monitorProduct.stockNumber = int(realBody.MainItem.Stock)
	m.monitorProduct.image = realBody.MainItem.Image.ItemCellImageName

	if m.Availability == false && monitorAvailability == true {
		fmt.Println("Item in Stock")
		go m.sendWebhook()
	}
	if m.Availability == true && monitorAvailability == false {
		fmt.Println("Item Out Of Stock")
	}
	m.Availability = monitorAvailability
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
			fmt.Printf("Site : %s, Product : %s Recovering from panic in printAllOperations error is: %v \n", m.Config.site, m.Config.skuName, r)
		}
	}()
	for _, letter := range m.monitorProduct.name {
		if string(letter) == `"` {
			m.monitorProduct.name = strings.Replace(m.monitorProduct.name, `"`, "", -1)
		}
	}
	for _, comp := range m.CurrentCompanies {
		fmt.Println(comp.Company)
		go webHookSend(comp, m.Config.site, m.Config.skuName, m.monitorProduct.name, m.monitorProduct.price, m.monitorProduct.stockNumber, "test", m.monitorProduct.image)
	}
	return nil
}
func webHookSend(c Company, site string, sku string, name string, price int, stockNumber int, time string, image string) {
	payload := strings.NewReader(fmt.Sprintf(`{
		"content": null,
		"embeds": [
		  {
			"title": "%s Monitor",
			"url": "https://www.newegg.com/Prada/p/%s",
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
				"name": "Stock Number",
				"value": "%d"
			  },
			  {
				"name": "Links",
				"value": "[ATC](https://secure.newegg.com/shopping/addtocart.aspx?submit=add&itemList=%s) | [Cart](https://secure.newegg.com/shop/cart)"
			  }
			],
			"footer": {
			  "text": "Prada#4873"
			},
			"timestamp": "2021-05-13 13:57:26.5157268",
			"thumbnail": {
			  "url": "http://c1.neweggimages.com/ProductImageOriginal/%s"
			}
		  }
		],
		"avatar_url": "%s"
	  }`, site, sku, c.Color, name, price, sku, stockNumber, sku, image, c.CompanyImage))
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
		url := fmt.Sprintf("http://localhost:7243/DB")
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
