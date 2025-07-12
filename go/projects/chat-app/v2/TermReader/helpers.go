package main

import (
	//"bytes"
	"strings"
)

func InsertByteAt(a []byte, b byte, startIdx int) []byte {
	al := len(a)
	if startIdx >= len(a) {
		// fmt.Printf("\r\x1b[2Kappending\n\r")
		//		fmt.Printf("\r\x1b[2k%s\n\r", string(append(a, b)))
		//	os.Exit(5)
		return append(a, b)
	} else if cap(a) >= len(a)+1 {
		// fmt.Printf("\r\x1b[2Kgrowing a\n\r")
		// os.Exit(6)
		a = a[0 : al+1]
		for i := al - 1; i >= startIdx; i-- {
			if i == startIdx {
				a[i] = b
			} else {
				a[i+1] = a[i]
			}
		}
		//		fmt.Printf("\r\x1b[2k%s\n\r", string(a))
		return a
	}
	// fmt.Printf("\r\x1b[2Kmaking new slice\n\r")
	// os.Exit(7)
	tmp := make([]byte, 0, (al+1)*2)
	tmp = append(tmp, a[0:startIdx]...)
	tmp = append(tmp, b)
	tmp = append(tmp, a[startIdx:]...)
	// fmt.Printf("\r\x1b[2k%s\n\r", string(tmp))
	return tmp
}

func InsertAt(a []byte, b []byte, startIdx int) []byte {
	al := len(a)
	bl := len(b)
	if startIdx >= len(a) {
		return append(a, b...)
	} else if cap(a) >= len(a)+len(b)+1 {
		a = a[0 : al+bl]
		if bl == 1 {
			for i := al - 1; i >= startIdx; i-- {
				a[i+1] = a[i]
				if i == startIdx {
					a[i] = b[0]
				}
			}
		} else {
			for i := al - 1 - bl; i >= startIdx; i-- {
				a[i+bl] = a[i]
			}
		}
		copy(a[startIdx:startIdx+bl], b)
		return a
	}
	tmp := make([]byte, 0, (al+bl)*2)
	tmp = append(tmp, a[0:startIdx]...)
	tmp = append(tmp, b...)
	tmp = append(tmp, a[startIdx:]...)
	return tmp
}

func InsertLineAt(a [][]byte, b [][]byte, startIdx int) [][]byte {
	al := len(a)
	bl := len(b)
	if startIdx >= len(a) {
		return append(a, b...)
	} else if cap(a) >= len(a)+len(b)+1 {
		a = a[0 : al+bl]
		if bl == 1 {
			for i := al - 1; i >= startIdx; i-- {
				a[i+1] = a[i]
				if i == startIdx {
					a[i] = b[0]
				}
			}
		} else {
			for i := al - 1 - bl; i >= startIdx; i-- {
				a[i+bl] = a[i]
			}
		}
		copy(a[startIdx:startIdx+bl], b)
		return a
	}
	tmp := make([][]byte, 0, (al+bl)*2)
	tmp = append(tmp, a[0:startIdx]...)
	tmp = append(tmp, b...)
	tmp = append(tmp, a[startIdx:]...)
	return tmp
}

func DeleteLineAt(a [][]byte, startIdx int, count int) [][]byte {
	aLen := len(a)
	if aLen == 0 || startIdx < 0 || (startIdx >= aLen) {
		return a
	} else if aLen == 1 && count > 0 {
		a[0] = make([]byte, 1)
		return a[:1]
	} else {
		if startIdx == aLen-1 {
			return a[:aLen-1]
		} else if (aLen - startIdx) < count {
			count = aLen - startIdx
		}
		for i := startIdx; i < aLen-count; i++ {
			a[i] = a[i+count]
		}
		return a[:aLen-count]
	}
}

func DeleteAt(a []byte, startIdx int, count int) []byte {
	if startIdx < 0 {
		return a
	}
	al := len(a)
	if al == 0 || startIdx < 0 {
		return a
	} else {
		if startIdx == al-1 {
			if al == 0 {
				return a
			} else {
				return a[:al-1]
			}
		} else {
			for i := startIdx; i < al-count; i++ {
				a[i] = a[i+count]
			}
		}
		return a[:al-count]
	}
}

func DeleteByteAt(a []byte, startIdx int) []byte {
	if startIdx < 0 {
		return a
	}
	GlobalLogger.Logln("DeleteByteAt called with:")
	al := len(a)
	GlobalLogger.Logln("\tByte Slice: %b\n\tstartIdx: %d\n\tSlice Length: %d", a, startIdx, al)
	// time.Sleep(1 * time.Millisecond)
	if startIdx >= al-1 {
		if al == 0 {
			return a
		} else {
			return a[:al-1]
		}
	} else {
		for i := startIdx; i < al-1; i++ {
			a[i] = a[i+1]
		}
		return a[:al-1]
	}
}

func Squeeze(s byte, buf []byte) []byte {
	newBuf := make([]byte, len(buf))
	bytesCopied := 0
	previousByteWasS := false
	for _, b := range buf {
		if b == s {
			if previousByteWasS {
				continue
			} else {
				previousByteWasS = true
				newBuf[bytesCopied] = b
				bytesCopied++
			}
		} else {
			previousByteWasS = false
			newBuf[bytesCopied] = b
			bytesCopied++
		}
	}
	return newBuf[0:bytesCopied]
}

func ProcessCmdArgs(cmdRaw []byte) (cmd string, cmdArgs []string) {
	cmdRaw = cmdRaw[1:]
	cmdRaw = Squeeze(' ', cmdRaw)
	cmdString := string(cmdRaw)
	cmdStrings := strings.Split(cmdString, " ")
	cmd = cmdStrings[0]
	if len(cmdStrings) > 1 {
		cmdArgs = cmdStrings[1:]
	} else {
		cmdArgs = make([]string, 0)
	}
	return cmd, cmdArgs
}

/*
Steps right in a utf-8 encoded byte string by one character
Moves to the start of the next utf-8 code point if in the middle of one
Moves one byte right if it cannot detect a valid utf-8 byte to prevent state freeze
*/
func StepRight(line []byte, curPos int) int {
	if curPos >= len(line) {
		return -1
	} else if ok, step := IsStartingByte(line[curPos]); ok {
		return step
	} else if IsContinuationByte(line[curPos]) {
		return StepRightMinimal(line, curPos)
	} else {
		return 1
	}
}

// minimally steps right until a utf-8 starting byte is found, including the one it starts on
func StepRightMinimal(line []byte, curPos int) int {
	for i := 0; i+curPos < len(line); i++ {
		if line[curPos+i]&CONTINUATION_BYTE_OPERATOR != CONTINUATION_BYTE_RESULT {
			return i
		}
	}
	return 1
}

/*
Steps left one character, or until beginning of line.
Returns how many index positions backward it needed to move to find a starting byte or the beginning of the line
*/
func StepLeft(line []byte, curPos int) int { // move left until a utf-8 starting byte is found
	step := -1
	if curPos-1 < 0 || curPos-1 >= len(line) {
		return 0
	} else {
		for i := curPos - 1; i >= 0; i-- {
			if ok, _ := IsStartingByte(line[i]); ok {
				break
			}
			step--
		}
		return step
	}
}

func StepLeftMinmal(line []byte, curPos int) int { // move left until a utf-8 starting byte is found, including the one it starts on
	step := 0
	if curPos < 0 || curPos >= len(line) {
		return 0
	} else {
		for i := curPos; i >= 0; i-- {
			if ok, _ := IsStartingByte(line[i]); ok {
				break
			}
			step--
		}
		return step
	}
}

func GetCharacterPositionByIndex(line []byte, curPos int) int { // Gets the display/character position of a point from a raw []byte location
	if curPos > len(line) {
		curPos = len(line)
	} else if curPos < 0 {
		curPos = 0
	}
	chars := 0
	curIndex := 0
	for curIndex < curPos {
		_, _, step := StepRightUntilValidByte(line, curIndex)
		curIndex += step
		chars += 1
	}
	return chars
}

func IsStartingByteByShifting(b byte) (bool, int) {
	switch {
	case b>>3 == 0b11110:
		return true, 4
	case b>>4 == 0b1110:
		return true, 3
	case b>>5 == 0b110:
		return true, 2
	case b>>7 == 0b0:
		return true, 1
	default:
		return false, -1
	}
}

func IsStartingByte(b byte) (bool, int) {
	switch {
	case b&SINGLE_BYTE_OPERATOR == SINGLE_BYTE_RESULT:
		return true, 1
	case b&TWO_BYTE_OPERATOR == TWO_BYTE_RESULT:
		return true, 2
	case b&THREE_BYTE_OPERATOR == THREE_BYTE_RESULT:
		return true, 3
	case b&FOUR_BYTE_OPERATOR == FOUR_BYTE_RESULT:
		return true, 4
	default:
		return false, -1
	}
}

func IsContinuationByte(b byte) bool {
	return b&CONTINUATION_BYTE_OPERATOR == CONTINUATION_BYTE_RESULT
}

func AbsInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// Returns the nuber of utf-8 characters in a line
func Utf8Len(line []byte) int {
	chars := 0
	i := 0
	for i < len(line) {
		_, _, step := StepRightUntilValidByte(line, i)
		if step < 1 {
			GlobalLogger.Logln("Step less than 1, breaking loop, char count is: %d, byte index was: %d", chars, i)
			break
		}
		i += step
		chars++
	}
	return chars
}

func ValidateUtf8(seq []byte) (bool, int) {
	if ok, seqLen := IsStartingByte(seq[0]); ok {
		if seqLen > 1 {
			for i := 1; i < seqLen; i++ {
				if i >= len(seq) {
					return false, i
				} else if !IsContinuationByte(seq[i]) {
					return false, i
				}
			}
			return true, seqLen
		} else {
			return true, 1
		}
	}
	return false, 1
}

/*
Validates the start of a utf-8 byte at position, and returns the length of the byte, or the length of the malformed sequence
*/
func ValidateUtf8At(line []byte, start int) (bool, int) {
	if start < 0 || start >= len(line) {
		return false, 0
	}
	if ok, seqLen := IsStartingByte(line[start]); ok {
		if seqLen > 1 {
			if seqLen+start >= len(line) {
				return false, len(line) - start
			}
			for i := 1; i < seqLen; i++ {
				if !IsContinuationByte(line[start+i]) {
					return false, i
				}
			}
			return true, seqLen
		} else {
			return true, 1
		}
	}
	return false, 1
}

/*
Will step right until it finds a valid utf-8 byte.
Will return one false if the starting byte is not a valid starting utf-8 byte
Will return a second false if it cannot find another valid starting byte within the the range line[curPos:]
Returns an int representing the total distance from the starting position to the next valid utf-8 byte start, or to EOL
*/
func StepRightUntilValidByte(line []byte, curPos int) (startOk bool, foundOk bool, totalStep int) {
	ok, step := ValidateUtf8At(line, curPos)
	if ok {
		GlobalLogger.Logln("Valid byte identified, stepped forward: %d", step)
		return ok, ok, step
	}

	totalStep += step
	for {
		if curPos+totalStep >= len(line) {
			GlobalLogger.Logln("No valid byte identified from %d to EOL", curPos)
			return false, false, len(line) - curPos
		}
		ok, step = ValidateUtf8At(line, curPos+totalStep)
		if ok {
			GlobalLogger.Logln("Valid byte identified at: %d, stepped forward: %d", curPos+totalStep, totalStep)
			return false, true, totalStep
		} else {
			totalStep += step
		}
	}
}

//func GetNthChar(line []byte, charPos int) (start int, end int) {
//	start, charLen := GetBytePositionByCharacter(line, charPos)
//	return start, start + charLen
//}

/*
Gets the byte intex of the start of the character position given (0 based index)
*/
func GetNthChar(line []byte, charPos int) (bytePos int, charLen int) {
	bytePos = 0
	chars := 0
	step := 0
	prevStep := 0
	for chars <= charPos {
		bytePos += prevStep
		_, _, step = StepRightUntilValidByte(line, bytePos)
		GlobalLogger.Logln("Validating byte at: %d", bytePos)
		prevStep = step
		chars++
		if bytePos+step >= len(line)-1 {
			break
		}
	}
	if bytePos >= len(line) {
		bytePos = len(line) - 1
	}
	return bytePos, step
}
