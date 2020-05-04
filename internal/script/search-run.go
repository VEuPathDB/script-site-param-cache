package script

import (
	"github.com/VEuPathDB/lib-go-wdk-api/v0/model/param"
	"github.com/VEuPathDB/lib-go-wdk-api/v0/model/record"
	"github.com/VEuPathDB/lib-go-wdk-api/v0/model/search"
	"github.com/VEuPathDB/script-site-param-cache/internal/out"
	"net/http"
	"time"

	"github.com/VEuPathDB/lib-go-rest-types/veupath/service/recordtypes"
	"github.com/VEuPathDB/script-site-param-cache/internal/util"
	"github.com/VEuPathDB/script-site-param-cache/internal/x"
)

var (
	disallowedParamNames = map[string]bool{
		"primaryKeys": false, // can't auto-populate this
	}
	disallowedParamTypes = map[param.Kind]bool{
		param.KindAnswer:  false, // we don't have any step ids
		param.KindDataset: false, // we don't have any dataset ids
	}
)

func (r *Runner) processSearch(
	record *record.Type,
	search *search.ValidatedSearch,
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

		var timing time.Duration
		res := util.PostRequest(fullUrl, &r.client, &timing, inputBody)

		if code := res.MustGetResponseCode(); code != http.StatusOK {
			out.PostRequestError(code, fullUrl, res.MustGetBody(),
				&search.SearchData, inputBody)
		}
	}))
}

func prepareSearchRequest(
	record *record.Type,
	search *search.FullSearch,
) (ret *recordtypes.OrganismSearchRequest, ok bool) {
	ret = recordtypes.NewOrganismSearchRequest()

	for i := range search.Parameters {

		tmp := &search.Parameters[i]

		if _, ok := disallowedParamNames[tmp.Name]; ok {
			out.WarnCannotRun("name: "+tmp.Name, search, record)
			return nil, false
		}

		if _, ok := disallowedParamTypes[tmp.Kind]; ok {
			out.WarnCannotRun("type: " + string(tmp.Kind), search, record)
			return nil, false
		}

		if tmp.Kind == param.KindMultiVocab {

			enum := tmp.AsKindMultiVocab()

			if enum.Vocabulary.Exists() {
				ret.SearchConfig.Parameters[tmp.Name] = `["yes"]`
				continue
			}

			switch enum.DisplayType {
			case param.DisplayTypeTreeBox:
				if val, ok := treeBoxParam(&enum); ok {
					ret.SearchConfig.Parameters[enum.Name] = val
					continue
				} else {
					return nil, false
				}
			case param.DisplayTypeTypeAhead, param.DisplayTypeCheckBox:
				if val, ok := enumParam(&enum); ok {
					ret.SearchConfig.Parameters[enum.Name] = val
					continue
				} else {
					return nil, false
				}
			}
		}

		if tmp.InitialDisplayValue != nil {
			ret.SearchConfig.Parameters[tmp.Name] = *tmp.InitialDisplayValue
		} else {
			ret.SearchConfig.Parameters[tmp.Name] = "1"
		}
	}

	ret.ReportConfig.Attributes = search.DefaultAttributes

	return ret, true
}

func treeBoxParam(param *param.Enum) (val string, ok bool) {
	voc := param.Vocabulary.GetAsTree()

	if voc.Data.Term == "@@fake@@" && len(voc.Children) > 0 {
		voc = &voc.Children[0]
	}

	return `["` + voc.Data.Term + `"]`, true
}

func enumParam(param *param.Enum) (val string, ok bool) {
	voc := param.Vocabulary.GetAsTable()

	return `["` + voc[0][0] + `"]`, true
}
