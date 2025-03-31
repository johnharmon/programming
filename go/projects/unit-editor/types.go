package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

func CleanLine(line string) []string {
	re := regexp.MustCompile(`\s+`)
	lineSections := strings.SplitN(re.ReplaceAllString(line, " "), " ", 2)
	cleanSections := make([]string, len(lineSections))
	for _, item := range lineSections {
		cleanSections = append(cleanSections, strings.TrimSpace(item))
	}
	return cleanSections
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
	Logger                 *UnitLogger
	LineRecords            []*LineRecord
	Lines                  []string
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
	Effects            map[string]int
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
	SidetoSideSpacingTight  float64
	FronttoBackSpacingTight float64
	SidetoSideSpacingLoose  float64
	FronttoBackSpacingLoose float64
	DefaultRanks            int
	PossibleFormations      []string
}

func (f *Formation) Unmarshal(formationInfo string) error {
	lineSections := CleanLine(formationInfo)
	formationStats := strings.Split(lineSections[1], ",")
	numFields := len(formationStats)
	if numFields < 6 {
		return fmt.Errorf("error, insufficient number of fields for formation")
	}
	f.SidetoSideSpacingTight, _ = strconv.ParseFloat(formationStats[0], 64)
	f.FronttoBackSpacingTight, _ = strconv.ParseFloat(formationStats[1], 64)
	f.SidetoSideSpacingLoose, _ = strconv.ParseFloat(formationStats[2], 64)
	f.FronttoBackSpacingLoose, _ = strconv.ParseFloat(formationStats[3], 64)
	f.DefaultRanks, _ = strconv.Atoi(formationStats[4])
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

func (h *Health) Unmarshal(healthInfo string) error {
	lineSections := CleanLine(healthInfo)
	healthStats := strings.Split(lineSections[1], ",")
	h.HP, _ = strconv.Atoi(strings.TrimSpace(healthStats[0]))
	h.MountHP, _ = strconv.Atoi(strings.TrimSpace(healthStats[1]))
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

func TrimValues(values []string) (tv []string) {
	for _, value := range values {
		tv = append(tv, strings.TrimSpace(value))
	}
	return tv
}

func CheckSetIntAttribute(attribute *int, sAttr string, attrName string, index int, ul *UnitLogger, errorFormat string, infoFormat string) (errs error) {
	iAttr, attrErr := strconv.Atoi(sAttr)
	if attrErr != nil {
		ul.FErrorf(errorFormat, attrName, sAttr, attrErr)
	}
	errors.Join(errs, attrErr)
	attribute = &iAttr
	ul.FDebugf(infoFormat, attrName, index, sAttr, string(iAttr))
	return errs
}

func SetStrAttribute(attribute *string, sAttr string, attrName string, ul *UnitLogger, lineNumber int) {
	attribute = &sAttr
	ul.FDebugf("Line: %d | Setting %s to %s\n", lineNumber, attrName, sAttr)

}

func CheckSetStrAttribute(attribute *string, sAttr string, attrName string, ul *UnitLogger, lineNumber int, errorFormat string, acceptedValues map[string]struct{}) {
	if _, ok := acceptedValues[sAttr]; ok {
		attribute = &sAttr
		ul.FDebugf("[INFO] Line: %d | Setting %s to %s\n", lineNumber, attrName, sAttr)
	} else {
		ul.FErrorf("[ERROR] Line: %d | Error setting %s attribute, unaccepted value: %s\n", lineNumber, attrName, sAttr)
	}

}

func (w *Weapon) Unmarshal(weaponInfo string, ul *UnitLogger, lr *LineRecord) error {
	conversionErrorFormat := fmt.Sprintf("Line: %d | error converting %%s value of %%s to %%s: %%s\n", lr.LineNumber)
	infoFormat := fmt.Sprintf("Line: %d | Attribute: \"%%s\" | Position: %%d | Converted \"%%s\" to %%s\n", lr.LineNumber)
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
			_ = CheckSetIntAttribute(&w.Attack, value, "Attack", index, ul, conversionErrorFormat, infoFormat)
		case 1:
			_ = CheckSetIntAttribute(&w.Charge, value, "Charge", index, ul, conversionErrorFormat, infoFormat)
		case 2:
			SetStrAttribute(&w.MissileType, value, "MissileType", ul, lr.LineNumber)
		case 3:
			_ = CheckSetIntAttribute(&w.MissileRange, value, "MissileRange", index, ul, conversionErrorFormat, infoFormat)
		case 4:
			_ = CheckSetIntAttribute(&w.MissileAmmo, value, "MissileAmmo", index, ul, conversionErrorFormat, infoFormat)
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
				_ = CheckSetIntAttribute(&w.MinDelay, value, "MinDelay", index, ul, conversionErrorFormat, infoFormat)
			} else {
				w.FireEffect = value
				SetStrAttribute(&w.FireEffect, value, "FireEffect", ul, lr.LineNumber)
			}
		case 10:
			if numFields == 11 {
				_ = CheckSetIntAttribute(&w.CompensationFactor, value, "CompensationFactor", index, ul, conversionErrorFormat, infoFormat)
			} else {
				_ = CheckSetIntAttribute(&w.MinDelay, value, "MinDelay", index, ul, conversionErrorFormat, infoFormat)
			}
		case 11:
			_ = CheckSetIntAttribute(&w.CompensationFactor, value, "CompensationFactor", index, ul, conversionErrorFormat, infoFormat)
		default:
			break
		}
		jsonBytes, _ := json.Marshal(w)
		sb := strings.Builder{}
		sb.Write(jsonBytes)
		jsonString := sb.String
		ul.FDebugf("Line: %d | Unmarshaled to %s\n", lr.LineNumber, jsonString)
		return nil
	}
}

type WeaponAttributes struct {
	Attributes map[string]bool
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
