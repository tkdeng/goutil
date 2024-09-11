package goutil

import (
	"errors"
	"strconv"
	"strings"

	regex "github.com/tkdeng/goregex"
)

// Contains returns true if an array contains a value
func Contains[T any](search []T, value T) bool {
	val := ToType[string](value)
	for _, v := range search {
		if ToType[string](v) == val {
			return true
		}
	}
	return false
}

// IndexOf returns the index of a value in an array
//
// returns -1 and an error if the value is not found
func IndexOf[T any](search []T, value T) (int, error) {
	val := ToType[string](value)
	for i, v := range search {
		if ToType[string](v) == val {
			return i, nil
		}
	}
	return -1, errors.New("array does not contain value: " + ToType[string](value))
}

// ContainsMap returns true if a map contains a value
func ContainsMap[T Hashable, J any](search map[T]J, value J) bool {
	val := ToType[string](value)
	for _, v := range search {
		if ToType[string](v) == val {
			return true
		}
	}
	return false
}

// IndexOfMap returns the index of a value in a map
//
// returns an error if the value is not found
func IndexOfMap[T Hashable, J any](search map[T]J, value J) (T, error) {
	val := ToType[string](value)
	for i, v := range search {
		if ToType[string](v) == val {
			return i, nil
		}
	}
	var blk T
	return blk, errors.New("map does not contain value: " + ToType[string](value))
}

// ContainsMapKey returns true if a map contains a key
func ContainsMapKey[T Hashable, J any](search map[T]J, key T) bool {
	/* for i := range search {
		if i == key {
			return true
		}
	}
	return false */

	_, ok := search[key]
	return ok
}

// TrimTabs trims exxess beginning tab characters from a multiline string
//
// @size: number of tabs to trim
func TrimTabs(size uint8, buf []byte, tabSize ...uint8) []byte {
	if size == 0 {
		size = 1
	}

	t := "2"
	if len(tabSize) != 0 {
		t = strconv.FormatUint(uint64(tabSize[0]), 10)
	}

	buf = regex.Comp(`^\r?\n`).RepStrLit(buf, []byte{})
	buf = regex.Comp(`(?m)^(\t|[ ]{`+t+`}){1,`+strconv.FormatUint(uint64(size), 10)+`}`).RepStrLit(buf, []byte{})

	return buf
}

/*
CloneBytes is a simple method for copying a stuborn []byte that wants to be a reference

	golang default:
		buf := make([]byte, 5)
		buf = []byte{'t', 'e', 's', 't', '1'}
		newBuf := buf
		newBuf[4] = 2
		string(buf) == string(newBuf)

	using this method:
		buf := make([]byte, 5)
		buf = []byte{'t', 'e', 's', 't', '1'}
		newBuf := goutil.CloneBytes(buf)
		newBuf[4] = 2
		string(buf) != string(newBuf)
*/
func CloneBytes(b []byte) []byte {
	buf := make([]byte, len(b))
	copy(buf, b)
	return buf
}

// Swap will swap the values of 2 variables
func Swap[T any](val1 *T, val2 *T) {
	v := *val1
	*val1 = *val2
	*val2 = v
}

// ToColumns sorts a list of strings into rows and columns
//
// The column size will be wrapped by words first, and wrap to the specisied size.
//
// This method is useful for printing help info to a terminal.
//
// example:
//
//	stty := 80 // import "golang.org/x/term"
//	if width, _, err := term.GetSize(0); err == nil {
//	  stty = width
//	}
//
//	goutil.ToColumns([][]string{
//	  {"--opt1", "info"}
//	  {"--opt2, -o", "more info"}
//	}, stty, "    ", "\n\n")
//
// @sep (optional): column separator || row separator
//   - if "\n" is included: this will assume a row separator.
//   - else: this will assume a column separator
func ToColumns(list [][]string, charSize int, sep ...string) string {
	maxColSize := []int{}

	colSep := "    "
	rowSep := "\n\n"

	for _, s := range sep {
		if strings.ContainsRune(s, '\n') {
			rowSep = s
		} else {
			colSep = s
		}
	}

	for r, col := range list {
		for len(maxColSize) < len(col) {
			maxColSize = append(maxColSize, 0)
		}

		for len(list[r]) < len(maxColSize) {
			list[r] = append(list[r], "")
		}

		for i, val := range col {
			if len(val) > maxColSize[i] {
				maxColSize[i] = len(val)
			}
		}
	}

	maxSize := 0
	for i, s := range maxColSize {
		if i == len(maxColSize)-1 {
			if s := maxSize + len(colSep) + 1; s > charSize {
				charSize = s
			}
		} else {
			maxColSize[i] += len(colSep)
		}

		maxSize += s
	}

	for r, col := range list {
		colSize := 0

		for i, val := range col {
			size := maxColSize[i]
			if i != len(maxColSize)-1 {
				size -= len(colSep)

				if len(val) < size {
					list[r][i] += strings.Repeat(" ", size-len(val))
				}

				list[r][i] += colSep

				colSize += len(list[r][i])
			} else if size := charSize - colSize; len(val) > size {
				p := 0

				for len(val) > size+p {
					s := size + p

					hasU := false
					for u := s; u > p; u-- {
						if val[u] == ' ' {
							s = u
							hasU = true
							break
						}
					}

					if hasU {
						val = val[:s] + "\n" + strings.Repeat(" ", colSize) + val[s+1:]
						p += colSize
					} else {
						val = val[:s] + "\n" + strings.Repeat(" ", colSize) + val[s:]
						p += colSize + 1
					}
					p += s
				}

				// fix long words
				valList := strings.Split(val, "\n")

				for v := 0; v < len(valList); v++ {
					if len(valList[v]) > charSize {
						valList = append(valList, strings.Repeat(" ", colSize)+valList[v][charSize:])
						valList[v] = valList[v][:charSize]
					}
				}

				list[r][i] = strings.Join(valList, "\n")
			}
		}
	}

	res := ""
	for _, col := range list {
		res += strings.Join(col, "") + rowSep
	}

	return res
}
