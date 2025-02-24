package main

import (
	"strings"
)

type PythonString string

type StringLike interface {
	~string
}

func (p PythonString) Split(splitString string) []PythonString {
	subStrings := strings.Split(string(p), splitString)
	pSubStrings := make([]PythonString, len(subStrings))
	for i, s := range subStrings {
		pSubStrings[i] = PythonString(s)
	}
	return pSubStrings
}

func (p PythonString) Replace(old string, new string) PythonString {
	newString := PythonString(strings.ReplaceAll(string(p), old, new))
	return newString
}

func (p PythonString) Strip(strip string) PythonString {
	return PythonString(strings.Trim(string(p), strip))
}

func (p PythonString) Concat(s ...string) PythonString {
	var sb strings.Builder
	sb.WriteString(string(p))
	for _, s := range s {
		sb.WriteString(s)
	}
	return PythonString(sb.String())
}

func (p PythonString) Capitalize() PythonString {
	first := strings.ToUpper(string(p[0]))
	rest := string(p[1:])
	var sb strings.Builder
	sb.WriteString(first)
	sb.WriteString(rest)
	return PythonString(sb.String())

}

func (p PythonString) Center(padding int) PythonString {
	spadding := strings.Repeat(" ", padding)
	var sb strings.Builder
	sb.WriteString(spadding)
	sb.WriteString(string(p))
	sb.WriteString(spadding)
	return PythonString(sb.String())
}

func (p PythonString) Count(sub string) int {
	return strings.Count(string(p), sub)
}

func (p PythonString) EndsWith(sub string) bool {
	sLen := len(sub)
	pLen := len(p)
	if sLen > pLen {
		return false
	}
	startIdx := pLen - sLen
	for i, s := range p[startIdx:] {
		if byte(s) != sub[i] {
			return false
		}
	}
	return true
}

func (p PythonString) ExpandTabs(size int) PythonString {
	col := 0
	var sb strings.Builder
	for _, r := range p {
		if r == '\t' {
			spaces := size - (col % size)
			sb.WriteString(strings.Repeat(" ", spaces))
			col += spaces
		} else {
			sb.WriteRune(r)
			if r == '\n' {
				col = 0
			}
			col++
		}
	}
	return PythonString(sb.String())
}

func (p PythonString) Find(sub string) int {
	matched := -1
	for i, r := range p {
		if byte(r) == sub[0] {
			matched = i
			for si, sr := range sub {
				if byte(sr) != p[i+si] {
					matched = -1
					break
				}
			}
			if matched != -1 {
				//return matched
				break
			}
		}
	}
	return matched
}
