//go:build rp2040

package golfclubs

import "machine"

var pwmGroups = []PWMGroup{
	machine.PWM0,
	machine.PWM1,
	machine.PWM2,
	machine.PWM3,
	machine.PWM4,
	machine.PWM5,
	machine.PWM6,
	machine.PWM7,
}

// getPWMGroup 基于序号获取 PWM 组
func getPWMGroup(i uint8) PWMGroup {
	//return pwmGroups[i]
	return machine.PWM1
}
