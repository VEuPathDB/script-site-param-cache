package script

import (
	"encoding/json"
	"github.com/VEuPathDB/lib-go-rest-types/veupath/service/recordtypes"
	"github.com/VEuPathDB/script-site-param-cache/internal/log"
)

const (
	_warnNogo = `Cannot run search with "%s" parameter.
              Search Name: %s
              Record Type: %s`

	_postReqErr = `POST request failed with code %d
                 URL:     %s
                 Message: %s
                 Search:  %s
                 Payload: %s`

	_getReqErr = `GET request failed with code %d
                URL:     %s
                Message: %s`

	_vocabParse = `Failed to parse vocabulary for %s
                 Search: %s`
)

func postReqError(
	code uint16,
	url string,
	response []byte,
	search *recordtypes.FullSearch,
	payload *recordtypes.OrganismSearchRequest,
) {
	sea, _ := json.Marshal(search)
	pay, _ := json.Marshal(payload)
	log.ErrorFmt(_postReqErr, code, url, string(response), sea, pay)
}

func getReqError(code uint16, url string, response []byte) {
	log.ErrorFmt(_getReqErr, code, url, string(response))
}

func warnNoGo(
	identifier string,
	search *recordtypes.Search,
	recordType *recordtypes.RecordType,
) {
	log.TraceFmt(_warnNogo, identifier, search.UrlSegment, recordType.UrlSegment)
}

func vocabParseErr(paramType string, search *recordtypes.Search) {
	b, _ := json.Marshal(search)
	log.ErrorFmt(_vocabParse, paramType, string(b))
}
