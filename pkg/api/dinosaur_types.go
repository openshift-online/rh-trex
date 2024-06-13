package api

import "gorm.io/gorm"

type Dinosaur struct {
	Meta
	Species string
	// This is to illustrate resource review in action
	// It passes integration tests as it's mocked
	// does not work for local envs pointing to integration AMS via proxy
	OrganizationId string
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
