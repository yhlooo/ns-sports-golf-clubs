package golfclubs

import (
	"log"
	"machine"
	"time"
)

const (
	// DefaultPulsesPerCircle 电机旋转一周默认所需脉冲数
	DefaultPulsesPerCircle uint32 = 1600
	// MaxSpeed 最大挥杆速度（单位： rpm ）
	MaxSpeed uint32 = 400
)

// New 创建一个 GolfClubs
func New(pwm, dir, en machine.Pin) *GolfClubs {
	return &GolfClubs{
		PWMPin: pwm,
		DirPin: dir,
		EnPin:  en,
	}
}

// GolfClubs 高尔夫球杆驱动器
type GolfClubs struct {
	// PWM 信号
	PWMPin machine.Pin
	// 方向控制
	DirPin machine.Pin
	// 脱机控制
	EnPin machine.Pin

	pwm   PWMGroup
	pwmCh uint8

	reverse         bool
	pulsesPerCircle uint32
}

// PWMGroup PWM 组
type PWMGroup interface {
	// Configure enables and configures this PWM.
	Configure(config machine.PWMConfig) error
	// Channel returns a PWM channel for the given pin. If pin does
	// not belong to PWM peripheral ErrInvalidOutputPin error is returned.
	// It also configures pin as PWM output.
	Channel(pin machine.Pin) (channel uint8, err error)
	// SetPeriod updates the period of this PWM peripheral in nanoseconds.
	// To set a particular frequency, use the following formula:
	//
	//	period = 1e9 / frequency
	//
	// Where frequency is in hertz. If you use a period of 0, a period
	// that works well for LEDs will be picked.
	//
	// SetPeriod will try not to modify TOP if possible to reach the target period.
	// If the period is unattainable with current TOP SetPeriod will modify TOP
	// by the bare minimum to reach the target period. It will also enable phase
	// correct to reach periods above 130ms.
	SetPeriod(period uint64) error
	// Top returns the current counter top, for use in duty cycle calculation.
	//
	// The value returned here is hardware dependent. In general, it's best to treat
	// it as an opaque value that can be divided by some number and passed to Set
	// (see Set documentation for more information).
	Top() uint32
	// Counter returns the current counter value of the timer in this PWM
	// peripheral. It may be useful for debugging.
	Counter() uint32
	// Period returns the used PWM period in nanoseconds.
	Period() uint64
	// SetInverting sets whether to invert the output of this channel.
	// Without inverting, a 25% duty cycle would mean the output is high for 25% of
	// the time and low for the rest. Inverting flips the output as if a NOT gate
	// was placed at the output, meaning that the output would be 25% low and 75%
	// high with a duty cycle of 25%.
	SetInverting(channel uint8, inverting bool)
	// Set updates the channel value. This is used to control the channel duty
	// cycle, in other words the fraction of time the channel output is high (or low
	// when inverted). For example, to set it to a 25% duty cycle, use:
	//
	//	pwm.Set(channel, pwm.Top() / 4)
	//
	// pwm.Set(channel, 0) will set the output to low and pwm.Set(channel,
	// pwm.Top()) will set the output to high, assuming the output isn't inverted.
	Set(channel uint8, value uint32)
	// Get current level (last set by Set). Default value on initialization is 0.
	Get(channel uint8) (value uint32)
	// SetTop sets TOP control register. Max value is 16bit (0xffff).
	SetTop(top uint32)
	// SetCounter sets counter control register. Max value is 16bit (0xffff).
	// Useful for synchronising two different PWM peripherals.
	SetCounter(ctr uint32)
	// Enable enables or disables PWM peripheral channels.
	Enable(enable bool)
	// IsEnabled returns true if peripheral is enabled.
	IsEnabled() (enabled bool)
}

// Config 配置
type Config struct {
	// 反向挥杆
	Reverse bool
	// 电机旋转一周所需脉冲数
	// 默认为 DefaultPulsesPerCircle
	PulsesPerCircle uint32
}

// Configure 初始配置
func (c *GolfClubs) Configure(cfg Config) error {
	c.reverse = cfg.Reverse
	c.pulsesPerCircle = cfg.PulsesPerCircle
	if c.pulsesPerCircle == 0 {
		c.pulsesPerCircle = DefaultPulsesPerCircle
	}

	// 配置 GPIO
	c.DirPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	c.DirPin.Low()
	c.EnPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	c.EnPin.High() // 先禁用

	// 获取 PWM 脚对应的 PWM 组
	pwmI, err := machine.PWMPeripheral(c.PWMPin)
	if err != nil {
		return err
	}
	c.pwm = getPWMGroup(pwmI)

	// 配置 PWM
	if err := c.pwm.Configure(machine.PWMConfig{
		Period: 1e9 / 1000,
	}); err != nil {
		return err
	}
	if c.pwmCh, err = c.pwm.Channel(c.PWMPin); err != nil {
		return err
	}
	c.pwm.Set(c.pwmCh, c.pwm.Top()/2)

	return nil
}

// SetReverse 设置是否反向挥杆
func (c *GolfClubs) SetReverse(reverse bool) {
	c.reverse = reverse
}

// SetPulsesPerCircle 设置电机旋转一周所需脉冲数
func (c *GolfClubs) SetPulsesPerCircle(pulsesPerCircle uint32) {
	c.pulsesPerCircle = pulsesPerCircle
}

// setDirFront 向前挥杆
func (c *GolfClubs) setDirFront() {
	c.DirPin.Set(c.reverse)
}

// setDirBack 向后挥杆
func (c *GolfClubs) setDirBack() {
	c.DirPin.Set(!c.reverse)
}

// Swing 挥杆一次
func (c *GolfClubs) Swing(speedPercent uint8) {
	c.hold()
	c.EnPin.Low()

	// 向后摆 38% 圈
	c.setDirBack()
	c.swingRaw(10, 38)
	c.hold()
	time.Sleep(time.Second)

	// 挥杆 76% 圈
	c.setDirFront()
	c.swingRaw(speedPercent/8, 3)
	c.swingRaw(speedPercent/4, 3)
	c.swingRaw(speedPercent/2, 6)
	c.swingRaw(uint8(uint32(speedPercent)*8/10), 13)
	c.swingRaw(speedPercent, 25)
	c.swingRaw(uint8(uint32(speedPercent)*8/10), 13)
	c.swingRaw(speedPercent/2, 6)
	c.swingRaw(speedPercent/4, 6)
	c.EnPin.High()
}

// swingRaw 挥杆
func (c *GolfClubs) swingRaw(speedPercent uint8, ringPercent uint8) {
	period := 1e9 * 60 / uint64(MaxSpeed*uint32(speedPercent)/100) / uint64(c.pulsesPerCircle)
	d := time.Duration(uint64(c.pulsesPerCircle*uint32(ringPercent)/100) * period)

	// 设置旋转速度
	if err := c.pwm.SetPeriod(period); err != nil {
		log.Printf("ERROR set pwm period to %d error: %v", period, err)
		return
	}

	// 挥
	c.pwm.Set(c.pwmCh, c.pwm.Top()/2)
	time.Sleep(d)
}

// hold 停住球杆
func (c *GolfClubs) hold() {
	c.pwm.Set(c.pwmCh, 0)
}
