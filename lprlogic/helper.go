package lprlogic

import (
	"git.kanosolution.net/kano/dbflex"
	"git.kanosolution.net/kano/kaos"
	"github.com/sebarcode/codekit"
)

type SearchFunction func(ctx *kaos.Context, payload codekit.M) (*dbflex.QueryParam, error)
