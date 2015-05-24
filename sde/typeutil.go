package sde

import (
	"fmt"
	"reflect"
	"strings"
)

/*
	Modifier represents a modifier that a Type has in it's attributes
*/
type Modifier struct {
	StackingPenalized string
	AttributeName     string
	ModifierValue     float64
	ModifierType      string
}

// ApplyTypeToType is a function that applies all the attributes given by tone to ttwo and returns the results as a pointer
// For now it doesn't work for skills.  I'll  likly write a different function for it.
func ApplyTypeToType(tone, ttwo SDEType) (*SDEType, error) {
	out := &tone

	mods := make(map[string]Modifier, 0)
	for k, _ := range ttwo.Attributes {
		if strings.Split(k, ".")[0] == "modifier" {
			index := strings.Split(k, ".")[1]
			if _, ok := mods[index]; !ok {
				// @TODO: Refactor this into it's own function
				mod := Modifier{}

				mod.StackingPenalized, _ = ttwo.Attributes["modifier."+index+".stackingPenalized"].(string)
				mod.AttributeName, _ = ttwo.Attributes["modifier."+index+".attributeName"].(string)
				mod.ModifierType, _ = ttwo.Attributes["modifier."+index+".modifierType"].(string)
				mod.ModifierValue, _ = ttwo.Attributes["modifier."+index+".modifierValue"].(float64)
				mods[index] = mod
			}
		}
	}

	for _, v := range mods {
		// @TODO: Refactor into it's own function
		switch v.ModifierType {
		case "ADD":
			out.Attributes[v.AttributeName] = modAdd(out.Attributes[v.AttributeName], v.ModifierValue)
		case "MULTIPLY":
			out.Attributes[v.AttributeName] = modMult(out.Attributes[v.AttributeName], v.ModifierValue)
		default:
			return nil, fmt.Errorf("Unsupported modifier: '%v'", v.ModifierType)
		}
	}

	return out, nil

}

// Generics are for pussies.
func modAdd(i interface{}, ii float64) interface{} {
	/*	if reflect.TypeOf(i) != reflect.TypeOf(ii) {
		log.Println("modAdd called with mismatched types.  WHY?", reflect.TypeOf(i), reflect.TypeOf(ii))
		return i
	}*/
	if i == nil {
		return nil
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Int:
		ic := i.(int)
		return interface{}(int(float64(ic) + ii))
	case reflect.Float64:
		ic := i.(float64)
		return interface{}(ic + ii)
	default:
		fmt.Println("modAdd called with unsupported type", reflect.TypeOf(i).Kind())
	}
	return nil
}

// Generics are for pussies.
func modMult(i interface{}, ii float64) interface{} {
	/*	if reflect.TypeOf(i) != reflect.TypeOf(ii) {
		log.Println("modMult called with mismatched types.  WHY?", reflect.TypeOf(i), reflect.TypeOf(ii))
		return i
	}*/
	if i == nil {
		return nil
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Int:
		ic := i.(int)
		return interface{}(int(float64(ic) * ii))
	case reflect.Float64:
		ic := i.(float64)
		return interface{}(ic * ii)
	default:
		fmt.Println("modMult called with unsupported type", reflect.TypeOf(i).Kind())
	}
	return nil
}
