package searchcafe

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("SearchCafe", SearchCafe)
}

func SearchCafe(w http.ResponseWriter, r *http.Request) {
	if _, bol := os.LookupEnv("SEARCH_API_KEY"); bol != true {
		fmt.Fprint(w, "SEARCH_API_KE not found")
		return
	}
	var d struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		fmt.Fprint(w, err)
		return
	}
	var middle_area string
	if d.Name == "" {
		middle_area = "Y030,Y035,Y040"
	} else {
		middle_area = d.Name
	}

	req, err := http.NewRequest("GET", "https://webservice.recruit.co.jp/hotpepper/gourmet/v1/", nil)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	q := req.URL.Query()
	q.Add("key", os.Getenv("SEARCH_API_KEY"))
	q.Add("budget", "B001,B002")
	q.Add("genre", "G014")
	q.Add("middle_area", middle_area)
	q.Add("special_category", "SPG5")
	q.Add("format", "json")
	req.URL.RawQuery = q.Encode()

	fmt.Println(req.URL.String())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("Status:", resp.Status)
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	fmt.Println(string(data))
}
