package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/elgs/gojq"
	"github.com/mattia-git/go-capmonster"
)

func TrimSpaceNewlineInString(s string) string {
	re := regexp.MustCompile(` +\r?\n +`)
	return re.ReplaceAllString(s, " ")
}

var datadomeCookie string

func main() {
	// for true {
	// 	sendRequest()
	// 	time.Sleep(5000 * (time.Millisecond))
	// }
	// for true {
	// 	testNewEndpoint()
	// }
	// // t := time.Now()
	// testing := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	// fmt.Println(testing)
	// fmt.Printf("%T", time.Now().UTC().Format("2006-01-02T15:04:05Z"))
	walmartNew()

}
func walmartNew() {

	url := "https://www.walmart.com/search/api/preso?prg=desktop&cat_id=0&facet=brand%3APanini%7C%7Cbrand%3ATopps%7C%7Cretailer%3AWalmart.com&grid=false&query=panini&soft_sort=false&sort=new"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("cookie", `DL=33473%2C%2C%2Cip%2C33473%2C%2C; vtc=Wc-5IJNXcCxYa0fpnFPpUA; TS013ed49a=01538efd7c0085dc072f274bb98d7f3e1623d95bcf8e3b594f2372d8c231a637702fc3fb2b74547f5075075afdf87b2113f04fcc69; TBV=7; _pxvid=7b2febbe-cfa6-11eb-92a4-0242ac120006; _gcl_au=1.1.531143447.1624047166; tb_sw_supported=true; _abck=atl1wiqokoeos5bao5pb_1985; rtoken=MDgyNTUyMDE4uRiNaOR42Li%2FaQqALWZv76zZdbjFdpv5ZbM4tPC%2B9Q26rzIZioWwa%2BTTC3UrNVsPwN1Tj6LkmULy%2B5RaADQRnytPo%2BMX7cdXA2k236eEKBAJi2JuCsXAPdkY4Xm5jUwd3vnArV%2BPLBRPWrgOcepPVULGpf7b2szLnKaieZQxg3Q6RNLx7McyDpt8hausTZOD5tHuax%2F7D%2BZhJoOfnhZIntW2NMKflYKyBHvRLZAdMPUv8JNE4Vga0o2bvXxYMc66rdV%2BOICiFSl6u%2B9soUDxKafMtLuzRYP8NYebsNtHWW1Bg93LiWTdF79xRZco9YtD6OcrmFydbBOCpJADN4cwu3igOjU1YJOpqjBdmMA1VDRajP6UfkCSQSh%2B2g5%2FNSnYes4oLiLLLGugVIO1xvxqvQ%3D%3D; SPID=12b5b2a8bec816faf2bff9663eec6a47a5dc4501089f1cff694a9251488b20f6d5ea803d5ca906a2ce35ad635aac7370myacc; CID=ad1e2798-ac54-4d33-9d9e-3803bca9a057; hasCID=1; customer=%7B%22firstName%22%3A%22Adrian%22%2C%22lastNameInitial%22%3A%22T%22%2C%22rememberme%22%3Atrue%7D; type=REGISTERED; WMP=4; oneapp_customer=true; hasCRT=1; CRT=7ebca80e-4e13-4929-bb1a-971f899f3979; cart-item-count=2; akavpau_p1=1624067154~id=2336dded40d2219f0ba35f145236c033; dtCookie=v_4_srv_36_sn_95D12B925BE6A9867771F833AB58D813_perc_100000_ol_0_mul_1_app-3Aea7c4b59f27d43eb_1; TS01af1d9b=01538efd7c67c354031bd3b3b17ef4818fcf24b60a1ddded8384c72ef31c1b351459853c5756406d97121527a55efe60ae430ad970; athrvi=RVI~h13a77e56-h798aed6; s_sess_2=c32_v%3DS2H%2Cnull%3B%20prop32%3D; next-day=1625257800|true|false|1625313600|1625196691; location-data=33473%3ABoynton%20Beach%3AFL%3A%3A8%3A1|2bn%3B%3B4.67%2C1uk%3B%3B4.7%2C4fz%3B%3B5.94%2C25h%3B%3B6.14%2C12u%3B%3B6.55%2C185%3B%3B6.7%2C5dj%3B%3B7.45%2C4k7%3B%3B8.09%2C4fy%3B%3B8.1%2C1uu%3B%3B9.56||7|1|1ydz%3B16%3B10%3B10.02%2C1ye0%3B16%3B11%3B10.64%2C1yoh%3B16%3B12%3B11.97%2C1ye2%3B16%3B13%3B13.82%2C1y3g%3B16%3B14%3B25.06; TB_Latency_Tracker_100=1; TB_Navigation_Preload_01=1; TB_SFOU-100=1; TB_DC_Flap_Test=0; bstc=T3CzSCsGQaqSFbN8NCvOqY; mobileweb=0; xpa=; xpm=3%2B1625196691%2BWc-5IJNXcCxYa0fpnFPpUA~ad1e2798-ac54-4d33-9d9e-3803bca9a057%2B0; TS01b0be75=01538efd7c0435e2a4f7f69e52557ffcfd7a06d5c2ac584f56a484847cbd8ee0dfb8388f866a7f09fea9bbbd83fcd772d9de594619; _pxff_cfp=1; _pxff_rf=1; _pxff_fp=1; com.wm.reflector="reflectorid:0000000000000000000000@lastupd:1625196720434@firstcreate:1624209371202"; wm_ul_plus=INACTIVE|1625200320462; auth=MTAyOTYyMDE4ZqpjZpvXLIEFVV49dSkA1Osdbi056ZywkLZXImEGSIHdy6UDTAuZJN54E6AYzSBTVzO4qGnXmG%2Fddq4ev%2FU8zjAxll6p58K75cKoQBqK4VPSEKCrJp1xdiMCoIH%2BXg3plrvBbkT8GAVfcnvLIPG4VyzvWvooXNnTyYN6TgAaFqE%2F8JEXyYQfnL0x8Syn1jvflNStVt84BmpJSfY%2Ftlm0IQ88%2FbreJ29hxGEqOcjsqKqJ00UMGQyiYLY97sfSUmPHCgi%2BXQu9OjpsIDyPAwsWQtgnd7gB%2FP2NDJHTrVlH9mCz0JdwmDDvu09GW0LIhDQfYMoSKvTkX3uFCS5L7pVExj9HIKPkiN811ppJbEo6FUw3I5nr37Zcxc7kizq5mWWY8p8Y2F9P%2FHM64p%2BbUtABoQ%3D%3D; akavpau_p8=1625197323~id=accaf851d51f7a43ca45ceaad24e1f62; _px3=a6b70673ec1bb0011b628afceca6d9f67100028a3a5aebd569907b79cc4dd070:IkHPQ07MIw0fDsajd5vBzqh7QRU1GMBjq41a6Q2fn6vFUyAxXr1sFNPzXY2wcYDbOZtBAEPnzjQP3HMEDaWObQ==:1000:ktznwSuWcv8hkmCfdby/Vh/+rHhrUMH8iNYuEDc4h5G7Um7TWVcDEIPAVE2SQILqRbpV3RledF7sflrJjmcQcjZ84Pk/aojXr9yNt6Y48M4YduhQhACD6GydVDH8PvHmOv6GBid5IxoRJdnEHynwoK9ecS6vqkINoy9Z3jescRSn5W/pOB7bbRTVFtgxky1iBtQ0CYhmgpnUX9xl82/qMg==; _uetsid=ff7404b0dae511ebb8b479e0d3604eca; _uetvid=8be1d8a0d07111eb887acdc793251c4f; _pxde=4f0d8d06dfc435f14cb4ae4bd080038957cd801cbb1204a661a568bd685cf03d:eyJ0aW1lc3RhbXAiOjE2MjUxOTY3MzIwNjIsImZfa2IiOjAsImlwY19pZCI6W119`)
	req.Header.Add("authority", "www.walmart.com")
	req.Header.Add("sec-ch-ua", `" Not;A Brand";v="99", "Google Chrome";v="91", "Chromium";v="91"`)
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("dnt", "1")
	req.Header.Add("upgrade-insecure-requests", "1")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36")
	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Add("service-worker-navigation-preload", "true")
	req.Header.Add("sec-fetch-site", "none")
	req.Header.Add("sec-fetch-mode", "navigate")
	req.Header.Add("sec-fetch-user", "?1")
	req.Header.Add("sec-fetch-dest", "document")
	req.Header.Add("accept-language", "en-US,en;q=0.9")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res.StatusCode)
	parser, err := gojq.NewStringQuery(string(body))
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(parser.Data)
	// for key, _ := range parser.Data.(map[string]interface{}) {
	// 	fmt.Println(key)
	// }
	products, err := parser.Query("items")
	if err != nil {
		fmt.Println(err)
	}
	for key, _ := range products.([]interface {}) {
		fmt.Println(parser.Query(fmt.Sprintf("items.[%d].id", key)))
	}
	fmt.Println(len(products.([]interface {})))
}
func testingDB() {
	fmt.Println("Sending request")
	url := "http://localhost:7243/DB"
	payload := strings.NewReader("{\n\t\"site\" : \"TARGET\",\n\t\"sku\" : \"16601601\"\n}")
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	//	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(res.StatusCode)
}
func sendRequest() {
	fmt.Println(datadomeCookie)
	url := "https://www.pokemoncenter.com/tpci-ecommweb-api/product/status/qgqvbkjsheyc2obqgu2dk=?format=zoom.nodatalinks"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("authority", "www.pokemoncenter.com")
	req.Header.Add("pragma", "no-cache")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("dnt", "1")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36")
	req.Header.Add("x-store-scope", "pokemon")
	req.Header.Add("accept", "*/*")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("referer", "https://www.pokemoncenter.com/product/290-80545/pokemon-tcg-champion-s-path-elite-trainer-box")
	req.Header.Add("accept-language", "en-US,en;q=0.9")
	req.Header.Add("cookie", "_sp_id.02f2=741ea5ed87ff1039.1621740658.1.1621740658.1621740658; _sp_ses.02f2=*; datadome="+datadomeCookie)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode == 403 {
		var datadome string
		var ddCookie string
		for _, v := range res.Cookies() {
			switch v.Name {
			case "datadome":
				datadome = v.Value
			default:

			}
		}
		currentUrl := getUrl(string(body))
		currentUrl = currentUrl + datadome
		fmt.Println(currentUrl)
		recaptchaToken := solveCaptcha(currentUrl, "6LccSjEUAAAAANCPhaM2c-WiRxCZ5CzsjR_vd8uX")
		u, _ := req.URL.Parse(currentUrl)
		//fmt.Println(u)
		q := u.Query()
		//fmt.Println(q)
		reqe, _ := http.NewRequest("GET", "https://geo.captcha-delivery.com/captcha/check", nil)
		query := reqe.URL.Query()
		query.Add("cid", q.Get("cid"))
		query.Add("icid", q.Get("initialCid"))
		query.Add("g-recaptcha-response", recaptchaToken)
		query.Add("hash", q.Get("hash"))
		query.Add("ua", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36")
		query.Add("referer", url)
		query.Add("parent_url", currentUrl)
		query.Add("x-forwarded-for", "")
		query.Add("s", "9817")
		reqe.URL.RawQuery = query.Encode()
		//fmt.Println(reqe.URL.RawQuery, "\n\n\n\n\n\n\n")
		//	fmt.Println(reqe.URL)
		response, err := http.DefaultClient.Do(reqe)
		if err != nil {
			fmt.Println(err)
		}
		body, err = ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
		}
		defer response.Body.Close()

		if response.StatusCode == 200 {

			// The dd cookie is in response body :)
			dd := make(map[string]interface{})
			// fmt.Println(dd)
			err = json.Unmarshal(body, &dd)
			if err != nil {
				fmt.Println(err)
			}
			cookie, ok := dd["cookie"]
			//	fmt.Println(cookie)
			if ok {
				ddCookie = cookie.(string)[9:115]
			}
			fmt.Println(ddCookie)
		} else {
			fmt.Println("failed fetching dd cookie")
			fmt.Println(response.StatusCode)
			body, _ := ioutil.ReadAll(res.Body)
			fmt.Println(string(body))
		}
		datadomeCookie = ddCookie
		fmt.Println(datadomeCookie)
		err = CallDataDome(currentUrl)
		if err != nil {
			fmt.Println(err)
		}

	} else {
		fmt.Println(res.StatusCode)
	}

}

func testNewEndpoint() {

	url := "https://search.mobile.walmart.com/v1/products-by-code/UPC/643690296973"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("authority", "www.walmart.com")
	req.Header.Add("pragma", "no-cache")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36")
	req.Header.Add("dnt", "1")
	req.Header.Add("accept", "*/*")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("accept-language", "en-US,en;q=0.9")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	fmt.Println(res.StatusCode)

}
func getUrl(body string) string {
	//	fmt.Println(body)
	//	fmt.Println(strings.Split(body, "'hsh':'"))
	hsh := strings.Split(strings.Split(body, "'hsh':'")[1], "'")[0]
	//	fmt.Println(hsh)
	Initalcid := strings.Split(strings.Split(body, "'cid':'")[1], "'")[0]
	//	fmt.Println(Initalcid)
	url := fmt.Sprintf("https://geo.captcha-delivery.com/captcha/?initialCid=%s&hash=%s&cid=", Initalcid, hsh)
	// fmt.Println(url)
	return url
}

func returnSplitted(s string) string {
	var returnString string
	returnString = strings.Split(s, " ")[0]
	return returnString
}

func solveCaptcha(url string, siteKey string) string {
	c := &capmonster.Client{APIKey: "8883585d0e1c3fab6202435a664a2b28"}
	key, err := c.SendRecaptchaV2(url, siteKey, time.Second*25)
	if err != nil {
		fmt.Println(err)
	} else {
		//	fmt.Println(key)
	}
	return key
}

func CallDataDome(captchaURL string) error {
	jsData := map[string]interface{}{
		"ttst":    34.56000000187487,
		"ifov":    false,
		"wdifts":  false,
		"wdifrm":  false,
		"wdif":    false,
		"br_h":    238,
		"br_w":    1414,
		"br_oh":   782,
		"br_ow":   1414,
		"nddc":    1,
		"rs_h":    900,
		"rs_w":    1440,
		"rs_cd":   30,
		"phe":     false,
		"nm":      false,
		"jsf":     false,
		"ua":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36",
		"lg":      "zh-CN",
		"pr":      2,
		"hc":      8,
		"ars_h":   785,
		"ars_w":   1440,
		"tz":      240,
		"str_ss":  true,
		"str_ls":  true,
		"str_idb": true,
		"str_odb": true,
		"plgod":   false,
		"plg":     2,
		"plgne":   true,
		"plgre":   true,
		"plgof":   false,
		"plggt":   false,
		"pltod":   false,
		"lb":      false,
		"eva":     33,
		"lo":      false,
		"ts_mtp":  0,
		"ts_tec":  false,
		"ts_tsa":  false,
		"vnd":     "Google Inc.",
		"bid":     "NA",
		"mmt":     "application/pdf,application/x-google-chrome-pdf",
		"plu":     "Chrome PDF Plugin,Chrome PDF Viewer",
		"hdn":     false,
		"awe":     false,
		"geb":     false,
		"dat":     false,
		"med":     "defined",
		"aco":     "probably",
		"acots":   false,
		"acmp":    "probably",
		"acmpts":  true,
		"acw":     "probably",
		"acwts":   false,
		"acma":    "maybe",
		"acmats":  false,
		"acaa":    "probably",
		"acaats":  true,
		"ac3":     "",
		"ac3ts":   false,
		"acf":     "probably",
		"acfts":   false,
		"acmp4":   "maybe",
		"acmp4ts": false,
		"acmp3":   "probably",
		"acmp3ts": false,
		"acwm":    "maybe",
		"acwmts":  false,
		"ocpt":    false,
		"vco":     "probably",
		"vcots":   false,
		"vch":     "probably",
		"vchts":   true,
		"vcw":     "probably",
		"vcwts":   true,
		"vc3":     "maybe",
		"vc3ts":   false,
		"vcmp":    "",
		"vcmpts":  false,
		"vcq":     "",
		"vcqts":   false,
		"vc1":     "probably",
		"vc1ts":   false,
		"dvm":     8,
		"sqt":     false,
		"so":      "landscape-primary",
		"wbd":     false,
		"wbdm":    true,
		"wdw":     true,
		"cokys":   "bG9hZFRpbWVzY3NpYXBwcnVudGltZQ==L=",
		"ecpc":    false,
		"lgs":     true,
		"lgsod":   false,
		"bcda":    true,
		"idn":     true,
		"capi":    false,
		"svde":    false,
		"vpbq":    true,
		"xr":      true,
		"bgav":    true,
		"rri":     true,
		"idfr":    true,
		"ancs":    true,
		"inlc":    true,
		"cgca":    true,
		"inlf":    true,
		"tecd":    true,
		"sbct":    true,
		"aflt":    true,
		"rgp":     true,
		"bint":    true,
		"spwn":    false,
		"emt":     false,
		"bfr":     false,
		"dbov":    false,
		"glvd":    "Apple",
		"glrd":    "Apple M1",
		"tagpu":   28.374999999869033,
		"prm":     true,
		"tzp":     "America/New_York",
		"cvs":     true,
		"usb":     "defined",
		"mp_cx":   441,
		"mp_cy":   223,
		"mp_tr":   true,
		"mp_mx":   -32,
		"mp_my":   48,
		"mp_sx":   467,
		"mp_sy":   355,
		"dcok":    ".captcha-delivery.com",
		"ewsi":    false,
	}

	events := []map[string]interface{}{
		{
			"source":  map[string]interface{}{"x": 0, "y": 110},
			"message": "scroll",
			"date":    time.Now().UnixNano()/int64(time.Millisecond) - 1000 - rand.Int63n(500),
			"id":      2,
		},
		{
			"source":  map[string]interface{}{"x": 594, "y": 12},
			"message": "mouse move",
			"date":    time.Now().UnixNano()/int64(time.Millisecond) - 800 - rand.Int63n(300),
			"id":      0,
		},
		{
			"source":  map[string]interface{}{"x": 0, "y": 25},
			"message": "scroll",
			"date":    time.Now().UnixNano()/int64(time.Millisecond) - 600 - rand.Int63n(100),
			"id":      2,
		},
		{
			"source":  map[string]interface{}{"x": 645, "y": 237},
			"message": "mouse move",
			"date":    time.Now().UnixNano()/int64(time.Millisecond) - 400 - rand.Int63n(50),
			"id":      0,
		},
	}

	eventCounters := map[string]interface{}{
		"mouse move":  rand.Intn(3),
		"mouse click": rand.Intn(3),
		"scroll":      rand.Intn(3),
		"touch start": 0,
		"touch end":   0,
		"touch move":  0,
		"key press":   0,
		"key down":    0,
		"key up":      0,
	}
	jsDataStr, _ := json.Marshal(jsData)
	eventsStr, _ := json.Marshal(events)
	eventCountersStr, _ := json.Marshal(eventCounters)
	u, _ := url.Parse(captchaURL)
	data := url.Values{
		"jsData":        {string(jsDataStr)},
		"events":        {string(eventsStr)},
		"eventCounters": {string(eventCountersStr)},
		"jsType":        {[]string{"le", "ch"}[rand.Intn(2)]},
		"cid":           {u.Query().Get("cid")},
		"ddk":           {"5B45875B653A484CC79E57036CE9FC"},
		"Referer":       {url.QueryEscape(captchaURL)},
		"request":       {url.QueryEscape(captchaURL[32:])},
		"responsePage":  {"origin"},
		"ddv":           {"4.1.50"},
	}
	req, _ := http.NewRequest("POST", "https://api-js.datadome.co/js/", strings.NewReader(data.Encode()))
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"89\", \"Chromium\";v=\"89\", \";Not A Brand\";v=\"99\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Origin", "https://www.pokemoncenter.com/")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://www.pokemoncenter.com/")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Println("Posting Data Status : ", resp.StatusCode)

	if resp.StatusCode == 200 {
		dd := make(map[string]interface{})
		err = json.Unmarshal(body, &dd)
		if err != nil {
			return err
		}
		cookie, ok := dd["cookie"]
		if ok {
			datadomeCookie = cookie.(string)[9:115]
			fmt.Println("Posting Data Cookie : ", datadomeCookie)
		}
	}
	return nil
}
