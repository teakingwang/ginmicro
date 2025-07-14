package config

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
	"os"
)

var Config *config

type config struct {
	Server   ServerConfig   `yaml:"server" json:"server"`
	Database DatabaseConfig `yaml:"database" json:"database"`
	Redis    RedisConfig    `yaml:"redis" json:"redis"`
	SMS      SMSConfig      `yaml:"sms" json:"sms"`
	JWT      JWTConfig      `yaml:"jwt" json:"jwt"`
	RocketMQ RocketMQConfig `yaml:"rocketmq" json:"rocketmq"`
	Consul   ConsulConfig   `yaml:"consul" json:"consul"` // 新增 Consul 配置结构体
	Kafka    KafkaConfig    `yaml:"kafka" json:"kafka"`   // 新增 Kafka 配置结构体
}

type KafkaConfig struct {
	Brokers []string `yaml:"brokers" json:"brokers"`
	Topic   string   `yaml:"topic" json:"topic"`
	GroupID string   `yaml:"groupId" json:"groupId"` // Kafka 消费者组 ID
}

type ServerConfig struct {
	User    ServiceConfig `yaml:"user" json:"user"`
	Order   ServiceConfig `yaml:"order" json:"order"`
	Gateway ServiceConfig `yaml:"gateway" json:"gateway"` // 新增 Gateway 服务配置
}

type ServiceConfig struct {
	Host     string `yaml:"host" json:"host"`
	GRPCPort string `yaml:"grpcPort" json:"grpcPort"` // gRPC 服务端口
	Name     string `yaml:"name" json:"name"`         // 服务名称
	HTTPPort string `yaml:"httpPort" json:"httpPort"` // HTTP 服务端口
}

type DatabaseConfig struct {
	Dialect  string `yaml:"dialect" json:"dialect"`
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	Database string `yaml:"database" json:"database"`
	User     string `yaml:"user" json:"user"`
	Password string `yaml:"password" json:"password"`
	Schema   string `yaml:"schema" json:"schema"`
	Level    string `yaml:"level" json:"level"`
}

type RedisConfig struct {
	Host         string `yaml:"host" json:"host"`
	Port         int    `yaml:"port" json:"port"`
	Password     string `yaml:"password" json:"password"`
	DB           int    `yaml:"db" json:"db"`
	PoolSize     int    `yaml:"poolSize" json:"poolSize"`
	MinIdleConns int    `yaml:"minIdleConns" json:"minIdleConns"`
}

type SMSConfig struct {
	Provider          string `yaml:"provider" json:"provider"`
	AccessKeyID       string `yaml:"accessKeyId" json:"accessKeyId"`
	AccessKeySecret   string `yaml:"accessKeySecret" json:"accessKeySecret"`
	SignName          string `yaml:"signName" json:"signName"`
	TemplateCode      string `yaml:"templateCode" json:"templateCode"`
	Region            string `yaml:"region" json:"region"`
	CodeExpireSeconds int    `yaml:"codeExpireSeconds" json:"codeExpireSeconds"`
}

type JWTConfig struct {
	Secret        string `yaml:"secret" json:"secret"`
	ExpireSeconds int    `yaml:"expireSeconds" json:"expireSeconds"`
}

type RocketMQConfig struct {
	NameServer    string `yaml:"nameServer" json:"nameServer"`       // RocketMQ NameServer 地址
	GroupName     string `yaml:"groupName" json:"groupName"`         // Producer 分组
	ProducerTopic string `yaml:"producerTopic" json:"producerTopic"` // Producer 使用的 topic
	ConsumerTopic string `yaml:"consumerTopic" json:"consumerTopic"` // Consumer 订阅的 topic
	RetryTimes    int    `yaml:"retryTimes" json:"retryTimes"`       // 发送失败重试次数
}

// 新增 Consul 配置结构体
type ConsulConfig struct {
	Address string `yaml:"address" json:"address"`
	Token   string `yaml:"token" json:"token"` // 可选，ACL Token
}

func LoadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	path := os.Getenv("MY_APP_CONFIG_PATH")
	if path == "" {
		path = "./resources"
	}
	viper.AddConfigPath(path)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	var cfg config
	if err := viper.Unmarshal(&cfg); err != nil {
		return err
	}
	Config = &cfg

	return nil
}

// 从 Consul 加载
func LoadConfigFromConsul(key string) error {
	consulAddr := os.Getenv("CONSUL_ADDR")
	if consulAddr == "" {
		consulAddr = "127.0.0.1:8500" // 默认地址
	}

	cli, err := api.NewClient(&api.Config{Address: consulAddr})
	if err != nil {
		return fmt.Errorf("create consul client failed: %w", err)
	}

	kv := cli.KV()
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return fmt.Errorf("get config from consul failed: %w", err)
	}
	if pair == nil {
		return fmt.Errorf("config not found in consul for key: %s", key)
	}

	var c config
	if err := json.Unmarshal(pair.Value, &c); err != nil {
		return fmt.Errorf("unmarshal config json failed: %w", err)
	}
	Config = &c
	return nil
}

func GetServiceID() string {
	if v := os.Getenv("SERVICE_ID"); v != "" {
		return v
	}
	return "default-service-id"
}

func GetServiceName() string {
	if v := os.Getenv("SERVICE_NAME"); v != "" {
		return v
	}
	return "default-service"
}

func GetServiceAddress() string {
	// 优先取环境变量
	if addr := os.Getenv("SERVICE_ADDRESS"); addr != "" {
		return addr
	}

	// fallback：读取配置文件
	return Config.Server.Gateway.Host
}
