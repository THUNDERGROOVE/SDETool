package sde

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/THUNDERGROOVE/SDETool/util/log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// SDEType holds hopefully all of the information you will need about a type.
// The interface in the Attributes map will always be either a float64, int or
// a string.  If the value is always going to a whole number pull an int out
// otherwise assume it's a float
type SDEType struct {
	parentSDE  *SDE                   `json:"parentSDE"`
	TypeID     int                    `json:"typeId"`
	TypeName   string                 `json:"typeName"`
	Attributes map[string]interface{} `json:"attributes"`
}

func (s *SDEType) ParentSDE() *SDE {
	return s.parentSDE
}

// GetAttributes grabs the attributes for the type and applied them.  This is
// used to speed up querries for simple lookups.
func (s *SDEType) GetAttributes() error {
	defer Debug(time.Now())
	if s == nil {
		log.Println("SDEType is nil")
		return errors.New("Type is nil")
	}
	if s.parentSDE == nil {
		if PrimarySDE == nil {
			log.Println("Primary SDE not set.  Returning error")
			return errors.New("No parent SDE set")
		}
		s.parentSDE = PrimarySDE
	}

	if s.parentSDE == nil {
		if PrimarySDE == nil {
			log.Println("Primary SDE not set.  Returning error")
			return errors.New("No parent SDE set")
		}
		s.parentSDE = PrimarySDE
	}

	rows, err := s.parentSDE.DB.Query(fmt.Sprintf("SELECT catmaAttributeName, catmaValueInt, catmaValueReal, catmaValueText FROM CatmaAttributes WHERE TypeID == '%v'", s.TypeID))
	if err != nil {
		return err
	}
	for rows.Next() {
		var catmaAttributeName string
		var catmaValueInt string
		var catmaValueReal string
		var catmaValueText string

		rows.Scan(&catmaAttributeName, &catmaValueInt, &catmaValueReal, &catmaValueText)
		if catmaValueInt != "None" {
			v, _ := strconv.Atoi(catmaValueInt)
			s.Attributes[catmaAttributeName] = v
		}
		if catmaValueReal != "None" {
			v, _ := strconv.ParseFloat(catmaValueReal, 64)
			s.Attributes[catmaAttributeName] = v
		}
		if catmaValueText != "None" {
			s.Attributes[catmaAttributeName] = catmaValueText
		}
	}
	return nil
}

// GetName returns the display name of a type.
func (s SDEType) GetName() string {
	if name, ok := s.Attributes["mDisplayName"]; ok {
		if v, kk := name.(string); kk {
			return v
		}
		log.Println("mDisplayName not string?", reflect.TypeOf(name), "typeID:", s.TypeID, "typeName:", s.TypeName)
	}

	return s.TypeName
}

// GetRoF gets the ROF of any weapon in rounds per minute.  Must vall GetAttributes first
func (s *SDEType) GetRoF() int {
	defer Debug(time.Now())

	// if v, ok := s.Attributes["mFireMode0.m_eFireMode"]; ok {
	// 	if v == "DWFM_SingleBurst" {
	bi := s.Attributes["m_BurstInfo.m_fBurstInterval"].(float64)
	fi := s.Attributes["mFireMode0.fireInterval"].(float64)
	return int((bi + fi + 0.01) * 10000)
	// 	}
	// }
	// if i, ok := s.Attributes["mFireMode0.fireInterval"]; ok {
	// 	interval := i.(float64)
	// 	return int(interval * 10000)
	// }
	return 0
}

// GetDPS returns the DPS of a type, if it can.
// Notice: CCP has some all kinds of fucked up shit with bursts and intervals
// don't expect these numbers to be accurate until I can finally fix all of it.
func (s *SDEType) GetDPS() float64 {
	defer Debug(time.Now())

	RoF := s.GetRoF()
	var damage float64
	if d, ok := s.Attributes["mFireMode0.instantHitDamage"]; ok {
		damage = d.(float64)
	} else {
		damage = 0
	}

	fmt.Printf("RoF: %v\n", RoF)
	return float64((damage * float64(RoF)) / 60)
}

// IsWeapon returns true if a type has a weapon tag.
func (s *SDEType) IsWeapon() bool {
	defer Debug(time.Now())

	for k, v := range s.Attributes {
		if strings.Contains(k, "tag.") {
			if v.(int) == 352335 {
				return true
			}
		}
	}
	return false
}

// IsAurum returns if the item is puchased with aurum.
// Be the soldiar of tomorrow, today with Aurum(C)(TM)(LOLCCP)
// Will also return true if it's a pack item or special edition.
// The only special edition suits that aren't filtered are unique.
func (s *SDEType) IsAurum() bool {
	defer Debug(time.Now())

	if strings.Contains(s.TypeName, "aur") || strings.Contains(s.TypeName, "promo") || strings.Contains(s.TypeName, "fit_") || strings.Contains(s.TypeName, "pack") || strings.Contains(s.TypeName, "harbinger") {
		return true
	}

	return false
}

// IsObtainable returns True if the item is consumable.
// The name is misleading but it should be used to check if an item is
// obtainable by a player.
func (s *SDEType) IsObtainable() bool {
	defer Debug(time.Now())

	if _, ok := s.Attributes["consumable"]; ok {
		return true
	}
	return false
}

// getFromTags is an internal method to return all types that share have a tag
// of the type provided.
func (s *SDEType) getFromTags(t SDEType) ([]*SDEType, error) {
	defer Debug(time.Now())

	types := make([]*SDEType, 0)
	rows, err := s.parentSDE.DB.Query(fmt.Sprintf("SELECT typeID FROM CatmaAttributes WHERE catmaValueInt == '%v';", t.TypeID))
	if err != nil {
		return types, err
	}
	for rows.Next() {
		var nTypeID int
		rows.Scan(&nTypeID)
		types = append(types, &SDEType{
			s.parentSDE,
			nTypeID,
			"",
			make(map[string]interface{})})
	}
	return types, nil
}

// HasTag returns true if the type contains a tag attribute by the tagID given.
func (s *SDEType) HasTag(tag int) bool {
	for k, v := range s.Attributes {
		if strings.Contains(k, "tag.") {
			if v, ok := v.(int); ok {
				if v == tag {
					return true
				}
			}
		}
	}
	return false
}

// GetSharedTagTypes returns a slice of SDETypes that share the 'main' tag
// of a type.
func (s *SDEType) GetSharedTagTypes() ([]*SDEType, error) {
	defer Debug(time.Now())

	types := make([]*SDEType, 0)
	if s.IsWeapon() {
		for k, v := range s.Attributes {
			if strings.Contains(k, "tag.") {
				tag, _ := s.parentSDE.GetType(v.(int))
				tag.GetAttributes()
				if strings.Contains(tag.TypeName, "tag_weapon_") { // if s is a scrambler rifle, return all scrambler rifles.
					types, err := s.getFromTags(tag)
					return types, err
				}
				if strings.Contains(tag.TypeName, "tag_core") { // Return all dropsuits since.
					types, err := s.getFromTags(tag)
					return types, err
				}
				if strings.Contains(tag.TypeName, "tag_module_") {
					types, err := s.getFromTags(tag)
					return types, err
				}
				if strings.Contains(tag.TypeName, "tag_amarr") {
					types, err := s.getFromTags(tag)
					return types, err
				}
				if strings.Contains(tag.TypeName, "tag_caldari") {
					types, err := s.getFromTags(tag)
					return types, err
				}
				if strings.Contains(tag.TypeName, "tag_gallente") {
					types, err := s.getFromTags(tag)
					return types, err
				}
				if strings.Contains(tag.TypeName, "tag_minmatar") {
					types, err := s.getFromTags(tag)
					return types, err
				}
			}
		}
	}
	return types, nil
}

// ToJSON returns a Marshaled and indented version of our SDEType.
func (s *SDEType) ToJSON() (string, error) {
	defer Debug(time.Now())

	v, err := json.MarshalIndent(s, "", "    ")
	return string(v), err
}

// PrintTags is a function to pretty print tags related to a type.
func (s *SDEType) PrintTags() {
	for k, v := range s.Attributes {
		if strings.Contains(k, "tag.") {
			tag, _ := s.parentSDE.GetType(v.(int))
			tag.GetAttributes()
			fmt.Println("-> Type has tag:", tag.GetName())
		}
	}
}

func (s *SDEType) IsFaction() bool {
	defer Debug(time.Now())
	if strings.Contains(s.GetName(), "Imperial") || strings.Contains(s.GetName(), "State") || strings.Contains(s.GetName(), "Federation") || strings.Contains(s.GetName(), "Republic") {
		return true
	}
	return false
}

// Lookup goes through our attributes in search of attributes that may
// be a TypeID and change it's value to an *SDEType.
func (s *SDEType) Lookup(depth int) {
	log.Println("Lookup() depth", depth)
	if depth <= 0 {
		return
	}
	for key, v := range s.Attributes {
		log.Println("Lookup checking", key)
		if i, ok := v.(int); ok {
			is := strconv.Itoa(i)
			if len(is) == 6 { // Must be a TypeID right? lol
				log.Println("Found TypeID candidate", i)
				t, err := s.parentSDE.GetType(i)
				if err != nil {
					log.Println(err.Error())
					continue
				}
				t.GetAttributes()
				t.Lookup(depth - 1)
				s.Attributes[key] = &t
			} else {
				log.Println("Attribute had len of", len(is))
			}
		}
	}
}

func (s SDEType) GetHighSlots() int {
	return s.slotFinder("IH")
}
func (s SDEType) GetLowSlots() int {
	return s.slotFinder("IL")
}
func (s SDEType) GetGrenadeSlots() int {
	return s.slotFinder("GS")
}
func (s SDEType) GetEquipmentSlots() int {
	return s.slotFinder("IE")
}
func (s SDEType) GetSidearmSlots() int {
	return s.slotFinder("WS")
}
func (s SDEType) GetPrimarySlots() int {
	return s.slotFinder("WP")
}
func (s SDEType) GetHeavySlots() int {
	return s.slotFinder("WH")
}

func (s *SDEType) slotFinder(slotType string) int {
	var c int
	for k, v := range s.Attributes {
		fmt.Println(k)
		if strings.Contains(k, "mModuleSlots.") {
			val := strings.Split(k, ".")
			if val[2] == "slotType" {
				if t, ok := v.(string); ok {
					if t == slotType {
						c += 1
					}
				}
			}
		}
	}
	return c
}

// Not documented for a reason. Don't ask.  Pretend this doesn't exist
func (s *SDEType) ESBA() {
	fmt.Println("ESBA")
	SBM := 342

	for i := 0; i < SBM; i++ {
		if s.Attributes[fmt.Sprintf("mSalvageBoxLootTable.%v.itemTypeSet", i)] == nil || s.Attributes[fmt.Sprintf("mSalvageBoxLootTable.%v.lootFreq", i)] == nil || s.Attributes[fmt.Sprintf("mSalvageBoxLootTable.%v.minQuantity", i)] == nil || s.Attributes[fmt.Sprintf("mSalvageBoxLootTable.%v.maxQuantity", i)] == nil {
			continue
		}
		id := s.Attributes[fmt.Sprintf("mSalvageBoxLootTable.%v.itemTypeSet", i)].(int)
		f := s.Attributes[fmt.Sprintf("mSalvageBoxLootTable.%v.lootFreq", i)].(int)
		mi := s.Attributes[fmt.Sprintf("mSalvageBoxLootTable.%v.minQuantity", i)].(int)
		ma := s.Attributes[fmt.Sprintf("mSalvageBoxLootTable.%v.maxQuantity", i)].(int)
		t, _ := s.parentSDE.GetType(id)
		t.GetAttributes()
		fmt.Printf("%v with frequency %v at least %v but no more than %v\n", t.GetName(), f, mi, ma)
	}
}

func (s *SDEType) genTemplateSlice(i int) []string {
	return make([]string, i-1)
}
