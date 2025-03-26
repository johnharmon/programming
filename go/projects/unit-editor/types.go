package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	SeaFaring          = "sea_faring"           // can board ships;can_swim = can swim across rivers
	HideForest         = "hide_forest"          // defines where the unit can hide
	HideImprovedForest = "hide_improved_forest" // defines where the unit can hide
	HideAnywhere       = "hide_anywhere"        // defines where the unit can hide
	CanSap             = "can_sap"              // Can dig tunnels under walls
	FrightenFoot       = "frighten_foot"        // Cause fear to certain nearby unit types
	FrightenMounted    = "frighten_mounted"     // Cause fear to certain nearby unit types
	CanRunAmok         = "can_run_amok"         // Unit may go out of control when riders lose control of animals
	GeneralUnit        = "general_unit"         // The unit can be used for a named character's bodyguard
	CantabrianCircle   = "cantabrian_circle"    // The unit has this special ability
	NoCustom           = "no_custom"            // The unit may not be selected in custom battles
	Command            = "command"              // The unit carries a legionary eagle, and gives bonuses to nearby units
	MercenaryUnit      = "mercenary_unit"       // The unit is s mercenary unit available to all factions
	IsPeasant          = "is_peasant"           // unknown
	Druid              = "druid"                // Can do a special morale raising chant
	PowerCharge        = "power_charge"         // unkown
	FreeUpkeepUnit     = "free_upkeep_unit"     // Unit can be supported free in a city
)

//	var UnitAttributes = map[string]string{
//		"SeaFaring":          "sea_faring",           // can board ships;can_swim : can swim across rivers
//		"HideForest":         "hide_forest",          // defines where the unit can hide
//		"HideImprovedForest": "hide_improved_forest", // defines where the unit can hide
//		"HideAnywhere":       "hide_anywhere",        // defines where the unit can hide
//		"CanSap":             "can_sap",              // Can dig tunnels under walls
//		"FrightenFoot":       "frighten_foot",        // Cause fear to certain nearby unit types
//		"FrightenMounted":    "frighten_mounted",     // Cause fear to certain nearby unit types
//		"CanRunAmok":         "can_run_amok",         // Unit may go out of control when riders lose control of animals
//		"GeneralUnit":        "general_unit",         // The unit can be used for a named character's bodyguard
//		"CantabrianCircle":   "cantabrian_circle",    // The unit has this special ability
//		"NoCustom":           "no_custom",            // The unit may not be selected in custom battles
//		"Command":            "command",              // The unit carries a legionary eagle, and gives bonuses to nearby units
//		"MercenaryUnit":      "mercenary_unit",       // The unit is s mercenary unit available to all factions
//		"IsPeasant":          "is_peasant",           // unknown
//		"Druid":              "druid",                // Can do a special morale raising chant
//		"PowerCharge":        "power_charge",         // unkown
//		"FreeUpkeepUnit":     "free_upkeep_unit",     // Unit can be supported free in a city
//
// }
type UnitAttributes struct {
	SeaFaring          string `unit:"sea_faring"`           // can board ships;can_swim : can swim across rivers
	HideForest         string `unit:"hide_forest"`          // defines where the unit can hide
	HideImprovedForest string `unit:"hide_improved_forest"` // defines where the unit can hide
	HideAnywhere       string `unit:"hide_anywhere"`        // defines where the unit can hide
	CanSap             string `unit:"can_sap"`              // Can dig tunnels under walls
	FrightenFoot       string `unit:"frighten_foot"`        // Cause fear to certain nearby unit types
	FrightenMounted    string `unit:"frighten_mounted"`     // Cause fear to certain nearby unit types
	CanRunAmok         string `unit:"can_run_amok"`         // Unit may go out of control when riders lose control of animals
	GeneralUnit        string `unit:"general_unit"`         // The unit can be used for a named character's bodyguard
	CantabrianCircle   string `unit:"cantabrian_circle"`    // The unit has this special ability
	NoCustom           string `unit:"no_custom"`            // The unit may not be selected in custom battles
	Command            string `unit:"command"`              // The unit carries a legionary eagle, and gives bonuses to nearby units
	MercenaryUnit      string `unit:"mercenary_unit"`       // The unit is s mercenary unit available to all factions
	IsPeasant          string `unit:"is_peasant"`           // unknown
	Druid              string `unit:"druid"`                // Can do a special morale raising chant
	PowerCharge        string `unit:"power_charge"`         // unkown
	FreeUpkeepUnit     string `unit:"free_upkeep_unit"`     // Unit can be supported free in a city

}

type BoolAttribute struct {
	Value  bool
	String string
}

type Unit struct {
	Type                   string            `unit:"type"`
	Dictionary             string            `unit:"dictionary"`
	Class                  string            `unit:"class"`
	VoiceType              string            `unit:"voice_type"`
	Accent                 string            `unit:"accent"`
	BannerFaction          string            `unit:"banner_faction"`
	BannerHoly             string            `unit:"banner_holy"`
	Soldier                *Soldier          `unit:"soldier"`
	Officer                string            `unit:"officer"`
	MountEffect            *MountEffect      `unit:"mount_effect"`
	Attributes             []string          `unit:"attributes"`
	Formation              *Formation        `unit:"formation"`
	StatHealth             *Health           `unit:"stat_health"`
	StatPrimary            *Weapon           `unit:"stat_pri"`
	StatPrimaryAttribute   *WeaponAttributes `unit:"stat_pri_attr"`
	StatSecondary          *Weapon           `unit:"stat_sec"`
	StatSecondaryAttribute *WeaponAttributes `unit:"stat_sec_attr"`
	StatPrimaryArmor       *Armor            `unit:"Stat_pri_armor"`
	StatSecondaryArmor     *Armor            `unit:"Stat_sec_armor"`
	StatHeat               *Heat             `unit:"stat_heat"`
	StatGround             *Ground           `unit:"stat_ground"`
	StatMental             string            `unit:"stat_mental"`
	StatChargeDistance     int               `unit:"stat_charge_dist"`
	StatFireDelay          int               `unit:"stat_fire_delay"`
	StatFood               string            `unit:"stat_food"`
	StatCost               string            `unit:"stat_cost"`
	ArmorUpgradeLevels     []int             `unit:"armor_upgrade_levels"`
	ArmorUpgradeModels     []string          `unit:"armor_upgrade_models"`
	Ownership              string            `unit:"ownership"`
	RecruitPriorityOffset  int               `unit:"recruit_priority_offset"`
}

type Soldier struct {
	Name      string
	Number    int
	Extras    int
	Collision float64
}

type MountEffect struct {
	Horse              int `unit:"horse"`
	Camel              int `unit:"camel"`
	Elephant           int `unit:"elephant"`
	ElephantCannon     int `unit:"elephant_cannon"`
	SimpleHorse        int `unit:"simple horse"`
	MountLightWolf     int `unit:"mount_light_wolf"`
	WargCamel          int `unit:"warg_camel"`
	SwanGuardHorse     int `unit:"swan guard horse"`
	Eorlingas          int `unit:"eorlingas"`
	NorthernHeavyHorse int `unit:"northern heavy horse"`
}

// type Attributes struct{}
type Formation struct {
	SidetoSideSpacingTight  float64
	FronttoBackSpacingTight float64
	SidetoSideSpacingLoose  float64
	FronttoBackSpacingLoose float64
	PossibleFormations      []string
}
type Health struct {
	HP      int
	unknown string
}
type Weapon struct {
	Attack             int
	Charge             int
	MissileType        string
	MissileRange       int
	MissileAmmo        int
	WeaponType         string
	TechType           string
	DamageType         string
	SoundType          string
	FireEffect         string
	MinDelay           int
	CompensationFactor int
}

func (w *Weapon) Unmarshal(weaponInfo string) error {
	re := regexp.MustCompile(`\s+`)
	lineSections := strings.SplitN(re.ReplaceAllString(weaponInfo, " "), " ", 2)
	weaponStats := strings.Split(lineSections[1], ",")
	numFields := len(weaponStats)
	if numFields < 11 {
		return fmt.Errorf("error parsing attack stats, too few fields")
	}
	switch numFields {
	case 11:
		w.MinDelay, _ = strconv.Atoi(strings.TrimSpace(weaponStats[10]))
		w.CompensationFactor, _ = strconv.Atoi(strings.TrimSpace(weaponStats[11]))
	default:
		w.FireEffect = weaponStats[9]
		w.MinDelay, _ = strconv.Atoi(strings.TrimSpace(weaponStats[10]))
		w.CompensationFactor, _ = strconv.Atoi(strings.TrimSpace(weaponStats[11]))
	}
	w.Attack, _ = strconv.Atoi(strings.TrimSpace(weaponStats[0]))
	w.Charge, _ = strconv.Atoi(strings.TrimSpace(weaponStats[1]))
	w.MissileType = weaponStats[2]
	w.MissileRange, _ = strconv.Atoi(strings.TrimSpace(weaponStats[3]))
	w.MissileAmmo, _ = strconv.Atoi(strings.TrimSpace(weaponStats[4]))
	w.WeaponType = weaponStats[5]
	w.TechType = weaponStats[6]
	w.DamageType = weaponStats[7]
	w.SoundType = weaponStats[8]
	return nil
}

type WeaponAttributes struct {
	AP         BoolAttribute `unit:"ap"`            // armour piercing. Only counts half of target's armour
	BP         BoolAttribute `unit:"bp"`            // body piercing. Missile can pass through men and hit those behind
	Spear      BoolAttribute `unit:"spear"`         // Used for long spears. Gives bonuses fighting cavalry, and penalties against infantry
	LongPike   BoolAttribute `unit:"long_pike"`     // Use very long pikes. Phalanx capable units only
	ShortPike  BoolAttribute `unit:"short_pike"`    // Use shorter than normal spears.
	Prec       BoolAttribute `unit:"prec"`          // Missile weapon is only thrown/ fired just before charging into combat
	Thrown     BoolAttribute `unit:"thrown"`        // The missile type if thrown rather than fired
	launching  BoolAttribute `unit:"launching"`     // attack may throw target men into the air
	Area       BoolAttribute `unit:"area"`          // attack affects an area, not just one man
	LightSpear BoolAttribute `unit:"light_spear"`   // The unit when braced has various protecting mechanisms versus cavalry charges from the frontk
	SpearBonus BoolAttribute `unit:"spear_bonus_x"` // attack bonus against cavalry. x = 2, 4, 6, 8, 10 or 12

}

type Armor struct {
	Armor        int
	DefenseSkill int
	Shield       int
	Sound        string
}

type ArmorEx struct {
}
type Heat struct{}
type Ground struct{}
type Mental struct{}
type Food struct{}
type Cost struct{}
