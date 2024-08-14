package supermodel

// State 地址状态
type State string

// 地址状态常量定义
const (
	Enabled  State = "enabled"
	Disabled State = "disabled"
)

// Endpoint 地址节点
type Endpoint struct {
	Addr   string // 地址，IP:Port
	State  State  // 状态
	Weight int    // 权重
}

// Cluster 集群节点
type Cluster struct {
	Name      string            // 集群名
	Labels    map[string]string // 集群标签
	Endpoints []Endpoint        // 地址节点列表
}

// Destination 目标节点
type Destination struct {
	Cluster string  // 目标集群名
	Percent float64 // 流量配比
}

// RoutePolicy 路由策略
type RoutePolicy struct {
	EnableScript   bool            // 是否启用脚本
	LabelSelectors []LabelSelector // 标签选择器列表
}
