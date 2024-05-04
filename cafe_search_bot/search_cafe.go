package searchcafe

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

type ResultsApi struct {
	Results Result `json:"results"`
}
type Result struct {
	Total int    `json:"results_available"`
	Shop  []Shop `json:"shop"`
}
type Shop struct {
	Id      string  `json:"id"`
	Name    string  `json:"name"`
	Address string  `json:"address"`
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
	Wifi    string  `json:"wifi"`
	Open    string  `json:"open"`
	Urls    struct {
		Pc string `json:"pc"`
	} `json:"urls"`
	Photo struct {
		Pc struct {
			Small string `json:"s"`
		} `json:"pc"`
	} `json:"photo"`
}

func init() {
	functions.HTTP("SearchCafe", SearchCafe)
}

func SearchCafe(w http.ResponseWriter, r *http.Request) {
	channelSecret := os.Getenv("LINE_CHANNEL_SECRET")
	bot, err := messaging_api.NewMessagingApiAPI(
		os.Getenv("LINE_CHANNEL_TOKEN"),
	)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	cb, err := webhook.ParseRequest(channelSecret, r)
	if err != nil {
		log.Printf("Cannot parse request: %+v\n", err)
		if errors.Is(err, webhook.ErrInvalidSignature) {
			w.WriteHeader(400)
			fmt.Fprint(w, err)
		} else {
			w.WriteHeader(500)
			fmt.Fprint(w, err)
		}
		return
	}

	log.Println("Handling events...")

	for _, event := range cb.Events {
		log.Printf("/callback called%+v...\n", event)

		switch e := event.(type) {
		case webhook.MessageEvent:
			switch message := e.Message.(type) {
			case webhook.TextMessageContent:
				res, err := innerLogic(message.Text)
				if err != nil {
					w.WriteHeader(500)
					fmt.Fprint(w, err)
				}
				if _, err = bot.ReplyMessage(
					&messaging_api.ReplyMessageRequest{
						ReplyToken: e.ReplyToken,
						Messages: []messaging_api.MessageInterface{
							messaging_api.TextMessage{
								Text: fmt.Sprintf("%#v", res),
							},
						},
					},
				); err != nil {
					log.Print(err)
				} else {
					log.Println("Sent text reply.")
				}
			default:
				log.Printf("Unsupported message content: %T\n", e.Message)
			}
		default:
			log.Printf("Unsupported message: %T\n", event)
		}
	}

}

func innerLogic(middle_area string) (ResultsApi, error) {
	var results ResultsApi
	if middle_area == "" {
		middle_area = "Y030,Y035,Y040"
	}

	req, err := http.NewRequest("GET", "https://webservice.recruit.co.jp/hotpepper/gourmet/v1/", nil)
	if err != nil {
		return results, err
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
		return results, err
	}
	defer resp.Body.Close()
	fmt.Println("Status:", resp.Status)
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return results, err
	}
	fmt.Println(string(data))
	err = json.Unmarshal(data, &results)
	if err != nil {
		return results, err
	}
	return results, nil
}
