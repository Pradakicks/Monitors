package Amazon

import (
	// "bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/bradhe/stopwatch"
	"github.com/nickname32/discordhook"
	"github.com/pkg/errors"

	Webhook "github.con/prada-monitors-go/helpers/discordWebhook"
	FetchProxies "github.con/prada-monitors-go/helpers/proxy"
	Types "github.con/prada-monitors-go/helpers/types"
)


func NewMonitormonitorFirst(sku string) *CurrentMonitor {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Site : %s, Product : %s Recovering from panic in printAllOperations error is: %v \n", sku, sku, r)
		}
	}()
	fmt.Println("AMAZON")
	m := CurrentMonitor{}
	m.Monitor.Availability = "OUT_OF_STOCK_ONLINE"
	m.Monitor.Config.Site = "Amazon"
	m.Monitor.Config.Sku = sku
	m.Monitor.Client = http.Client{Timeout: 10 * time.Second}
	proxyList := FetchProxies.Get()

	// time.Sleep(15000 * (time.Millisecond))
	// go m.Monitor.CheckStop()
	// time.Sleep(3000 * (time.Millisecond))

	i := true
	jar, _ := cookiejar.New(nil)
	m.Monitor.Client = http.Client{
		Jar: jar,
	}
	m.getCookies()
	m.getCookies()
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
				Jar:       jar,
				Transport: defaultTransport,
			}
			m.monitor()
			// time.Sleep(500 * (tixme.Millisecond))
		} else {
			fmt.Println(m.Monitor.Config.Sku, "STOPPED STOPPED STOPPED")
			i = false
		}

	}
	return &m
}

func (m *CurrentMonitor) monitorFirst() error {
	watch := stopwatch.Start()
	fmt.Println("HERE")
	oid := "Qmw1Lc22Jrfu5q1XTiZ3PLJlYAVCN7fK5bRy%2Bh2OP4qSsJMfsqZ8Zfy9TuBqSqSMDI9W%2FgWKiDqyhLI8g4tzsXV3wYysyGu9o6tja6qhDybdTdgyT8D7OKIWXUWkZJcWNwfxhU8E9t9ZqT%2FAq3zz2w%3D%3D"
	payload := strings.NewReader(fmt.Sprintf(`marketplaceId=ATVPDKIKX0DER&asin=%v&customerId=&sessionId=%v&accessoryItemAsin=B002M40VJM&accessoryItemOfferingId=%v&languageOfPreference=en_US&accessoryItemQuantity=1&accessoryItemPrice=9.99&accessoryMerchantId=ATVPDKIKX0DER&accessoryProductGroupId=8652000`, m.Monitor.Config.Sku, m.sid, oid))
	// fmt.Println(fmt.Sprintf(`marketplaceId=ATVPDKIKX0DER&asin=%v&customerId=&sessionId=%v&accessoryItemAsin=B002M40VJM&accessoryItemOfferingId=%v&languageOfPreference=en_US&accessoryItemQuantity=1&accessoryItemPrice=9.99&accessoryMerchantId=ATVPDKIKX0DER&accessoryProductGroupId=8652000`, m.Monitor.Config.Sku, m.sid, oid))
	req, err := http.NewRequest("POST", "https://smile.amazon.com/gp/product/features/aloha-ppd/udp-ajax-handler/attach-add-to-cart.html", payload)
	if err != nil {
		fmt.Println(err)
		fmt.Println(errors.Cause(err))
		return nil
	}
	req.Header.Add("authority", "smile.amazon.com")
    req.Header.Add("sec-ch-ua", `"Chromium";v="92", " Not A;Brand";v="99", "Google Chrome";v="92"`)
    req.Header.Add("rtt", "0")
    req.Header.Add("sec-ch-ua-mobile", "?0")
    req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36")
    req.Header.Add("content-type", "application/x-www-form-urlencoded")
    req.Header.Add("accept", "/")
    req.Header.Add("x-requested-with", "XMLHttpRequest")
    req.Header.Add("downlink", "10")
    req.Header.Add("ect", "4g")
    req.Header.Add("origin", "https://smile.amazon.com/")
    req.Header.Add("sec-fetch-site", "same-origin")
    req.Header.Add("sec-fetch-mode", "cors")
    req.Header.Add("sec-fetch-dest", "empty")
    req.Header.Add("referer", "https://smile.amazon.com/Adhesive-Organizer-Holder-Durable-Management/dp/B07Y4ZYRQ3/ref=pd_ybh_a_6?_encoding=UTF8&psc=1&refRID=RKBCY923SXTSR97XRD2Q")
    req.Header.Add("accept-language", "en-US,en;q=0.9")
	// req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36")
	req.Header.Add("cookie", fmt.Sprintf("session-id=%s", m.sid))

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
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		fmt.Println(errors.Cause(err))
		return nil
	}
	var realBody AmazonResponse
	json.Unmarshal(body, &realBody)
	fmt.Println(realBody)
	fmt.Println(res.StatusCode)
	fmt.Println(string(body))

	return nil
}

func (m *CurrentMonitor) sendWebhookmonitorFirst(sku string, name string, price float64, link string, image string, Qty int) error {
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

func (m *CurrentMonitor) webHookSendmonitorFirst(c Types.Company, site string, name string, price float64, link string, currentTime time.Time, image string, Qty int) {
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
func AmazonmonitorFirst(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "Amazon Monitor")
	fmt.Println("Amazon")
	var currentMonitor Types.MonitorResponse
	_ = json.NewDecoder(r.Body).Decode(&currentMonitor)
	fmt.Println(currentMonitor)
	go NewMonitor(currentMonitor.Sku, currentMonitor.SkuName)
	json.NewEncoder(w).Encode(currentMonitor)
}

func (m *CurrentMonitor) getCookiesmonitorFirst() error {
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
				m.sid = strings.Split(strings.Split(v[0], "session-id=")[1],";")[0]
				fmt.Println(strings.Split(strings.Split(v[0], "session-id=")[1],";")[0])
			}
			
			
		}
	}

	return nil
}
