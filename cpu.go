package goutil

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/host"
)

type CPU struct {
	// HighTemp is the highest temperature for the WaitToCool method to detect that the cpu is too hot
	//
	// note: if set to strict mode, LowTemp will take this variables place
	//
	// default: 64 (celsius)
	HighTemp uint16

	// LowTemp is the lowest temperature for the WaitToCool method to wait for, before deciding that the cpu has colled down,
	// if it was previously throttled by HighTemp (except in strict mode, where LowTemp takes the place of HighTemp)
	//
	// default: 56 (celsius)
	LowTemp uint16

	// when Logging is enabled, the WaitToCool method will log info to the console to report when it is waiting for the cpu to cool down
	//
	// you can set this var to false to disable this feature
	//
	// default: true
	Logging bool
}

// GetTemp returns the average cpu temperature in celsius
func (cpu *CPU) GetTemp() uint16 {
	temps, err := host.SensorsTemperatures()
	if err != nil {
		return 0
	}

	var i float64
	var temp float64
	for _, t := range temps {
		if strings.HasSuffix(t.SensorKey, "_input") {
			i++
			temp += t.Temperature
		}
	}

	temp = math.Round(temp / i)
	if temp < 0 || uint16(temp) > 1000 {
		return 0
	}
	return uint16(temp)
}

// WaitToCool makes your function wait for the cpu to cool down
//
// by default, if the temperature > HighTemp, it will wait until the temperature <= LowTemp
//
// in strict mode, this will run if temperature > LowTemp
//
// HighTemp = 64
// LowTemp = 56
func (cpu *CPU) WaitToCool(strict bool) {
	if cpu.HighTemp == 0 {
		cpu.HighTemp = 64
	}

	if cpu.LowTemp == 0 {
		cpu.LowTemp = 56
	}

	flagTemp := cpu.HighTemp
	if strict {
		flagTemp = cpu.LowTemp
	}

	if temp := cpu.GetTemp(); temp > flagTemp {
		if cpu.Logging {
			fmt.Println("CPU Too Hot!")
			fmt.Println("Waiting for it to cool down...")
			fmt.Print("CPU Temp:", strconv.Itoa(int(temp))+"°C", "          \r")
		}
		for {
			time.Sleep(10 * time.Second)
			temp := cpu.GetTemp()
			if cpu.Logging {
				fmt.Print("CPU Temp:", strconv.Itoa(int(temp))+"°C", "          \r")
			}
			if temp <= cpu.LowTemp {
				break
			}
		}
		if cpu.Logging {
			fmt.Println("\nCPU Temperature Stable!")
		}
	}

	time.Sleep(300 * time.Millisecond)
}
