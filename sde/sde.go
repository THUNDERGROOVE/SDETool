package sde

import (
	"encoding/gob"
	"fmt"
	"os"
	"strings"
)

var (
	ErrTypeDoesNotExist = fmt.Errorf("sde: type does not exist")
)

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

type SDE struct {
	Version  string
	Official bool
	Types    map[int]SDEType
}

func (s *SDE) GetType(id int) (sdetype *SDEType, err error) {
	if v, ok := s.Types[id]; ok {
		return &v, nil
	} else {
		return nil, ErrTypeDoesNotExist
	}

}

func (s *SDE) Search(ss string) (sdetypes []*SDEType, err error) {
	out := make([]*SDEType, 0)
	for _, v := range s.Types {
		if strings.Contains(v.GetName(), ss) || strings.Contains(v.TypeName, ss) {
			out = append(out, &v)
		}
	}
	return out, nil
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
