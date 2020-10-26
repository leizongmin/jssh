package jsshcmd

import (
	"github.com/leizongmin/go/typeutil"
	"github.com/leizongmin/jssh/internal/jsexecutor"
	"net"
)

func JsFnNetworkinterfaces(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			return ctx.ThrowError(err)
		}
		list := make([]typeutil.H, 0)
		for _, addr := range addrs {
			ip, ok := addr.(*net.IPNet)
			if ok {
				family := "IPv4"
				if ip.IP.To16() != nil {
					family = "IPv6"
				}
				list = append(list, typeutil.H{
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

		return jsexecutor.AnyToJSValue(ctx, list)
	}
}
