package meta

import (
	"context"
	"net"
	"os"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	uuid "github.com/satori/go.uuid"
)

var (
	startedAt time.Time
	server    *Server
	ipv4      = getIpv4()
)

func init() {
	ctx := gctx.GetInitCtx()
	hostName, _ := os.Hostname()
	startedAt = time.Now()
	server = &Server{
		ServiceName: tryGetCfgString(ctx, "meta.service_name", "undefined-server"),
		HostName:    hostName,
		InstanceId:  tryGetCfgString(ctx, "meta.uuid", uuid.NewV4().String()),
	}
}

func tryGetCfgString(ctx context.Context, pattern string, def string) string {
	cfg, err := g.Cfg().Get(ctx, pattern, def)
	if err != nil {
		return ""
	}
	return cfg.String()
}

func getIpv4() string {
	ip, err := getIp()
	if err != nil {
		g.Log().Warningf(context.Background(), "failed to get ipv4: %v", err)
		return ""
	}
	return ip.To4().String()
}

func getIp() (v4 net.IP, err error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return
	}
	var addresses []net.Addr
	for _, inf := range interfaces {
		if inf.Flags|net.FlagUp > 0 && inf.Flags|net.FlagRunning > 0 {
			addresses, err = inf.Addrs()
			if err != nil {
				return
			}
			for _, addr := range addresses {
				if ipAddr, ok := addr.(*net.IPNet); ok &&
					ipAddr.IP != nil &&
					ipAddr.IP.To4() != nil &&
					!ipAddr.IP.IsLoopback() {
					v4 = ipAddr.IP
					return
				}
			}
		}
	}

	return
}
