package util

import (
	"encoding/json"
	"github.com/Foxcapades/Go-ChainRequest"
	"github.com/Foxcapades/Go-ChainRequest/request/header"
	req "github.com/Foxcapades/Go-ChainRequest/simple"
	"github.com/VEuPathDB/script-site-param-cache/internal/log"
	"github.com/VEuPathDB/script-site-param-cache/internal/x"
	"net/http"
	"time"
)

const (
	reqEndText = "%s Request %s%s\nUrl:  %s\nTime: %s"
)

func PostRequest(
	url string,
	client *http.Client,
	timing *time.Duration,
	body interface{},
) (res creq.Response) {
	start := time.Now()
	defer func() { printRequestStats(http.MethodPost, url, start, timing, res) }()
	printRequestStart(http.MethodPost, url, body)

	res = req.PostRequest(url).SetHeader(header.CONTENT_TYPE, "application/json").
		SetHttpClient(client).MarshalBody(body, req.MarshallerFunc(json.Marshal)).
		Submit()
	return
}

func GetRequest(url string, timing *time.Duration, client *http.Client) (res creq.Response) {
	start := time.Now()
	defer func() { printRequestStats(http.MethodGet, url, start, timing, res) }()
	printRequestStart(http.MethodGet, url, nil)

	res = req.GetRequest(url).SetHttpClient(client).Submit()
	return
}

func printRequestStart(method, url string, body interface{}) {
	log.TraceFmt("Starting %s Request\nURL: %s", method, url)
	if body != nil {
		log.TraceFn(func() []interface{} {
			return []interface{}{"Payload:", string(x.JsonMarshal(body))}
		})
	}
}

func printRequestStats(method, url string, start time.Time, timing *time.Duration, res creq.Response) {
	dur := time.Now().Sub(start)
	*timing = dur

	if res.GetError() != nil {
		log.DebugFmt(reqEndText, method, "failed", "", url, dur)
		return
	}

	code := res.MustGetResponseCode()
	succ := ""

	if code < 200 || code > 299 {
		succ = " unsuccessfully"
	} else {
		succ = " successfully"
	}

	log.DebugFmt(reqEndText, method, "completed", succ, url, dur)
	log.TraceFn(func() []interface{} {
		return []interface{}{
			"Message Body:", string(res.MustGetBody())}
	})
}
