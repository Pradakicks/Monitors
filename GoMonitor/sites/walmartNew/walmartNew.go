package WalmartNew

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/bradhe/stopwatch"
	"github.com/elgs/gojq"
	MonitorLogger "github.con/prada-monitors-go/helpers/logging"
	FetchProxies "github.con/prada-monitors-go/helpers/proxy"
	Types "github.con/prada-monitors-go/helpers/types"
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
	Query               string
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

func NewMonitor(query string, sku string) *Monitor {
	fmt.Println("Walmart New Monitor Testing Product", query, sku)
	m := Monitor{}
	m.Availability = false
	// var err error
	m.Client = http.Client{Timeout: 10 * time.Second}
	m.Config.site = "WalmartNew"
	m.Config.startDelay = 3000
	m.Query = query
	m.Config.sku = sku
	m.Config.skuName = sku
	// 	m.file, err = os.Create("./testing.txt")
	m.Client = http.Client{
		Timeout: 60 * time.Second,
	}
	m.Config.discord = "https://discord.com/api/webhooks/833902775196450818/wCEYgKpT7eJaOtNERfwe5AlietWcFomGp10zTP3JyDEc8Kk3f2ujqY-BpdXPmpQYANiT"
	m.monitorProduct.name = "Testing Product"
	m.monitorProduct.stockNumber = 10

	proxyList := FetchProxies.Get()
	fmt.Println("Walmart New Monitor Testing Product", query, sku)
	// fmt.Println(timeout)
	//m.Availability = "OUT_OF_STOCK"
	//fmt.Println(m)
	fmt.Println(fmt.Sprintf("https://www.walmart.com/search/api/preso?%s", m.Query))
	i := true
	time.Sleep(5000 * (time.Millisecond))
	go m.checkStop()
	time.Sleep(3000 * (time.Millisecond))
	for i {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Recovering from panic in printAllOperations error is: %v \n", r)
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
		} else {
			fmt.Println(m.Config.sku, "STOPPED STOPPED STOPPED")
			i = false
		}

	}
	return &m
}
func rangeIn(low, hi int) int {
    return low + rand.Intn(hi-low)
}
func (m *Monitor) monitor() error {
	watch := stopwatch.Start()
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovering from panic in printAllOperations error is: %v \n", r)
		}
	}()

	url := fmt.Sprintf("https://www.walmart.com/search/api/preso?%s", m.Query)
	fmt.Println(url)
	req, _ := http.NewRequest("GET", url, nil)
	// req.Header.Add("cookie", `vtc=RXvPjiy9rlmT0_HqtWSkow; TBV=f; DL=33473%2C%2C%2Cip%2C33473%2C%2C; _gcl_au=1.1.241378710.1626465262; tb_sw_supported=true; _ga=GA1.2.866450145.1626480791; GCRT=0894391d-1384-42c8-8c5e-3aa0f6a3c93e; hasGCRT=1; wm_mystore=Fe26.2**b6fdc3ed75cfec5a86f7500843caf9d46920ec44d7f799b931b0935cafb86d09*Q8B4tDXg7igLG6UC9SFeog*8MIkpb5IrhsSY8JRU4ru6fvuYetQPke2ZdPnenObtgL5zyGERYbmHIjQ01szjqV4Wcp-2H20k1H2CRKFg2fZ1mMsCPasLdBK4f87fouVz3NnYSqgks4A1gzGu5bJtExX9NnRRsacu-c7Px1R-jHUQw**f29ff4b5d91b6c2a53b54ef2e169a9605bd1911690c4bfb8b0b8f1b5250c66b3*HmlpBCu52yoYed92RwUYO_CPNLEvD7kFbzpXTujR5UY; s_vi=[CS]v1|307910EA444113A9-60000C762BDE3647[CE]; __gads=ID=fa090174c9a64a18:T=1626481110:S=ALNI_MbUTAHJB0JJ9sUY3dfbrdmg9ah1qg; s_pers=%20s_fid%3D5DD1589A478FCCE6-2D8C526D83D3C278%7C1689553143877%3B%20s_v%3DY%7C1626482943882%3B%20s_cmpstack%3D%255B%255B%2527wmt%2527%252C%25271626481143900%2527%255D%255D%7C1784247543900%3B%20gpv_p11%3DPersistent%2520Cart%7C1626482943902%3B%20gpv_p44%3Dno%2520value%7C1626482943905%3B%20s_vs%3D1%7C1626482943909%3B; s_pers_2=+s_fid%3D5DD1589A478FCCE6-2D8C526D83D3C278%7C1689553143877%3B+s_v%3DY%7C1626482943882%3B+s_cmpstack%3D%255B%255B%2527wmt%2527%252C%25271626481143900%2527%255D%255D%7C1784247543900%3B+gpv_p11%3DPersistent%2520Cart%7C1626482943902%3B+gpv_p44%3Dno%2520value%7C1626482943905%3B+s_vs%3D1%7C1626482943909%3BuseVTC%3DN%7C1689596345; _abck=b8ndxd0chveomw8fjgvx_1998; rtoken=MDgyNTUyMDE4uRiNaOR42Li%2FaQqALWZv76zZdbjFdpv5ZbM4tPC%2B9Q26rzIZioWwa%2BTTC3UrNVsPwN1Tj6LkmULy%2B5RaADQRnytPo%2BMX7cdXA2k236eEKBAJi2JuCsXAPdkY4Xm5jUwd7CKH8eYyQq%2FrrBUDG3GmJpiXMaA8pakndcAM6we6Md0Sa2wDcJ%2FBxXYOLm9%2FME3cnND9OuLLEGNK1EqJcC1faaFohiGsBMDNSb6rb3%2F1hnDT7vSCO135TrrDhhwQ9O4T%2BIrUprBkM7EVGb1FBEzSfn2JWJ4foK9TnSim%2B3ecw2L5s3aK5jNm88C3uC1qxAs1Mu1te4%2FI9oqks53lWamUMr4RZkRA3dyP%2BjmXewsMLZAYd2QezKuINSm%2BwhOWEZn1dlHS9R09XTXrN0pxpvHayg%3D%3D; SPID=e1c98e20fb84e2a4dd2efa88af3689bd4d7da42334c9a532e1fa8d76acf112df93787ead8e1e8e574ccaf39dea3bdd4ccprof; CID=ad1e2798-ac54-4d33-9d9e-3803bca9a057; hasCID=1; customer=%7B%22firstName%22%3A%22Adrian%22%2C%22lastNameInitial%22%3A%22T%22%2C%22rememberme%22%3Atrue%7D; type=REGISTERED; WMP=4; oneapp_customer=true; cart-item-count=2; TS013ed49a=01538efd7ce2244f64106808eff1b2d9dc6334a964dae90d0ac3af729a65b4873bd5697ff428564f31abf2bd4d514d53db03110e04; TB_Latency_Tracker_100=1; TB_Navigation_Preload_01=1; TB_SFOU-100=1; athrvi=RVI~h344dd504-h2aea8d39-h798aed6; s_sess_2=c32_v%3DS2H%2Cnull%3B%20prop32%3DS2H-V%2CS2H; com.wm.reflector="reflectorid:0000000000000000000000@lastupd:1626843310925@firstcreate:1626483269227"; next-day=1626897600|true|false|1626955200|1626843310; location-data=33473%3ABoynton%20Beach%3AFL%3A%3A8%3A1|2bn%3B%3B4.67%2C1uk%3B%3B4.7%2C4fz%3B%3B5.94%2C25h%3B%3B6.14%2C12u%3B%3B6.55%2C185%3B%3B6.7%2C5dj%3B%3B7.45%2C4k7%3B%3B8.09%2C4fy%3B%3B8.1%2C1uu%3B%3B9.56||7|1|1ydz%3B16%3B10%3B10.02%2C1ye0%3B16%3B11%3B10.64%2C1yoh%3B16%3B12%3B11.97%2C1ye2%3B16%3B13%3B13.82%2C1y3g%3B16%3B14%3B25.06; wm_ul_plus=INACTIVE|1626846911000; TB_DC_Flap_Test=0; bstc=WUYguFe1m0kjrZidR4mLXI; mobileweb=0; xpa=; xpm=3%2B1626843310%2BRXvPjiy9rlmT0_HqtWSkow~ad1e2798-ac54-4d33-9d9e-3803bca9a057%2B0; auth=MTAyOTYyMDE4ZqpjZpvXLIEFVV49dSkA1Osdbi056ZywkLZXImEGSIHdy6UDTAuZJN54E6AYzSBTVzO4qGnXmG%2Fddq4ev%2FU8zjAxll6p58K75cKoQBqK4VPSEKCrJp1xdiMCoIH%2BXg3plrvBbkT8GAVfcnvLIPG4V2MZbsr3i13w3j4OuDNfeQDNxtSTOS1hoyKAFxdCiAGja9JZkMvNMuBZk%2FC9ia9BVA88%2FbreJ29hxGEqOcjsqKqJ00UMGQyiYLY97sfSUmPHCgi%2BXQu9OjpsIDyPAwsWQtQHp2cX6P%2F2A5hDcWMIRyMtbf8dslv0ny2JFI%2BXsNGVHWR7H6AUgDne%2FuivR6MM3kTkOjZ36avZqLmfPcLgm2CJDVOEALs9wT40OUOZ1T2ykx043eLnK0EtE8s4eqciIQ%3D%3D; TS01b0be75=01538efd7cde2a6624b7f77e4a7534bb3c55af61da99fadd5d346e26c730d2bd1a3f541d76f0cf78f9c812fba503c73031c8e425c6; akavpau_p8=1626844420~id=bf3719bf3f35085a9b7b1dd661d8e233; _internal.verticalId=default; _internal.verticalTheme=default; _uetsid=47b3ed80e9d311eb8f83f5bfd4d34304; _uetvid=9d3166d0e66f11eb84a2c14d43f5e5ca`)
	req.Header.Add("authority", "www.walmart.com")
	req.Header.Add("pragma", "no-cache")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("sec-ch-ua", `" Not;A Brand";v="99", "Google Chrome";v="91", "Chromium";v="91"`)
	req.Header.Add("accept", "application/json")
	req.Header.Add("dnt", "1")
	req.Header.Add("wm_client_ip", "")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("user-agent", "insomnia/2021.4.1")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("referer", "https://www.walmart.com/browse/trading-cards-by-brand/panini/4171_4191_9807313_6249075_1619679?cat_id=4171_4191_9807313_6249075_1619679&facet=retailer%3AWalmart.com&sort=new")
	req.Header.Add("accept-language", "en-US,en;q=0.9")
	// req.Header.Add("X-PX-ORIGINAL-TOKEN", "3:54d80e223ff20d77a6dbb92f93469c00fcfd44e7c81f546fdabbd012f25acef1:DSzLSuIb2qonIbfjTNitkRI/DMGe6aaJyyYUBgWdIc/ggfMuVpg+oxqSvO5BHO7vDCBtH9/uKNMJ9Z7b0EmvDg==:1000:No7RZwgvEUE4miQH3U06HAhE5kssPXFVwI1/5xJQOG70IvQtBmTTmEmeFnxRtuGDonFtUp6nJ+RlIHI7Uy14infBQvq2TuxOYzn9vAtn4zO9WgA+BnUscGjaznJHpFbgg2KVa6IEnt3TBE1uJGuLC2A7b+giVi9Gs000xWdgcqoyYDrsFaTT7WZM22coge5nz3ORwXxXMWdGTMw1mjcTRA==")
	// req.Header.Add("WM_SITE_MODE", "0")
	// req.Header.Add("mobile-platform", "ios")
	// req.Header.Add("Accept-Language", "en-us")
	// req.Header.Add("X-PX-AUTHORIZATION", "3")
	// req.Header.Add("ACCESS_KEY", "532c28d5412dd75bf975fb951c740a30")
	// req.Header.Add("User-Agent", "insomnia/2021.4.1")
	// req.Header.Add("mobile-app-version", "21.14.0")
	// req.Header.Add("Accept", "*/*")
	// // req.Header.Add("did", "3e012308a1c74e34adc3134caace295600000000")
	req.Header.Add("Connection", "close")
	req.Header.Add("Host", "www.walmart.com")
	req.Close = true

	res, err := m.Client.Do(req)

	if err != nil {
		fmt.Println(err)
		go MonitorLogger.LogError(m.Config.site, m.Config.sku, err)
		return nil
	}
	defer res.Body.Close()
	defer func() {
		watch.Stop()
		fmt.Printf("Walmart New - Status Code : %d Cache : %s : Milliseconds elapsed: %v\n", res.StatusCode, res.Header["Cache-Status"], watch.Milliseconds())
	}()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		go MonitorLogger.LogError(m.Config.site, m.Config.sku, err)
		return nil
	}
	if res.StatusCode != 200 {
		if res.StatusCode == 412 || res.StatusCode == 444 {
			fmt.Println("Blocked by PX 12341")
			fmt.Println("Blocked by PX1234")
			fmt.Println("Blocked by PX 1234")
			fmt.Println("Blocked by PX 1234")
		}
		return nil
	}

	var newList []string
	parser, err := gojq.NewStringQuery(string(body))
	if err != nil {
		fmt.Println(err)
	}
	products, err := parser.Query("items")
	if err != nil {
		fmt.Println(err)
	}
	for key := range products.([]interface{}) {
		var isPresent bool
		pid, err := parser.Query(fmt.Sprintf("items.[%d].usItemId", key))
		if err != nil {
			fmt.Println(err)
			fmt.Println(parser.Query(fmt.Sprintf("items.[%d]", key)))
		}
		sellerName, err := parser.Query(fmt.Sprintf("items.[%d].sellerName", key))
		if err != nil {
			fmt.Println(err)
		}
		title, err := parser.Query(fmt.Sprintf("items.[%d].title", key))
		if err != nil {
			fmt.Println(err)
		}
		image, err := parser.Query(fmt.Sprintf("items.[%d].imageUrl", key))
		if err != nil {
			fmt.Println(err)
		}
		stockNum, err := parser.Query(fmt.Sprintf("items.[%d].quantity", key))
		if err != nil {
			fmt.Println(err)
		}
		offerId, err := parser.Query(fmt.Sprintf("items.[%d].primaryOffer.offerId", key))
		if err != nil {
			fmt.Println(err)
		}
		price, err := parser.Query(fmt.Sprintf("items.[%d].primaryOffer.offerPrice", key))
		if err != nil {
			fmt.Println(err)
			price = float64(00000)
		}
		newList = append(newList, pid.(string))
		for _, v := range m.products {
			if v == pid {
				isPresent = true
			}
		}
		if !isPresent {
			m.products = append(m.products, pid.(string))
			// fmt.Println(pid.(string), offerId.(string), price.(float64), title.(string), image.(string), stockNum.(float64), sellerName.(string))
			go m.sendWebhook(pid.(string), offerId.(string), price.(float64), title.(string), image.(string), stockNum.(float64), sellerName.(string))

		}
	}
	fmt.Println(len(m.products))

	m.products = newList
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
func (m *Monitor) sendWebhook(sku string, offerId string, price float64, productName string, image string, stockNum float64, sellerName string) error {
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
	t := time.Now().UTC().Format("2006-01-02T15:04:05Z")

	for _, comp := range m.CurrentCompanies {
		fmt.Println(comp.Company)
		go m.webHookSend(comp, m.Config.site, sku, productName, int(price), offerId, t, image, int(stockNum), sellerName)
	}
	return nil
}
func (m *Monitor) webHookSend(c Company, site string, sku string, name string, price int, offerId string, time string, image string, stockNum int, sellerName string) {
	payload := strings.NewReader(fmt.Sprintf(`{
		"content": null,
		"embeds": [
		  {
			"title": "%s Monitor",
			"url": "https://www.walmart.com/ip/prada/%s",
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
				"name": "OfferId",
				"value": "%s"
			  },
			  {
				"name": "Stock Quantity",
				"value": "%d",
				"inline": true
			  },
			  {
				"name": "Seller Name",
				"value": "%s",
				"inline": true
			  },
			  {
				"name": "Links",
				"value": "[Product](https://www.walmart.com/ip/prada/%s) | [ATC](https://affil.walmart.com/cart/buynow?items=%s) | [Checkout](https://www.walmart.com/checkout/) | [Cart](https://www.walmart.com/cart)"
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
	  }`, site, sku, c.Color, name, price, sku, offerId, stockNum, sellerName, sku, sku, time, image, c.CompanyImage))
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
				fmt.Printf("Site : %s, Product : %s  CHECK STOP Recovering from panic in printAllOperations error is: %v \n", m.Config.site, m.Config.sku, r)
			}
		}()

		fmt.Println("TESTING PAYLOADS", strings.ToUpper(m.Config.site), m.Config.sku)
		getDBPayload := strings.NewReader(fmt.Sprintf(`{
			"site" : "%s",
			"sku" : "%s"
		  }`, strings.ToUpper(m.Config.site), m.Config.sku))
		url := "http://104.249.128.207:7243/DB"
		req, err := http.NewRequest("POST", url, getDBPayload)
		if err != nil {
			fmt.Println(err)
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println(err)

		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
		}
		var currentObject ItemInMonitorJson
		err = json.Unmarshal(body, &currentObject)
		if err != nil {
			fmt.Println(err)
		}
		m.stop = currentObject.Stop
		m.CurrentCompanies = currentObject.Companies
		fmt.Println(m.CurrentCompanies)
		time.Sleep(5000 * (time.Millisecond))
	}
	return nil
}

func WalmartNew(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "Walmart New Monitor")
	fmt.Println("Walmart New Monitor")
	var currentMonitor Types.MonitorResponse
	_ = json.NewDecoder(r.Body).Decode(&currentMonitor)
	// fmt.Println(currentMonitor)
	go NewMonitor(currentMonitor.SkuName, currentMonitor.Sku)
	json.NewEncoder(w).Encode(currentMonitor)
}