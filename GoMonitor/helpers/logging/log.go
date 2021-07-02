package MonitorLogger

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type itemLogger struct {
	site  string
	sku   string
	Error string
}

func LogError(site string, sku string, Error error) {

	url := fmt.Sprintf("https://monitors-9ad2c-default-rtdb.firebaseio.com/Errors/%s/%s.json", site, sku)

	itemErr := itemLogger{
		site:  site,
		sku:   sku,
		Error: Error.Error(),
	}

	payload, err := json.Marshal(itemErr)
	if err != nil {
		fmt.Println(err)
		return
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(payload)))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Connection", "close")
	req.Close = true

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	defer func(){
	message := fmt.Sprintf("Logger For %s Sku %s : Status : %d", site, sku, res.StatusCode)
		fmt.Println(message)
	}()
	

}
