package supermodel

// ServiceModel 服务模型
type ServiceModel struct {
	Domain             string    // 域名
	DefaultDestination string    // 默认路由目标
	Clusters           []Cluster // 集群节点列表
}

// RouteModel 路由模型
type RouteModel struct {
	Domain string      // 域名
	Policy RoutePolicy // 路由策略
}
