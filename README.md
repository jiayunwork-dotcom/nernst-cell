# nernst-cell — 电极电位与极化核算

nernst-cell 是一个电化学核算工具：用户输入标准电位 E°、反应电子数 n、温度与氧化/还原活度，按 Nernst 方程算出平衡电位 E 与 Nernst 斜率（mV/decade）；再输入交换电流密度与传递系数，按 Butler–Volmer 方程算出过电位下的电流 i(η)，并与 Tafel 极限对照。能力边界：Nernst 平衡电位 + BV 极化（单电子/多电子按 α_a + α_c = n 约定，默认 α = 0.5、n = 1），不做物质输运、不做阻抗谱。所有计算走同一个 CODATA 常数与同一套温度换算，跨接口结果一致。

## 用法

```bash
go run . -http :8080
```

打开 http://localhost:8080 进入交互页面，或直接用 curl 调用 API：

```bash
curl -s -X POST http://localhost:8080/api/nernst \
  -H 'Content-Type: application/json' \
  --data-binary @example/cu-conc.json
```

`example/cu-conc.json` 是铜浓差电池算例：两极同为 Cu/Cu²⁺，E° 抵消只剩浓度比，算出的电位为正，浓侧为正极。页面左上角可一键加载该算例。

## 关键约定

- **Nernst 符号书写**：`E = E° + (RT/(nF))·ln(a_ox / a_red)`，还原态写在分母，与 E° 表（还原反应式）一致。
- **常数**：R = 8.314462618 J/(mol·K)、F = 96485.33212 C/mol，均为 CODATA 值；温度以摄氏输入，内部换算为开尔文。
- **浓差电池**：两极同种、E° 相消，只剩 `E = (RT/(nF))·ln(a_浓 / a_稀)`，电位为正且浓侧为正极。
- **Butler–Volmer**：`i = i0·(exp(α_a·F·η/(RT)) − exp(−α_c·F·η/(RT)))`，其中 `α_c = n − α_a`；η = 0 时净电流为 0。
- **Tafel**：阳极 η 足够正时 `η = a + b·log10|i|`，`b = 2.303·RT/(αF)`；该式是 BV 在 `|η| ≫ RT/(αF)` 的极限，API 同时返回 BV 电流与 Tafel 参考电流以便对照。
- **校验**：n ≤ 0、T ≤ 0、活度 ≤ 0、i0 ≤ 0、α 越界（≤0 或 ≥n）均返回 error JSON，不做静默错值。

## 构建与测试

```bash
go build ./...
go test ./...
```

## 许可

MIT。
