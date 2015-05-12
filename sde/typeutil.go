package sde

import (
	"fmt"
	"github.com/THUNDERGROOVE/SDETool/util/log"
	"reflect"
	"strings"
	"time"
)

type Modifier struct {
	StackingPenalized string
	AttributeName     string
	ModifierValue     float64
	ModifierType      string
}

func ApplyTypeToType(tone, ttwo SDEType) (*SDEType, error) {
	defer Debug(time.Now())
	out := &tone

	mods := make(map[string]Modifier, 0)
	for k, _ := range ttwo.Attributes {
		if strings.Split(k, ".")[0] == "modifier" {
			index := strings.Split(k, ".")[1]
			if _, ok := mods[index]; !ok {
				mod := Modifier{}

				mod.StackingPenalized, _ = ttwo.Attributes["modifier."+index+".stackingPenalized"].(string)
				mod.AttributeName, _ = ttwo.Attributes["modifier."+index+".attributeName"].(string)
				mod.ModifierType, _ = ttwo.Attributes["modifier."+index+".modifierType"].(string)
				mod.ModifierValue, _ = ttwo.Attributes["modifier."+index+".modifierValue"].(float64)
				mods[index] = mod
			} else {
				log.Println("Skipping attribute; Already indexed")
			}
		}
	}

	for _, v := range mods {
		log.Printf("Applying %v: %v to %v with mode %v\n", v.AttributeName, v.ModifierValue, out.GetName(), v.ModifierType)
		switch v.ModifierType {
		case "ADD":
			out.Attributes[v.AttributeName] = modAdd(out.Attributes[v.AttributeName], v.ModifierValue)
		case "MULTIPLY":
			out.Attributes[v.AttributeName] = modMult(out.Attributes[v.AttributeName], v.ModifierValue)
		default:
			log.Println("Unsupported modifier:", v.ModifierType)
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
		log.Println("modAdd called with nil interface.")
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
		log.Println("modMult called with nil interface.")
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
