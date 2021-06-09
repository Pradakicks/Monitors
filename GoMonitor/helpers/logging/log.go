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
	Error error
}

func LogError(site string, sku string, Error error) {

	url := fmt.Sprintf("https://monitors-9ad2c-default-rtdb.firebaseio.com/Errors/%s/%s.json", site, sku)

	itemErr := itemLogger{
		site:  site,
		sku:   sku,
		Error: Error,
	}

	payload, err := json.Marshal(itemErr)
	if err != nil {
		fmt.Println(err)
	}

	req, _ := http.NewRequest("POST", url, strings.NewReader(string(payload)))

	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Connection", "close")
	req.Close = true

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	message := fmt.Sprintf("Logger For %s Sku %s : Status : %d", site, sku, res.StatusCode)
	fmt.Println(message)

}
