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
	// t := time.Now()
	testing := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	fmt.Println(testing)
	fmt.Printf("%T", time.Now().UTC().Format("2006-01-02T15:04:05Z"))

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
