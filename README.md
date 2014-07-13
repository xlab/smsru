Go-lang API bindings for [sms.ru](http://sms.ru)
------------------------------------------------

#### Installation

```
go get github.com/xlab/smsru
```

#### Implemented features

* Error handling;
* Simple auth by token;
* Status report listener.

#### API coverage

* [sms/send], [sms/status], [sms/cost];
* [my/balance], [my/limit], [my/senders];
* [auth/get_token].

[sms/send]: http://sms.ru/?panel=api&subpanel=method&show=sms/send
[sms/status]: http://sms.ru/?panel=api&subpanel=method&show=sms/status
[sms/cost]: http://sms.ru/?panel=api&subpanel=method&show=sms/cost
[my/balance]: http://sms.ru/?panel=api&subpanel=method&show=my/balance
[my/limit]: http://sms.ru/?panel=api&subpanel=method&show=my/limit
[my/senders]: http://sms.ru/?panel=api&subpanel=method&show=my/senders
[auth/get_token]: http://sms.ru/?panel=api&subpanel=method&show=auth/get_token

#### License

MIT
