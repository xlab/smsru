package smsru

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testToken = "ddc5b4d1-4563-2344-658c-63ba15f677a2"
const devPhone = "79269965690" // your dev phone, the cost should be 0.

var simpleSms = &Sms{
	To:   devPhone,
	Text: "Привет, Xlab!",
	Test: true,
}

var multiSms = &Sms{
	Multiple: map[string]string{
		devPhone:      "Земля",
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

func TestSmsCost(t *testing.T) {
	api := NewApi(testToken)
	cost, n, err := api.SmsCost(simpleSms)
	assert.NoError(t, err)
	assert.Condition(t, func() bool { return cost == 0 }, "cost of messages should be 0")
	assert.Condition(t, func() bool { return n == 1 }, "n of messages should be 1")

	// API bug on the serverside
	// cost, n, err = api.SmsCost(multiSms)
	// assert.NoError(t, err)
	// assert.Condition(t, func() bool { return cost > 1 }, "cost of messages should be > 0")
	// assert.Condition(t, func() bool { return n == 2 }, "n of messages should be 2")
}

func TestMyBalance(t *testing.T) {
	api := NewApi(testToken)
	balance, err := api.MyBalance()
	assert.NoError(t, err)
	assert.Condition(t, func() bool { return balance > 0 }, "balance should be > 0")
}

func TestMyLimit(t *testing.T) {
	api := NewApi(testToken)
	limit, _, err := api.MyLimit()
	assert.NoError(t, err)
	assert.Condition(t, func() bool { return limit > 0 }, "limit should be > 0")
}

func TestMySenders(t *testing.T) {
	api := NewApi(testToken)
	senders, err := api.MySenders()
	assert.NoError(t, err)
	assert.NotEmpty(t, senders)
}

func TestToken(t *testing.T) {
	token, err := Token()
	assert.NoError(t, err)
	assert.Condition(t, func() bool { return len(token)/2 == 16 }, "size of token should be 16 bit")
}
