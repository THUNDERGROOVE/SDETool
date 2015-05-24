package sde

import (
	"encoding/gob"
	"os"
)

type SDE struct {
	Version  string
	Official bool
	Types    map[int]SDEType
}

func Load(filename string) (*SDE, error) {
	if f, err := os.OpenFile(filename, os.O_RDONLY, 0777); err != nil {
		return nil, err
	} else {
		s := &SDE{}
		dec := gob.NewDecoder(f)
		if err := dec.Decode(s); err != nil {
			return nil, err
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
		if err := enc.Encode(s); err != nil {
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

func (s *SDEType) GetName() string {
	if v, ok := s.Attributes["mDisplayName"]; ok {
		return v.(string)
	}
	return s.TypeName
}

func init() {
	gob.Register(SDE{})
	gob.Register(SDEType{})
}
