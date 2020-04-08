package script

import (
	"encoding/json"
	"net/http"

	R "github.com/Foxcapades/Go-ChainRequest/simple"

	"github.com/VEuPathDB/lib-go-rest-types/veupath/service/recordtypes"
	"github.com/VEuPathDB/script-site-param-cache/internal/util"
	"github.com/VEuPathDB/script-site-param-cache/internal/x"
)

func (r *Runner) processShortSearch(
	record *recordtypes.RecordType,
	sSearch *recordtypes.Search,
) {
	fullUrl := r.url.RecordSearchUrl(record.UrlSegment, sSearch.UrlSegment)
	r.push(fullUrl)
	r.wp.Submit(x.PanicCatcher(func() {
		r.start(fullUrl)
		defer r.pop(fullUrl)
		search := new(recordtypes.FullSearch)

		res := util.GetRequest(fullUrl, &r.client)

		if code := res.MustGetResponseCode(); code != http.StatusOK {
			getReqError(code, fullUrl, res.MustGetBody())
			return
		}

		res.MustUnmarshalBody(&search, R.UnmarshallerFunc(json.Unmarshal))

		if r.opts.SearchEnabled() {
			r.processSearch(record, search)
		}
	}))
}
