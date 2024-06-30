package lprlogic

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"git.kanosolution.net/kano/dbflex"
	"git.kanosolution.net/kano/kaos"
	"github.com/ariefdarmawan/byter"
	"github.com/sebarcode/codekit"
	"github.com/sebarcode/laporan/lprmodel"
)

type ReportRunner struct {
	eder   byter.Byter
	prefix string
}

func NewReportRunner(enc byter.Byter, prefix string) *ReportRunner {
	rr := new(ReportRunner)
	rr.eder = enc
	rr.prefix = prefix
	return rr
}

func (rr *ReportRunner) Run(ctx *kaos.Context, report *lprmodel.ReportConfig, payload codekit.M, searchFn SearchFunction) (*lprmodel.ReportResult, error) {
	if ctx == nil {
		return nil, errors.New("ctx is null")
	}

	if report == nil {
		return nil, errors.New("report config is null")
	}

	var (
		req *http.Request
		err error
	)
	origRequest := ctx.Data().Get("http_request", new(http.Request)).(*http.Request)

	getUrl := report.GetUrl
	if rr.prefix != "" {
		getUrl, err = url.JoinPath(rr.prefix, getUrl)
		if err != nil {
			return nil, fmt.Errorf("report get url: %s, %s", getUrl, err.Error())
		}
	}

	qp := new(dbflex.QueryParam)
	switch report.SearchType {
	case lprmodel.SearchFunction:
		if searchFn == nil {
			return nil, fmt.Errorf("search function is nil")
		}
		qp, err = searchFn(ctx, payload)
		if err != nil {
			return nil, errors.New("searchFn to queryParam: " + err.Error())
		}

	case lprmodel.SearchQueryParam:
		err = byter.Cast(rr.eder, payload, qp, nil)
		if err != nil {
			return nil, errors.New("payload to QueryParam: " + err.Error())
		}
	}

	payloadBytes, err := rr.eder.Encode(qp)
	if err != nil {
		return nil, errors.New("byter encode: " + err.Error())
	}

	reqReader := bytes.NewReader(payloadBytes)
	if report.GetMethod == "" {
		report.GetMethod = http.MethodPost
	}
	req, err = http.NewRequest(report.GetMethod, getUrl, reqReader)
	if err != nil {
		return nil, errors.New("invoke rest: " + err.Error())
	}
	for key, vals := range origRequest.Header {
		for _, val := range vals {
			req.Header.Add(key, val)
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call: %s, %s ", getUrl, err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("call: %s, %s ", getUrl, resp.Status)
	}
	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("get response: %s, %s ", getUrl, err.Error())
	}
	if len(bs) == 0 {
		return nil, fmt.Errorf("empty response: %s", getUrl)
	}

	dataset := []codekit.M{}
	err = rr.eder.DecodeTo(bs, &dataset, nil)
	if err != nil {
		return nil, fmt.Errorf("get response: %s ", err.Error())
	}

	resDs := make([]codekit.M, len(dataset))
	if len(report.ViewSetups) == 0 {
		resDs = dataset
	} else {
		for i, d := range dataset {
			ds := codekit.M{}
			for k, v := range d {
				ds.Set(k, v)
			}
			resDs[i] = ds
		}
	}

	return &lprmodel.ReportResult{
		Data: resDs,
	}, nil
}
