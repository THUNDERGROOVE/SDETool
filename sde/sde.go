package sde

import (
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strings"
)

var (
	ErrTypeDoesNotExist    = fmt.Errorf("sde: type does not exist")
	ErrCacheNotInitialized = fmt.Errorf("sde: cache not initialized")
	ErrSDEIsNil            = fmt.Errorf("sde: SDE struct was nil")
	ErrTypeIsNil           = fmt.Errorf("sde: SDEType struct was nil")
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
		s.Cache = nil // Kill the cache.  Don't want it to be stale do we?
		if err := enc.Encode(s); err != nil {
			return err
		}
	}
	return nil
}

// SDE is a struct that owns every type for a given SDE.
type SDE struct {
	Version  string
	Official bool
	Types    map[int]*SDEType
	Cache    *Cache
}

// Cache is a struct that is included within SDE.
//
// Whenever an SDE file is loaded we populate this and whenever an SDE is
// saved we make the pointer nil.  The struct is supposed to provide
// faster lookups for things like TypeName and mDisplayName
//
// @TODO: Finish implementing Cache relating things
type Cache struct {
	TypeNameLookup map[string]*SDEType
	//DisplayNameLookup map[string]*SDEType Some types have overlapping names
}

// GetType returns a pointer to an SDEType or nil and an error
func (s *SDE) GetType(id int) (sdetype *SDEType, err error) {
	if s == nil {
		return nil, ErrSDEIsNil
	}
	if v, ok := s.Types[id]; ok {
		if v == nil {
			return nil, ErrTypeIsNil
		}
		v.Parent = s
		return v, nil
	} else {
		return nil, ErrTypeDoesNotExist
	}

}

// Search checks for the existance of ss in mDisplayName or TypeName in every
// type and returns a slice of pointers to SDETypes
func (s *SDE) Search(ss string) (sdetypes []*SDEType, err error) {
	out := make([]*SDEType, 0)
	for _, v := range s.Types {
		if strings.Contains(strings.ToLower(v.GetName()), strings.ToLower(ss)) ||
			strings.Contains(strings.ToLower(v.TypeName), strings.ToLower(ss)) {
			out = append(out, v)
		}
	}
	return out, nil
}

// GetTypeByTag @TODO: Finish implementing
func (s *SDE) GetTypesByTag(tag int) ([]*SDEType, error) { return nil, nil }

// GetTypeByName @TODO: Finish implementing
func (s *SDE) GetTypeByName(name string) ([]*SDEType, error) { return nil, nil }

// VerifySDEPrint prints the entire list of types/typeids to check for DB
// corruption
func (s *SDE) VerifySDEPrint() {
	for k, v := range s.Types {
		fmt.Printf("  [%v][%v] %v at %p\n", k, v.TypeID, v.GetName(), v)
	}
}

// FindTypeThatReferences returns any time that refers to the given type
//
// Suprising how fast this method runs
//
// @TODO:
//	When our caching system is finished update this to not iterate all ~3400
// types lol
// Not sure how I would cache referencing type attributes.. Hmmm
func (s *SDE) FindTypesThatReference(t *SDEType) ([]*SDEType, error) {
	out := make([]*SDEType, 0)
	for _, v := range s.Types {
		for _, attr := range v.Attributes {
			switch tid := attr.(type) {
			case int:
				if tid == t.TypeID && !sdeslicecontains(out, tid) {
					out = append(out, v)
				}
			}
		}
	}
	return out, nil
}

// Size estimates the memory usage of the SDE instance.
func (s *SDE) Size() int {
	base := int(reflect.ValueOf(*s).Type().Size())
	for _, v := range s.Types {
		vv := int(reflect.ValueOf(*v).Type().Size())
		for _, a := range v.Attributes {
			switch reflect.TypeOf(a).Kind() {
			case reflect.String:
				vv += len(a.(string))
				fallthrough
			default:
				vv += int(reflect.ValueOf(a).Type().Size())
			}
		}
		base += vv
	}
	return base
}

// Internal methods

func (s *SDE) cleanCache() {
	s.Cache = new(Cache)
	s.Cache.TypeNameLookup = make(map[string]*SDEType)
}

// Use whenever possible.  Benchmarks have shown it takes roughly the same
// amount of time to generate the cache as it does to perform one SDEType
// level lookup.  Let alone one that looks into SDEType.Attributes
func (s *SDE) generateCache() {
	s.cleanCache()
	for _, v := range s.Types {
		s.Cache.TypeNameLookup[v.TypeName] = v
	}
}

func (s *SDE) DoCaching(c bool) {
	if c {
		s.generateCache()
	} else {
		s.cleanCache()
	}
}
func (s *SDE) lookupByTypeName(typeName string) (*SDEType, error) {
	if s.Cache != nil { // Fast lookup
		if v, ok := s.Cache.TypeNameLookup[typeName]; ok {
			return v, nil
		} else {
			return nil, ErrTypeDoesNotExist
		}
	}
	// Default to slow lookup if cache is nil

	for _, v := range s.Types {
		if v.TypeName == typeName {
			return v, nil
		}
	}
	return nil, ErrTypeDoesNotExist
}

// SDEType is a struct representing a single individual type in an SDE.
type SDEType struct {
	TypeID     int
	TypeName   string
	Attributes map[string]interface{}
	Parent     *SDE
}

// GetName returns the string value of Attributes["mDisplayName"] if it exists.
// Otherwise we return TypeName
func (s *SDEType) GetName() string {
	if v, ok := s.Attributes["mDisplayName"]; ok {
		return v.(string)
	}
	return s.TypeName
}

// GetAttribute checks if the type has the attribute and returns it.
// If it doesn't exist we lookup the weapons projectile type
func (s *SDEType) GetAttribute(attr string) interface{} {
	if v, ok := s.Attributes[attr]; ok {
		return v
	} else {
		if v, ok := s.Attributes["mFireMode0.projectileType"]; ok {
			v, _ := v.(int) // Ditch ok.  If it fails we get 0 and t will be nil
			if s.Parent == nil {
				log.Printf("GetAttributes can't lookup projectile.  Parent is nil\n")
				return nil
			}
			t, _ := s.Parent.GetType(v)
			// Ditching error because we don't return an error.
			// I don't want to break SDETool things yet
			if t == nil {
				log.Printf("Got nil type.  Returning nil\n")
				return nil
			}
			if v, ok := t.Attributes[attr]; ok {
				return v
			}
		}
	}
	return nil
}

// CompareTo prints the differences between two types
func (s *SDEType) CompareTo(t *SDEType) {
	if s.TypeID != t.TypeID {
		fmt.Printf("TYPEID CHANGE: %v: %v\n", s.TypeID, t.TypeID)
	}
	if s.TypeName != t.TypeName {
		fmt.Printf("TYPENAME CHANGE: %v: %v\n", s.TypeName, s.TypeName)
	}
	for key, value := range s.Attributes {
		if v, ok := t.Attributes[key]; ok {
			if value != v {
				fmt.Printf("CHANGE: %v: %v\n", value, v)
			}
		} else {
			fmt.Printf("ADD: %v: %v\n", key, value)
		}
	}
	for key, value := range t.Attributes {
		if _, ok := s.Attributes[key]; ok {
		} else {
			fmt.Printf("REMOVE: %v: %v\n", key, value)
		}
	}
}

/*
	Helpers
*/

func sdeslicecontains(s []*SDEType, tid int) bool {
	for _, v := range s {
		if v.TypeID == tid {
			return true
		}
	}
	return false
}

func init() {
	gob.Register(SDE{})
	gob.Register(SDEType{})
}
