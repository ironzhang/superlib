package parameter

import (
	"os"
	"path"

	"github.com/ironzhang/tlog"

	"github.com/ironzhang/superlib/fileutil"
)

// Parameter 配置参数
type Parameter struct {
	AgentServer   string // 服务器地址
	ResourcePath  string // 资源路径
	WatchInterval int    // 订阅间隔，单位：秒
}

// Param 全局配置参数
var Param Parameter

func init() {
	Param = readParameter()
}

func getDefaultResourcePath() string {
	home, err := os.UserHomeDir()
	if err == nil {
		return path.Join(home, ".superdns")
	}
	return "/var/superdns"
}

func getDefaultParameter() Parameter {
	return Parameter{
		AgentServer:   "127.0.0.1:1789",
		ResourcePath:  getDefaultResourcePath(),
		WatchInterval: 1,
	}
}

func readParameter() Parameter {
	param := getDefaultParameter()

	const path = "/etc/superdns.conf"
	if fileutil.FileExist(path) {
		err := fileutil.ReadTOML(path, &param)
		if err != nil {
			tlog.Errorw("read toml", "path", path, "error", err)
		}
	}

	if param.WatchInterval <= 0 {
		param.WatchInterval = 1
	}

	return param
}
