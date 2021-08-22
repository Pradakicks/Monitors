package RestirMonitor

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/bradhe/stopwatch"
	"github.com/nickname32/discordhook"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	Webhook "github.con/prada-monitors-go/helpers/discordWebhook"
	helper "github.con/prada-monitors-go/helpers/mongo"
	FetchProxies "github.con/prada-monitors-go/helpers/proxy"
	Types "github.con/prada-monitors-go/helpers/types"
)

type CurrentMonitor struct {
	Monitor  Types.Monitor
	products []ProductIdentifier
}

type Product struct {
	ArrivalYmd  string `json:"arrival_ymd,omitempty"`
	BrandID     string `json:"brand_id,omitempty"`
	BrandName   string `json:"brand_name,omitempty"`
	BrandNameJp string `json:"brand_name_jp,omitempty"`
	Code        string `json:"code,omitempty"`
	// DeliveryItemChargeAmount                 int64         `json:"delivery_item_charge_amount,omitempty"`
	DeliveryItemChargeAmountStr              string        `json:"delivery_item_charge_amount_str,omitempty"`
	DeliveryItemChargeCode                   string        `json:"delivery_item_charge_code,omitempty"`
	DeliveryItemChargeInternationalAmount    int64         `json:"delivery_item_charge_international_amount,omitempty"`
	DeliveryItemChargeInternationalAmountStr string        `json:"delivery_item_charge_international_amount_str,omitempty"`
	Icons                                    []interface{} `json:"icons,omitempty"`
	ID                                       int64         `json:"id,omitempty"`
	IsAsk                                    bool          `json:"is_ask,omitempty"`
	IsComingSoon                             bool          `json:"is_coming_soon,omitempty"`
	IsCouponApplied                          bool          `json:"is_coupon_applied,omitempty"`
	IsDisplayOnly                            bool          `json:"is_display_only,omitempty"`
	IsInternationalShippingPossible          bool          `json:"is_international_shipping_possible,omitempty"`
	IsLimitedSale                            bool          `json:"is_limited_sale,omitempty"`
	IsLottery                                bool          `json:"is_lottery,omitempty"`
	IsNew                                    bool          `json:"is_new,omitempty"`
	IsOpen                                   bool          `json:"is_open,omitempty"`
	IsSecret                                 bool          `json:"is_secret,omitempty"`
	IsSoftClose                              bool          `json:"is_soft_close,omitempty"`
	ItemLink                                 string        `json:"item_link,omitempty"`
	ListImageURL                             string        `json:"list_image_url,omitempty"`
	ListLink                                 string        `json:"list_link,omitempty"`
	ListMoImageURL                           string        `json:"list_mo_image_url,omitempty"`
	Name                                     string        `json:"name,omitempty"`
	// OffRate                                  int64         `json:"off_rate,omitempty"`
	// Price                                    int64         `json:"price,omitempty"`
	// PriceProper                              int64         `json:"price_proper,omitempty"`
	PriceProperStr string `json:"price_proper_str,omitempty"`
	// PriceProperSum                           int64         `json:"price_proper_sum,omitempty"`
	PriceProperSumStr string `json:"price_proper_sum_str,omitempty"`
	PriceStr          string `json:"price_str,omitempty"`
	// PriceSum                                 int64         `json:"price_sum,omitempty"`
	PriceSumStr string `json:"price_sum_str,omitempty"`
	// ReceivePointCampaignSum                  int64         `json:"receive_point_campaign_sum,omitempty"`
	// ReceivePointSum                          int64         `json:"receive_point_sum,omitempty"`
	Stock int64 `json:"stock,omitempty"`
}

type ProductIdentifier struct {
	ID       int64  `json:"id,omitempty"`
	ItemLink string `json:"item_link,omitempty"`
}
type RestirSearchResponse struct {
	Message string `json:"message,omitempty"`
	Result  struct {
		Page struct {
			PageNo int64 `json:"page_no,omitempty"`
			Size   int64 `json:"size,omitempty"`
			Start  int64 `json:"start,omitempty"`
			Total  int64 `json:"total,omitempty"`
		} `json:"page,omitempty"`
		Records []Product `json:"records,omitempty"`
	} `json:"result,omitempty"`
	Status int64 `json:"status,omitempty"`
}

var RestirCollection = helper.ConnectDBRestir()

func NewMonitor(sku string, collection *mongo.Collection) *CurrentMonitor {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Site : %s, Product : %s Recovering from panic in printAllOperations error is: %v \n", sku, sku, r)
		}
	}()
	fmt.Println("Restir Monitor ", sku)
	m := CurrentMonitor{}
	// m.Monitor.Availability = "Copy"
	m.Monitor.Config.Site = "Restir"
	m.Monitor.Config.Sku = sku
	m.Monitor.Collection = collection
	m.Monitor.Client = http.Client{Timeout: 10 * time.Second}
	proxyList := FetchProxies.Get()

	go func() {
		for !m.Monitor.Stop {
			go func() {
				cur, err := m.Monitor.Collection.Find(context.TODO(), bson.M{})

				if err != nil {
					fmt.Println(err)
					fmt.Println(errors.Cause(err))

					fmt.Println(err)
					fmt.Println(errors.Cause(err))

					fmt.Println(err)
					fmt.Println(errors.Cause(err))

					fmt.Println(err)
					fmt.Println(errors.Cause(err))

				}
				defer cur.Close(context.TODO())
				var testingArr []ProductIdentifier
				for cur.Next(context.TODO()) {

					// create a value into which the single document can be decoded
					var prod Product
					// & character returns the memory address of the following variable.
					err := cur.Decode(&prod) // decode similar to deserialize process.
					if err != nil {
						fmt.Println(err)
						fmt.Println(errors.Cause(err))

					}

					// deletedResult, err := m.collection.DeleteOne(context.TODO(), bson.M{"id": prod.ID})
					// fmt.Println(deletedResult, prod.ID)
					// if err != nil {
					// 	fmt.Println(err)
					// fmt.Println(errors.Cause(err))

					// }
					// add item our array
					testingArr = append(testingArr, ProductIdentifier{ID: prod.ID})
				}

				if err := cur.Err(); err != nil {
					fmt.Println(err)
					fmt.Println(errors.Cause(err))

				}
				m.products = testingArr
				time.Sleep(10 * time.Second)
				fmt.Println("Restarting", len(m.products), len(testingArr))
			}()
			time.Sleep(10 * time.Second)

		}
	}()

	time.Sleep(5000 * (time.Millisecond))
	go m.Monitor.CheckStop()
	time.Sleep(3000 * (time.Millisecond))

	i := true
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
			m.Monitor.Client.Transport = defaultTransport
			go m.monitor()
			time.Sleep(500 * (time.Millisecond))
		} else {
			fmt.Println(m.Monitor.Config.Sku, "STOPPED STOPPED STOPPED")
			i = false
		}

	}
	return &m
}

func (m *CurrentMonitor) monitor() error {
	watch := stopwatch.Start()

	url := "https://www.restir.com/brand/lifestyle/b/20710"

	payload := strings.NewReader("{\"item_category_code\":\"lifestyle\",\"brand_id\":\"20710\",\"master_category_code\":null,\"category_code\":null}")

	req, _ := http.NewRequest("POST", url, payload)

	// req.Header.Add("cookie", "_mkra_stck=a983840e3e72eae8c8f1682800d8312f%253A1629244102.857295; AUI=EB746501D88FF65C7DCC0F503CAF728AF893D68F782B89E62B4D2BD83587432E; _front_session=b9c454496c9735162286d2178a5592f4")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Pragma", "no-cache")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("sec-ch-ua", `"Chromium";v="92", " Not A;Brand";v="99", "Google Chrome";v="92"`)
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("DNT", "1")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36")
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	req.Header.Add("Origin", "https://www.restir.com")
	req.Header.Add("Sec-Fetch-Site", "same-origin")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Referer", "https://www.restir.com/brand/lifestyle/b/20710")
	req.Header.Add("Accept-Language", "en-US,en;q=0.9")
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
		fmt.Printf("Restir %s - Code : %d Milli elapsed: %v\n", m.Monitor.Config.Sku, res.StatusCode, watch.Milliseconds())
	}()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		fmt.Println(errors.Cause(err))
		return nil
	}
	if res.StatusCode != 200 {
		time.Sleep(10 * time.Second)
		return nil
	}

	var jsonResponse RestirSearchResponse

	err = json.Unmarshal([]byte(body), &jsonResponse)
	if err != nil {
		fmt.Println(err)
		fmt.Println(errors.Cause(err))
		return nil
	}

	for _, currentProduct := range jsonResponse.Result.Records {
		isPresent := false

		for _, identifer := range m.products {

			// if identifer.URL == "nfl/green-bay-packers/bart-starr-green-bay-packers-1960-topps-number-51-card-bvg-55/o-1372+t-03820298+p-71703716874+z-9-1863514192" && currentProduct.URL == "nfl/green-bay-packers/bart-starr-green-bay-packers-1960-topps-number-51-card-bvg-55/o-1372+t-03820298+p-71703716874+z-9-1863514192" {
			// 	fmt.Println("HERER ", identifer, currentProduct.ProductID, currentProduct.URL)
			// }

			if identifer.ID == currentProduct.ID {
				isPresent = true
			}
		}

		if !isPresent {
			fmt.Println("New Item in Stock", currentProduct.Name)
			fmt.Println(isPresent, currentProduct.Name, currentProduct.ID)
			m.products = append(m.products, ProductIdentifier{ID: currentProduct.ID, ItemLink: currentProduct.ItemLink})
			link := fmt.Sprintf("https://www.restir.com/%s", currentProduct.ItemLink)
			// image := fmt.Sprintf("https://fanatics.frgimages.com/%s", currentProduct.ImageSelector.DefaultImage.Image.Src)
			go m.sendWebhook(currentProduct.BrandID, currentProduct.Name, currentProduct.PriceProperStr, link, currentProduct.ListImageURL)
			result, err := m.Monitor.Collection.InsertOne(context.TODO(), currentProduct)
			if err != nil {
				fmt.Println(err)
				fmt.Println(errors.Cause(err))

			} else {
				fmt.Println("New Restir Product ", result)
			}
		}

	}

	m.Monitor.AvailabilityBool = m.Monitor.CurrentAvailabilityBool
	return nil
}

func (m *CurrentMonitor) sendWebhook(sku string, name string, price string, link string, image string) error {
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
	t := time.Now().UTC()

	// for _, comp := range m.Monitor.CurrentCompanies {
	// fmt.Println(comp.Company)
	currentCompany := Types.Company{Company: "Test", Webhook: "https://discord.com/api/webhooks/797249480410923018/NPL3ktXS78z5EHo_cpYyrtFl_2iB0ARgz9IW5kwAZA-UkiseiinnBmUPJZlGgxw8TZiW", Color: "1752220", CompanyImage: "https://cdn.discordapp.com/attachments/802755133582475315/842627264482508820/unknown.png"}

	go m.webHookSend(currentCompany, sku, m.Monitor.Config.Site, name, price, link, t, image)
	// go m.webHookSend(comp, sku, m.Monitor.Config.Site, name, price, link, t, image)
	// }
	return nil
}

func (m *CurrentMonitor) webHookSend(c Types.Company, sku string, site string, name string, price string, link string, currentTime time.Time, image string) {
	Color, err := strconv.Atoi(c.Color)
	if err != nil {
		fmt.Println(err)
		fmt.Println(errors.Cause(err))
		return
	}
	var currentFields []*discordhook.EmbedField
	currentFields = append(currentFields, &discordhook.EmbedField{
		Name:   "Product Name",
		Value:  name,
		Inline: false,
	})
	currentFields = append(currentFields, &discordhook.EmbedField{
		Name:   "Product Identifer",
		Value:  sku,
		Inline: true,
	})
	currentFields = append(currentFields, &discordhook.EmbedField{
		Name:   "Price",
		Value:  price,
		Inline: true,
	})
	currentFields = append(currentFields, &discordhook.EmbedField{
		Name:  "Links",
		Value: "[Cart](https://www.restir.com/cart) | [Login](https://www.restir.com/login?nextPathname=/account)",
	})

	var discordParams discordhook.WebhookExecuteParams = discordhook.WebhookExecuteParams{
		Content: "",
		Embeds: []*discordhook.Embed{
			{
				Title:  "Restir New Products" + " Monitor",
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

func Restir(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)
	w.Header().Set("Content-Type", "application/json")
	var currentMonitor Types.MonitorResponse
	_ = json.NewDecoder(r.Body).Decode(&currentMonitor)
	fmt.Println(currentMonitor)
	go NewMonitor(currentMonitor.Sku, RestirCollection)
	json.NewEncoder(w).Encode(currentMonitor)
}
