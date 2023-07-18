package xnet

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func Local() (ip string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}
	return "", errors.New("could not find any ip")
}

// GetPeerIP 获取对端ip
func GetPeerIP(ctx context.Context) string {
	clientIP := GrpcHeaderValue(ctx, "client-ip")
	if clientIP != "" {
		return clientIP
	}
	// 从grpc里取对端ip
	pr, ok2 := peer.FromContext(ctx)
	if !ok2 {
		return ""
	}
	if pr.Addr == net.Addr(nil) {
		return ""
	}
	addSlice := strings.Split(pr.Addr.String(), ":")
	if len(addSlice) > 1 {
		return addSlice[0]
	}
	return ""
}

// GrpcHeaderValue 获取context value
func GrpcHeaderValue(ctx context.Context, key string) string {
	if key == "" {
		return ""
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	// 小写
	return strings.Join(md.Get(key), ";")
}
