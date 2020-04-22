package out

import (
	"encoding/json"
	"github.com/VEuPathDB/lib-go-rest-types/veupath/service/recordtypes"
	"github.com/VEuPathDB/script-site-param-cache/internal/log"
	"strings"
)

const (
	warnNogo = `Cannot run search with "%s" parameter.
		Search Name: %s
		Record Type: %s`

	postReqErr = `POST request failed with code %d
		URL:     %s
		Message: %s
		Search:  %s
		Payload: %s`

	getReqErr = `GET request failed with code %d
		URL:     %s
		Message: %s`

	getSearchErr = `Search lookup failed with code %d
		Service Url: %s
		App Url:     %s
		Response:    %s`

	vocabParse = `Failed to parse vocabulary for %s
		Search: %s`
)


// PostRequestError prints a formatted error message about a
// search submission that failed.
func PostRequestError(
	code uint16,
	url string,
	response []byte,
	search *recordtypes.FullSearch,
	payload *recordtypes.OrganismSearchRequest,
) {
	sea, _ := json.Marshal(search)
	pay, _ := json.Marshal(payload)
	log.ErrorFmt(postReqErr, code, trimTok(url), string(response), sea, pay)
}

// GetRequestError prints a formatted error message about an
// arbitrary GET request that failed.
func GetRequestError(code uint16, url string, response []byte) {
	log.ErrorFmt(getReqErr, code, trimTok(url), string(response))
}

// GetSearchError prints a formatted error message about a
// search lookup request that failed.
func GetSearchError(code uint16, url string, response []byte) {
	url = trimTok(url)
	log.ErrorFmt(getSearchErr, code, url, convertUrl(url), string(response))
}

func WarnCannotRun(
	identifier string,
	search *recordtypes.Search,
	recordType *recordtypes.RecordType,
) {
	log.TraceFmt(warnNogo, identifier, search.UrlSegment, recordType.UrlSegment)
}

func VocabParseErr(paramType string, search *recordtypes.Search) {
	b, _ := json.Marshal(search)
	log.ErrorFmt(vocabParse, paramType, string(b))
}

func trimTok(url string) string {
	if ind := strings.Index(url, "?"); ind > -1 {
		return url[:ind]
	} else {
		return url
	}
}

const (
	siteAppSegment = "/app/search/"
	siteSkip1      = "/service/record-types/"
	siteSkip2      = "/searches/"
)
func convertUrl(svcUrl string) (appUrl string) {
	ind1 := strings.Index(svcUrl, siteSkip1)
	ind2 := strings.Index(svcUrl, siteSkip2)

	if ind1 == -1 || ind2 == -1 {
		return "!!invalid service url"
	}

	appUrl = svcUrl[:ind1]
	appUrl += siteAppSegment
	appUrl += svcUrl[ind1 + len(siteSkip1):ind2]
	appUrl += "/"
	appUrl += svcUrl[ind2 + len(siteSkip2):]

	return
}