package script

import (
	"fmt"
	"time"

	"github.com/VEuPathDB/lib-go-wdk-api/v0/except"
	"github.com/VEuPathDB/lib-go-wdk-api/v0/model/record"
	"github.com/VEuPathDB/lib-go-wdk-api/v0/model/search"
	log "github.com/sirupsen/logrus"

	"github.com/VEuPathDB/script-site-param-cache/internal/out"
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

// log messages
const (
	logSkipSearch = "Skipping search \"%s\", it is marked as excluded."
	logSearchLookup = "Successful search GET"
)

// retrieves the full search details json for a given
// search and record type.  Optionally runs the search if
// search running is enabled.
func (r *Runner) processShortSearch(
	record *record.Type,
	sSearch *search.ShortSearch,
) {
	if ok := exclusions[sSearch.UrlSegment]; ok {
		log.Infof(logSkipSearch, sSearch.FullName)
		r.stats.SearchDetailSkipped()
		return
	}

	api := r.api.RecordApiFor(record.UrlSegment)
	fullUrl := api.UrlBuilder().Search(sSearch.UrlSegment)
	r.push(fullUrl)
	r.wp.Submit(x.PanicCatcher(func() {
		r.start(fullUrl)
		defer r.pop(fullUrl)
		defer func() {
			if rec := recover(); rec != nil {
				r.stats.SearchDetailFailed()
				out.GetSearchError(0, fullUrl, []byte(fmt.Sprint(rec)))
			}
		}()

		start := time.Now()
		sch, err := api.GetSearch(sSearch.UrlSegment)
		timing := time.Since(start)
		r.stats.RecordTiming(sSearch.FullName, timing)
		if err != nil {
			r.stats.SearchDetailFailed()

			if cst, ok := err.(except.HttpRequestError); ok {
				var code uint16
				if cst.ResponseCode().Exists() {
					code = uint16(cst.ResponseCode().Get())
				}

				out.GetSearchError(code, fullUrl, cst.ResponseBody())
			} else {
				out.GetSearchError(0, fullUrl, []byte(err.Error()))
			}

			return
		}

		log.WithFields(log.Fields{
			"time": timing,
			"search": sSearch.FullName,
		}).Debug(logSearchLookup)
		r.stats.SearchDetailSuccess()
		if r.opts.SearchEnabled() {
			r.processSearch(record, &sch)
		}
	}))
}


