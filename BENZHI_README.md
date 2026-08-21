# nernst-cell — 电极电位与极化核算

nernst-cell 是电化学核算工具：输入标准电位、电子数、温度与氧化/还原活度，按 Nernst 方程计算平衡电位与 mV/decade 斜率；再输入交换电流密度与传递系数，按 Butler–Volmer 计算过电位下的电流并给出 Tafel 对照。内置铜浓差电池算例。

## 构建 / 运行 / 测试

```text
go build ./...     # 编译
go run . -http :8080   # 启动 Web 控制台，页面可加载 example/cu-conc.json
go test ./...      # 测试
```

## 评测镜像

本目录评测专用文件（勿覆盖项目自带 Dockerfile/README）：

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md`（本文件）

两种架构都要构建并进容器验证：

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/arm64
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -it <image-name>:latest
```
