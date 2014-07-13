package smsru

import (
	"bufio"
	"errors"
	"io/ioutil"
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

func Token() (token string, err error) {
	var cli http.Client
	var resp *http.Response
	u, _ := url.Parse(baseURL)
	u.Path = authGetToken
	if resp, err = cli.Get(u.String()); err != nil {
		return
	}
	defer resp.Body.Close()
	var tok []byte
	if tok, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}
	if len(tok) < 1 {
		err = ErrBadServerReply
		return
	}
	token = string(tok)
	return
}

func (a *Api) smsUrl(sms *Sms) (u *url.URL) {
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

	u, _ = url.Parse(baseURL)
	u.RawQuery = val.Encode()
	return
}

func (a *Api) SmsSend(sms *Sms) (ids []string, err error) {
	u := a.smsUrl(sms)
	u.Path = smsSend

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

func (a *Api) SmsStatus(id string) (err error) {
	val := make(url.Values)
	val.Add("id", id)
	u, _ := url.Parse(baseURL)
	u.Path = smsSend
	u.RawQuery = val.Encode()
	_, err = a.get(u)
	return
}

func (a *Api) SmsCost(sms *Sms) (cost float32, n int, err error) {
	u := a.smsUrl(sms)
	u.Path = smsCost
	var lines []string
	if lines, err = a.get(u); err != nil {
		return
	}
	if len(lines) < 2 {
		err = ErrBadServerReply
		return
	}
	var c float64
	if c, err = strconv.ParseFloat(lines[0], 32); err != nil {
		return
	}
	cost = float32(c)
	var num uint64
	if num, err = strconv.ParseUint(lines[1], 10, 8); err != nil {
		return
	}
	n = int(num)
	return
}

func (a *Api) MyBalance() (balance float32, err error) {
	u, _ := url.Parse(baseURL)
	u.Path = myBalance
	var lines []string
	if lines, err = a.get(u); err != nil {
		return
	}
	if len(lines) < 1 {
		err = ErrBadServerReply
		return
	}
	var bal float64
	if bal, err = strconv.ParseFloat(lines[0], 32); err != nil {
		return
	}
	balance = float32(bal)
	a.Balance = balance
	return
}

func (a *Api) MyLimit() (limit int, n int, err error) {
	u, _ := url.Parse(baseURL)
	u.Path = myLimit
	var lines []string
	if lines, err = a.get(u); err != nil {
		return
	}
	if len(lines) < 2 {
		err = ErrBadServerReply
		return
	}
	var lim, num uint64
	if lim, err = strconv.ParseUint(lines[0], 10, 0); err != nil {
		return
	}
	limit = int(lim)
	if num, err = strconv.ParseUint(lines[0], 10, 0); err != nil {
		return
	}
	n = int(num)
	return
}

func (a *Api) MySenders() (senders []string, err error) {
	u, _ := url.Parse(baseURL)
	u.Path = mySenders
	var lines []string
	if lines, err = a.get(u); err != nil {
		return
	}
	if len(lines) < 1 {
		err = ErrBadServerReply
		return
	}
	for _, l := range lines {
		senders = append(senders, l)
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
		i, err = strconv.ParseInt(sc.Text(), 10, 16)
		if err != nil {
			return
		}
		if u.Path == smsStatus {
			if err = Status(int(i)); err != nil {
				return
			}
		} else {
			if err = Error(int(i)); err != nil {
				return
			}
		}
	} else {
		return nil, ErrBadServerReply
	}
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return
}
