package api

import (
	"github.com/openshift-online/rh-trex/pkg/util"
	"gorm.io/gorm"
)

var (
	DinosaurTypeName = util.GetBaseType(Dinosaur{})
)

type Dinosaur struct {
	Meta
	Species string
}

type DinosaurList []*Dinosaur
type DinosaurIndex map[string]*Dinosaur

func (l DinosaurList) Index() DinosaurIndex {
	index := DinosaurIndex{}
	for _, o := range l {
		index[o.ID] = o
	}
	return index
}

func (d *Dinosaur) BeforeCreate(tx *gorm.DB) error {
	d.ID = NewID()
	return nil
}

type DinosaurPatchRequest struct {
	Species *string `json:"species,omitempty"`
}
