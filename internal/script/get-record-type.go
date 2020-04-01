package script

import (
	"encoding/json"
	"net/http"

	R "github.com/Foxcapades/Go-ChainRequest/simple"

	"github.com/VEuPathDB/lib-go-rest-types/veupath/service/recordtypes"
	"github.com/VEuPathDB/script-site-param-cache/internal/log"
)

func (r *Runner) processRecordType(rType string) {
	fullUrl := r.url.RecordTypeUrl(rType)
	r.push(fullUrl)
	r.wp.Submit(func() {
		r.start(fullUrl)
		defer r.pop(fullUrl)

		record  := new(recordtypes.RecordType)

		log.TraceFmt("Looking up searches for record type %s", rType)

		res       := R.GetRequest(fullUrl).Submit()
		code, err := res.GetResponseCode()

		if err != nil {
			log.ErrorFmt(err.Error())
			return
		}

		if code != http.StatusOK {
			getReqError(code, fullUrl, res.MustGetBody())
			return
		}

		err = res.UnmarshalBody(record, R.UnmarshallerFunc(json.Unmarshal))

		if err != nil {
			log.ErrorFmt(err.Error())
			return
		}

		for i := range record.Searches {
			r.processShortSearch(record, &record.Searches[i])
		}
	})
}
