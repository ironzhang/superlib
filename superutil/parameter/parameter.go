package parameter

import (
	"os"
	"path/filepath"

	"github.com/ironzhang/tlog"

	"github.com/ironzhang/superlib/fileutil"
)

const (
	defaultAgentServer            = "127.0.0.1:1789" // 默认代理地址
	defaultAgentTimeout           = 2                // 2秒
	defaultAgentKeepAliveTTL      = 10 * 60          // 10分钟
	defaultAgentKeepAliveInterval = 10               // 10秒
	defaultWatchInterval          = 1                // 1秒
)

// Agent 代理参数
type Agent struct {
	Server            string // 服务器地址
	SkipError         bool   // 忽略接口调用错误
	Timeout           int    // 接口调用超时，单位：秒
	KeepAliveTTL      int    // 保活过期时间，单位：秒
	KeepAliveInterval int    // 保活间隔时间，单位：秒
}

// Watch 订阅参数
type Watch struct {
	ResourcePath  string // 资源路径
	WatchInterval int    // 订阅间隔，单位：秒
}

// Parameter 配置参数
type Parameter struct {
	Agent Agent // 代理参数
	Watch Watch // 订阅参数
}

// Param 全局配置参数
var Param Parameter

func init() {
	Param = loadParameter()
}

func getSuperNamePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "/var/supername"
	}
	return filepath.Join(home, ".supername")
}

func getDefaultResourcePath() string {
	return filepath.Join(getSuperNamePath(), "resource")
}

func getUserSuperNameConfPath() string {
	return filepath.Join(getSuperNamePath(), "supername.conf")
}

func getGlobalSuperNameConfPath() string {
	return "/etc/supername.conf"
}

func getDefaultParameter() Parameter {
	return Parameter{
		Agent: Agent{
			Server:            defaultAgentServer,
			SkipError:         true,
			Timeout:           defaultAgentTimeout,
			KeepAliveTTL:      defaultAgentKeepAliveTTL,
			KeepAliveInterval: defaultAgentKeepAliveInterval,
		},
		Watch: Watch{
			ResourcePath:  getDefaultResourcePath(),
			WatchInterval: defaultWatchInterval,
		},
	}
}

func readParameter() Parameter {
	param := getDefaultParameter()
	paths := []string{getGlobalSuperNameConfPath(), getUserSuperNameConfPath()}
	for _, path := range paths {
		if fileutil.FileExist(path) {
			err := fileutil.ReadTOML(path, &param)
			if err != nil {
				tlog.Errorw("read toml", "path", path, "error", err)
			}
			return param
		}
	}
	return param
}

func loadParameter() Parameter {
	param := readParameter()
	if param.Agent.Timeout < 0 {
		param.Agent.Timeout = defaultAgentTimeout
	}
	if param.Agent.KeepAliveInterval <= 0 {
		param.Agent.KeepAliveInterval = defaultAgentKeepAliveInterval
	}
	if param.Agent.KeepAliveTTL < 3*param.Agent.KeepAliveInterval {
		param.Agent.KeepAliveTTL = 3 * param.Agent.KeepAliveInterval
	}
	if param.Watch.WatchInterval <= 0 {
		param.Watch.WatchInterval = defaultWatchInterval
	}
	return param
}
