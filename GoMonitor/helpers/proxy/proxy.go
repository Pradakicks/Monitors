package FetchProxies

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	Types "github.con/prada-monitors-go/helpers/types"
)

func Get() []string {
	var proxyList = make([]string, 0)
	url := "http://localhost:7243/PROXY"
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
	var proxies Types.ProxyResponseType
	err = json.Unmarshal(body, &proxies)
	if err != nil {
		log.Fatal(err)
		res.Body.Close()
	}
	//	fmt.Println(proxies)
	for _, proxy := range proxies.Proxies {
		//	fmt.Println(proxy)
		proxyList = append(proxyList, proxy)
	}
	return proxyList
}
