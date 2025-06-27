package main

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
	al := len(a)
	if al == 0 || startIdx < 0 || al == 1 {
		return a
	} else {
		if startIdx == al-1 {
			return a[:al-1]
		}
		for i := startIdx; i < al-count; i++ {
			a[i] = a[i+count]
		}
		return a[:al-count]
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
