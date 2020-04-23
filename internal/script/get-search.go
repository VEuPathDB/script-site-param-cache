package script

import (
	"encoding/json"
	R "github.com/Foxcapades/Go-ChainRequest/simple"
	"github.com/VEuPathDB/script-site-param-cache/internal/log"
	"github.com/VEuPathDB/script-site-param-cache/internal/out"
	"net/http"

	"github.com/VEuPathDB/lib-go-rest-types/veupath/service/recordtypes"
	"github.com/VEuPathDB/script-site-param-cache/internal/util"
	"github.com/VEuPathDB/script-site-param-cache/internal/x"
)

// retrieves the full search details json for a given
// search and record type.  Optionally runs the search if
// search running is enabled.
func (r *Runner) processShortSearch(
	record *recordtypes.RecordType,
	sSearch *recordtypes.Search,
) {
	if ok := exclusions[sSearch.UrlSegment]; ok {
		log.DebugFmt("Skipping search \"$s\", it is marked as excluded.")
		r.stats.SearchDetailSkipped()
		return
	}

	fullUrl := r.url.RecordSearchUrl(record.UrlSegment, sSearch.UrlSegment)
	r.push(fullUrl)
	r.wp.Submit(x.PanicCatcher(func() {
		r.start(fullUrl)
		defer r.pop(fullUrl)
		search := new(recordtypes.FullSearch)

		res := util.GetRequest(fullUrl, &r.client)

		if code := res.MustGetResponseCode(); code != http.StatusOK {
			out.GetSearchError(code, fullUrl, res.MustGetBody())
			r.stats.SearchDetailFailed()
			return
		}

		res.MustUnmarshalBody(&search, R.UnmarshallerFunc(json.Unmarshal))
		r.stats.SearchDetailSuccess()

		if r.opts.SearchEnabled() {
			r.processSearch(record, search)
		}
	}))
}

// Searches that cannot be looked up due to the guest user
// not having the necessary prerequisites for those searches
// to be available
var exclusions = map[string]bool{
	"GenesByUserDatasetAntisense": true,
	"GenesByRNASeqUserDataset": true,
	"DatasetsByReferenceName": true,
}

