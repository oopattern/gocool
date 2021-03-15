package config

// 阿波罗配置中心: http://localhost:8070/
// 用户名: apollo
// 密码: admin

// consul UI页面: http://localhost:8500/ui

// skywalking backend port
const (
	SkyWalkingEndPoint = "127.0.0.1:11800"
	MetricsEndPoint = "0.0.0.0:9095"
	GatewayEndPoint = ":8006"
	AgolloEndPoint = "localhost:8080"
	ConsulEndPoint = "localhost:8500"
)