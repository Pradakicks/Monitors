package Webhook

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/andersfylling/snowflake"
	"github.com/nickname32/discordhook"
	"github.com/pkg/errors"
)

type ErrorWebhook struct {
	Global     bool   `json:"global"`
	Message    string `json:"message"`
	RetryAfter int64  `json:"retry_after"`
}
func SendWebhook(webhook string, webhookData *discordhook.WebhookExecuteParams){

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Send Webhook Class : Webhook In Use : %s ,Recovering from panic in printAllOperations error is: %v \n", webhook, r)
		}
	}()

	s := strings.Split(webhook, "/")
	// s[4] =
	fmt.Println(s)
	Int, err := strconv.Atoi(s[5])
	if err != nil {
		fmt.Println(err)
		fmt.Println(errors.Cause(err))

		return
	}
	wa, err := discordhook.NewWebhookAPI(snowflake.NewSnowflake(uint64(Int)), s[6], true, nil)
	if err != nil {
		fmt.Println(err)
		fmt.Println(errors.Cause(err))

		return
	}


	msg, err := wa.Execute(context.TODO(), webhookData, nil, "")
	if err != nil {
		fmt.Println("Error with webhook", err)
		if strings.Contains(err.Error(), "You are being") {
			if strings.Contains(err.Error(), "You are being blocked from accessing our API temporarily due to exceeding our rate limits frequently"){
				time.Sleep(240 * time.Second)
			}
			var errorWebhookMessage ErrorWebhook
			err = json.Unmarshal([]byte(err.Error()), &errorWebhookMessage)
			if err != nil {
				fmt.Println(err)
				fmt.Println(errors.Cause(err))
				return
			}
			if errorWebhookMessage.RetryAfter > 500 {
				errorWebhookMessage.RetryAfter = errorWebhookMessage.RetryAfter / 10
			}
			fmt.Println("Retrying IN : " + string(errorWebhookMessage.RetryAfter))
			time.Sleep(time.Duration(errorWebhookMessage.RetryAfter) * time.Second)
			go SendWebhook(webhook, webhookData)
		} else {
			// fmt.Println(link, Color, site, image, c.CompanyImage)
			// for _, v := range currentFields {
			// 	fmt.Println(v.Name, v.Value)
			// }

		}
		return
	}

	fmt.Println("Discord Webhook Client : ", msg.ID)
}