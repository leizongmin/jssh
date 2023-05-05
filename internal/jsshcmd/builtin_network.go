package jsshcmd

import (
	"net"

	"github.com/leizongmin/jssh/internal/jsexecutor"
	"github.com/leizongmin/jssh/internal/utils"
)

func jsFnNetworkinterfaces(global utils.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		interfaces, err := net.Interfaces()
		if err != nil {
			return ctx.ThrowError(err)
		}

		ret := make(utils.H)
		for _, item := range interfaces {
			addrs, err := item.Addrs()
			if err != nil {
				return ctx.ThrowError(err)
			}

			list := make([]utils.H, 0)
			for _, addr := range addrs {
				ip, ok := addr.(*net.IPNet)
				if ok {
					family := "IPv4"
					if ip.IP.To16() != nil {
						family = "IPv6"
					}
					list = append(list, utils.H{
						"address":     ip.IP.String(),
						"netmask":     ip.Mask.String(),
						"family":      family,
						"cidr":        ip.String(),
						"internal":    ip.IP.IsLoopback(),
						"multicast":   ip.IP.IsMulticast(),
						"unspecified": ip.IP.IsUnspecified(),
					})
				}
			}

			ret[item.Name] = utils.H{
				"index": item.Index,
				"mac":   item.HardwareAddr.String(),
				"list":  list,
			}
		}

		return jsexecutor.AnyToJSValue(ctx, ret)
	}
}
