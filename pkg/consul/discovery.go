// pkg/discovery/consul_discovery.go
package consul

import (
	"fmt"
	"google.golang.org/grpc/connectivity"
	"sync"

	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
)

type Discovery interface {
	// GetGRPCConn 获取指定服务的 gRPC 连接
	GetGRPCConn(serviceName string) (*grpc.ClientConn, error)
}

type consulDiscovery struct {
	client     *api.Client
	connCache  map[string]*grpc.ClientConn
	cacheMutex sync.Mutex
	consulAddr string
}

func NewConsulDiscovery(addr string) (Discovery, error) {
	config := api.DefaultConfig()
	config.Address = addr

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &consulDiscovery{
		client:     client,
		connCache:  make(map[string]*grpc.ClientConn),
		consulAddr: addr,
	}, nil
}

func (d *consulDiscovery) GetGRPCConn(serviceName string) (*grpc.ClientConn, error) {
	d.cacheMutex.Lock()
	if conn, ok := d.connCache[serviceName]; ok {
		if conn.GetState() != connectivity.Shutdown {
			d.cacheMutex.Unlock()
			return conn, nil
		}
		// 删除已关闭连接
		delete(d.connCache, serviceName)
	}
	d.cacheMutex.Unlock()

	// 查找服务实例
	entries, _, err := d.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, fmt.Errorf("consul discover error: %w", err)
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("no healthy instance for service: %s", serviceName)
	}

	svc := entries[0].Service
	target := fmt.Sprintf("%s:%d", svc.Address, svc.Port)

	// 建立新连接
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s: %w", target, err)
	}

	d.cacheMutex.Lock()
	d.connCache[serviceName] = conn
	d.cacheMutex.Unlock()

	return conn, nil
}
