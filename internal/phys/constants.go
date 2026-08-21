// Package phys 提供电化学核算共用的物理常数、温度换算与
// 热电压等基础量。R 与 F 取 CODATA 推荐值，全仓唯一来源，
// 避免各包各自写死导致数值不一致。
package phys

// GasConstant 是摩尔气体常数 R，CODATA 2018 推荐值，
// 单位 J/(mol·K)。
const GasConstant = 8.314462618

// FaradayConstant 是法拉第常数 F，CODATA 2018 推荐值，
// 单位 C/mol。数值取 96485.33212，截断后与 2018 推荐值一致
// 到实验不确定度之内。
const FaradayConstant = 96485.33212

// R 返回摩尔气体常数，单位 J/(mol·K)。
// 保留直接字母别名，便于公式侧阅读：R*T/F。
func R() float64 {
	return GasConstant
}

// F 返回法拉第常数，单位 C/mol。
func F() float64 {
	return FaradayConstant
}
