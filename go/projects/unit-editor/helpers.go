package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func TrimValues(values []string) (tv []string) {
	for _, value := range values {
		tv = append(tv, strings.TrimSpace(value))
	}
	return tv
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

func GetFieldInfo(field string, neededFields int, ul *UnitLogger, lr *LineRecord) (fieldStats []string, numFields int, fieldErrors error) {
	lineSections := CleanLine(field)
	fieldStats = strings.Split(lineSections[1], ",")
	numFields = len(fieldStats)
	if numFields < neededFields {
		fieldErrors = errors.Join(fieldErrors, fmt.Errorf("Line: "))
		ul.FErrorf("Line: %d | Error, too few values in field: %s | Required: %d, provided: %d", lr.LineNumber, lineSections[0], neededFields, numFields)
	}
	return fieldStats, numFields, fieldErrors

}

func CheckSetIntAttribute(attribute *int, sAttr string, attrName string, index int, ul *UnitLogger, errorFormat string, infoFormat string, lineNumber int) (errs error) {
	iAttr, attrErr := strconv.Atoi(sAttr)
	if attrErr != nil {
		ul.FErrorf(errorFormat, attrName, sAttr, attrErr)
	}
	errors.Join(errs, attrErr)
	*attribute = iAttr
	ul.FDebugf(infoFormat, lineNumber, attrName, index, sAttr, sAttr)
	return errs
}

func CheckSetFloatAttribute(attribute *float64, sAttr string, attrName string, index int, ul *UnitLogger, errorFormat string, infoFormat string, lineNumber int) (errs error) {
	fAttr, attrErr := strconv.ParseFloat(sAttr, 64)
	if attrErr != nil {
		ul.FErrorf(errorFormat, attrName, sAttr, attrErr)
	}
	errors.Join(errs, attrErr)
	*attribute = fAttr
	ul.FDebugf(infoFormat, attrName, index, sAttr, sAttr)
	return errs
}

func SetStrAttribute(attribute *string, sAttr string, attrName string, ul *UnitLogger, lineNumber int) {
	*attribute = sAttr
	ul.FDebugf("Line: %d | Setting %s to %s\n", lineNumber, attrName, sAttr)

}

func CheckSetStrAttribute(attribute *string, sAttr string, attrName string, ul *UnitLogger, lineNumber int, errorFormat string, acceptedValues map[string]struct{}) {
	if _, ok := acceptedValues[sAttr]; ok {
		*attribute = sAttr
		ul.FDebugf("[INFO] Line: %d | Setting %s to %s\n", lineNumber, attrName, sAttr)
	} else {
		ul.FErrorf("[ERROR] Line: %d | Error setting %s attribute, unaccepted value: %s\n", lineNumber, attrName, sAttr)
	}

}
