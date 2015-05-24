package sde

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	ErrTypeDoesNotExist = fmt.Errorf("sde: type does not exist")
)

// Load loads an encoding/gob encoded SDE object from file
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

// LoadReader returns an SDE pointer given an io.Reader to read from
func LoadReader(r io.Reader) (*SDE, error) {
	s := &SDE{}
	dec := gob.NewDecoder(r)
	if err := dec.Decode(s); err != nil {
		return nil, err
	}
	return s, nil
}

// Save saves a provided SDE object to disk
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

/*
	SDE is a struct that owns every type for a given SDE.
 		@TODO:
		Add more old methods:
			GetTypeByName
			GetTypeByTag
			...
		Add lookups:
			TypeName
			Attrribute["mDiplsayName"]
			Use a map that isn't gobed and generate on load(use goroutine)
*/
type SDE struct {
	Version  string
	Official bool
	Types    map[int]*SDEType
}

// GetType returns a pointer to an SDEType or nil and an error
func (s *SDE) GetType(id int) (sdetype *SDEType, err error) {
	if v, ok := s.Types[id]; ok {
		return v, nil
	} else {
		return nil, ErrTypeDoesNotExist
	}

}

// Search checks for the existance of ss in mDisplayName or TypeName in every type and returns
// a slice of pointers to SDETypes
func (s *SDE) Search(ss string) (sdetypes []*SDEType, err error) {
	out := make([]*SDEType, 0)
	for _, v := range s.Types {
		if strings.Contains(strings.ToLower(v.GetName()), strings.ToLower(ss)) || strings.Contains(strings.ToLower(v.TypeName), strings.ToLower(ss)) {
			fmt.Printf("Appending %v to slice.\nAddress: %p\n", v.GetName(), &v)
			out = append(out, v)
		}
	}
	return out, nil
}

/*
	SDEType is a struct representing a single individual type in an SDE.
	@TODO:
		Add old methods.
		Make some cleaner way than before of checking for the existance of *.*... atributes:
		Options:
			1) Substruct them out and create a parser for each(yuck)
			2) Map[string]interface{} parser(ehh)
*/
type SDEType struct {
	TypeID     int
	TypeName   string
	Attributes map[string]interface{}
}

// GetName returns the string value of Attributes["mDisplayName"] if it exists.  Otherwise we return TypeName
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
