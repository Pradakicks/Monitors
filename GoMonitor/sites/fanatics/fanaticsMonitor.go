package FanaticsMonitor

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
	Auction struct {
		AuctionBidCount      interface{} `json:"auctionBidCount,omitempty"`
		AuctionEndTime       interface{} `json:"auctionEndTime,omitempty"`
		AuctionLastBidAmount interface{} `json:"auctionLastBidAmount,omitempty"`
		AuctionLastBidTime   interface{} `json:"auctionLastBidTime,omitempty"`
		AuctionPdpURL        interface{} `json:"auctionPdpUrl,omitempty"`
		AuctionStartingBid   interface{} `json:"auctionStartingBid,omitempty"`
	} `json:"auction,omitempty"`
	ClickBeacon       interface{}   `json:"clickBeacon,omitempty"`
	Colors            []interface{} `json:"colors,omitempty"`
	DerivedExclusions int64         `json:"derivedExclusions,omitempty"`
	FewLeft           bool          `json:"fewLeft,omitempty"`
	ImageSelector     struct {
		AdditionalImages []struct {
			Image struct {
				Alt         interface{} `json:"alt,omitempty"`
				Href        interface{} `json:"href,omitempty"`
				ScreenWidth interface{} `json:"screenWidth,omitempty"`
				Src         string      `json:"src,omitempty"`
				Title       interface{} `json:"title,omitempty"`
			} `json:"image,omitempty"`
		} `json:"additionalImages,omitempty"`
		DefaultImage struct {
			Image struct {
				Alt         interface{} `json:"alt,omitempty"`
				Href        interface{} `json:"href,omitempty"`
				ScreenWidth interface{} `json:"screenWidth,omitempty"`
				Src         string      `json:"src,omitempty"`
				Title       interface{} `json:"title,omitempty"`
			} `json:"image,omitempty"`
		} `json:"defaultImage,omitempty"`
	} `json:"imageSelector,omitempty"`
	IsDailyDeal               bool        `json:"isDailyDeal,omitempty"`
	IsGtgt                    bool        `json:"isGtgt,omitempty"`
	IsJerseyAssuranceEligible bool        `json:"isJerseyAssuranceEligible,omitempty"`
	IsShopRunnerEligible      bool        `json:"isShopRunnerEligible,omitempty"`
	IsSponsored               interface{} `json:"isSponsored,omitempty"`
	PotentialDiscount         struct {
		CalculatedPercentOff         int64 `json:"calculatedPercentOff,omitempty"`
		IsCouponEligibleWithCart     bool  `json:"isCouponEligibleWithCart,omitempty"`
		IsFreeShippingEligible       bool  `json:"isFreeShippingEligible,omitempty"`
		IsProductEligibleForDiscount bool  `json:"isProductEligibleForDiscount,omitempty"`
		SomeItemsExcluded            bool  `json:"someItemsExcluded,omitempty"`
	} `json:"potentialDiscount,omitempty"`
	Price struct {
		Clearance struct {
			Money struct {
				AdmCC             string `json:"admCC,omitempty"`
				AdmCV             string `json:"admCV,omitempty"`
				UserCC            string `json:"userCC,omitempty"`
				UserCurrencyValue string `json:"userCurrencyValue,omitempty"`
			} `json:"money,omitempty"`
		} `json:"clearance,omitempty"`
		DiscountAmount interface{} `json:"discountAmount,omitempty"`
		DiscountPrice  interface{} `json:"discountPrice,omitempty"`
		MarkdownAmount struct {
			Money struct {
				AdmCC             string `json:"admCC,omitempty"`
				AdmCV             string `json:"admCV,omitempty"`
				UserCC            string `json:"userCC,omitempty"`
				UserCurrencyValue string `json:"userCurrencyValue,omitempty"`
			} `json:"money,omitempty"`
		} `json:"markdownAmount,omitempty"`
		QualifierPromo struct {
			DiscountAmount      interface{} `json:"discountAmount,omitempty"`
			ItemDiscountedTotal interface{} `json:"itemDiscountedTotal,omitempty"`
		} `json:"qualifierPromo,omitempty"`
		Regular struct {
			Money struct {
				AdmCC             string `json:"admCC,omitempty"`
				AdmCV             string `json:"admCV,omitempty"`
				UserCC            string `json:"userCC,omitempty"`
				UserCurrencyValue string `json:"userCurrencyValue,omitempty"`
			} `json:"money,omitempty"`
		} `json:"regular,omitempty"`
		Sale struct {
			Money struct {
				AdmCC             string `json:"admCC,omitempty"`
				AdmCV             string `json:"admCV,omitempty"`
				UserCC            string `json:"userCC,omitempty"`
				UserCurrencyValue string `json:"userCurrencyValue,omitempty"`
			} `json:"money,omitempty"`
		} `json:"sale,omitempty"`
	} `json:"price,omitempty"`
	ProductID     string      `json:"productId,omitempty"`
	ProductTags   []string    `json:"productTags,omitempty"`
	ReadyToShip   bool        `json:"readyToShip,omitempty"`
	ShipDetailsID string      `json:"shipDetailsId,omitempty"`
	Title         string      `json:"title,omitempty"`
	TopSeller     interface{} `json:"topSeller,omitempty"`
	URL           string      `json:"url,omitempty"`
	ViewBeacon    interface{} `json:"viewBeacon,omitempty"`
}

type FanticsSearchResponse struct {
	Data struct {
		Search struct {
			Products          []Product `json:"products,omitempty"`
			TotalProductCount int64     `json:"totalProductCount,omitempty"`
		} `json:"search,omitempty"`
	} `json:"data,omitempty"`
}

type ProductIdentifier struct {
	ProductID string `json:"productId,omitempty"`
	URL string `json:"url,omitempty"`
}

var FanaticsCollection = helper.ConnectDBFanatics()

func NewMonitor(sku string, collection *mongo.Collection) *CurrentMonitor {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Site : %s, Product : %s Recovering from panic in printAllOperations error is: %v \n", sku, sku, r)
		}
	}()
	fmt.Println("New Fanatics Monitor ", sku)
	m := CurrentMonitor{}
	// m.Monitor.Availability = "Copy"
	m.Monitor.Config.Site = "FanaticsNewProducts"
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
					testingArr = append(testingArr, ProductIdentifier{ProductID: prod.ProductID})
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

	url := "https://www.fanatics.com/graphql"

		payload := strings.NewReader("{\"query\":\"\\n  query {\\n    search (\\n        queryString: \\\"\\\",\\n        pageSize: 25,\\n        pageNumber: 1,\\n        sortBy: \\\"NewestArrivals\\\",\\n        filters: [{key: \\\"o\\\", value: \\\"25\\\"},{key: \\\"d\\\", value: \\\"8004\\\"},{key: \\\"d\\\", value: \\\"8073\\\"}],\\n        siteId: 510005,\\n        coupon: {\\n    isBMSM: false,\\n    isPercentOff: false,\\n    percentOff: \\\"0\\\",\\n    threshold: \\\"48\\\",\\n    shippingDiscountThreshold : \\\"48\\\",\\n    orderDiscountThreshold: \\\"48\\\",\\n    cartTotal: \\\"0\\\"\\n  },\\n        dataServiceSettings: {isSmartExclusionsEligibleForDiscount:true,fewLeftThreshold:20,showFreeShippingforOversizeItems:false,isLastPurchasedEnabled:false}\\n      ) {\\n      breadcrumb {\\n   name\\n   key\\n   count\\n   link {\\n     href\\n   }\\n}\\nproducts {\\n  title\\n  productId\\n  url\\n  topSeller {\\n  type\\n  value\\n  valueUrl {\\n    href\\n  }\\n}\\n  fewLeft\\n  shipDetailsId\\n  isShopRunnerEligible\\n  isJerseyAssuranceEligible\\n  isGtgt\\n  price {\\n  regular {\\n    money {\\n  userCC\\n  userCurrencyValue\\n  admCC\\n  admCV\\n}\\n  }\\n  sale {\\n    money {\\n  userCC\\n  userCurrencyValue\\n  admCC\\n  admCV\\n}\\n  }\\n  clearance {\\n    money {\\n  userCC\\n  userCurrencyValue\\n  admCC\\n  admCV\\n}\\n  }\\n  markdownAmount {\\n    money {\\n  userCC\\n  userCurrencyValue\\n  admCC\\n  admCV\\n}\\n  }\\n  discountPrice {\\n    money {\\n  userCC\\n  userCurrencyValue\\n  admCC\\n  admCV\\n}\\n  }\\n  discountAmount {\\n    money {\\n  userCC\\n  userCurrencyValue\\n  admCC\\n  admCV\\n}\\n  }\\n  qualifierPromo {\\n    itemDiscountedTotal {\\n      money {\\n  userCC\\n  userCurrencyValue\\n  admCC\\n  admCV\\n}\\n    }\\n    discountAmount {\\n      money {\\n  userCC\\n  userCurrencyValue\\n  admCC\\n  admCV\\n}\\n    }\\n  }\\n}\\n  potentialDiscount {\\n  isCouponEligibleWithCart\\n  isFreeShippingEligible\\n  isProductEligibleForDiscount\\n  someItemsExcluded\\n  calculatedPercentOff\\n}\\n  imageSelector {\\n  defaultImage {\\n    image {\\n  src\\n  href\\n  alt\\n  title\\n  screenWidth\\n}\\n  }\\n  additionalImages {\\n    image {\\n  src\\n  href\\n  alt\\n  title\\n  screenWidth\\n}\\n  }\\n}\\n  colors {\\n    productId\\n    productUrl\\n    code\\n    name\\n    image {\\n      src\\n    }\\n  }\\n  isSponsored\\n  viewBeacon\\n  clickBeacon\\n  productTags\\n  auction {\\n    auctionEndTime\\n    auctionPdpUrl\\n    auctionStartingBid\\n    auctionLastBidTime\\n    auctionLastBidAmount\\n    auctionBidCount\\n  }\\n  derivedExclusions\\n  isDailyDeal\\n  readyToShip\\n}\\ntotalProductCount\\nseoContext {\\n  canonicalUrl\\n  paginationLinks {\\n    rel\\n    url\\n  }\\n  shouldIndex\\n  shouldFollow\\n}\\nleftNavigation {\\n  yourSelections {\\n  name\\n  key\\n  valueId\\n  link {\\n    title\\n    href\\n    alias\\n    text\\n    type\\n  }\\n  removeLink\\n}\\n  facets {\\n  title\\n  type\\n  isCollapsed\\n  hasPills\\n  links {\\n    alias\\n    type\\n    title\\n    href\\n    text,\\n    isFeatured,\\n    isPill,\\n    productCount\\n  }\\n}\\n}\\n    }\\n  }\\n\",\"variables\":{}}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("cookie", `platform1=e; u_loc=en-US; st=510005; ac=USD; uc=USD; priv=%7B%22acc%22%3Afalse%2C%22fcc%22%3Afalse%2C%22tcc%22%3Afalse%2C%22pc%22%3Atrue%2C%22ecc%22%3Afalse%7D; cme=; _s=www.fanatics.com; vid=6eba57c0-eefb-11eb-be63-11440f8736b1; akacd_PR_Iris_permanent=3804857721~rv=21~id=2efca3b5999196d3fb76a0af4899f3e5; akacd_PR_Iris_Assets=2177452799~rv=40~id=324c6c0fd5e7bc9fe527b7bbc007e435; bm_sz=7CF9E379A3A8C9FFE18230A96ECF4074~YAAQxGfJF2BF1eR6AQAAd3Lk6AzXdise0zyhaE15DODZVCyk6ONjwHjSw2w3dL0nmcHf5mCn/DPsf9Gvjzx0YSCf38T2ebZlZ+BII2wJqTqBQcwcYwQM2PuWbQpBFOVg+X+9VCy1Vevmcey3EgM17PZME+Ar2HFZ40XZaxnYPlIMAZoWex1uoQtUvDExmKONtdSzcPGScEfOXV0TmNp8DjFRVz0KexYlaK7bg9pORDhtizXwKxyZDjdnm9duftcDDCi+BE1u7n2r2Ejd5tBcgFOj1RrPPKT5pldSDWwFUFS1EAc0Wg==~3552050~4469059; civ=1.1.0-rc-20210721.38758; pu=true; s_fid=5B6BE393B532FB64-38A9CF7CE8FFA017; s_cc=true; s_fuid=65876091034991972844083022672580026391; _gcl_au=1.1.2087603185.1627404924; sr_browser_id=16578ab5-b496-44be-b267-657e71909242; _fbp=fb.1.1627404924191.1702496856; _hjid=7071c4e8-90be-4756-a406-8e35aa1b9b4a; _abck=A83C57B2681EB7822551F166212FF1C8~0~YAAQxGfJF79F1eR6AQAABnnk6AZApaMRNn5YYSNoSegvENduAk2TXFS+o3I2e/9ZCdiFYW+o5QdSXm88ISBFDLafwLl0TC5+Aj9dLqt9iM+ibHvKIZ0JsJa1jFNIoq2Im5658mPghnaIZyihut9Gvyulw7R5pscHVmvUQ1323ggwiujm1SAUNojiYoT74bIBHJcJywFRfTeSwGCynaZKcdl2mTeWJH+ZQWHBVc7canXVBBNquNf8MZYCgdkBqm7udP6uITUwqjpaJQlSbnmgw+LSsWIIl0bMVaFHkwoJeamssZv+4jXCYBpJsd5XPOTLrjTtXgkIs/IA3CAbccgwmWFIUsAV3RzmwlfU8GIhWjE8r+w0k+XwrGtmd6L7rRSCynAnsyyWLGWZAxfmfkH5LHME+CETxLVmChY=~-1~||-1||~-1; _pbjs_userid_consent_data=3524755945110770; _pubcid=9d70ff0b-8db3-41cb-bd7d-6f40e9762e2b; _fssid=c73732d0-5cf3-4130-92e4-6b0ab7f014a0; __qca=P0-1916065638-1627404925843; pbjs-unifiedid=%7B%22TDID%22%3A%22fc92c060-2891-42aa-9e41-a867087f0aa8%22%2C%22TDID_LOOKUP%22%3A%22FALSE%22%2C%22TDID_CREATED_AT%22%3A%222021-07-27T16%3A55%3A27%22%7D; panoramaId_expiry=1628009727512; _cc_id=a9477b8aecc517906c56aa2f99e3e934; panoramaId=674d68094828dd8383b9ec9e86db16d53938c7de1104d398604218b98096a7df; __gads=ID=f037743e755b39dd:T=1627404926:S=ALNI_MYWnQgjV4bGw8MxPoyN_rtUrGeEnw; csl=/nfl/collectibles-and-memorabilia-trading-cards/o-2461+d-08339937-64441858+z-9-1762906313?pageSize=72&sortOption=NewestArrivals; cqe=%5B%225057%3AA%3A0%3A0%22%2C%225447%3AB%3A1%3A1%22%2C%225662%3AB%3A1%3A1%22%2C%225747%3AB%3A1%3A1%22%2C%225778%3AA%3A0%3A1%22%5D; ist=12c2a794-8918-45f1-9c69-ed206a7f9040; _fsloc=?i=US&c=Boynton Beach&s=FL; sa=sid%3D12c2a794-8918-45f1-9c69-ed206a7f9040%7Cfpr%3Dwww.fanatics.com; va=%7B%22cc%22%3A0%2C%22ct%22%3A0%2C%22cpi%22%3A%5B%5D%2C%22nv%22%3Afalse%2C%22el%22%3Afalse%2C%22ch%22%3A%22%22%2C%22ci%22%3A%22rsh%22%2C%22lic%22%3A%22rsh%22%7D; xsrft=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyZXFIb3N0Ijoid3d3LmZhbmF0aWNzLmNvbSIsInZpc2l0b3JJZCI6IjZlYmE1N2MwLWVlZmItMTFlYi1iZTYzLTExNDQwZjg3MzZiMSIsImlhdCI6MTYyNzQxMzI3NSwiZXhwIjoxNjI5MDQwMzg4ODkxfQ.E8pJQrRg6OPmLAFCtjfFj1P0DtPMVkNFF5wNRo_wKHM; xsrfp=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyZXFIb3N0Ijoid3d3LmZhbmF0aWNzLmNvbSIsInZpc2l0b3JJZCI6IjZlYmE1N2MwLWVlZmItMTFlYi1iZTYzLTExNDQwZjg3MzZiMSIsImlhdCI6MTYyNzQxMzI3NSwiZXhwIjoxNjI5MDQwMzg4ODkxfQ.0a9wnXiPxMePFkNgqyMjFOm-v5mRSSSyaR36oN9MM18; eci=111c4a6ba3347efd; AWSALB=kGkDutIRmiuOpadSLlSwYqwbcZQPcIXeCuYP9lu++lgwkJjdslz8mgtjEdVXEnmgwdm1RiW9HYklgnOtSPFMV7PFBcjaeATcvM289dpTISiajGpNIP/GdavS/bPgmK6yIW2dOJVxRR1CatANXBdUtE7G3atVurfDow68dTSlBpENDb4XVD4GcKiqZU/IAcgiOSuG4XTf6Pht+SqO4xUzLJe7v4BPUxscaOhI8O8yoDFACtePiR9YvTNLBV3h8CA=; AWSALBCORS=kGkDutIRmiuOpadSLlSwYqwbcZQPcIXeCuYP9lu++lgwkJjdslz8mgtjEdVXEnmgwdm1RiW9HYklgnOtSPFMV7PFBcjaeATcvM289dpTISiajGpNIP/GdavS/bPgmK6yIW2dOJVxRR1CatANXBdUtE7G3atVurfDow68dTSlBpENDb4XVD4GcKiqZU/IAcgiOSuG4XTf6Pht+SqO4xUzLJe7v4BPUxscaOhI8O8yoDFACtePiR9YvTNLBV3h8CA=; vrc=74be2097a84a9f87; ak_bmsc=96AFB889B8E5FFF4FFFAEF1E7E6156D4~000000000000000000000000000000~YAAQdOVwaBv6npl6AQAANfBj6QyJWTS3jNVO9axYVkbrKbGfaPli5nr/cCfGrT0r3AZaHwAAfUqrSfqlxZSDUUfoYLz3CaqzpTNQ/X0Jp4588WPHQ3haK6A2YTzBhHSN1c+9vfl9c6XocWKML5x+tFHnozpRmkoN7aXOJIr42U3DUUs/qkZdwxSiViZgBNUzLhynsxCZj6gHrRwhF5nyw+sXyvJE61i+jK+QRHGaLWb6auxxvreoujdYPo+l/cJSKzl0NLv1ExFllNbdvrS1u6NgxB0mDe6I6ZcHI57mnbl0Zs7jpV2jNLKNdrI8Jhz+98kfurimTwAPC/6xVxbr9+uy5XBEVNNvO1s9mV0emrCaElF5257/YgKjGRAi3KI28oFh5cEmNKFM8yUKXrYEY0f/jFMjexQBHIVuTCPxwVN4aJor7oX5PODV79VTk6+G3egbWJ7QFd9nDUnYe0BgoV5qaVulJpvQsxYQOZ8dnpav4/Q1YFPqWerMPQ==; cto_bidid=eARrP19DWXIyakFDWU0yQXZLMzJZZEhVVlNqclNJU3A5YTVVUGJraXAyMXpMZFBZOWlCVCUyRkJCd2JQanU5ejZVSjB1aVc1Sm83MDM1RUVQaDNFbUlRS0phJTJGSiUyRlhoSFpBQ3Z2Q3VPcmZ5eEREamVBcyUzRA; cto_bundle=x1Y2nV9aQU56S2F1ZUY5WGdPd2I5SmN0SmpJWXlMJTJGQnRkZERTUjRqRGVhVTBVSzM2M2pVaU5QQWxtOSUyRnEyJTJCNFdqVzJmT0FvNE5YJTJCVWlWemkxdFRUSDRYWDh2RDc1eWJadmRxQ29oU25vcXlNY0swNlVxczFlbkdnS1JwcFpGeUNZZlVVWFpLJTJCZEczWnRxMCUyRiUyQnNIY3BoJTJCUEdnJTNEJTNE; _4c_=lVTBbts4EP2VQodcGlmiRElkAKOwXRubRROncbrXgJLGNlGJFEjKqlPk35e0JNvbXro6CPPeDDnzhhz%2B9Lo9CO8OpVGGURyRMI7jW%2B87HLV399MrGvc%2FuF%2BrKu%2FO2xvT6Lsg6LpusmWCGV7oSSHrQGyroJBVBYXheQXaZ6L0a6ilYjmvOPONYiUXO79gqtSB9COcoo%2BlH5I4pjTO%2FBRjjEhCPr751EdZGtEwjVH8qWE72PA3mGbRjZbKrBvDpZg%2BQgfazJTiB1Zp79YrZAm2QEQn2SSx2LxZZLOE1m6ULNvCvJpj42I6yD%2Fo8rt1lHDgBbx2vDR7tzgl4YXdA9%2FtjaPD5EQ3yoLIWh0XpezOywihF%2FK8imbIsl%2BA7VpYqxdg9VOrij3TUFqvda1VCer%2B8xX4h1Ut9PjJyn7py3VoM8Y9g2mVsI1ctNrIGtT1XgbqhWyF6bnHts57UnOxYMrRj5Z35pjJ4aVSUj2A1jblQP0ltRGsHuG9MKAEqzbArIIXUPXg6NWNQO52UN6LAfYllQNyep5hC0qdSh6pb89fzuhYgzAPYPbyvKhiYt2a2XzEStbSHf9abLgZ8%2BrXLS9XXGnzZJWdGuXojVVkQy94z5vGNc51SB1%2FYf%2BT1x3WYH5rNFRVf%2Bk%2Bc91U7AhjWK5kp09yFntbGXyg7nJIOzneAyusqS6CXZ19yddjM7B23H53rB9enl%2Fny9li%2FXg1eOcghKJJJLOJABPkgdZnj729AQr%2B3vhokgY6ySglmNIwS9Ms%2BzT7Op%2BiG9uvaTJP58uYxvMkjlbzFPsxmdHFKlssyWo1C1F2M%2Fu6nLor3BzcDDirkgWrnAb7ZNjK3axLcVW5937r%2FRgeE0rTjJAU2Uk09uUgKQ7dZyMUL4dXxSuiJEQ0TnyM49zHOQOfRXjrR4ATVqS2%2BND18bRnnCAS0ghnGbWbNNWwB7qktDtFmBA8pET4nLI5jNHoTyrsD3ZU9T8Xa23GeBcYYuqKSn5RYRkbe%2BBjJ8rt1jlinxS09HHKEp9leeTbOBymDIqCpN6V1CTMSBbGo1TSZ39%2F%2Fxc%3D; RT="z=1&dm=www.fanatics.com&si=e7b8a77b-fa9d-4bbc-9b08-5b2f56535250&ss=krmfu8yi&sl=2&tt=6or&bcn=%2F%2F17d09917.akstat.io%2F&ld=bh6"; s_sq=fanaticsdev%3D%2526c.%2526a.%2526activitymap.%2526page%253Ddlp%25253A%252520NFL%252520Trading%252520Cards%25252C%252520NFL%252520Trading%252520Card%252520Set%252520%25257C%252520Fanatics%2526link%253DNewest%252520Items%2526region%253DBODY%2526pageIDType%253D1%2526.activitymap%2526.a%2526.c%2526pid%253Ddlp%25253A%252520NFL%252520Trading%252520Cards%25252C%252520NFL%252520Trading%252520Card%252520Set%252520%25257C%252520Fanatics%2526pidt%253D1%2526oid%253DfunctionLr%252528%252529%25257B%25257D%2526oidt%253D2%2526ot%253DLI`)
	req.Header.Add("authority", "www.fanatics.com")
	req.Header.Add("pragma", "no-cache")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("sec-ch-ua", `"Chromium";v="92", " Not A;Brand";v="99", "Google Chrome";v="92"`)
	req.Header.Add("dnt", "1")
	req.Header.Add("x-frg-tq", "0.0,505700:A,544700:B,566200:B,574700:B,577800:A")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36")
	req.Header.Add("x-xsrf-token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyZXFIb3N0Ijoid3d3LmZhbmF0aWNzLmNvbSIsInZpc2l0b3JJZCI6IjZlYmE1N2MwLWVlZmItMTFlYi1iZTYzLTExNDQwZjg3MzZiMSIsImlhdCI6MTYyNzQxMzI3NSwiZXhwIjoxNjI5MDQwMzg4ODkxfQ.E8pJQrRg6OPmLAFCtjfFj1P0DtPMVkNFF5wNRo_wKHM")
	req.Header.Add("x-frg-ss", "12c2a794-8918-45f1-9c69-ed206a7f9040")
	req.Header.Add("x-sec-clge-req-type", "ajax")
	req.Header.Add("accept", "application/json")
	req.Header.Add("x-frg-st", "510005")
	req.Header.Add("content-type", "application/json;charset=UTF-8")
	req.Header.Add("x-frg-promo", "49SHIP")
	req.Header.Add("x-frg-ci", "rsh")
	req.Header.Add("x-frg-si", "6eba57c0-eefb-11eb-be63-11440f8736b1")
	req.Header.Add("origin", "https://www.fanatics.com")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("referer", "https://www.fanatics.com/nfl/collectibles-and-memorabilia-trading-cards/o-2461+d-08339937-64441858+z-9-1762906313?pageSize=72&sortOption=NewestArrivals")
	req.Header.Add("accept-language", "en-US,en;q=0.9")
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
		fmt.Printf("Fanatics New Products %s - Code : %d Milli elapsed: %v\n", m.Monitor.Config.Sku, res.StatusCode, watch.Milliseconds())
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

	var jsonResponse FanticsSearchResponse

	err = json.Unmarshal([]byte(body), &jsonResponse)
	if err != nil {
		fmt.Println(err)
		fmt.Println(errors.Cause(err))
		return nil
	}

	for _, currentProduct := range jsonResponse.Data.Search.Products {
		isPresent := false

		for _, identifer := range m.products {

			// if identifer.URL == "nfl/green-bay-packers/bart-starr-green-bay-packers-1960-topps-number-51-card-bvg-55/o-1372+t-03820298+p-71703716874+z-9-1863514192" && currentProduct.URL == "nfl/green-bay-packers/bart-starr-green-bay-packers-1960-topps-number-51-card-bvg-55/o-1372+t-03820298+p-71703716874+z-9-1863514192" {
			// 	fmt.Println("HERER ", identifer, currentProduct.ProductID, currentProduct.URL)
			// }

			if identifer.ProductID == currentProduct.ProductID {
				isPresent = true
			}
		}

		if !isPresent {
			fmt.Println("New Item in Stock", currentProduct.Title)
			fmt.Println(isPresent, currentProduct.Title, currentProduct.ProductID)
			m.products = append(m.products, ProductIdentifier{ProductID: currentProduct.ProductID, URL: currentProduct.URL})
			// link := fmt.Sprintf("https://www.fanatics.com/%s", currentProduct.URL)
			// image := fmt.Sprintf("https://fanatics.frgimages.com/%s", currentProduct.ImageSelector.DefaultImage.Image.Src)
			// go m.sendWebhook(currentProduct.ProductID, currentProduct.Title, currentProduct.Price.Regular.Money.AdmCV, link, image)
			result, err := m.Monitor.Collection.InsertOne(context.TODO(), currentProduct)
			if err != nil {
				fmt.Println(err)
				fmt.Println(errors.Cause(err))

			} else {
				fmt.Println("New Fanatic Product ", result)
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

	for _, comp := range m.Monitor.CurrentCompanies {
		fmt.Println(comp.Company)
		go m.webHookSend(comp, sku, m.Monitor.Config.Site, name, price, link, t, image)
	}
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
		Value:  "$" + price,
		Inline: true,
	})
	currentFields = append(currentFields, &discordhook.EmbedField{
		Name:  "Links",
		Value: "[Cart](https://www.fanatics.com/cart) | [Login](https://www.fanatics.com/login?nextPathname=/account)",
	})

	var discordParams discordhook.WebhookExecuteParams = discordhook.WebhookExecuteParams{
		Content: "",
		Embeds: []*discordhook.Embed{
			{
				Title:  "Fanatics New Products" + " Monitor",
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

func FanaticsNewProducts(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)
	w.Header().Set("Content-Type", "application/json")
	var currentMonitor Types.MonitorResponse
	_ = json.NewDecoder(r.Body).Decode(&currentMonitor)
	fmt.Println(currentMonitor)
	go NewMonitor(currentMonitor.Sku, FanaticsCollection)
	json.NewEncoder(w).Encode(currentMonitor)
}