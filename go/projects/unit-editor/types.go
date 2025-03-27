package main

import (
	"encoding/json"
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
	LevelError LogLevel  
	LevelWarn LogLevel 
	LevelInfo LogLevel 
	LevelDebug LogLevel  
)

var LogLevel = 0

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

func DebugLogger() ( *UnitLogger) {
	logBuffer := &bytes.Buffer{}
	ul := &UnitLogger {
	debugStream: logBuffer,
	infoStream: logBuffer,
	warnStream: logBuffer,
	errorStream: logBuffer,
	}
	return ul
}

type UnitLogger struct {
	logLevel  int
	debugStream *bytes.Buffer
	infoStream  *bytes.Buffer
	warnStream  *bytes.Buffer
	errorStream *bytes.Buffer
	defaultDebugFormat string 
	defaultInfoFormat string 
	defaultWarnFormat string 
	defaultErrorFormat string 
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
	Logger *UnitLogger
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
		effectInt, err := ParseModifier(effectValue)
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
	Attack             int `json:"attack"`
	Charge             int `json:"charge"`
	MissileType        string `json:"missile_type"`
	MissileRange       int `json:"missile_range"`
	MissileAmmo        int `json:"missile_ammo"`
	WeaponType         string `json:"weapon_type"`
	TechType           string `json:"tech_type"`
	DamageType         string `json:"damage_type"`
	SoundType          string `json:"sound_type"`
	FireEffect         string `json:"fire_effect"`
	MinDelay           int `json:"min_delay"`
	CompensationFactor int `json:"compensation_factor"`
}

func (w *Weapon) Unmarshal(weaponInfo string, ul *UnitLogger, lr *LineRecord) error {
	conversionErrorFormat := fmt.Sprintf("line: %d | error converting %%s value of %%s to %%s: %%s\n")
	lineSections := CleanLine(weaponInfo)
	weaponStats := strings.Split(lineSections[1], ",")
	numFields := len(weaponStats)
	if numFields < 11 {
		ul.FErrorf("error parsing attack stats, too few fields")
	}
	switch numFields {
	case 11:
		md := strings.TrimSpace(weaponStats[9])
		cf := strings.TrimSpace(weaponStats[10])
		w.MinDelay, delayErr = strconv.Atoi(md)
		w.CompensationFactor, cfErr = strconv.Atoi(cf)
		if delayErr != nil {
		ul.FErrorf(conversionErrorFormat, lr.LineNumber, "MinDelay", md, delayErr)
		return delayErr
		}
		if cfErr != nil {
		ul.FErrorf(conversionErrorFormat, lr.LineNumber, "CompensationFactor", cf, cfErr)
		return cfErr
		}
	default:
		w.FireEffect = weaponStats[9]
		md := strings.TrimSpace(weaponStats[10])
		cf := strings.TrimSpace(weaponStats[11])
		w.MinDelay, delayErr = strconv.Atoi(md)
		w.CompensationFactor, cfErr = strconv.Atoi(cf)
		if delayErr != nil {
		ul.FErrorf(conversionErrorFormat, lr.LineNumber, "MinDelay", md, delayErr)
		return delayErr
		}
		if cfErr != nil {
		ul.FErrorf(conversionErrorFormat, lr.LineNumber, "CompensationFactor", cf, cfErr)
		return cfErr
		}
	}
	atk := strings.TrimSpace(weaponStats[0])
	chg := strings.TrimSpace(weaponStats[1])
	mr := strings.TrimSpace(weaponStats[3])
	ma := strings.TrimSpace(weaponStats[4])
	w.Attack, atkErr = strconv.Atoi(atk)
	w.Charge, chgErr = strconv.Atoi(crg)
	w.MissileRange, rangeErr = strconv.Atoi(mr)
	w.MissileAmmo, ammoErr = strconv.Atoi(atk)
	if atkErr != nil {
		ul.FErrorf(conversionErrorFormat, "Attack", atk, atkErr)
		return atkErr
	}
	if chgErr != nil {
		ul.FErrorf(conversionErrorFormat, "Charge", chg, chgErr)
		return chgErr
	}
	if rangeErr != nil {
		ul.FErrorf(conversionErrorFormat, "MissileRange", mr, rangeErr)
		return rangeErr
	}
	if ammoErr != nil {
		ul.FErrorf(conversionErrorFormat, "MissileAmmo", ma, ammoErr)
		return ammoErr
	}
	w.MissileType := strings.TrimSpace(weaponStats[2])
	w.WeaponType := strings.TrimSpace(weaponStats[5])
	w.TechType := strings.TrimSpace(weaponStats[6])
	w.DamageType := strings.TrimSpace(weaponStats[7])
	w.SoundType := strings.TrimSpace(weaponStats[8])
	jsonBytes, _ := json.Marshal(w)
	sb := strings.Builder{}
	sb.WriteByte(jsonBytes)
	jsonString := sb.String
	ul.FDebugf("line: %d | Unmarshaled to %s\n", lr.LineNumber, jsonString)
	return nil
}

type WeaponAttributes struct {
	Attributes map[string][bool]
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
