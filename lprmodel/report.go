package lprmodel

import (
	"time"

	"git.kanosolution.net/kano/dbflex"
	"git.kanosolution.net/kano/dbflex/orm"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SearchType string

const (
	SearchNone       SearchType = "JSON"
	SearchQueryParam SearchType = "QueryParam"
	SearchFunction   SearchType = "Function" // function prototype: func (ctx *kaos.Context, payload interface{}) (*dbflex.QueryParam, error)
)

type Report struct {
	orm.DataModelBase `bson:"-" json:"-"`
	ID                string `bson:"_id" json:"_id" key:"1" form_read_only_edit:"1" form_section:"General" form_section_auto_col:"2"`
	Name              string `form_required:"1" form_section:"General"`
	SearchType        SearchType
	FunctionName      string
	GetMethod         string
	GetUrl            string
	Created           time.Time `form_kind:"datetime" form_read_only:"1" grid:"hide" form_section:"Time Info" form_section_auto_col:"2"`
	LastUpdate        time.Time `form_kind:"datetime" form_read_only:"1" grid:"hide" form_section:"Time Info"`
}

func (o *Report) TableName() string {
	return "Reports"
}

func (o *Report) FK() []*orm.FKConfig {
	return orm.DefaultRelationManager().FKs(o)
}

func (o *Report) ReverseFK() []*orm.ReverseFKConfig {
	return orm.DefaultRelationManager().ReverseFKs(o)
}

func (o *Report) SetID(keys ...interface{}) {
	o.ID = keys[0].(string)
}

func (o *Report) GetID(dbflex.IConnection) ([]string, []interface{}) {
	return []string{"_id"}, []interface{}{o.ID}
}

func (o *Report) PreSave(dbflex.IConnection) error {
	if o.ID == "" {
		o.ID = primitive.NewObjectID().Hex()
	}
	if o.Created.IsZero() {
		o.Created = time.Now()
	}
	o.LastUpdate = time.Now()
	return nil
}

func (o *Report) PostSave(dbflex.IConnection) error {
	return nil
}
func (o Report) Indexes() []dbflex.DbIndex {
	return []dbflex.DbIndex{}
}
