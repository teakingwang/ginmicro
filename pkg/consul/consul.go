package consul

import (
	"context"
	"fmt"
	"github.com/teakingwang/ginmicro/config"
	"net"

	"github.com/hashicorp/consul/api"
	"github.com/teakingwang/ginmicro/pkg/logger"
	"time"
)

// GetLocalIP 获取本机非回环 IP 地址
func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "127.0.0.1", nil
}

type ConsulClient struct {
	client *api.Client
	kv     *api.KV
}

func NewConsulClient(addr string) (*ConsulClient, error) {
	config := api.DefaultConfig()
	config.Address = addr
	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consul client: %w", err)
	}

	return &ConsulClient{
		client: client,
		kv:     client.KV(),
	}, nil
}

// PutKV 将数据写入 Consul KV，key 如 config/app.yaml
func (c *ConsulClient) PutKV(key string, value []byte) error {
	p := &api.KVPair{
		Key:   key,
		Value: value,
	}
	_, err := c.kv.Put(p, nil)
	if err != nil {
		return fmt.Errorf("failed to put kv: %w", err)
	}
	return nil
}

// GetKV 从 Consul KV 读取数据
func (c *ConsulClient) GetKV(key string) ([]byte, error) {
	pair, _, err := c.kv.Get(key, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get kv: %w", err)
	}
	if pair == nil {
		return nil, fmt.Errorf("key %s not found", key)
	}
	return pair.Value, nil
}

type ServiceRegistration struct {
	ID      string
	Name    string
	Address string
	Port    int
	Client  *api.Client
}

func (c *ConsulClient) RegisterService(id, name, address string, port int, tags []string) error {
	// 如果 address 是 0.0.0.0 或 127.0.0.1，使用 host.docker.internal
	// 因为 Consul 运行在 Docker 容器中，需要访问宿主机的服务
	effectiveAddress := address
	if address == "0.0.0.0" || address == "" || address == "127.0.0.1" {
		effectiveAddress = "host.docker.internal"
	}

	// 检查 Consul 是否在 Docker 中运行，如果是，使用网关 IP 让容器访问宿主机
	// 通过检查 Consul 地址是否为 localhost/127.0.0.1 来判断
	consulAddr := config.Config.Consul.Address
	checkAddress := effectiveAddress
	if consulAddr == "127.0.0.1:8500" || consulAddr == "localhost:8500" {
		// Consul 在本地，但可能是在 Docker 中
		// 使用 Docker 网关 IP 让容器可以访问宿主机上的服务
		// 通常 Docker bridge 网络的网关是 xxx.xxx.xxx.1
		checkAddress = "host.docker.internal"
		logger.Info("Consul appears to be running in Docker, using host.docker.internal for health check")
	}

	logger.Info("Registering service %s with ID %s at %s:%d (check addr: %s:%d)", name, id, address, port, effectiveAddress, port)
	registration := &api.AgentServiceRegistration{
		ID:      id,
		Name:    name,
		Address: effectiveAddress,
		Port:    port,
		Tags:    tags,
		Check: &api.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d/%s", checkAddress, port, name),
			Interval:                       "10s",
			Timeout:                        "5s",
			DeregisterCriticalServiceAfter: "30s",
		},
	}

	return c.client.Agent().ServiceRegister(registration)
}

func (c *ConsulClient) DeregisterService(id string) error {
	return c.client.Agent().ServiceDeregister(id)
}

// Example: Watch KV key changes (simple blocking query)
func (c *ConsulClient) WatchKey(ctx context.Context, key string, waitIndex uint64) ([]byte, uint64, error) {
	opts := &api.QueryOptions{
		WaitIndex: waitIndex,
		WaitTime:  5 * time.Minute,
	}
	pair, meta, err := c.kv.Get(key, opts.WithContext(ctx))
	if err != nil {
		return nil, waitIndex, err
	}
	if pair == nil {
		return nil, waitIndex, fmt.Errorf("key %s not found", key)
	}
	return pair.Value, meta.LastIndex, nil
}
