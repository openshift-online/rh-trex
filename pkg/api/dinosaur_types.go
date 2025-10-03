package api

import (
	"gorm.io/gorm"
	coreapi "github.com/openshift-online/rh-trex-core/api"
)

type Dinosaur struct {
	coreapi.Meta
	Species string `json:"species" gorm:"index"`
}

// GetMeta returns the metadata for the dinosaur (implements MetaProvider)
func (d *Dinosaur) GetMeta() *coreapi.Meta {
	return &d.Meta
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
