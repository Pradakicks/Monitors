package FetchProxies

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)


func Get() []string {
	var proxyList = make([]string, 0)
	url := "https://monitors-9ad2c-default-rtdb.firebaseio.com/proxy.json"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
		res.Body.Close()
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
		res.Body.Close()
	}
	var proxies []string
	err = json.Unmarshal(body, &proxies)
	if err != nil {
		log.Fatal(err)
		res.Body.Close()
	}
//	fmt.Println(proxies)
	for _, proxy := range proxies {
	//	fmt.Println(proxy)
		proxyList = append(proxyList, proxy)
	}
	return proxyList
}
