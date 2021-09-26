package Amazon

import (
	// "bytes"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/bradhe/stopwatch"
	"github.com/nickname32/discordhook"
	"github.com/pkg/errors"

	Webhook "github.con/prada-monitors-go/helpers/discordWebhook"
	FetchProxies "github.con/prada-monitors-go/helpers/proxy"
	Types "github.con/prada-monitors-go/helpers/types"
)

type CurrentMonitor struct {
	Monitor Types.Monitor
	sid     string
	csrf    string
	oid     string
}

type AmazonResponse struct {
	CartSubtotalString       string        `json:"cartSubtotalString"`
	CoreFreeShippingMessage  string        `json:"coreFreeShippingMessage"`
	FormattedTotalPrice      string        `json:"formattedTotalPrice"`
	IncludedAsins            []interface{} `json:"includedAsins"`
	ItemQuantity             string        `json:"itemQuantity"`
	ItemQuantityString       string        `json:"itemQuantityString"`
	TotalPrice               string        `json:"totalPrice"`
	TotalPriceInBaseCurrency string        `json:"totalPriceInBaseCurrency"`
}

type CookieMonitorResponse struct {
	Sid string `json:"sid"`
	// cookies string `json:"cookies,omitempty"`
	Csrf string `json:"csrf"`
}

// var client http.Client

var sid string
var csrf string

func NewMonitor(sku string, oid string) *CurrentMonitor {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Site : %s, Product : %s Recovering from panic in printAllOperations error is: %v \n", sku, sku, r)
		}
	}()
	fmt.Println("AMAZON")
	m := CurrentMonitor{}
	m.Monitor.AvailabilityBool = false
	m.Monitor.Config.Site = "Amazon"
	m.Monitor.Config.Sku = sku
	m.oid = oid
	m.Monitor.Client = http.Client{Timeout: 10 * time.Second}
	proxyList := FetchProxies.Get()

	// time.Sleep(15000 * (time.Millisecond))
	// go m.Monitor.CheckStop()
	// time.Sleep(3000 * (time.Millisecond))

	i := true
	// jar, _ := cookiejar.New(nil)
	// m.Monitor.Client = http.Client{
	// 	Jar: jar,
	// }
	// m.getCookies()
	// m.getCookies()
	m.getHome()
	for i {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Site : %s, Product : %s Recovering from panic in printAllOperations error is: %v \n", m.Monitor.Config.Site, m.Monitor.Config.Sku, r)
			}
		}()
		if !m.Monitor.Stop {
			currentProxy := m.Monitor.GetProxy(proxyList)
			splittedProxy := strings.Split(currentProxy, ":")
			proxy := Types.Proxy{splittedProxy[0], splittedProxy[1], splittedProxy[2], splittedProxy[3]}
			//	fmt.Println(proxy, proxy.ip)
			prox1y := fmt.Sprintf("http://%s:%s@%s:%s", proxy.UserAuth, proxy.UserPass, proxy.Ip, proxy.Port)
			proxyUrl, err := url.Parse(prox1y)
			if err != nil {
				fmt.Println(errors.Cause(err))
				return nil
			}
			defaultTransport := &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			}
			m.Monitor.Client = http.Client{
				// Jar:       jar,
				Transport: defaultTransport,
			}

			go m.monitor()
			time.Sleep(100 * (time.Millisecond))
		} else {
			fmt.Println(m.Monitor.Config.Sku, "STOPPED STOPPED STOPPED")
			i = false
		}

	}
	return &m
}

func (m *CurrentMonitor) monitor() error {
	watch := stopwatch.Start()
	// m.sid = sid
	// m.csrf = csrf
	// oid := "f9FL%2FHvb8mR%2Fj3J361zTxy2J87piw0thjhoprF8qOayYYFTZsg1cIfDAEcG5D8tXKxNQiQsTBT%2F53D4%2FhW1Kgt1ruMYpP6XrgF5pie3MZfHXZ0xLMVgGaTUrNS9DOrabFGjD3TIzGMTMlDoz%2FKmCzA%3D%3D"
	// fmt.Println("HERE", fmt.Sprintf("session-id=%s;", m.sid), fmt.Sprintf(`"{\"items\":[{\"asin\":\"%s\",\"offerListingId\":\"%s\",\"quantity\":1}]}"`, m.Monitor.Config.Sku, m.oid))
	// payload := strings.NewReader(fmt.Sprintf(`marketplaceId=ATVPDKIKX0DER&asin=%v&customerId=&sessionId=%v&accessoryItemAsin=B002M40VJM&accessoryItemOfferingId=%v&languageOfPreference=en_US&accessoryItemQuantity=1&accessoryItemPrice=9.99&accessoryMerchantId=ATVPDKIKX0DER&accessoryProductGroupId=8652000`, m.Monitor.Config.Sku, m.sid, oid))
	// fmt.Println(fmt.Sprintf(`marketplaceId=ATVPDKIKX0DER&asin=%v&customerId=&sessionId=%v&accessoryItemAsin=B002M40VJM&accessoryItemOfferingId=%v&languageOfPreference=en_US&accessoryItemQuantity=1&accessoryItemPrice=9.99&accessoryMerchantId=ATVPDKIKX0DER&accessoryProductGroupId=8652000`, m.Monitor.Config.Sku, m.sid, oid))
	url := "https://data.amazon.com/api/marketplaces/ATVPDKIKX0DER/cart/carts/retail/items?ref=aod_dpdsk_used_1"

	payload := strings.NewReader(fmt.Sprintf("{\"items\":[{\"asin\":\"%s\",\"offerListingId\":\"%s\",\"quantity\":1}]}", m.Monitor.Config.Sku, m.oid))
	// fmt.Println(m.sid, m.csrf, m.Monitor.Config.Sku, m.oid)
	req, _ := http.NewRequest("POST", url, payload)

	// req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Pragma", "no-cache")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("sec-ch-ua", `" Not;A Brand";v="99", "Google Chrome";v="91", "Chromium";v="91"`)
	req.Header.Add("DNT", "1")
	// req.Header.Add("x-api-csrf-token", "1@g+J9lHIdBuX8o3AjxeE7866vU4nadfUdUycwHS/gmkL2AAAADAAAAABhQWnMcmF3AAAAABVX8CwXqz4nuL9RKX///w==@NLD_Y47GLU")
	req.Header.Add("x-api-csrf-token", m.csrf)
	req.Header.Add("Accept-Language", "en-US")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Add("Content-Type", `application/vnd.com.amazon.api+json; type="cart.add-items.request/v1"`)
	req.Header.Add("Accept", `application/vnd.com.amazon.api+json; type="cart.add-items/v1"`)
	req.Header.Add("Origin", "https://www.amazon.com")
	req.Header.Add("Sec-Fetch-Site", "same-site")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Referer", "https://www.amazon.com/")
	// req.Header.Add("cookie", "session-id=147-2666751-6954768;")//X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36")
	req.Header.Add("cookie", fmt.Sprintf("session-id=%s;", m.sid))
	req.Header.Set("Connection", "close")
	req.Close = true
	res, err := m.Monitor.Client.Do(req)
	if err != nil {
		fmt.Println(err)
		fmt.Println(errors.Cause(err))
		return nil
	}
	defer res.Body.Close()
	defer func() {
		watch.Stop()
		fmt.Printf("Amazon %s - InStock : %t - Code : %d Milli elapsed: %v\n", m.Monitor.Config.Sku, m.Monitor.AvailabilityBool, res.StatusCode, watch.Milliseconds())
	}()
	// body, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// 	fmt.Println(errors.Cause(err))
	// 	return nil
	// }

	// var realBody AmazonResponse
	// json.Unmarshal(body, &realBody)
	// fmt.Println(realBody)

	// fmt.Println(res.StatusCode)
	// fmt.Println(string(body))

	if res.StatusCode == 200 {
		m.Monitor.CurrentAvailabilityBool = true
	} else if res.StatusCode == 422 {
		m.Monitor.CurrentAvailabilityBool = false
	} else if res.StatusCode == 404 {
		m.getHome()
	}

	if !m.Monitor.AvailabilityBool && m.Monitor.CurrentAvailabilityBool {
		fmt.Println("Item In Stock", m.Monitor.Config.Sku)
		currentC := Types.Company{}
		currentC.Company = "Vibris"
		currentC.Webhook = "https://discord.com/api/webhooks/797249480410923018/NPL3ktXS78z5EHo_cpYyrtFl_2iB0ARgz9IW5kwAZA-UkiseiinnBmUPJZlGgxw8TZiW"
		currentC.Color = "15277667"
		currentC.CompanyImage = "https://cdn.discordapp.com/attachments/843905652790263838/871833603770810489/3.png"
		t := time.Now().UTC()
		go m.sendRestockNotification(m.oid, m.Monitor.Config.Sku, "Amazon Product")
		go m.webHookSend(currentC, "Amazon", "Test Product", 999, fmt.Sprintf("https://www.amazon.com/gp/product/%s", m.Monitor.Config.Sku), t, "https://cdn.discordapp.com/attachments/843905652790263838/871833603770810489/3.png", 1)
		time.Sleep(60000 * time.Millisecond)

	}

	m.Monitor.AvailabilityBool = m.Monitor.CurrentAvailabilityBool

	return nil
}

func (m *CurrentMonitor) sendWebhook(sku string, name string, price float64, link string, image string, Qty int) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Site : %s, Product : %s Recovering from panic in printAllOperations error is: %v \n", m.Monitor.Config.Site, m.Monitor.Config.Sku, r)
		}
	}()
	for _, letter := range m.Monitor.MonitorProduct.Name {
		if string(letter) == `"` {
			m.Monitor.MonitorProduct.Name = strings.Replace(m.Monitor.MonitorProduct.Name, `"`, "", -1)
		}
	}
	fmt.Println("Testing Here : ", m.Monitor.MonitorProduct.Name, "Here")
	if strings.HasSuffix(m.Monitor.MonitorProduct.Name, "                       ") {
		m.Monitor.MonitorProduct.Name = strings.Replace(m.Monitor.MonitorProduct.Name, "                       ", "", -1)
	}
	fmt.Println("Testing Here : ", m.Monitor.MonitorProduct.Name, "Here")
	// now := time.Now()
	// currentTime := strings.Split(now.String(), "-0400")[0]
	t := time.Now().UTC()

	for _, comp := range m.Monitor.CurrentCompanies {
		fmt.Println(comp.Company)
		go m.webHookSend(comp, m.Monitor.Config.Site, name, price, link, t, image, Qty)
	}
	// payload := strings.NewReader("{\"content\":null,\"embeds\":[{\"title\":\"Target Monitor\",\"url\":\"https://discord.com/developers/docs/resources/channel#create-message\",\"color\":507758,\"fields\":[{\"name\":\"Product Name\",\"value\":\"%s\"},{\"name\":\"Product Availability\",\"value\":\"In Stock\\u0021\",\"inline\":true},{\"name\":\"Stock Number\",\"value\":\"%s\",\"inline\":true},{\"name\":\"Links\",\"value\":\"[Product](https://www.walmart.com/ip/prada/%s)\"}],\"footer\":{\"text\":\"Prada#4873\"},\"timestamp\":\"2021-04-01T18:40:00.000Z\",\"thumbnail\":{\"url\":\"https://cdn.discordapp.com/attachments/815507198394105867/816741454922776576/pfp.png\"}}],\"avatar_url\":\"https://cdn.discordapp.com/attachments/815507198394105867/816741454922776576/pfp.png\"}")
	return nil
}

func (m *CurrentMonitor) webHookSend(c Types.Company, site string, name string, price float64, link string, currentTime time.Time, image string, Qty int) {
	Color, err := strconv.Atoi(c.Color)
	if err != nil {
		fmt.Println(err)
		fmt.Println(errors.Cause(err))
		return
	}
	Quantity := strconv.Itoa(Qty)
	Price := fmt.Sprintf("%f", price)
	var currentFields []*discordhook.EmbedField
	currentFields = append(currentFields, &discordhook.EmbedField{
		Name:   "Product Name",
		Value:  name,
		Inline: false,
	})
	currentFields = append(currentFields, &discordhook.EmbedField{
		Name:   "Price",
		Value:  Price,
		Inline: true,
	})
	currentFields = append(currentFields, &discordhook.EmbedField{
		Name:   "Stock #",
		Value:  Quantity,
		Inline: true,
	})
	currentFields = append(currentFields, &discordhook.EmbedField{
		Name:  "Links",
		Value: "[In Development](" + link + ")",
	})
	var discordParams discordhook.WebhookExecuteParams = discordhook.WebhookExecuteParams{
		Content: "",
		Embeds: []*discordhook.Embed{
			{
				Title:  site + " Monitor",
				URL:    link,
				Color:  Color,
				Fields: currentFields,
				Footer: &discordhook.EmbedFooter{
					Text: "Prada#4873",
				},
				Timestamp: &currentTime,
				Thumbnail: &discordhook.EmbedThumbnail{
					URL: image,
				},
				Provider: &discordhook.EmbedProvider{
					URL: c.CompanyImage,
				},
			},
		},
	}
	go Webhook.SendWebhook(c.Webhook, &discordParams)
}
func Amazon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "Amazon Monitor")
	fmt.Println("Amazon")
	var currentMonitor Types.MonitorResponse
	_ = json.NewDecoder(r.Body).Decode(&currentMonitor)
	fmt.Println(currentMonitor)
	go NewMonitor(currentMonitor.Sku, currentMonitor.SkuName)
	// go NewMonitor(currentMonitor.Sku, currentMonitor.SkuName)
	json.NewEncoder(w).Encode(currentMonitor)
}

func AmazonSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "Amazon Cookies")
	fmt.Println("Amazon Cookies")
	fmt.Println("Amazon Cookies")
	fmt.Println("Amazon Cookies")
	fmt.Println("Amazon Cookies")
	var currentBody CookieMonitorResponse
	_ = json.NewDecoder(r.Body).Decode(&currentBody)
	fmt.Println(currentBody)
	// sid = currentBody.Sid
	// csrf = currentBody.Csrf
	// fmt.Println(sid, csrf)
	// fmt.Println(r.Response.Body)
	json.NewEncoder(w).Encode(currentBody)
}

func (m *CurrentMonitor) getCookies() error {
	watch := stopwatch.Start()
	fmt.Println("COOKIES GET")
	req, err := http.NewRequest("GET", "https://smile.amazon.com/gp/mobile/udp/ajax-handlers/reftag.html?ref_=dp_atch_abb_i", nil)
	if err != nil {
		fmt.Println(err)
		fmt.Println(errors.Cause(err))
		return nil
	}
	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Add("accept-language", "en-US,en;q=0.9")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("pragma", "no-cache")
	req.Header.Add("sec-ch-ua", "\" Not;A Brand\";v=\"99\", \"Google Chrome\";v=\"91\", \"Chromium\";v=\"91\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-fetch-dest", "document")
	req.Header.Add("sec-fetch-mode", "navigate")
	req.Header.Add("sec-fetch-site", "none")
	req.Header.Add("sec-fetch-user", "?1")
	req.Header.Add("upgrade-insecure-requests", "1")
	// req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36")
	res, err := m.Monitor.Client.Do(req)
	if err != nil {
		fmt.Println(err)
		fmt.Println(errors.Cause(err))
		return nil
	}
	defer res.Body.Close()
	defer func() {
		watch.Stop()
		// fmt.Printf("Home Depot %s - Code : %d Milli elapsed: %v\n", m.Monitor.Config.Sku, res.StatusCode, watch.Milliseconds())
	}()
	// body, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// 	fmt.Println(errors.Cause(err))
	// 	return nil
	// }
	// if res.StatusCode != 200 {
	// 	time.Sleep(10 * time.Second)
	// 	return nil

	// }
	fmt.Println(res.StatusCode)
	for i, v := range res.Header {
		if i == "Set-Cookie" {
			fmt.Println(v)
			if strings.Contains(v[0], "session-id=") {
				m.sid = strings.Split(strings.Split(v[0], "session-id=")[1], ";")[0]
				fmt.Println(strings.Split(strings.Split(v[0], "session-id=")[1], ";")[0])
			}

		}
	}

	return nil
}
func (m *CurrentMonitor) getHome() error {
	watch := stopwatch.Start()
	fmt.Println("COOKIES GET HOME")
	url := fmt.Sprintf("https://www.amazon.com/gp/product/%s", m.Monitor.Config.Sku)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		fmt.Println(errors.Cause(err))
		return nil
	}
	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Add("accept-language", "en-US,en;q=0.9")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("pragma", "no-cache")
	req.Header.Add("sec-ch-ua", "\" Not;A Brand\";v=\"99\", \"Google Chrome\";v=\"91\", \"Chromium\";v=\"91\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-fetch-dest", "document")
	req.Header.Add("sec-fetch-mode", "navigate")
	req.Header.Add("sec-fetch-site", "none")
	req.Header.Add("sec-fetch-user", "?1")
	req.Header.Add("upgrade-insecure-requests", "1")
	// req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36")
	res, err := m.Monitor.Client.Do(req)
	if err != nil {
		fmt.Println(err)
		fmt.Println(errors.Cause(err))
		return nil
	}
	defer res.Body.Close()
	defer func() {
		watch.Stop()
		// fmt.Printf("Home Depot %s - Code : %d Milli elapsed: %v\n", m.Monitor.Config.Sku, res.StatusCode, watch.Milliseconds())
	}()
	// // body, err := ioutil.ReadAll(res.Body)
	// // if err != nil {
	// // 	fmt.Println(err)
	// // 	fmt.Println(errors.Cause(err))
	// // 	return nil
	// // }
	// // if res.StatusCode != 200 {
	// // 	time.Sleep(10 * time.Second)
	// // 	return nil

	// // }
	// fmt.Println(res.StatusCode)
	// // for i, v := range res.Header {
	// // 	if i == "Set-Cookie" {
	// // 		fmt.Println(v)
	// // 		if strings.Contains(v[0], "session-id=") {
	// // 			m.sid = strings.Split(strings.Split(v[0], "session-id=")[1], ";")[0]
	// // 			fmt.Println(strings.Split(strings.Split(v[0], "session-id=")[1], ";")[0])
	// // 		}

	// // 	}
	// // }
	fmt.Println(res.StatusCode)
	if res.StatusCode != 200 {
		time.Sleep(time.Millisecond * time.Duration(500))
		m.getHome()
	}
	var page *goquery.Document

	if res != nil {
		page, err = goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			fmt.Println(err)
			return errors.New("Error Parsing Doc")
		}
		res.Body.Close()
	}

	var ok bool
	if sid, ok = page.Find("input[id='session-id']").Attr("value"); !ok {
		fmt.Println("Something missing 1")
		m.getCSRF()
		return nil
		// return errors.New("Error Parsing Doc")
	} 
	// m.si
	fmt.Println(sid, "TESTING")
	fmt.Println(sid, "TESTING")
	m.sid = sid
	m.getCSRF()
	return nil
}
func (m *CurrentMonitor) getCSRF() error {
	watch := stopwatch.Start()
	url := fmt.Sprintf("https://www.amazon.com/gp/aod/ajax/ref=auto_load_aod?asin=%s&pc=dp", m.Monitor.Config.Sku)
	fmt.Println("CSRF Fetch", fmt.Sprintf("session-id=%s;", m.sid), url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		fmt.Println(errors.Cause(err))
		return nil
	}
	// req.Header.Add("cookie", fmt.Sprintf("session-id=%s;", m.sid))
	req.Header.Add("authority", "www.amazon.com")
	req.Header.Add("pragma", "no-cache")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("rtt", "50")
	req.Header.Add("downlink", "10")
	req.Header.Add("ect", "4g")
	req.Header.Add("sec-ch-ua", `"Chromium";v="92", " Not A;Brand";v="99", "Google Chrome";v="92"`)
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("dnt", "1")
	req.Header.Add("upgrade-insecure-requests", "1")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36")
	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Add("sec-fetch-site", "none")
	req.Header.Add("sec-fetch-mode", "navigate")
	req.Header.Add("sec-fetch-user", "?1")
	req.Header.Add("sec-fetch-dest", "document")
	req.Header.Add("accept-language", "en-US,en;q=0.9")
	req.Header.Add("cookie", fmt.Sprintf("session-id=%s;", sid))
	// req.Header.Add("cookie", "ubid-main=132-7199499-9402615; lc-main=en_US; s_vnum=2055786480819%26vn%3D1; s_fid=085844D920902910-0DB536AAEB684C2B; aws-target-data=%7B%22support%22%3A%221%22%7D; aws-target-visitor-id=1623792154921-241206.34_0; aws-ubid-main=172-3048843-7260845; remember-account=false; regStatus=registered; awsc-color-theme=light; s_vn=1655867630274%26vn%3D1; s_dslv=1624331690698; s_nr=1624331690710-Repeat; x-main=ZIE86NWu1jB5HxfuG0qVxaU08jqV6sGBp2DA9tsu928OHZhsC2qT9C43prcuqzeA; at-main=Atza|IwEBIKXGflTuoexerWYqIsnUV3p9wik3OPXv7v2EJ3qT5QZcN1hIG0FguU0zlYCut1lKB29gF7TV70OO9ZGqDwC_sKWGnpSdyfcy8xrCYEC1WLUHAGY63bQFmRfeqMrtKIz05d9Ik6rmDVoHURLgrsefT9dx6Muxe1vdUXrANNcx6P5PiM0N5Gd4WnE4P_asSodtZakEim9GXwYYRRNjQO8GVox2; sess-at-main=\"Xdkcueks4+wSSK6zT0h1mv7AkUcMz53Pk31NyOT7CfY=\"; sst-main=Sst1|PQEN0NQUJ894Zthy9VGLmzrRCQaDdSIA3GIXq5hTe4BPG3e0w-PyJrrCn_fZ5Wc8qK4u5c57XsmADeArS5YQBweKHWSB8T0wSfCPlrL4KhpM1LWuQsyUTASxbFkeIsMLy1PZfL3MsDjRHsFhrJfkjWkqdb2t9WbYMgGGY-p-WjhzNC8eSXpqujMbzXkA7oBsf189_3sy3u3o645DZ2uhRSKh-UNygWmTtd5YJLcvi9O95BT-T4s2JvMN8QIvC9KSYsBKwCQd4q0ENdNI9nucKNTyQjOZ0kzWSIwd6xCTqFaqNqM; i18n-prefs=USD; session-id-apay=145-1370814-2780545; session-id=131-7297596-5151307; session-id-time=2082787201l; aws-userInfo-signed=eyJ0eXAiOiJKV1MiLCJrZXlSZWdpb24iOiJ1cy1lYXN0LTEiLCJhbGciOiJFUzM4NCIsImtpZCI6ImFiMWE2OTgwLTQzMDAtNGQ3Yy1iMzRlLTYzZWFiODJhYTA4NyJ9.eyJzdWIiOiIiLCJzaWduaW5UeXBlIjoiUFVCTElDIiwiaXNzIjoiaHR0cDpcL1wvc2lnbmluLmF3cy5hbWF6b24uY29tXC9zaWduaW4iLCJrZXliYXNlIjoiQ2kwZDI5eU9cL2pPbXo2aTRMdU9xNE1DcHNldUxqaFlOWjhVMDFGUXNPZXM9IiwiYXJuIjoiYXJuOmF3czppYW06OjEwNTMwMzA1MzQxNzpyb290IiwidXNlcm5hbWUiOiJQcmFkYUtpY2tzIn0.Oxi7dsCC28WdyqXBxEl1VICR_NEHAmWRvOjDr6cKWFgpa-SsEKS9bu9F4ZVwhBzbZNs4Q7uuzjER6nbbo0yPlIQpMXiDTemlCICY89VIIGCOdKC_YvHNInN7OMlj9MUj; aws-userInfo=%7B%22arn%22%3A%22arn%3Aaws%3Aiam%3A%3A105303053417%3Aroot%22%2C%22alias%22%3A%22%22%2C%22username%22%3A%22PradaKicks%22%2C%22keybase%22%3A%22Ci0d29yO%2FjOmz6i4LuOq4MCpseuLjhYNZ8U01FQsOes%5Cu003d%22%2C%22issuer%22%3A%22http%3A%2F%2Fsignin.aws.amazon.com%2Fsignin%22%2C%22signinType%22%3A%22PUBLIC%22%7D; csd-key=eyJ3YXNtVGVzdGVkIjp0cnVlLCJ3YXNtQ29tcGF0aWJsZSI6dHJ1ZSwid2ViQ3J5cHRvVGVzdGVkIjpmYWxzZSwidiI6MSwia2lkIjoiOTI5ZDY4Iiwia2V5IjoiR0dRS3hoZ1l1VXVPbjNjQmRTbDdkTy80a2Z1MjI2WnpKTGdXS09OQng4QlRDRGtReDVWYnFxU3htcXQrSzRkdDU0TWh3bENpbWFGQm5ReEJwRjIwblBsU1pqZGgyNVpCcTVjemI2aEwzSDMwcFVKZXZtUVJKaVZrZThNc29lcTFrdDh1TGp3ODhmaHpyTUtnaWlKQlhvUWlQaFNPSFZxbXFqTUg1bFpKR2lnWGFhMUl0YXliTHpxVUU2OTZ0bWVuU1k1OFRVNDU2Z3NMWEdBU2lUTnZyK29Ja1RhckFBN2w1cUxROUp1Z2JvRHI2bzBzb2puMjVTMGZSWE1lcUNuTUM4SWFQT082OTg1UVZGTkIxd0pjelRycFVzQjZFeEg4cXBjdVBIYlJUeEs1UzNxSnY5K1Z2VUR0a29RNnVVTkErQzJYand2cEFWVkFFa1V3djZaT3JBPT0ifQ==; session-token=\"rKnsqVM2fzGbIN1bkI8rxANamaoKq2H3kT6Orvj6eLMK1KsuRsseHXOIvweaC6QZiauxj6vf3pod003PxBvKzkgNnshnB74Oh/CARufVe6dmS+K/LsdFPbjjLIQpj3kHtTqa8ABB05EQTVgTPeF6LuefVnH3LzLol0AaiFRlFDnkVA/nunNWilHMHE10US7NMf2iZmWZgBbP71bQrfjRrQ==\"; csm-hit=tb:SARMPC92MC0MJX8GH41R+b-29F64KJKMRJ6M6QSV7Y9|1631583069455&t:1631583069455&adb:adblk_no")
	res, err := m.Monitor.Client.Do(req)
	if err != nil {
		fmt.Println(err)
		fmt.Println(errors.Cause(err))
		return nil
	}
	defer res.Body.Close()
	defer func() {
		watch.Stop()
		// fmt.Printf("Home Depot %s - Code : %d Milli elapsed: %v\n", m.Monitor.Config.Sku, res.StatusCode, watch.Milliseconds())
	}()
	var page *goquery.Document

	if res != nil {
		page, err = goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			fmt.Println(err)
			return errors.New("Error Parsing Doc")
		}
		res.Body.Close()
	}
	var ok bool
	if csrf, ok = page.Find("input[id='aod-atc-csrf-token']").Attr("value"); !ok {
		fmt.Println("Something missing 2")
		m.getCSRF()
		return nil
	}
	fmt.Println(csrf, "TESTING")
	m.csrf = csrf

	return nil
}
func (m *CurrentMonitor) sendRestockNotification(oid string, sku string, title string) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Restock Notification : Site : %s, Product : %s Recovering from panic in printAllOperations error is: %v \n", m.Monitor.Config.Site, m.Monitor.Config.Sku, r)
		}
	}()

	url := "http://159.203.179.167:3030/"

	t := (time.Now().UTC().UnixNano() / 1e6)

	var jsonData = []byte(fmt.Sprintf(`{
		"site": "Amazon",
		"offerId" : "%s",
		"sku" : "%s",
		"title" : "%s",
		"time" : %d
	}`, oid, sku, title, t))
	fmt.Println(string(jsonData))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
	fmt.Println(string(body))
}
