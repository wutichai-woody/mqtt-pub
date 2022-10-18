package mqtt

import "context"

type PoolConfig struct {
	Url      string
	ClientId string
	Username string
	Password string
	UseSSL   bool
	Retained bool
	QoS      int
	Hostname string
	MaxTotal int
	MaxIdle  int
	MinIdle  int
}

type PoolManager struct {
	maxTotal int
	maxIdle  int
	minIdle  int
}

type MqttPool interface {
	BorrowObject(ctx context.Context) (interface{}, error)
	ReturnObject(ctx context.Context, object interface{}) error
	StoppedChannel() chan struct{}
	Hostname() string
}
