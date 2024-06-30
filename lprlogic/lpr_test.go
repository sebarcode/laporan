package lprlogic_test

import (
	"context"
	"net/http"
	"testing"

	"git.kanosolution.net/kano/dbflex"
	"git.kanosolution.net/kano/kaos"
	"github.com/ariefdarmawan/byter"
	"github.com/sebarcode/codekit"
	"github.com/sebarcode/laporan/lprlogic"
	"github.com/sebarcode/laporan/lprmodel"
	"github.com/sebarcode/logger"
	"github.com/smartystreets/goconvey/convey"
)

var (
	lgr = logger.NewLogEngine(true, false, "", "", "")
)

func TestLaporan(t *testing.T) {
	convey.Convey("laporan", t, func() {
		var (
			ctx *kaos.Context
			rpt *lprmodel.ReportConfig
		)

		ctx = kaos.NewContext(context.Background(), lgr, nil, nil, kaos.NewSharedData(), nil)
		rpt = &lprmodel.ReportConfig{
			ID:         "laporan",
			Name:       "Item Master",
			SearchType: lprmodel.SearchQueryParam,
			GetMethod:  http.MethodPost,
			GetUrl:     "fico/paymentterm/find",
		}

		param, _ := codekit.ToM(dbflex.NewQueryParam())
		lpr := lprlogic.NewReportRunner(byter.NewByter(""), "https://bis-dev.kanosolution.app/v1")
		res, err := lpr.Run(ctx, rpt, param, nil)
		convey.So(err, convey.ShouldBeNil)

		convey.Convey("validate result", func() {
			convey.So(len(res), convey.ShouldBeGreaterThan, 0)
			convey.Printf(" len of result is %d", len(res))
		})
	})
}
