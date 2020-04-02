package script

import (
	"encoding/json"
	"net/http"

	R "github.com/Foxcapades/Go-ChainRequest/simple"

	"github.com/VEuPathDB/lib-go-rest-types/veupath/service/recordtypes"
	"github.com/VEuPathDB/script-site-param-cache/internal/log"
)

func (r *Runner) processShortSearch(
	record  *recordtypes.RecordType,
	sSearch *recordtypes.Search,
) {
	fullUrl := r.url.RecordSearchUrl(record.UrlSegment, sSearch.UrlSegment)
	r.push(fullUrl)
	r.wp.Submit(func() {
		r.start(fullUrl)
		defer r.pop(fullUrl)
		search  := new(recordtypes.FullSearch)

		log.TraceFmt("Fetching full search data for search %s on record-type %s",
			sSearch.UrlSegment, record.UrlSegment)

		res := R.GetRequest(fullUrl).SetHttpClient(&r.client).Submit()
		code, err := res.GetResponseCode()

		if err != nil {
			log.Error(err.Error())
			return
		}

		if code != http.StatusOK {
			getReqError(code, fullUrl, res.MustGetBody())
			return
		}

		err = res.UnmarshalBody(&search, R.UnmarshallerFunc(json.Unmarshal))

		if err != nil {
			log.Error(err.Error())
			return
		}

		if r.opts.RunSearches {
			r.processSearch(record, search)
		}
	})
}