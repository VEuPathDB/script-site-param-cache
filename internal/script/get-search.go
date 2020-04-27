package script

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	R "github.com/Foxcapades/Go-ChainRequest/simple"
	"github.com/VEuPathDB/lib-go-rest-types/veupath/service/recordtypes"
	"github.com/VEuPathDB/lib-go-wdk-api/v0/service/recordTypes"
	"github.com/VEuPathDB/lib-go-wdk-api/v0/service/recordTypes/searches"
	"github.com/VEuPathDB/script-site-param-cache/internal/log"
	"github.com/VEuPathDB/script-site-param-cache/internal/out"

	"github.com/VEuPathDB/script-site-param-cache/internal/util"
	"github.com/VEuPathDB/script-site-param-cache/internal/x"
)

// Searches that cannot be looked up due to the guest user
// not having the necessary prerequisites for those searches
// to be available
var exclusions = map[string]bool{
	"GenesByUserDatasetAntisense": true,
	"GenesByRNASeqUserDataset": true,
	"DatasetsByReferenceName": true,
}

// retrieves the full search details json for a given
// search and record type.  Optionally runs the search if
// search running is enabled.
func (r *Runner) processShortSearch(
	record *recordTypes.RecordTypeResponse,
	sSearch *searches.SearchResponse,
) {
	if ok := exclusions[sSearch.UrlSegment]; ok {
		log.InfoFmt("Skipping search \"%s\", it is marked as excluded.", sSearch.FullName)
		r.stats.SearchDetailSkipped()
		return
	}

	fullUrl := r.url.RecordSearchUrl(record.UrlSegment, sSearch.UrlSegment)
	r.push(fullUrl)
	r.wp.Submit(x.PanicCatcher(func() {
		r.start(fullUrl)
		defer r.pop(fullUrl)
		defer func() {
			if r := recover(); r != nil {
				out.GetSearchError(0, fullUrl, []byte(fmt.Sprint(r)))
			}
		}()
		search := new(recordtypes.FullSearch)

		var timing time.Duration
		res := util.GetRequest(fullUrl, &timing, &r.client)

		if code := res.MustGetResponseCode(); code != http.StatusOK {
			out.GetSearchError(code, fullUrl, res.MustGetBody())
			r.stats.SearchDetailFailed()
			return
		}

		r.stats.RecordTiming(sSearch.FullName, timing)

		res.MustUnmarshalBody(&search, R.UnmarshallerFunc(json.Unmarshal))
		r.stats.SearchDetailSuccess()

		if r.opts.SearchEnabled() {
			r.processSearch(record, search)
		}
	}))
}


