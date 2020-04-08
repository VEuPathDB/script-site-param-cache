package script

import (
	"encoding/json"
	"net/http"

	R "github.com/Foxcapades/Go-ChainRequest/simple"

	"github.com/VEuPathDB/lib-go-rest-types/veupath/service/recordtypes"
	"github.com/VEuPathDB/script-site-param-cache/internal/util"
	"github.com/VEuPathDB/script-site-param-cache/internal/x"
)

func (r *Runner) processRecordType(rType string) {
	fullUrl := r.url.RecordTypeUrl(rType)

	r.push(fullUrl)
	r.wp.Submit(x.PanicCatcher(func() {
		defer r.pop(fullUrl)
		r.start(fullUrl)

		record := new(recordtypes.RecordType)

		res := util.GetRequest(fullUrl, &r.client)
		if code := res.MustGetResponseCode(); code != http.StatusOK {
			getReqError(code, fullUrl, res.MustGetBody())
			return
		}

		res.MustUnmarshalBody(record, R.UnmarshallerFunc(json.Unmarshal))

		for i := range record.Searches {
			r.processShortSearch(record, &record.Searches[i])
		}
	}))
}
