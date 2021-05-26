package TargetNewTradingCards

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

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
	products            []string
	keywords            []string
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
type targetNewProduct struct {
	Typename string `json:"__typename"`
	Item     struct {
		CartAddOnThreshold int64 `json:"cart_add_on_threshold"`
		ChokingHazard      []struct {
			Code string `json:"code"`
		} `json:"choking_hazard"`
		Dpci       string `json:"dpci"`
		Enrichment struct {
			BuyURL string `json:"buy_url"`
			Images struct {
				AlternateImageUrls []string `json:"alternate_image_urls"`
				PrimaryImageURL    string   `json:"primary_image_url"`
			} `json:"images"`
		} `json:"enrichment"`
		Fulfillment               struct{} `json:"fulfillment"`
		MerchandiseClassification struct {
			ClassID      int64 `json:"class_id"`
			DepartmentID int64 `json:"department_id"`
		} `json:"merchandise_classification"`
		PrimaryBrand struct {
			CanonicalURL string `json:"canonical_url"`
			FacetID      string `json:"facet_id"`
			Name         string `json:"name"`
		} `json:"primary_brand"`
		ProductDescription struct {
			BulletDescriptions []string `json:"bullet_descriptions"`
			SoftBullets        struct {
				Bullets []string `json:"bullets"`
			} `json:"soft_bullets"`
			Title string `json:"title"`
		} `json:"product_description"`
		ProductVendors []struct {
			ID         string `json:"id"`
			VendorName string `json:"vendor_name"`
		} `json:"product_vendors"`
		RelationshipType     string `json:"relationship_type"`
		RelationshipTypeCode string `json:"relationship_type_code"`
	} `json:"item"`
	OriginalTcin string `json:"original_tcin"`
	Price        struct {
		CurrentRetail             float64 `json:"current_retail"`
		FormattedCurrentPrice     string  `json:"formatted_current_price"`
		FormattedCurrentPriceType string  `json:"formatted_current_price_type"`
	} `json:"price"`
	Promotions        []interface{} `json:"promotions"`
	RatingsAndReviews struct {
		Statistics struct {
			Rating struct {
				Average           int64 `json:"average"`
				Count             int64 `json:"count"`
				SecondaryAverages []struct {
					ID    string `json:"id"`
					Label string `json:"label"`
					Value int64  `json:"value"`
				} `json:"secondary_averages"`
			} `json:"rating"`
		} `json:"statistics"`
	} `json:"ratings_and_reviews"`
	Tcin string `json:"tcin"`
}

var file os.File

func NewMonitor(sku string, keywords []string) *Monitor {
	// fmt.Println("TESTING", sku)
	m := Monitor{}
	m.Availability = false
	// var err error
	m.Client = http.Client{Timeout: 5 * time.Second}
	m.Config.site = "Target New"
	m.Config.startDelay = 3000
	m.Config.sku = sku
	m.Config.skuName = sku
	// 	m.file, err = os.Create("./testing.txt")
	m.Client = http.Client{Timeout: 60 * time.Second}
	m.Config.discord = "https://discord.com/api/webhooks/833902775196450818/wCEYgKpT7eJaOtNERfwe5AlietWcFomGp10zTP3JyDEc8Kk3f2ujqY-BpdXPmpQYANiT"
	m.monitorProduct.name = "Testing Product"
	m.monitorProduct.stockNumber = 10
	m.keywords = keywords
	fmt.Println(keywords)

	proxyList := FetchProxies.Get()

	// fmt.Println(timeout)
	//m.Availability = "OUT_OF_STOCK"
	//fmt.Println(m)
	i := true
	go m.checkStop()
	time.Sleep(3000 * (time.Millisecond))
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
			go m.monitor()
			time.Sleep(300 * (time.Millisecond))
			// fmt.Println(m.Availability)
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
			fmt.Printf("Recovering from panic in printAllOperations error is: %v \n", r)
		}
	}()
	// url := "https://httpbin.org/ip"

	// req, _ := http.NewRequest("GET", url, nil)

	// res, _ := m.Client.Do(req)

	// defer res.Body.Close()
	// body, _ := ioutil.ReadAll(res.Body)

	// fmt.Println(res)
	// fmt.Println(string(body))

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

	//	fmt.Println(string(body))
	fmt.Println(res.StatusCode)
	if res.StatusCode != 200 {
		return nil
	}

	var realBody map[string]interface{}
	err = json.Unmarshal([]byte(body), &realBody)
	if err != nil {
		fmt.Println(err)

		return nil
	}
	products := realBody["data"].(map[string]interface{})["search"].(map[string]interface{})["products"].([]interface{})
	for _, value := range products {
		//	fmt.Println(value)
		var isPresent bool
		// var currentProduct targetNewProduct
		tcin := value.(map[string]interface{})["tcin"].(string)
		price := int(value.(map[string]interface{})["price"].(map[string]interface{})["current_retail"].(float64))
		productName := value.(map[string]interface{})["item"].(map[string]interface{})["product_description"].(map[string]interface{})["title"].(string)
		link := value.(map[string]interface{})["item"].(map[string]interface{})["enrichment"].(map[string]interface{})["buy_url"].(string)
		image := value.(map[string]interface{})["item"].(map[string]interface{})["enrichment"].(map[string]interface{})["images"].(map[string]interface{})["primary_image_url"].(string)
		for _, v := range m.products {
			if v == tcin {
				isPresent = true

			}
		}
		if isPresent == false {
			m.products = append(m.products, tcin)
			go m.sendWebhook(tcin, link, price, productName, image)
			// for _, kw := range m.keywords{
			// 		// fmt.Println(kw, productName)
			// 		if strings.Contains(strings.ToUpper(productName), strings.ToUpper(kw)) {
			// 			// m.products = append(m.products, tcin)
			// 			// go m.sendWebhook(tcin, link, price, productName, image)
			// 		}
			// 	}

		}
		// fmt.Println(m.products)
	}
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
func (m *Monitor) sendWebhook(tcin string, link string, price int, productName string, image string) error {
	// payload := strings.NewReader("{\"content\":null,\"embeds\":[{\"title\":\"Target Monitor\",\"url\":\"https://discord.com/developers/docs/resources/channel#create-message\",\"color\":507758,\"fields\":[{\"name\":\"Product Name\",\"value\":\"%s\"},{\"name\":\"Product Availability\",\"value\":\"In Stock\\u0021\",\"inline\":true},{\"name\":\"Stock Number\",\"value\":\"%s\",\"inline\":true},{\"name\":\"Links\",\"value\":\"[Product](https://www.walmart.com/ip/prada/%s)\"}],\"footer\":{\"text\":\"Prada#4873\"},\"timestamp\":\"2021-04-01T18:40:00.000Z\",\"thumbnail\":{\"url\":\"https://cdn.discordapp.com/attachments/815507198394105867/816741454922776576/pfp.png\"}}],\"avatar_url\":\"https://cdn.discordapp.com/attachments/815507198394105867/816741454922776576/pfp.png\"}")
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
	// payload := strings.NewReader("{\"content\":null,\"embeds\":[{\"title\":\"Target Monitor\",\"url\":\"https://discord.com/developers/docs/resources/channel#create-message\",\"color\":507758,\"fields\":[{\"name\":\"Product Name\",\"value\":\"%s\"},{\"name\":\"Product Availability\",\"value\":\"In Stock\\u0021\",\"inline\":true},{\"name\":\"Stock Number\",\"value\":\"%s\",\"inline\":true},{\"name\":\"Links\",\"value\":\"[Product](https://www.walmart.com/ip/prada/%s)\"}],\"footer\":{\"text\":\"Prada#4873\"},\"timestamp\":\"2021-04-01T18:40:00.000Z\",\"thumbnail\":{\"url\":\"https://cdn.discordapp.com/attachments/815507198394105867/816741454922776576/pfp.png\"}}],\"avatar_url\":\"https://cdn.discordapp.com/attachments/815507198394105867/816741454922776576/pfp.png\"}")
	for _, comp := range m.CurrentCompanies {
		fmt.Println(comp.Company)
		fmt.Println(comp.Company)
		go webHookSend(comp, m.Config.site, tcin, m.monitorProduct.name, price, "test", image, link)
	}
	return nil
}
func webHookSend(c Company, site string, sku string, name string, price int, time string, image string, link string) {
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
				"name": "Tcin",
				"value": "%s",
				"inline": true
			  },
			  {
				"name": "Links",
				"value": "[Product](%s) | [Cart](https://www.target.com/co-cart)"
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
	  }`, site, link, c.Color, name, price, sku, link, image, c.CompanyImage))
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
func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
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
