package consul

import (
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/teakingwang/ginmicro/pkg/logger"
	"time"
)

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
	logger.Info("Registering service %s with ID %s at %s:%d", name, id, address, port)
	registration := &api.AgentServiceRegistration{
		ID:      id,
		Name:    name,
		Address: address,
		Port:    port,
		Tags:    tags,
		Check: &api.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d", address, port), // 不加 grpc:// 前缀
			Interval:                       "10s",
			Timeout:                        "3s",
			DeregisterCriticalServiceAfter: "1m",
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
