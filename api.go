package smsru

import (
	"bufio"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const baseURL = "http://sms.ru"

const (
	smsSend      = "sms/send"
	smsStatus    = "sms/status"
	smsCost      = "sms/cost"
	myBalance    = "my/balance"
	myLimit      = "my/limit"
	mySenders    = "my/senders"
	authGetToken = "auth/get_token"
	authCheck    = "auth/check"
	stoplistAdd  = "stoplist/add"
	stoplistDel  = "stoplist/del"
	stoplistGet  = "stoplist/get"
)

var ErrBadServerReply = errors.New("smsru: bad server reply")

type Api struct {
	Token string

	Password string
	Login    string
	Balance  float32
}

func NewApi(token string) *Api {
	return &Api{
		Token: token,
	}
}

type Sms struct {
	From string
	To   string
	Text string

	Time      time.Time
	Translit  bool
	Test      bool
	PartnerId string
	Multiple  map[string]string
}

func (a *Api) SmsSend(sms *Sms) (ids []string, err error) {
	val := make(url.Values)
	if len(sms.Multiple) > 0 {
		for to, text := range sms.Multiple {
			val.Add("multi["+to+"]", text)
		}
	} else {
		val.Add("to", sms.To)
		val.Add("text", sms.Text)
	}
	if len(sms.From) > 0 {
		val.Add("from", sms.From)
	}
	if sms.Time.After(time.Now()) {
		val.Add("time", strconv.FormatInt(sms.Time.Unix(), 10))
	}
	if sms.Translit {
		val.Add("translit", "1")
	}
	if sms.Test {
		val.Add("test", "1")
	}
	if len(sms.PartnerId) > 0 {
		val.Add("partner_id", sms.PartnerId)
	}

	u, _ := url.Parse(baseURL)
	u.Path = smsSend
	u.RawQuery = val.Encode()
	var lines []string
	if lines, err = a.get(u); err != nil {
		return
	}
	if len(lines) < 1 || (!sms.Test && len(lines) < 2) {
		err = ErrBadServerReply
		return
	}

	var i int
	for ; i < len(sms.Multiple); i++ {
		ids = append(ids, lines[i])
	}
	if len(ids) < 1 {
		ids = append(ids, lines[0])
		i++
	}
	if !sms.Test {
		balanceStr := lines[i]
		idx := strings.Index(balanceStr, "=")
		if idx < 0 {
			err = ErrBadServerReply
			return
		}
		var val float64
		if val, err = strconv.ParseFloat(balanceStr[idx+1:], 32); err != nil {
			return
		}
		a.Balance = float32(val)
	}
	return
}

func (a *Api) auth(param url.Values) url.Values {
	param.Add("api_id", a.Token)
	return param
}

func (a *Api) get(u *url.URL) (lines []string, err error) {
	var cli http.Client
	var resp *http.Response

	// send request
	u.RawQuery = a.auth(u.Query()).Encode()
	if resp, err = cli.Get(u.String()); err != nil {
		return
	}
	defer resp.Body.Close()

	sc := bufio.NewScanner(resp.Body)
	if sc.Scan() {
		var i int64
		i, err = strconv.ParseInt(sc.Text(), 10, 8)
		if err != nil {
			return
		}
		if err = Error(int(i)); err != nil {
			return
		}
	} else {
		return nil, ErrBadServerReply
	}
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return
}
