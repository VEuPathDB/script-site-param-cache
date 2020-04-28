package util

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"

	"github.com/Foxcapades/Go-ChainRequest"
	"github.com/Foxcapades/Go-ChainRequest/request/header"
	req "github.com/Foxcapades/Go-ChainRequest/simple"

	"github.com/VEuPathDB/script-site-param-cache/internal/x"
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

func printRequestStart(method, url string, body interface{}) {
	log.Tracef("Starting %s Request\nURL: %s", method, url)
	if log.IsLevelEnabled(log.TraceLevel) {
		if body != nil {
			log.Tracef("Payload:", string(x.JsonMarshal(body)))
		}
	}
}

func printRequestStats(method, url string, start time.Time, timing *time.Duration, res creq.Response) {
	dur := time.Now().Sub(start)
	*timing = dur

	if res.GetError() != nil {
		log.Debugf(reqEndText, method, "failed", "", url, dur)
		return
	}

	code := res.MustGetResponseCode()
	succ := ""

	if code < 200 || code > 299 {
		succ = " unsuccessfully"
	} else {
		succ = " successfully"
	}

	log.Debugf(reqEndText, method, "completed", succ, url, dur)
	log.Tracef("Message Body:", string(res.MustGetBody()))
}
