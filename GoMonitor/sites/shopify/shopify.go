package Shopify

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bradhe/stopwatch"
	"github.com/nickname32/discordhook"
	"github.com/pkg/errors"

	Webhook "github.con/prada-monitors-go/helpers/discordWebhook"
	FetchProxies "github.con/prada-monitors-go/helpers/proxy"
	Types "github.con/prada-monitors-go/helpers/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	products            []ProductsItem
	stop                bool
	CurrentCompanies    []Company
	collection          *mongo.Collection
	zerosArray          []string
	Keywords            []string
	NegKeywords         []string
	BaseUrl             string
}
type Product struct {
	name        string
	stockNumber string
	productId   string
	price       int
	image       string
	link        string
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
type ShopifyJson struct {
	Products []Types.ShopifyNewProduct
}
type ProductsItem struct {
	ID          int64
	UpdatedTime string
	Variants    []Types.Variant
	Store       string
	Handle      string
}

func NewMonitor(sku string, collection *mongo.Collection) *Monitor {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Site : %s, Product : %s Recovering from panic in printAllOperations error is: %v \n", sku, sku, r)
		}
	}()
	fmt.Println("TESTING")
	m := Monitor{}
	m.Availability = "OUT_OF_STOCK_ONLINE"
	m.collection = collection
	// var err error
	//	m.Client = http.Client{Timeout: 10 * time.Second}
	m.Config.site = "Shopify"
	m.Config.startDelay = 3000
	m.Config.sku = sku
	// 	m.file, err = os.Create("./testing.txt")
	m.Client = http.Client{Timeout: 10 * time.Second}
	m.monitorProduct.name = "Testing Product"
	m.monitorProduct.stockNumber = ""

	proxyList := FetchProxies.Get()
	m.Keywords = append(m.Keywords, "jordan")
	m.Keywords = append(m.Keywords, "air")
	m.Keywords = append(m.Keywords, "yeezy")
	m.Keywords = append(m.Keywords, "adidas")
	m.Keywords = append(m.Keywords, "Nike")
	// m.Keywords = append(m.Keywords, "New")
	m.Keywords = append(m.Keywords, "Balance")
	m.NegKeywords = append(m.NegKeywords, "Fleece")
	m.NegKeywords = append(m.NegKeywords, "Tee")
	m.NegKeywords = append(m.NegKeywords, "Accessories")
	m.NegKeywords = append(m.NegKeywords, "Hat")
	m.NegKeywords = append(m.NegKeywords, "Pant")
	m.NegKeywords = append(m.NegKeywords, "Track Jacket")
	m.NegKeywords = append(m.NegKeywords, "Ultraboost")
	m.NegKeywords = append(m.NegKeywords, "Sock")
	m.NegKeywords = append(m.NegKeywords, "Presto")
	m.NegKeywords = append(m.NegKeywords, "Jeans")
	m.NegKeywords = append(m.NegKeywords, "T-Shirt")
	m.NegKeywords = append(m.NegKeywords, "Snapback")
	m.NegKeywords = append(m.NegKeywords, "Leggings")
	m.NegKeywords = append(m.NegKeywords, "Mesh Top")
	m.NegKeywords = append(m.NegKeywords, "Max Aura")
	m.NegKeywords = append(m.NegKeywords, "Leggings")
	m.NegKeywords = append(m.NegKeywords, "Air Structure")
	m.NegKeywords = append(m.NegKeywords, "Huarache")
	m.NegKeywords = append(m.NegKeywords, "Cargo Short")
	m.NegKeywords = append(m.NegKeywords, "Court Polo")
	m.NegKeywords = append(m.NegKeywords, "Jacket")
	m.NegKeywords = append(m.NegKeywords, "Dress")
	m.NegKeywords = append(m.NegKeywords, "Short")
	m.NegKeywords = append(m.NegKeywords, "Sweatshirt")
	m.NegKeywords = append(m.NegKeywords, "Sweatshirt")
	// fmt.Println(timeout)
	//m.Availability = "OUT_OF_STOCK"
	//fmt.Println(m)
	// time.Sleep(15000 * (time.Millisecond))
	// go m.checkStop()
	// time.Sleep(3000 * (time.Millisecond))

	// Close the cursor once finished
	/*A defer statement defers the execution of a function until the surrounding function returns.
	simply, run cur.Close() process but after cur.Next() finished.*/
	go func() {
		for {
			go func() {
				cur, err := collection.Find(context.TODO(), bson.M{"store": m.Config.sku})

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
				var testingArr []ProductsItem
				for cur.Next(context.TODO()) {

					// create a value into which the single document can be decoded
					var prod Types.ShopifyNewProduct
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
					testingArr = append(testingArr, ProductsItem{
						ID:          prod.ID,
						UpdatedTime: prod.UpdatedAt,
						Variants:    prod.Variants,
						Store:       m.Config.sku,
						Handle:      prod.Handle,
					})
				}

				if err := cur.Err(); err != nil {
					fmt.Println(err)
					fmt.Println(errors.Cause(err))

				}
				// fmt.Println(len(m.products), len(testingArr))
				m.products = testingArr

				// fmt.Println(len(m.products))
				time.Sleep(10 * time.Second)
				fmt.Println("Restarting")
			}()
			time.Sleep(10 * time.Second)

		}
	}()

	// time.Sleep(1000000 * time.Second)

	i := true
	for i {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Site : %s, Product : %s Recovering from panic in printAllOperations error is: %v \n", m.Config.site, m.Config.sku, r)
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
				// fmt.Println(err)
				fmt.Println(errors.Cause(err))

				return nil
			}
			defaultTransport := &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			}
			m.Client.Transport = defaultTransport
			go m.monitor()
			time.Sleep(500 * (time.Millisecond))
		} else {
			fmt.Println(m.Config.sku, "STOPPED STOPPED STOPPED")
			i = false
		}

	}
	return &m
}

func (m *Monitor) monitor() error {
	// start := time.Now()

	watch := stopwatch.Start()
	t := time.Now().UTC().UnixNano()
	var url string
	var baseUrl string
	switch m.Config.sku {
	case "ShopNiceKicks":
		baseUrl = fmt.Sprintf("https://%s.com", m.Config.sku)
		url = fmt.Sprintf("%s/products.json?limit=%d", baseUrl, t)
		break
	case "travisscott":
		baseUrl = fmt.Sprintf("https://shop.%s.com", m.Config.sku)
		url = fmt.Sprintf("%s/products.json?limit=%d", baseUrl, t)
		break
	case "deadstock":
		baseUrl = fmt.Sprintf("https://www.%s.ca", m.Config.sku)
		url = fmt.Sprintf("%s/products.json?limit=%d", baseUrl, t)
		break
	case "mountaindew":
		baseUrl = fmt.Sprintf("https://store.%s.com", m.Config.sku)
		url = fmt.Sprintf("%s/products.json?limit=%d", baseUrl, t)
		break
	default:
		baseUrl = fmt.Sprintf("https://www.%s.com", m.Config.sku)
		url = fmt.Sprintf("%s/products.json?limit=%d", baseUrl, t)
	}
	m.BaseUrl = baseUrl
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		fmt.Println(errors.Cause(err))

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
	// res, err := http.DefaultClient.Do(req)
	res, err := m.Client.Do(req)
	if err != nil {
		fmt.Println(err)
		fmt.Println(errors.Cause(err))

		return nil
	}
	// elapsed := time.Since(start)
	// fmt.Println("Time Elapsed Request : ", elapsed)
	defer res.Body.Close()
	var newList []ProductsItem
	defer func() {
		watch.Stop()
		fmt.Printf("Shopify Store %s - Code : %d Cache: %s Length : %d :Old %d Milli elapsed: %v\n", m.Config.sku, res.StatusCode, res.Header["X-Cache"], len(newList), len(m.products), watch.Milliseconds())
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

	var jsonResponse ShopifyJson

	err = json.Unmarshal([]byte(body), &jsonResponse)
	if err != nil {
		fmt.Println(err)
		fmt.Println(errors.Cause(err))
		return nil
	}
	// elapsed = time.Since(start)
	// fmt.Println("Time Elapsed BODY : ", elapsed)
	// fmt.Println("Length of Products in ResponseWriter", len(jsonResponse.Products))
	for _, value := range jsonResponse.Products {
		var isPresent bool
		var containsKeywords bool
		var containsNegativeKeywords bool
		value.Store = m.Config.sku

		for _, v := range m.products {
			if v.Handle == value.Handle && v.Store == value.Store {
				isPresent = true
				newList = append(newList, ProductsItem{ID: value.ID, UpdatedTime: value.UpdatedAt, Variants: value.Variants, Store: value.Store, Handle: value.Handle})
				if v.UpdatedTime != value.UpdatedAt { // This monitor all changes in a product
					fmt.Println("Restock!", v.UpdatedTime, value.UpdatedAt)
					go func(value Types.ShopifyNewProduct) {
						go func(value Types.ShopifyNewProduct) {
							filter := bson.M{"id": value.ID, "handle": value.Handle, "store": value.Store}
							var testing Types.ShopifyNewProduct
							replacedResult := m.collection.FindOneAndDelete(context.TODO(), filter).Decode(&testing)
							if replacedResult != nil {
								fmt.Printf("remove fail %v\n", replacedResult)
								fmt.Printf("remove fail %v\n", replacedResult)
								fmt.Printf("remove fail %v\n", replacedResult)
								fmt.Println(errors.Cause(replacedResult))
							} else {
								fmt.Println("Remove Result", testing.ID, testing.Store, testing.Handle, replacedResult)
								var testingAgaain Types.ShopifyNewProduct
								isPresent := m.collection.FindOne(context.TODO(), filter).Decode(&testingAgaain)
								fmt.Println("Finding", value.Handle, value.Store, "After", testingAgaain, isPresent)

								result, err := m.collection.InsertOne(context.TODO(), value)
								fmt.Println("Insterting", value.Handle, value.Store, "Now")
								if err != nil {
									fmt.Println("Adding Fail", err)
									// fmt.Println(errors.Cause(err))
								} else {
									fmt.Println("Instered ", result)
								}

							}
						}(value)
						restockVariants := m.restockedVariants(v.Variants, value.Variants)
						if len(restockVariants) != 0 {

							fmt.Println("Restocked Variants", len(restockVariants))

							link := fmt.Sprintf("https://www.%s.com/products/%s", m.Config.sku, value.Handle)
							var testCompany Company
							testCompany.Webhook = "https://webhooks.aycd.io/webhooks/api/v1/send/8028/d1464662-73f6-4971-83c3-609e923d170e"
							testCompany.Color = "1752220"
							testCompany.CompanyImage = "https://cdn.discordapp.com/attachments/802755133582475315/842627264482508820/unknown.png"
							testCompany.Company = "Testing"
							t := time.Now().UTC()
							var price string
							var image string
							if len(value.Variants) < 1 || len(value.Images) == 0 {
								price = "N/A"
								image = testCompany.CompanyImage
							} else {
								price = value.Variants[0].Price
								image = value.Images[0].Src
							}
							go m.webHookSend(testCompany, m.Config.sku, value.Title, price, link, t, image, restockVariants)
						}
					}(value)
					continue
				} else {
					continue
				}

			}
		}

		// Add Keywords
		for _, val := range m.Keywords {
			if strings.Contains(value.Title, val) || strings.Contains(value.Handle, val) {
				containsKeywords = true
				continue
			}
		}
		
		for _, val := range m.NegKeywords {
			if strings.Contains(value.Title, val) || strings.Contains(value.Handle, val) {
				containsNegativeKeywords = true
				continue
			}
		}
		if !isPresent && containsKeywords && !containsNegativeKeywords {
			newList = append(newList, ProductsItem{ID: value.ID, UpdatedTime: value.UpdatedAt, Variants: value.Variants, Store: value.Store, Handle: value.Handle})
			go func(value Types.ShopifyNewProduct) {
				result, err := m.collection.InsertOne(context.TODO(), value)
				if err != nil {
					// fmt.Println(err)
					// fmt.Println(errors.Cause(err))

				} else {
					fmt.Println("New Product ", value.ID, value.UpdatedAt, value.Store, value.Handle, len(value.Images), len(value.Variants))
					fmt.Println(result)
					link := fmt.Sprintf("https://www.%s.com/products/%s", m.Config.sku, value.Handle)
					var testCompany Company
					testCompany.Webhook = "https://webhooks.aycd.io/webhooks/api/v1/send/8028/d1464662-73f6-4971-83c3-609e923d170e"
					testCompany.Color = "1752220"
					testCompany.CompanyImage = "https://cdn.discordapp.com/attachments/802755133582475315/842627264482508820/unknown.png"
					testCompany.Company = "Testing"
					var price string
					var image string
					if len(value.Variants) == 0 || len(value.Images) == 0 {
						price = "N/A"
						image = testCompany.CompanyImage
					} else {
						price = value.Variants[0].Price
						image = value.Images[0].Src
					}
					go m.sendWebhook(m.Config.sku, value.Title, price, link, image, value.Variants)
				}
			}(value)
		}

	}
	// elapsed = time.Since(start)
	// fmt.Println("Time Elapsed Loop : ", elapsed)
	// fmt.Println("New List", len(newList), len(m.products))
	m.products = newList
	return nil
}

func (m *Monitor) getProxy(proxyList []string) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Site : %s, Product : %s Recovering from panic in printAllOperations error is: %v \n", m.Config.site, m.Config.sku, r)
		}
	}()
	if m.Config.proxyCount+1 == len(proxyList) {
		m.Config.proxyCount = 0
	}
	m.Config.proxyCount++
	return proxyList[m.Config.proxyCount]
}

func (m *Monitor) sendWebhook(site string, name string, price string, link string, image string, restockedVariants []Types.Variant) error {
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
	t := time.Now().UTC()


	for _, comp := range m.CurrentCompanies {
		fmt.Println(comp.Company)
		go m.webHookSend(comp, site, name, price, link, t, image, restockedVariants)
	}
	// payload := strings.NewReader("{\"content\":null,\"embeds\":[{\"title\":\"Target Monitor\",\"url\":\"https://discord.com/developers/docs/resources/channel#create-message\",\"color\":507758,\"fields\":[{\"name\":\"Product Name\",\"value\":\"%s\"},{\"name\":\"Product Availability\",\"value\":\"In Stock\\u0021\",\"inline\":true},{\"name\":\"Stock Number\",\"value\":\"%s\",\"inline\":true},{\"name\":\"Links\",\"value\":\"[Product](https://www.walmart.com/ip/prada/%s)\"}],\"footer\":{\"text\":\"Prada#4873\"},\"timestamp\":\"2021-04-01T18:40:00.000Z\",\"thumbnail\":{\"url\":\"https://cdn.discordapp.com/attachments/815507198394105867/816741454922776576/pfp.png\"}}],\"avatar_url\":\"https://cdn.discordapp.com/attachments/815507198394105867/816741454922776576/pfp.png\"}")
	return nil
}

func (m *Monitor) webHookSend(c Company, site string, name string, price string, link string, currentTime time.Time, image string, restockedVariants []Types.Variant) {
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
		Name:   "Price",
		Value:  price,
		Inline: true,
	})

	for _, variants := range restockedVariants {
		var currentName string
		var value string

		if variants.Available {
			currentName = variants.Title
			value = fmt.Sprintf("[ATC](%s/cart/%d:1)", m.BaseUrl, variants.ID)
		} else {
			// Add Inventory Quantity
			currentName = variants.Title
			value = "OOS"
		}
		currentFields = append(currentFields, &discordhook.EmbedField{
			Name:   currentName,
			Value:  value,
			Inline: true,
		})
	}
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
		url := "http://104.249.128.207:7243/DB"
		req, err := http.NewRequest("POST", url, getDBPayload)
		if err != nil {
			fmt.Println(err)
			fmt.Println(errors.Cause(err))

			return nil
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println(err)
			fmt.Println(errors.Cause(err))

			res.Body.Close()
			return nil

		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			fmt.Println(errors.Cause(err))

			res.Body.Close()
			return nil
		}
		var currentObject ItemInMonitorJson
		err = json.Unmarshal(body, &currentObject)
		if err != nil {
			fmt.Println(err)
			fmt.Println(errors.Cause(err))

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

func (m *Monitor) restockedVariants(oldVariants []Types.Variant, newVariants []Types.Variant) []Types.Variant {
	var returnArray []Types.Variant
	for _, newVariant := range newVariants {
		isPresent := false
		for _, oldVariant := range oldVariants {
			if newVariant.ID == oldVariant.ID {
				fmt.Println("Variant Compare ", newVariant.Available, newVariant.UpdatedAt, oldVariant.Available, oldVariant.UpdatedAt)
				isPresent = true
				if newVariant.Available != oldVariant.Available && newVariant.Available {
					returnArray = append(returnArray, newVariant)
				}
			}

		}

		if !isPresent && newVariant.Available {
			fmt.Println("New Variant", newVariant.Available, newVariant.CreatedAt, newVariant.ID, len(oldVariants))
			returnArray = append(returnArray, newVariant)
		}
	}

	return returnArray
}
