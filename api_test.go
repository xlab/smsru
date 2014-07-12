package smsru

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testToken = "ddc5b4d1-4563-2344-658c-63ba15f677a2"

var simpleSms = &Sms{
	To:   "79269965690",
	Text: "Привет, Xlab!",
	Test: true,
}

var multiSms = &Sms{
	Multiple: map[string]string{
		"79269965690": "Земля",
		"79853329310": "Воздух",
	},
	Test: true,
}

func TestSmsSend(t *testing.T) {
	api := NewApi(testToken)
	ids, err := api.SmsSend(simpleSms)
	if !assert.NoError(t, err) {
		return
	}
	log.Printf("Sent sms: %s", ids)
	log.Printf("Balance: %f", api.Balance)
	ids, err = api.SmsSend(multiSms)
	if !assert.NoError(t, err) {
		return
	}
	log.Printf("Sent multiple sms: %s", ids)
	log.Printf("Balance: %f", api.Balance)
}
