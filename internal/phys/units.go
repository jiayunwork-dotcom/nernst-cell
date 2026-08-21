package phys

// Volt 表示电位值，单位 V。包内用裸 float64 计算，
// 类型包装只用于区分语义、便于校验与格式化。
type Volt float64

// Millivolt 表示电位值，单位 mV。
type Millivolt float64

// CurrentDensity 表示电流密度，单位 A/cm²。
type CurrentDensity float64

// V 构造伏特值。
func V(x float64) Volt { return Volt(x) }

// MV 构造毫伏值。
func MV(x float64) Millivolt { return Millivolt(x) }

// ToMV 把伏特换算为毫伏。
func (v Volt) ToMV() Millivolt { return Millivolt(float64(v) * 1000) }

// ToV 把毫伏换算为伏特。
func (m Millivolt) ToV() Volt { return Volt(float64(m) / 1000) }

// Add 返回两个电位的和。
func (v Volt) Add(o Volt) Volt { return Volt(float64(v) + float64(o)) }

// Sub 返回两个电位的差。
func (v Volt) Sub(o Volt) Volt { return Volt(float64(v) - float64(o)) }

// Scaled 返回电位乘以系数后的值。
func (v Volt) Scaled(k float64) Volt { return Volt(float64(v) * k) }
