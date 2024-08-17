package goutil

import (
	"math"
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/AspieSoft/go-regex-re2/v2"
)

// SysFreeMemory returns the amount of memory available in megabytes
func SysFreeMemory() float64 {
	in := &syscall.Sysinfo_t{}
	err := syscall.Sysinfo(in)
	if err != nil {
		return 0
	}

	// If this is a 32-bit system, then these fields are
	// uint32 instead of uint64.
	// So we always convert to uint64 to match signature.
	return math.Round(float64(uint64(in.Freeram)*uint64(in.Unit))/1024/1024*100) / 100
}

// FormatMemoryUsage converts bytes to megabytes
func FormatMemoryUsage(b uint64) float64 {
	return math.Round(float64(b)/1024/1024*100) / 100
}

var regIsAlphaNumeric *regex.Regexp = regex.Comp(`^[A-Za-z0-9]+$`)

// MapArgs will convert a bash argument array ([]string) into a map (map[string]string)
//
// When @args is left blank with no values, it will default to os.Args[1:]
//
// -- Arg Conversions:
//
// "--Key=value" will convert to "key:value"
//
// "--boolKey" will convert to "boolKey:true"
//
// "-flags" will convert to "f:true, l:true, a:true, g:true, s:true" (only if its alphanumeric [A-Za-z0-9])
// if -flags is not alphanumeric (example: "-test.paniconexit0" "-test.timeout=10m0s") it will be treated as a --flag (--key=value --boolKey)
//
// keys that match a number ("--1" or "-1") will start with a "-" ("--1=value" -> "-1:value", "-1" -> -1:true)
// this prevents a number key from conflicting with an index key
//
// everything else is given a number value index starting with 0
//
// this method will not allow --args to have their values modified after they have already been set
func MapArgs(args ...[]string) map[string]string {
	if len(args) == 0 {
		args = append(args, os.Args[1:])
	}

	argMap := map[string]string{}
	i := 0

	for _, argList := range args {
		for _, arg := range argList {
			if strings.HasPrefix(arg, "--") {
				arg = arg[2:]
				if strings.ContainsRune(arg, '=') {
					data := strings.SplitN(arg, "=", 2)
					if _, err := strconv.Atoi(data[0]); err == nil {
						if argMap["-"+data[0]] == "" {
							argMap["-"+data[0]] = data[1]
						}
					} else {
						if argMap[data[0]] == "" {
							argMap[data[0]] = data[1]
						}
					}
				} else {
					if _, err := strconv.Atoi(arg); err == nil {
						if argMap["-"+arg] == "" {
							argMap["-"+arg] = "true"
						}
					} else {
						if argMap[arg] == "" {
							argMap[arg] = "true"
						}
					}
				}
			} else if strings.HasPrefix(arg, "-") {
				arg = arg[1:]
				if regIsAlphaNumeric.Match([]byte(arg)) {
					flags := strings.Split(arg, "")
					for _, flag := range flags {
						if _, err := strconv.Atoi(flag); err == nil {
							if argMap["-"+flag] == "" {
								argMap["-"+flag] = "true"
							}
						} else {
							if argMap[flag] == "" {
								argMap[flag] = "true"
							}
						}
					}
				} else {
					if strings.ContainsRune(arg, '=') {
						data := strings.SplitN(arg, "=", 2)
						if _, err := strconv.Atoi(data[0]); err == nil {
							if argMap["-"+data[0]] == "" {
								argMap["-"+data[0]] = data[1]
							}
						} else {
							if argMap[data[0]] == "" {
								argMap[data[0]] = data[1]
							}
						}
					} else {
						if _, err := strconv.Atoi(arg); err == nil {
							if argMap["-"+arg] == "" {
								argMap["-"+arg] = "true"
							}
						} else {
							if argMap[arg] == "" {
								argMap[arg] = "true"
							}
						}
					}
				}
			} else {
				argMap[strconv.Itoa(i)] = arg
				i++
			}
		}
	}

	return argMap
}
