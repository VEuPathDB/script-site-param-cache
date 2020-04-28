package out

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"strings"

	"github.com/VEuPathDB/lib-go-rest-types/veupath/service/recordtypes"
	"github.com/VEuPathDB/lib-go-wdk-api/v0/model/record"
	"github.com/VEuPathDB/lib-go-wdk-api/v0/model/search"
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

	getSearchErr = `Search lookup failed with code %d
		Service Url: %s
		App Url:     %s
		Response:    %s`
)


// PostRequestError prints a formatted error message about a
// search submission that failed.
func PostRequestError(
	code uint16,
	url string,
	response []byte,
	search *search.FullSearch,
	payload *recordtypes.OrganismSearchRequest,
) {
	sea, _ := json.Marshal(search)
	pay, _ := json.Marshal(payload)
	log.Errorf(postReqErr, code, trimTok(url), string(response), sea, pay)
}

// GetSearchError prints a formatted error message about a
// search lookup request that failed.
func GetSearchError(code uint16, url string, response []byte) {
	url = trimTok(url)
	log.Errorf(getSearchErr, code, url, convertUrl(url), string(response))
}

func WarnCannotRun(
	identifier string,
	search *search.FullSearch,
	recordType *record.Type,
) {
	log.Tracef(warnNogo, identifier, search.UrlSegment, recordType.UrlSegment)
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