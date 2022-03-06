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
type KeyWordMonitor struct {
	Endpoint string   `json:"endpoint"`
	Keywords []string `json:"keywords"`
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
	Config                  Config
	MonitorProduct          Product
	Availability            string
	AvailabilityBool        bool
	CurrentAvailability     string
	CurrentAvailabilityBool bool
	Client                  http.Client
	File                    *os.File
	Products                []ProductsItem
	Stop                    bool
	CurrentCompanies        []Company
	Collection              *mongo.Collection
	ZerosArray              []string
	Keywords                []string
	NegKeywords             []string
	BaseUrl                 string
}
type Product struct {
	Name           string
	StockNumber    string
	StockNumberInt int
	ProductId      string
	Price          int
	Image          string
	Link           string
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

type MonitorResponse struct {
	Site          string `json:"site"`
	Sku           string `json:"sku"`
	PriceRangeMin int    `json:"priceRangeMin"`
	PriceRangeMax int    `json:"priceRangeMax"`
	SkuName       string `json:"skuName"`
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
		url := "http://localhost:7243/DB"
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
