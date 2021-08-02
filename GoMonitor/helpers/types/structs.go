package Types

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)
type ProxyResponseType struct {
	Proxies []string `json:"proxies"`
}
type ShopifyNewProduct struct {
	Store            string     `json:"store,omitempty"`
	BodyHTML  string `json:"body_html,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	Handle    string `json:"handle,omitempty"`
	ID        int64  `json:"id,omitempty"`
	Images    []struct {
		CreatedAt  string        `json:"created_at,omitempty"`
		Height     int64         `json:"height,omitempty"`
		ID         int64         `json:"id,omitempty"`
		Position   int64         `json:"position,omitempty"`
		ProductID  int64         `json:"product_id,omitempty"`
		Src        string        `json:"src,omitempty"`
		UpdatedAt  string        `json:"updated_at,omitempty"`
		VariantIds []interface{} `json:"variant_ids,omitempty"`
		Width      int64         `json:"width,omitempty"`
	} `json:"images,omitempty"`
	Options []struct {
		Name     string   `json:"name,omitempty"`
		Position int64    `json:"position,omitempty"`
		Values   []string `json:"values,omitempty"`
	} `json:"options,omitempty"`
	ProductType string   `json:"product_type,omitempty"`
	PublishedAt string   `json:"published_at,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Title       string   `json:"title,omitempty"`
	UpdatedAt   string   `json:"updated_at,omitempty"`
	Variants    []Variant `json:"variants,omitempty"`
	Vendor string `json:"vendor,omitempty"`
}

type Variant struct {
	Available        bool        `json:"available,omitempty"`
	CompareAtPrice   string      `json:"compare_at_price,omitempty"`
	CreatedAt        string      `json:"created_at,omitempty"`
	FeaturedImage    interface{} `json:"featured_image,omitempty"`
	Grams            int64       `json:"grams,omitempty"`
	ID               int64       `json:"id,omitempty"`
	Option1          string      `json:"option1,omitempty"`
	Option2          string      `json:"option2,omitempty"`
	Option3          interface{} `json:"option3,omitempty"`
	Position         int64       `json:"position,omitempty"`
	Price            string      `json:"price,omitempty"`
	ProductID        int64       `json:"product_id,omitempty"`
	RequiresShipping bool        `json:"requires_shipping,omitempty"`
	Sku              string      `json:"sku,omitempty"`
	Taxable          bool        `json:"taxable,omitempty"`
	Title            string      `json:"title,omitempty"`
	UpdatedAt        string      `json:"updated_at,omitempty"`
}
type ShopifyProductJS struct {
	Available            bool     `json:"available,omitempty"`
	Store            string     `json:"store,omitempty"`
	CompareAtPrice       int64    `json:"compare_at_price,omitempty"`
	CompareAtPriceMax    int64    `json:"compare_at_price_max,omitempty"`
	CompareAtPriceMin    int64    `json:"compare_at_price_min,omitempty"`
	CompareAtPriceVaries bool     `json:"compare_at_price_varies,omitempty"`
	CreatedAt            string   `json:"created_at,omitempty"`
	Description          string   `json:"description,omitempty"`
	FeaturedImage        string   `json:"featured_image,omitempty"`
	Handle               string   `json:"handle,omitempty"`
	ID                   int64    `json:"id,omitempty"`
	Images               []string `json:"images,omitempty"`
	Media                []struct {
		Alt          interface{} `json:"alt,omitempty"`
		AspectRatio  float64     `json:"aspect_ratio,omitempty"`
		Height       int64       `json:"height,omitempty"`
		ID           int64       `json:"id,omitempty"`
		MediaType    string      `json:"media_type,omitempty"`
		Position     int64       `json:"position,omitempty"`
		PreviewImage struct {
			AspectRatio float64 `json:"aspect_ratio,omitempty"`
			Height      int64   `json:"height,omitempty"`
			Src         string  `json:"src,omitempty"`
			Width       int64   `json:"width,omitempty"`
		} `json:"preview_image,omitempty"`
		Src   string `json:"src,omitempty"`
		Width int64  `json:"width,omitempty"`
	} `json:"media,omitempty"`
	Options []struct {
		Name     string   `json:"name,omitempty"`
		Position int64    `json:"position,omitempty"`
		Values   []string `json:"values,omitempty"`
	} `json:"options,omitempty"`
	Price       int64    `json:"price,omitempty"`
	PriceMax    int64    `json:"price_max,omitempty"`
	PriceMin    int64    `json:"price_min,omitempty"`
	PriceVaries bool     `json:"price_varies,omitempty"`
	PublishedAt string   `json:"published_at,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Title       string   `json:"title,omitempty"`
	Type        string   `json:"type,omitempty"`
	URL         string   `json:"url,omitempty"`
	Variants    []struct {
		Available           bool        `json:"available,omitempty"`
		Barcode             string      `json:"barcode,omitempty"`
		CompareAtPrice      int64       `json:"compare_at_price,omitempty"`
		FeaturedImage       interface{} `json:"featured_image,omitempty"`
		ID                  int64       `json:"id,omitempty"`
		InventoryManagement string      `json:"inventory_management,omitempty"`
		InventoryPolicy     string      `json:"inventory_policy,omitempty"`
		InventoryQuantity   int64       `json:"inventory_quantity,omitempty"`
		Name                string      `json:"name,omitempty"`
		Option1             string      `json:"option1,omitempty"`
		Option2             string      `json:"option2,omitempty"`
		Option3             interface{} `json:"option3,omitempty"`
		Options             []string    `json:"options,omitempty"`
		Price               int64       `json:"price,omitempty"`
		PublicTitle         string      `json:"public_title,omitempty"`
		RequiresShipping    bool        `json:"requires_shipping,omitempty"`
		Sku                 string      `json:"sku,omitempty"`
		Taxable             bool        `json:"taxable,omitempty"`
		Title               string      `json:"title,omitempty"`
		Weight              int64       `json:"weight,omitempty"`
	} `json:"variants,omitempty"`
	Vendor string `json:"vendor,omitempty"`
}

type SiteCollectMongo struct {
	Store string
	Products []ShopifyNewProduct
}

type Config struct {
	Sku              string
	SkuName          string // Only for new Egg
	StartDelay       int
	Discord          string
	Site             string
	PriceRangeMax    int
	PriceRangeMin    int
	ProxyCount       int
	IndexMonitorJson int
}
type Monitor struct {
	Config              Config
	MonitorProduct      Product
	Availability        string
	AvailabilityBool        bool
	CurrentAvailability string
	CurrentAvailabilityBool bool
	Client              http.Client
	File                *os.File
	Products            []ProductsItem
	Stop                bool
	CurrentCompanies    []Company
	Collection          *mongo.Collection
	ZerosArray          []string
	Keywords            []string
	NegKeywords         []string
	BaseUrl             string
}
type Product struct {
	Name        string
	StockNumber string
	ProductId   string
	Price       int
	Image       string
	Link        string
}
type Proxy struct {
	Ip       string
	Port     string
	UserAuth string
	UserPass string
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
	Products []ShopifyNewProduct
}
type ProductsItem struct {
	ID          int64
	UpdatedTime string
	Variants    []Variant
	Store       string
	Handle      string
}

func (m *Monitor) GetProxy(proxyList []string) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Site : %s, Product : %s Recovering from panic in printAllOperations error is: %v \n", m.Config.Site, m.Config.Sku, r)
		}
	}()
	if m.Config.ProxyCount+1 == len(proxyList) {
		m.Config.ProxyCount = 0
	}
	m.Config.ProxyCount++
	return proxyList[m.Config.ProxyCount]
}

func (m *Monitor) CheckStop() error {
	for !m.Stop {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Site : %s, Product : %s  CHECK STOP Recovering from panic in printAllOperations error is: %v \n", m.Config.Site, m.Config.Sku, r)
			}
		}()
		getDBPayload := strings.NewReader(fmt.Sprintf(`{
			"site" : "%s",
			"sku" : "%s"
		  }`, strings.ToUpper(m.Config.Site), m.Config.Sku))
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
		fmt.Println(string(body))

		if strings.Contains(string(body), "Failed to ") {
			time.Sleep(2000 * (time.Millisecond))
			continue
		}

		var currentObject ItemInMonitorJson

		err = json.Unmarshal(body, &currentObject)
		
		if err != nil {
			fmt.Println(err)
			fmt.Println(errors.Cause(err))
			res.Body.Close()
			return nil
		}
		
		m.Stop = currentObject.Stop
		m.CurrentCompanies = currentObject.Companies
		fmt.Println(m.CurrentCompanies)
		res.Body.Close()
		time.Sleep(5000 * (time.Millisecond))
	}
	return nil
}