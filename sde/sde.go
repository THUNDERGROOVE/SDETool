package sde

import (
	"encoding/gob"
	"os"
)

type SDE struct {
	Types map[int]SDEType
}

func Load(filename string) (*SDE, error) {
	if f, err := os.OpenFile(filename, os.O_RDONLY, 0777); err != nil {
		return nil, err
	} else {
		s := &SDE{}
		dec := gob.NewDecoder(f)
		if err := dec.Decode(s); err != nil {
			return err
		}
		return s, nil
	}
	return nil, nil
}

func Save(filename string, s *SDE) error {
	if f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0777); err != nil {
		return err
	} else {
		enc := gob.NewEncoder(f)
		if err := dec.Encode(s); err != nil {
			return err
		}
	}
	return nil
}

type SDEType struct {
	TypeID     int
	TypeName   string
	Attributes map[string]interface{}
}

func init() {
	gob.Register(SDE{})
	gob.Register(SDEType{})
}
