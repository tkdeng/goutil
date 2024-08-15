package goutil

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha512"
	"encoding/base64"
	"math"
	"strconv"
	"time"

	"github.com/AspieSoft/go-regex-re2/v2"
)

// RandBytes generates random bytes using crypto/rand
//
// @exclude[0] allows you can to pass an optional []byte to ensure that set of chars
// will not be included in the output string
//
// @exclude[1] provides a replacement string to put in place of the unwanted chars
//
// @exclude[2:] is currently ignored
func RandBytes(size uint, exclude ...[]byte) []byte {
	b := make([]byte, size)
	rand.Read(b)
	b = []byte(base64.URLEncoding.EncodeToString(b))

	if len(exclude) >= 2 {
		if exclude[0] == nil || len(exclude[0]) == 0 {
			b = regex.Comp(`[^\w_-]`).RepStrLit(b, exclude[1])
		} else {
			b = regex.Comp(`[`+regex.Escape(string(exclude[0]))+`]`).RepStrLit(b, exclude[1])
		}
	} else if len(exclude) >= 1 {
		if exclude[0] == nil || len(exclude[0]) == 0 {
			b = regex.Comp(`[^\w_-]`).RepStrLit(b, []byte{})
		} else {
			b = regex.Comp(`[`+regex.Escape(string(exclude[0]))+`]`).RepStrLit(b, []byte{})
		}
	}

	for uint(len(b)) < size {
		a := make([]byte, size)
		rand.Read(a)
		a = []byte(base64.URLEncoding.EncodeToString(a))

		if len(exclude) >= 2 {
			if exclude[0] == nil || len(exclude[0]) == 0 {
				a = regex.Comp(`[^\w_-]`).RepStrLit(a, exclude[1])
			} else {
				a = regex.Comp(`[`+regex.Escape(string(exclude[0]))+`]`).RepStrLit(a, exclude[1])
			}
		} else if len(exclude) >= 1 {
			if exclude[0] == nil || len(exclude[0]) == 0 {
				a = regex.Comp(`[^\w_-]`).RepStrLit(a, []byte{})
			} else {
				a = regex.Comp(`[`+regex.Escape(string(exclude[0]))+`]`).RepStrLit(a, []byte{})
			}
		}

		b = append(b, a...)
	}

	return b[:size]
}

// URandBytes tries to generate a unique random bytes
//
// This method uses the current microsecond and crypto random bytes to generate unique keys.
// This method also only returns alphanumeric characters [A-Za-z0-9]
//
// @unique (optional): add a list pointer, to keep track of what keys were already used.
// This method will automattically append new keys to the list.
// If the same key is generated twice, the function will try again (using recursion).
func URandBytes(size uint, unique ...*[][]byte) []byte {
	if size < 8 {
		size = 8
	}

	var b []byte
	if size < 12 {
		b = []byte(strconv.FormatUint(uint64(time.Now().UnixMicro())/100000000000, 36))
	} else {
		b = []byte(strconv.FormatUint(uint64(time.Now().UnixMicro())/100000000, 36))
	}

	for uint(len(b)) < size {
		a := make([]byte, size)
		rand.Read(a)
		a = bytes.TrimRight([]byte(base64.URLEncoding.EncodeToString(a)), "=")

		b = append(b, a...)
	}

	b = bytes.ReplaceAll(b, []byte{'-'}, []byte{})
	b = bytes.ReplaceAll(b, []byte{'_'}, []byte{})
	b = b[:size]

	if len(unique) != 0 {
		if Contains(*unique[0], b) {
			time.Sleep(1 * time.Microsecond)
			return URandBytes(size, unique[0])
		}
	}

	*unique[0] = append(*unique[0], b)
	return b
}

var uuidGenLastTime int64

// GenUUID generates a Unique Identifier using a custom build method
//
// Notice: This feature is currently in beta
//
// @size: (minimum: 8) the bit size for the last part of the uuid
// (note: other parts may vary)
//
// @timezone: optionally add a timezone string to the uuid
// (note: you could also pass random info into here for a more complex algorithm)
//
// This method uses the following data:
//   - A hash of the current year and day of year
//   - A hash of the current timezone
//   - A hash of the current unix time (in seconds)
//   - A hash of the current unix time in nanoseconds and a random number
//
// The returned value is url encoded and will look something like this: xxxx-xxxx-xxxx-xxxxxxxx
func GenUUID(size uint, timezone ...string) string {
	for time.Now().UnixNano() <= uuidGenLastTime {
		time.Sleep(3 * time.Nanosecond)
	}
	uuidGenLastTime = time.Now().UnixNano()

	if size < 8 {
		size = 8
	}

	uuid := [][]byte{{}, {}, {}, {}}

	// year
	{
		s := int(math.Min(float64(size/4), 8))
		if s < 4 {
			s = 4
		}

		sm := s / 2
		if s%2 != 0 {
			sm++
		}

		b := sha1.Sum([]byte(strconv.Itoa(time.Now().Year())))
		uuid[0] = []byte(base64.URLEncoding.EncodeToString(b[:]))[:sm]
		b = sha1.Sum([]byte(strconv.Itoa(time.Now().YearDay())))
		uuid[0] = append(uuid[0], []byte(base64.URLEncoding.EncodeToString(b[:]))[:sm]...)
		uuid[0] = uuid[0][:s]
	}

	// time zone
	{
		s := int(math.Min(float64(size/8), 8))
		if s < 4 {
			s = 4
		}

		if len(timezone) != 0 {
			sm := s / len(timezone)
			if s%2 != 0 {
				sm++
			}

			for _, zone := range timezone {
				b := sha1.Sum([]byte(zone))
				uuid[1] = append(uuid[1], []byte(base64.URLEncoding.EncodeToString(b[:]))[:sm]...)
			}
			uuid[1] = uuid[1][:s]
		} else {
			z, _ := time.Now().Zone()
			b := sha1.Sum([]byte(z))
			uuid[1] = []byte(base64.URLEncoding.EncodeToString(b[:]))[:s]
		}
	}

	// unix time
	{
		s := int(math.Min(float64(size/2), 16))
		if s < 4 {
			s = 4
		}

		b := sha1.Sum([]byte(strconv.Itoa(int(time.Now().Unix()))))
		uuid[2] = []byte(base64.URLEncoding.EncodeToString(b[:]))[:s]
	}

	// random
	{
		s := uint(math.Min(float64(size/4), 64))
		if s < 4 {
			s = 4
		}

		b := sha512.Sum512([]byte(strconv.Itoa(int(time.Now().UnixNano()))))
		uuid[3] = []byte(base64.URLEncoding.EncodeToString(b[:]))[:s]
		// uuid[3] = append(uuid[3], []byte(base64.URLEncoding.EncodeToString(RandBytes(size)))[:size-s]...)
		uuid[3] = append(uuid[3], []byte(base64.URLEncoding.EncodeToString(URandBytes(size)))[:size-s]...)
	}

	if len(uuid[1]) == 0 {
		uuid = append(uuid[:1], uuid[2:]...)
	}

	for i := range uuid {
		uuid[i] = bytes.ReplaceAll(uuid[i], []byte{'-'}, []byte{'0'})
		uuid[i] = bytes.ReplaceAll(uuid[i], []byte{'_'}, []byte{'1'})
	}

	return string(bytes.Join(uuid, []byte{'-'}))
}
