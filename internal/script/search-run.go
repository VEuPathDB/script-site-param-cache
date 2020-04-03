package script

import (
	"encoding/json"
	"net/http"

	"github.com/VEuPathDB/lib-go-rest-types/veupath/service/recordtypes"
	"github.com/VEuPathDB/script-site-param-cache/internal/util"
	"github.com/VEuPathDB/script-site-param-cache/internal/x"
)

var (
	disallowedParamNames = map[string]bool{
		"primaryKeys": false, // can't auto-populate this
	}
	disallowedParamTypes = map[string]bool{
		"input-step":    false, // we don't have any step ids
		"input-dataset": false, // we don't have any dataset ids
	}
)

func (r *Runner) processSearch(
	record *recordtypes.RecordType,
	search *recordtypes.FullSearch,
) {
	fullUrl := r.url.RecordSearchStdReportUrl(record.UrlSegment, search.SearchData.UrlSegment)

	r.push(fullUrl)

	r.wp.Submit(x.PanicCatcher(func() {
		r.start(fullUrl)
		defer r.pop(fullUrl)

		inputBody, ok := prepareSearchRequest(record, &search.SearchData)
		if !ok {
			return
		}

		res := util.PostRequest(fullUrl, &r.client, inputBody)

		if code := res.MustGetResponseCode(); code != http.StatusOK {
			postReqError(code, fullUrl, res.MustGetBody(), search, inputBody)
		}
	}))
}

func prepareSearchRequest(
	record *recordtypes.RecordType,
	search *recordtypes.Search,
) (out *recordtypes.OrganismSearchRequest, ok bool) {
	out = recordtypes.NewOrganismSearchRequest()

	for i := range search.Parameters {

		tmp := &search.Parameters[i]

		if _, ok := disallowedParamNames[tmp.Name]; ok {
			warnNoGo("name: "+tmp.Name, search, record)
			return nil, false
		}

		if _, ok := disallowedParamTypes[tmp.Type]; ok {
			warnNoGo("type: "+tmp.Type, search, record)
			return nil, false
		}

		if tmp.Type == "multi-pick-vocabulary" {

			if tmp.Vocabulary == nil {
				out.SearchConfig.Parameters[tmp.Name] = `["yes"]`
				continue
			}

			if tmp.DisplayType != nil {
				switch *tmp.DisplayType {
				case "treeBox":
					if val, ok := treeBoxParam(tmp, search); ok {
						out.SearchConfig.Parameters[tmp.Name] = val
						continue
					} else {
						return nil, false
					}
				case "typeAhead", "checkBox":
					if val, ok := enumParam(tmp, search); ok {
						out.SearchConfig.Parameters[tmp.Name] = val
						continue
					} else {
						return nil, false
					}
				}
			}
		}

		if len(tmp.InitialDisplayValue) > 0 {
			out.SearchConfig.Parameters[tmp.Name] = tmp.InitialDisplayValue
		} else {
			out.SearchConfig.Parameters[tmp.Name] = "1"
		}
	}

	out.ReportConfig.Attributes = search.DefaultAttributes

	return out, true
}

func treeBoxParam(
	param *recordtypes.Parameter,
	search *recordtypes.Search,
) (val string, ok bool) {
	voc := new(recordtypes.EnumParamTermNode)
	err := json.Unmarshal(param.Vocabulary, voc)

	if err != nil {
		vocabParseErr("tree box", search)
		return "", false
	}

	if voc.Data.Term == "@@fake@@" && len(voc.Children) > 0 {
		voc = &voc.Children[0]
	}

	return `["` + voc.Data.Term + `"]`, true
}

func enumParam(
	param *recordtypes.Parameter,
	search *recordtypes.Search,
) (val string, ok bool) {
	voc := make([][3]string, 0, 15)
	err := json.Unmarshal(param.Vocabulary, &voc)

	if err != nil {
		vocabParseErr("enum param", search)
		return "", false
	}

	return `["` + voc[0][0] + `"]`, true
}
