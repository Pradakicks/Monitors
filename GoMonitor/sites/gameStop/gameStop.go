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
	go m.checkStop()
	time.Sleep(3000 * (time.Millisecond))
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
			fmt.Println("Gamestop : ", m.Availability, m.Config.sku, m.statusCode)
			time.Sleep(250 * (time.Millisecond))
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

	url := fmt.Sprintf("https://www.gamestop.com/on/demandware.store/Sites-gamestop-us-Site/default/Product-Variation?dwvar_%s_condition=New&pid=%s&quantity=1&redesignFlag=true&rt=productDetailsRedesign", m.Config.sku, m.Config.sku)
	req, _ := http.NewRequest("GET", url, nil)
	gmeRef := fmt.Sprintf("https://www.gamestop.com/products/prada/%s.html", m.Config.sku)
	//	req.Header.Add("cookie", "dwac_a78eb5a11975e4c9cbc84f55ad=hsID1-Tt0AwdpVxt3843SGPN-twDq54LmqE%253D%7Cdw-only%7C%7C%7CUSD%7Cfalse%7CUS%252FCentral%7Ctrue; cqcid=ab5UPFZX9yyAfmyTsKUgaaHFzv; cquid=%7C%7C; customergroups=Everyone; sid=hsID1-Tt0AwdpVxt3843SGPN-twDq54LmqE; dwanonymous_420142ceefb9f0c103b3815e84e9fcef=ab5UPFZX9yyAfmyTsKUgaaHFzv; userInfo=S%3DN%7CR%3DN%7CC%3D0; __cq_dnt=0; dw_dnt=0; dwsid=Wlq6vy6WB_fCXkDfe9xQG79VQ7TJ0EnwRVpVHu76k74TZNmN_l1gaPvcKEW4cHBDjoAZPqq1X-8Fep-T3_Al6A%3D%3D; akaas_ChatThrottling=2147483647~rv%3D96~id%3D63cbc0fee3058ac39b466275a70905af~rn%3D; _abck=87E16D9492B4FF25B48E128780C8B242~-1~YAAQFjPKF3WLVop5AQAA%2FdqAigUCb4OWPBCZD53lDEvJlGuVt0EWx8aijJXv78XmQNE3Jcn0HhA7uuDlcW46NkbVrSezW78SqsyDvh%2BcuXwNI1HodryzQYa8myrq7UClFXdgMQl0LGO3hMu2u8K1rKbX4i2th5MajKHLGz12mYraytOUtFQfgOc79XF6DkO2cdHzPUY7S0KlwrdPyQLJ1gnYTX6%2Bwv8VTblTlANXfdMOCuTJwglqgpDHjfT0CAzItmKKyzFA7q45VESCP4m%2F2TYBY4dP0gdLf62je%2F%2FND27XoxnDLyHTI%2Fpoz4Tu62E8VEBhj1OUASBncq7bVeqYdF3hK%2Btd6ANpSRnEEJpAZ%2FEQWFU6oSOny%2BYMjazZR8NKN13Y6GSeFSPJ2K52~-1~-1~-1")
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
	//	req.Header.Add("$cookie", "akaas_ChatThrottling=2147483647~rv=14~id=f9ecbc4735b5a359f6c7d9ce04a4478e~rn=; _caid=a9fc9999-b8eb-497c-84eb-f5e30de59b5c; RES_TRACKINGID=831656757534380; ResonanceSegment=1; notice_behavior=implied,eu; _gcl_au=1.1.403930278.1621428291; cquid=||; customergroups=Everyone; userInfo=S=N|R=N|C=0; __cq_dnt=0; dw_dnt=0; cqcid=abbwdRaFf3tfLzthJuEsJWoOX0; dwanonymous_420142ceefb9f0c103b3815e84e9fcef=abbwdRaFf3tfLzthJuEsJWoOX0; _mibhv=anon-1611193773687-9267464732_6874; _gid=GA1.2.139986699.1621428292; _ga=GA1.2.470299063.1621428292; _fbp=fb.1.1621428291733.400088248; __cq_uuid=ac2nyrc9myCuCaJd7go98GvLl7; sto__vuid=9b0a62ad062fe7634881c08be1357815; QuantumMetricUserID=3255d30ede3179288cea36f3fcb75bee; _aeaid=0b21af47-dd3a-427e-a8dc-d4ff1db632ed; aeatstartmessage=true; lastVisitedCgid=pc; BVBRANDID=21f8293a-6961-4947-8643-e3c4da07b51c; salsify_session_id=3f8d9c2e-46b0-4bdd-8873-d026aef0f5c7; __cq_seg=0~0.16\u00211~-0.30\u00212~-0.23\u00213~-0.78\u00214~0.16\u00215~0.03\u00216~-0.42\u00217~0.10\u00218~-0.11\u00219~0.01; AnonymousSaveForLaterCookie=a2b82aafce4c3ecb32c055217c; spCheckout=true; BVImplsfcc=9014_3_0; bm_sz=6E4624F205C8D9D157959B7DC72F043C~YAAQBjPKFwdSvDR5AQAA9/1/igv0dp9v2oGrh6RLT8nLJ78toTLM9AKEnOEnDk+IKGsIUTfOpoNOy40aYKU6awzO49Duc3SYBHm85I6w68tP7KFbZN0NfTStPy8tcLMmFVtb7P8WoscTtYZSPlFoMb5ER+cGpta3vFwp2d/KvQHl20Bqxf4i1hDpKpVeTBe3rg==; _gaexp=GAX1.2.IodINR6fQb--45IDXrQ7eQ.18851.2\u0021kQ74Gh4ZQ3ic2rdDBej41Q.18860.0\u0021UZIDaZ6yRSCK79It9Iq0ew.18860.0; _cavisit=1798a800045|; BVBRANDSID=562c3483-d621-4e1f-a87b-2be56fd4ee24; RES_SESSIONID=950727001230365; _gat_UA-10897913-30=1; dwac_a78eb5a11975e4c9cbc84f55ad=xtZB7DghmV3O0XUSJMQIyAQnNAviCQQxDKo%3D|dw-only|||USD|false|US%2FCentral|true; sid=xtZB7DghmV3O0XUSJMQIyAQnNAviCQQxDKo; dwsid=cPacK9ND5wH-rPX9-0yTF7IH6cwJa3ptPNdh2US1PhJDd_0DB0MPdHxCeLFbx5mswTKEhw64haEzQfkJ0qFu5w==; bm_mi=C40316A27AF356148E35871AAFE5AC3C~RVY/QOu/xX5cZD1z7xMu9AqHST30b/wkMXYdSXKDSoTGAiZXlU3CxYFbLqkYQROZaFgGXkavMtfE2LSTGMdAoiYQd9ww4nSDZc/e2xZntTKaYuruK1bsn/PTyF7jtPryhNwyJlgJJABg583jhKxwcdaJ2uRxcUVXUBFyOXrFGDvFYmJhPXVHpR1V1inoQfdj3NFdcOSk/Ip71dJEVpirpR3mJi34Pw0cwMJ0DI75yGvxp2V9XB2bvv2TdCrg38lj0A6j3rrQsBG++DdFb+WXhSZCBnuFIOXtEiKuS2Z77ikHN+1elTdjAVPYUiGfVo9a; ak_bmsc=7C4EDD042E8147FA924491DE2161CF0417CA3306844B00002A87A6606EF3E77A~plEzQxzJ892lN2n9XR27jfZkGrTNUrOlaBBC+n7OxU3Ypbokp31vwzRV0AT2+U5l/BiEQJtSGeNqMgqGXPdD2LPWIQcq3cE2cJU+x2QP65irmHevt16mFk3X9AM/unop8cufXssMTVMkNHureuWkT8RTcHdIwP9ngax267K/NkF/tChl/3LDevJ7spk0JDwStMgXK7bMgqNjHa7MOC/g9tbfe69ghPEwADvtd0cPzKnJrLqAldhHQK1y8esRoZlM3n; sto__session=1621526317850; QuantumMetricSessionID=f62e7111fa6c11abd96e7c92429b0358; _uetsid=00fa3ab0b8a011eba0a725a2ac83ba10; _uetvid=00fa3f90b8a011eb804207a6f4a89e19; __cq_bc=%7B%22bcpk-gamestop-us%22%3A%5B%7B%22id%22%3A%2211146493%22%2C%22sku%22%3A%22295368%22%7D%2C%7B%22id%22%3A%2211107421%22%7D%2C%7B%22id%22%3A%2211095775%22%7D%5D%7D; stc118903=env:1621526315%7C20210620155835%7C20210520162845%7C2%7C1084111:20220520155845|uid:1621428291596.933694301.3955164.118903.1117792204.:20220520155845|srchist:1084110%3A1%3A20210619124451%7C1084111%3A1621526315%3A20210620155835:20220520155845|tsa:1621526315835.1379394640.781375.174304299060569.42:20210520162845; sto__count=1; _abck=97165556E8D94CA58EB68395D61C1EE0~-1~YAAQBjPKFwVTvDR5AQAAoUuAigWGeft0uIYVC3ciasUfslYO2AS9c9po2kBv2W+Wj8B36dfFq+nEc5UO6EwyjQSVV++LcRREW4hIsy+iv/f4fTQ+Lpdxlb1bvi9b+lGXNJqcJvGZ1L/17XkU9PBfEhYD8Gp6SBUSE/JjiCbQijN8FIPnW8eC59aDWPE+y9QNcMH1QxgORCYd7NEOWGXz7HlWCO1pChv+ySD6sINUrZ7DLztY5NQBAJ430yNJ3xW4hLvMICs8XWoT50Zak7cSCGA3u7er9Vl9uDk54vIRmvQqKwvn5KipZICgRKAm8Es6/MA7PAYgCxNf/YFgAJV4pogrko8JmyOlfhaC74+VpBKKFhQaiXfjfMO5ToAYF5UdI1bDuw06goWXhGD1Cr0+bRImVOH44vcNsmkYx4fxKgJ1/hHmXPdh~-1~-1~-1")

	res, _ := m.Client.Do(req)
	m.statusCode = res.StatusCode
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

	//	fmt.Println(m.monitorProduct.link, m.monitorProduct.image)
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
