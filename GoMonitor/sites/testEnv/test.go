package main

import (
	"fmt"
	"strings"
	"net/http"
	"io/ioutil"
)

func main() {

	url := "https://www.homedepot.com/product-information/model?opname=mediaPriceInventory"

	payload := strings.NewReader("{\"operationName\":\"mediaPriceInventory\",\"variables\":{\"excludeInventory\":false,\"itemIds\":[\"312422229\"],\"storeId\":\"224\"},\"query\":\"query mediaPriceInventory($excludeInventory: Boolean = false, $itemIds: [String!]!, $storeId: String!) {\\n  mediaPriceInventory(itemIds: $itemIds, storeId: $storeId) {\\n    productDetailsList {\\n      itemId\\n      imageLocation\\n      onlineInventory @skip(if: $excludeInventory) {\\n        enableItem\\n        totalQuantity\\n        __typename\\n      }\\n      pricing {\\n        value\\n        original\\n        message\\n        mapAboveOriginalPrice\\n        __typename\\n      }\\n      storeInventory @skip(if: $excludeInventory) {\\n        enableItem\\n        totalQuantity\\n        __typename\\n      }\\n      __typename\\n    }\\n    __typename\\n  }\\n}\\n\"}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:89.0) Gecko/20100101 Firefox/89.0")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("X-Experience-Name", "general-merchandise")
	req.Header.Add("apollographql-client-name", "general-merchandise")
	req.Header.Add("apollographql-client-version", "0.0.0")
	req.Header.Add("x-debug", "false")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Pragma", "no-cache")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("TE", "Trailers")

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