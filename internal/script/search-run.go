package script

import (
	"encoding/json"
	"net/http"

	"github.com/Foxcapades/Go-ChainRequest/request/header"
	R "github.com/Foxcapades/Go-ChainRequest/simple"

	"github.com/VEuPathDB/lib-go-rest-types/veupath/service/recordtypes"
	"github.com/VEuPathDB/script-site-param-cache/internal/log"
)

const (
	warnNogo = `Cannot run search with "%s" parameter.
              Search Name: %s
              Record Type: %s`
)

var (
	disallowedParamNames = map[string]bool{
		"primaryKeys": false, // can't auto-populate this
	}
	disallowedParamTypes = map[string]bool{
		"input-step":    false, // we don't have any step ids
		"input-dataset": false, // we don't have any dataset ids
	}
)

func (r *Runner) processSearch(
	record *recordtypes.RecordType,
	search *recordtypes.FullSearch,
) {
	fullUrl := r.url.RecordSearchStdReportUrl(record.UrlSegment, search.SearchData.UrlSegment)
	r.push(fullUrl)
	r.wp.Submit(func() {
		r.start(fullUrl)
		defer r.pop(fullUrl)
		log.TraceFmt("Running search %s for record type %s", search.SearchData.UrlSegment, record.UrlSegment)

		inputBody, ok := prepareSearchRequest(record, &search.SearchData)
		if !ok {
			return
		}


		res := R.PostRequest(fullUrl).
			SetHeader(header.CONTENT_TYPE, "application/json").
			MarshalBody(inputBody, R.MarshallerFunc(json.Marshal)).
			Submit()

		code, err := res.GetResponseCode()

		if err != nil {
			log.Error(err)
			return
		}

		if code != http.StatusOK {
			body, _ := res.GetBody()
			b, _ := json.Marshal(search)
			c, _ := json.Marshal(inputBody)
			log.ErrorFmt("Request failed with code %d\nURL:     %s\nMessage: %s\nSearch:  %s\nPayload: %s",
				code, fullUrl, string(body), b, c)
		} else {
			log.TraceFmt("Finished record-types/%s/searches/%s/reports/standard %d",
				record.UrlSegment, search.SearchData.UrlSegment, code)
		}
	})
}

func prepareSearchRequest(
	record *recordtypes.RecordType,
	search *recordtypes.Search,
) (out *recordtypes.OrganismSearchRequest, ok bool) {
	out = recordtypes.NewOrganismSearchRequest()

	for i := range search.Parameters {

		tmp := &search.Parameters[i]

		if _, ok := disallowedParamNames[tmp.Name]; ok {
			log.WarnFmt(warnNogo, "name: " + tmp.Name, search.UrlSegment,
				record.UrlSegment)
			return nil, false
		}

		if _, ok := disallowedParamTypes[tmp.Type]; ok {
			log.WarnFmt(warnNogo, "type: " + tmp.Type, search.UrlSegment,
				record.UrlSegment)
			return nil, false
		}

		if len(tmp.InitialDisplayValue) > 0 {
			out.SearchConfig.Parameters[tmp.Name] = tmp.InitialDisplayValue
		} else if tmp.Type == "multi-pick-vocabulary" {
			out.SearchConfig.Parameters[tmp.Name] = `["yes"]`
		} else {
			out.SearchConfig.Parameters[tmp.Name] = "yes"
		}
	}

	out.ReportConfig.Attributes = search.DefaultAttributes

	return out, true
}
