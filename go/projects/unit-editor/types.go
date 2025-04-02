package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	DefaultErrorFormat = "Line: %d | error converting %s value of %s to %s: %s\n"
	DefaultInfoFormat  = "Line: %d | Attribute: \"%s\" | Position: %d | Converted \"%s\" to %s\n"
)

type LogLevel int

const (
	LevelNone LogLevel = iota
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
)

type logger interface {
	SDebugf(format string, fields ...any) string
	SInfof(format string, fields ...any) string
	SWarnf(format string, fields ...any) string
	SErrorf(format string, fields ...any) string
	FDebugf(format string, fields ...any)
	FInfof(format string, fields ...any)
	FWarnf(format string, fields ...any)
	FErrorf(format string, fields ...any)
}

func DebugLogger(w io.Writer) *UnitLogger {
	ul := &UnitLogger{
		debugStream: w,
		infoStream:  w,
		warnStream:  w,
		errorStream: w,
	}
	return ul
}

type UnitLogger struct {
	Unit               *Unit
	logLevel           LogLevel
	debugStream        io.Writer
	infoStream         io.Writer
	warnStream         io.Writer
	errorStream        io.Writer
	defaultDebugFormat string
	defaultInfoFormat  string
	defaultWarnFormat  string
	defaultErrorFormat string
}

func NewUnitLogger(logLevel LogLevel, output io.Writer, discard io.Writer) (ul *UnitLogger) {
	ul = &UnitLogger{logLevel: logLevel}
	for level := 0; level <= int(logLevel); level++ {
		switch {
		case level == 1 && level <= int(logLevel):
			ul.errorStream = output
		case level == 2 && level <= int(logLevel):
			ul.warnStream = output
		case level == 3 && level <= int(logLevel):
			ul.infoStream = output
		case level == 3 && level <= int(logLevel):
			ul.debugStream = output
		}
	}
	return ul
}

func (ul UnitLogger) SDebugf(format string, fields ...any) string {
	return fmt.Sprintf(format, fields...)
}
func (ul UnitLogger) SInfof(format string, fields ...any) string {
	return fmt.Sprintf(format, fields...)
}
func (ul UnitLogger) SWarnf(format string, fields ...any) string {
	return fmt.Sprintf(format, fields...)
}
func (ul UnitLogger) SErrorf(format string, fields ...any) string {
	return fmt.Sprintf(format, fields...)
}

func (ul *UnitLogger) FDebugf(format string, fields ...any) {
	fmt.Fprintf(ul.debugStream, format, fields...)
}
func (ul *UnitLogger) FInfof(format string, fields ...any) {
	fmt.Fprintf(ul.infoStream, format, fields...)
}
func (ul *UnitLogger) FWarnf(format string, fields ...any) {
	fmt.Fprintf(ul.warnStream, format, fields...)
}
func (ul *UnitLogger) FErrorf(format string, fields ...any) {
	fmt.Fprintf(ul.errorStream, format, fields...)
}

func ParseModifier(modifier string) (int, error) {
	// parses a string: '+/-'<int> into an actual integer
	switch modifier[0] {
	case '-':
		result, err := strconv.Atoi(modifier[1:])
		if err != nil {
			return 0, err
		}
		return 0 - result, nil
	default:
		result, err := strconv.Atoi(modifier[1:])
		if err != nil {
			return 0, err
		}
		return 0 + result, nil
	}

}

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

func ProcessUnit()

type UnitField interface {
	Marshal() string
}

type LineRecord struct {
	LineNumber int
	Raw        string
	Unit       *Unit
	FieldValue UnitField
	FieldName  string
	Comment    bool
	Empty      bool
}

func ParseLineRecord(lr *LineRecord) (err error) {
	return err
}

func UnmarshalLineRecord(line string, lineNumber int, unit *Unit) (lr *LineRecord) {
	lr = &LineRecord{}
	lr.LineNumber = lineNumber
	lr.Raw = line
	lr.Unit = unit
	err := ParseLineRecord(lr)
	if err != nil {

	}
	return lr
}

func (lr *LineRecord) Unmarshal(line string, lineNumber int) {
	lr.LineNumber = lineNumber
	lr.Raw = line

}

type UnitMetadata struct {
	Type      string
	LineStart int
	LineEnd   int
}

type UnitLog struct {
	Logs     []string
	RawLogs  []byte
	Metadata *UnitMetadata
}

type UnitAttributes struct {
	SeaFaring          string `unit:"sea_faring" json:"sea_faring"`                     // can board ships;can_swim : can swim across rivers
	HideForest         string `unit:"hide_forest" json:"hide_forest"`                   // defines where the unit can hide
	HideImprovedForest string `unit:"hide_improved_forest" json:"hide_improved_forest"` // defines where the unit can hide
	HideAnywhere       string `unit:"hide_anywhere" json:"hide_anywhere"`               // defines where the unit can hide
	CanSap             string `unit:"can_sap" json:"can_sap"`                           // Can dig tunnels under walls
	FrightenFoot       string `unit:"frighten_foot" json:"frighten_foot"`               // Cause fear to certain nearby unit types
	FrightenMounted    string `unit:"frighten_mounted" json:"frighten_mounted"`         // Cause fear to certain nearby unit types
	CanRunAmok         string `unit:"can_run_amok" json:"can_run_amok"`                 // Unit may go out of control when riders lose control of animals
	GeneralUnit        string `unit:"general_unit" json:"general_unit"`                 // The unit can be used for a named character's bodyguard
	CantabrianCircle   string `unit:"cantabrian_circle" json:"cantabrian_circle"`       // The unit has this special ability
	NoCustom           string `unit:"no_custom" json:"no_custom"`                       // The unit may not be selected in custom battles
	Command            string `unit:"command" json:"command"`                           // The unit carries a legionary eagle, and gives bonuses to nearby units
	MercenaryUnit      string `unit:"mercenary_unit" json:"mercenary_unit"`             // The unit is s mercenary unit available to all factions
	IsPeasant          string `unit:"is_peasant" json:"is_peasant"`                     // unknown
	Druid              string `unit:"druid" json:"druid"`                               // Can do a special morale raising chant
	PowerCharge        string `unit:"power_charge" json:"power_charge"`                 // unkown
	FreeUpkeepUnit     string `unit:"free_upkeep_unit" json:"free_upkeep_unit"`         // Unit can be supported free in a city

}

type BoolAttribute struct {
	Value  bool
	String string
}

type Unit struct {
	Logger                 *UnitLogger
	LineRecords            []*LineRecord
	Lines                  []string
	Type                   string            `unit:"type" json:"type"`
	Dictionary             string            `unit:"dictionary" json:"dictionary"`
	Class                  string            `unit:"class" json:"class"`
	VoiceType              string            `unit:"voice_type" json:"voice_type"`
	Accent                 string            `unit:"accent" json:"accent"`
	BannerFaction          string            `unit:"banner_faction" json:"banner_faction"`
	BannerHoly             string            `unit:"banner_holy" json:"banner_holy"`
	Soldier                *Soldier          `unit:"soldier" json:"soldier"`
	Officer                string            `unit:"officer" json:"officer"`
	MountEffect            *MountEffect      `unit:"mount_effect" json:"mount_effect"`
	Attributes             []string          `unit:"attributes" json:"attributes"`
	Formation              *Formation        `unit:"formation" json:"formation"`
	StatHealth             *Health           `unit:"stat_health" json:"stat_health"`
	StatPrimary            *Weapon           `unit:"stat_pri" json:"stat_pri"`
	StatPrimaryAttribute   *WeaponAttributes `unit:"stat_pri_attr" json:"stat_pri_attr"`
	StatSecondary          *Weapon           `unit:"stat_sec" json:"stat_sec"`
	StatSecondaryAttribute *WeaponAttributes `unit:"stat_sec_attr" json:"stat_sec_attr"`
	StatPrimaryArmor       *Armor            `unit:"Stat_pri_armor" json:"Stat_pri_armor"`
	StatSecondaryArmor     *Armor            `unit:"Stat_sec_armor" json:"Stat_sec_armor"`
	StatHeat               *Heat             `unit:"stat_heat" json:"stat_heat"`
	StatGround             *Ground           `unit:"stat_ground" json:"stat_ground"`
	StatMental             string            `unit:"stat_mental" json:"stat_mental"`
	StatChargeDistance     int               `unit:"stat_charge_dist" json:"stat_charge_dist"`
	StatFireDelay          int               `unit:"stat_fire_delay" json:"stat_fire_delay"`
	StatFood               string            `unit:"stat_food" json:"stat_food"`
	StatCost               string            `unit:"stat_cost" json:"stat_cost"`
	ArmorUpgradeLevels     []int             `unit:"armor_upgrade_levels" json:"armor_upgrade_levels"`
	ArmorUpgradeModels     []string          `unit:"armor_upgrade_models" json:"armor_upgrade_models"`
	Ownership              string            `unit:"ownership" json:"ownership"`
	RecruitPriorityOffset  int               `unit:"recruit_priority_offset" json:"recruit_priority_offset"`
}

type Soldier struct {
	Name      string
	Number    int
	Extras    int
	Collision float64
}

type MountEffect struct {
	Effects            map[string]int
	Horse              int `unit:"horse json:"horse"`
	Camel              int `unit:"camel" json:"camel"`
	Elephant           int `unit:"elephant" json:"elephant"`
	ElephantCannon     int `unit:"elephant_cannon" json:"elephant_cannon"`
	SimpleHorse        int `unit:"simple horse" json:"simple horse"`
	MountLightWolf     int `unit:"mount_light_wolf" json:"mount_light_wolf"`
	WargCamel          int `unit:"warg_camel" json:"warg_camel"`
	SwanGuardHorse     int `unit:"swan guard horse" json:"swan guard horse"`
	Eorlingas          int `unit:"eorlingas" json:"eorlingas"`
	NorthernHeavyHorse int `unit:"northern heavy horse" json:"northern heavy horse"`
}

func (me *MountEffect) Unmarshal(effectInfo string) error {
	lineSections := CleanLine(effectInfo)
	effectStats := strings.Split(lineSections[1], ",")
	for _, effect := range effectStats {
		effects := strings.SplitN(effect, " ", 2)
		effectKey := effects[0]
		effectValue := effects[2]
		effectInt, _ := ParseModifier(effectValue)
		me.Effects[effectKey] = effectInt
	}
	return nil
}

// type Attributes struct{}
type Formation struct {
	SidetoSideSpacingTight  float64  `json:"sideTight"`
	FronttoBackSpacingTight float64  `json:"frontTight"`
	SidetoSideSpacingLoose  float64  `json:"sideLoose"`
	FronttoBackSpacingLoose float64  `json:"frontLoose"`
	DefaultRanks            int      `json:"defaultRanks"`
	PossibleFormations      []string `json:"possibleFormations"`
}

func (f *Formation) Unmarshal(formationInfo string, ul *UnitLogger, lr *LineRecord) (fieldErrors error) {
	formationStats, numFields, fieldErrors := GetFieldInfo(formationInfo, 6, ul, lr)
	if numFields < 6 {
		return fmt.Errorf("error, insufficient number of fields for formation")
	}
	fieldErrors = errors.Join(fieldErrors, CheckSetFloatAttribute(&f.SidetoSideSpacingTight, formationStats[0], "sideToSideTight", 0, ul, DefaultErrorFormat, DefaultInfoFormat, lr.LineNumber))
	fieldErrors = errors.Join(fieldErrors, CheckSetFloatAttribute(&f.FronttoBackSpacingTight, formationStats[1], "frontToBackTight", 1, ul, DefaultErrorFormat, DefaultInfoFormat, lr.LineNumber))
	fieldErrors = errors.Join(fieldErrors, CheckSetFloatAttribute(&f.SidetoSideSpacingLoose, formationStats[2], "sideToSideLoose", 2, ul, DefaultErrorFormat, DefaultInfoFormat, lr.LineNumber))
	fieldErrors = errors.Join(fieldErrors, CheckSetFloatAttribute(&f.FronttoBackSpacingLoose, formationStats[3], "frontToBackLoose", 3, ul, DefaultErrorFormat, DefaultInfoFormat, lr.LineNumber))
	fieldErrors = errors.Join(fieldErrors, CheckSetIntAttribute(&f.DefaultRanks, formationStats[3], "frontToBackLoose", 3, ul, DefaultErrorFormat, DefaultInfoFormat, lr.LineNumber))
	f.PossibleFormations = append(f.PossibleFormations, formationStats[5])
	if numFields > 6 {
		f.PossibleFormations = append(f.PossibleFormations, formationStats[6])
	}
	return nil
}

type Health struct {
	HP      int
	MountHP int
}

func (h *Health) Unmarshal(healthInfo string, ul *UnitLogger, lr *LineRecord) error {
	lineSections := CleanLine(healthInfo)
	healthStats := strings.Split(lineSections[1], ",")
	_ = CheckSetIntAttribute(&h.HP, healthStats[0], "Health", 0, ul, DefaultErrorFormat, DefaultInfoFormat, lr.LineNumber)
	_ = CheckSetIntAttribute(&h.MountHP, healthStats[0], "Health", 0, ul, DefaultErrorFormat, DefaultInfoFormat, lr.LineNumber)
	return nil
}

type Weapon struct {
	Attack             int    `json:"attack"`
	Charge             int    `json:"charge"`
	MissileType        string `json:"missile_type"`
	MissileRange       int    `json:"missile_range"`
	MissileAmmo        int    `json:"missile_ammo"`
	WeaponType         string `json:"weapon_type"`
	TechType           string `json:"tech_type"`
	DamageType         string `json:"damage_type"`
	SoundType          string `json:"sound_type"`
	FireEffect         string `json:"fire_effect"`
	MinDelay           int    `json:"min_delay"`
	CompensationFactor int    `json:"compensation_factor"`
}

func (w *Weapon) Unmarshal(weaponInfo string, ul *UnitLogger, lr *LineRecord) (fieldErrors error) {
	lineSections := CleanLine(weaponInfo)
	weaponStats := strings.Split(lineSections[1], ",")
	numFields := len(weaponStats)
	if numFields < 11 {
		ul.FErrorf("error parsing attack stats, too few fields")
	}
	weaponStats = TrimValues(weaponStats)
	for index, value := range weaponStats {
		switch index {
		case 0:
			err := CheckSetIntAttribute(&w.Attack, value, "Attack", index, ul, DefaultErrorFormat, DefaultInfoFormat, lr.LineNumber)
			if err != nil {
				fieldErrors = errors.Join(fieldErrors, err)
			}
		case 1:
			err := CheckSetIntAttribute(&w.Charge, value, "Charge", index, ul, DefaultErrorFormat, DefaultInfoFormat, lr.LineNumber)
			if err != nil {
				fieldErrors = errors.Join(fieldErrors, err)
			}
		case 2:
			SetStrAttribute(&w.MissileType, value, "MissileType", ul, lr.LineNumber)
		case 3:
			err := CheckSetIntAttribute(&w.MissileRange, value, "MissileRange", index, ul, DefaultErrorFormat, DefaultInfoFormat, lr.LineNumber)
			if err != nil {
				fieldErrors = errors.Join(fieldErrors, err)
			}
		case 4:
			err := CheckSetIntAttribute(&w.MissileAmmo, value, "MissileAmmo", index, ul, DefaultErrorFormat, DefaultInfoFormat, lr.LineNumber)
			if err != nil {
				fieldErrors = errors.Join(fieldErrors, err)
			}
		case 5:
			SetStrAttribute(&w.WeaponType, value, "WeaponType", ul, lr.LineNumber)
		case 6:
			SetStrAttribute(&w.TechType, value, "TechType", ul, lr.LineNumber)
		case 7:
			SetStrAttribute(&w.DamageType, value, "DamageType", ul, lr.LineNumber)
		case 8:
			SetStrAttribute(&w.SoundType, value, "SoundType", ul, lr.LineNumber)
		case 9:
			if numFields == 11 {
				err := CheckSetIntAttribute(&w.MinDelay, value, "MinDelay", index, ul, DefaultErrorFormat, DefaultInfoFormat, lr.LineNumber)
				if err != nil {
					fieldErrors = errors.Join(fieldErrors, err)
				}
			} else {
				w.FireEffect = value
				SetStrAttribute(&w.FireEffect, value, "FireEffect", ul, lr.LineNumber)
			}
		case 10:
			if numFields == 11 {
				err := CheckSetIntAttribute(&w.CompensationFactor, value, "CompensationFactor", index, ul, DefaultErrorFormat, DefaultInfoFormat, lr.LineNumber)
				if err != nil {
					fieldErrors = errors.Join(fieldErrors, err)
				}
			} else {
				err := CheckSetIntAttribute(&w.MinDelay, value, "MinDelay", index, ul, DefaultErrorFormat, DefaultInfoFormat, lr.LineNumber)
				if err != nil {
					fieldErrors = errors.Join(fieldErrors, err)
				}
			}
		case 11:
			err := CheckSetIntAttribute(&w.CompensationFactor, value, "CompensationFactor", index, ul, DefaultErrorFormat, DefaultInfoFormat, lr.LineNumber)
			if err != nil {
				fieldErrors = errors.Join(fieldErrors, err)
			}
		default:
			break
		}
	}
	jsonBytes, _ := json.Marshal(w)
	sb := strings.Builder{}
	sb.Write(jsonBytes)
	jsonString := sb.String()
	ul.FDebugf("Line: %d | Unmarshaled to %s\n", lr.LineNumber, jsonString)
	return fieldErrors
}

type WeaponAttributes struct {
	Attributes map[string]bool
	AP         BoolAttribute `unit:"ap" json:"ap"`                       // armour piercing. Only counts half of target's armour
	BP         BoolAttribute `unit:"bp" json:"bp"`                       // body piercing. Missile can pass through men and hit those behind
	Spear      BoolAttribute `unit:"spear" json:"spear"`                 // Used for long spears. Gives bonuses fighting cavalry, and penalties against infantry
	LongPike   BoolAttribute `unit:"long_pike" json:"long_pike"`         // Use very long pikes. Phalanx capable units only
	ShortPike  BoolAttribute `unit:"short_pike" json:"short_pike"`       // Use shorter than normal spears.
	Prec       BoolAttribute `unit:"prec" json:"prec"`                   // Missile weapon is only thrown/ fired just before charging into combat
	Thrown     BoolAttribute `unit:"thrown" json:"thrown"`               // The missile type if thrown rather than fired
	Launching  BoolAttribute `unit:"launching" json:"launching"`         // attack may throw target men into the air
	Area       BoolAttribute `unit:"area" json:"area"`                   // attack affects an area, not just one man
	LightSpear BoolAttribute `unit:"light_spear" json:"light_spear"`     // The unit when braced has various protecting mechanisms versus cavalry charges from the frontk
	SpearBonus BoolAttribute `unit:"spear_bonus_x" json:"spear_bonus_x"` // attack bonus against cavalry. x = 2, 4, 6, 8, 10 or 12
}

type Armor struct {
	Armor        int
	DefenseSkill int
	Shield       int
	Sound        string
	FieldName    string
}

func (a *Armor) Unmarshal(armorInfo string, ul *UnitLogger, lr *LineRecord) (fieldErrors error) {
	armorStats, numFields, err := GetFieldInfo(armorInfo, 4, ul, lr)
	ul.FInfof("[INFO] Line: %d | %d fields detected", lr.LineNumber, numFields)
	if err != nil {
		fieldErrors = errors.Join(fieldErrors, err)
	}
	armorErr := CheckSetIntAttribute(&a.Armor, armorStats[0], "armor", 0, ul, DefaultErrorFormat, DefaultInfoFormat, lr.LineNumber)
	defErr := CheckSetIntAttribute(&a.Armor, armorStats[1], "defense skill", 1, ul, DefaultErrorFormat, DefaultInfoFormat, lr.LineNumber)
	shieldErr := CheckSetIntAttribute(&a.Armor, armorStats[2], "shield", 2, ul, DefaultErrorFormat, DefaultInfoFormat, lr.LineNumber)
	SetStrAttribute(&a.Sound, armorStats[3], "sound", ul, lr.LineNumber)
	if armorErr != nil || defErr != nil || shieldErr != nil {
		fieldErrors = errors.Join(fieldErrors, armorErr, defErr, shieldErr)
	}
	return fieldErrors
}

type ArmorEx struct{}
type Heat struct{}
type Ground struct{}
type Mental struct{}
type Food struct{}
type Cost struct{}
