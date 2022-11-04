package mqtt

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"techberry-go/common/v2/core/components/mqtt"
	"techberry-go/common/v2/facade/pdk"
	"time"

	pool "github.com/jolestar/go-commons-pool/v2"
)

var (
	mqttPools   map[string]MqttPool
	once        sync.Once
	logger      pdk.Logger
	stoppedChan chan struct{}
)

type MqttFactory struct {
	url         string
	clientId    string
	username    string
	password    string
	useSSL      bool
	config      *PoolConfig
	objectPool  *pool.ObjectPool
	stoppedChan chan struct{}
}

func init() {
	once.Do(func() {
		mqttPools = make(map[string]MqttPool)
	})
}

func New() *PoolManager {
	return &PoolManager{
		maxTotal: 100,
		minIdle:  1,
		maxIdle:  2,
	}
}

func (pMgr *PoolManager) GetMqttPool(_logger pdk.Logger, config *PoolConfig) MqttPool {
	logger = _logger
	replacer := strings.NewReplacer("://", "_", ".", "_", ":", "_")
	url := replacer.Replace(config.Url)
	clientId := ""
	if config.ClientId != "" {
		clientId = config.ClientId
	} else {
		rgen := rand.New(rand.NewSource(time.Now().UnixNano()))
		// Generate random client id if one wasn't supplied.
		b := make([]byte, 16)
		rgen.Read(b)
		clientId = string(b)
	}
	key := fmt.Sprintf("%s_%s", url, clientId)
	if v, ok := mqttPools[key]; ok {
		Debugf(logger, "-> Mqtt pool key : %s exists.\n", key)
		return v
	} else {
		Debugf(logger, "-> Create mqtt pool for key : %s.\n", key)
	}
	ctx := context.Background()
	factory := &MqttFactory{
		url:         config.Url,
		clientId:    config.ClientId,
		username:    config.Username,
		password:    config.Password,
		useSSL:      config.UseSSL,
		config:      config,
		stoppedChan: make(chan struct{}),
	}
	// Default value
	if config.MaxTotal == 0 {
		pMgr.maxTotal = 100
		config.MaxTotal = 100
	}
	if config.MinIdle == 0 {
		pMgr.minIdle = 1
		config.MinIdle = 1
	}
	if config.MaxIdle == 0 {
		pMgr.maxIdle = 2
		config.MaxIdle = 2
	}
	objConfig := &pool.ObjectPoolConfig{
		MaxTotal: config.MaxTotal,
		MaxIdle:  config.MaxIdle,
		MinIdle:  config.MinIdle,
	}
	p := pool.NewObjectPool(ctx, factory, objConfig)
	p.Config.MaxTotal = config.MaxTotal
	p.Config.MaxIdle = config.MaxIdle
	p.Config.MinIdle = config.MinIdle
	registerListener(factory)
	factory.objectPool = p
	mqttPools[key] = factory
	return factory
}

func (pMgr *PoolManager) Close() {
	for _, v := range mqttPools {
		t := v.(MqttPool)
		if t.StoppedChannel() != nil {
			close(t.StoppedChannel())
		}
	}
}

func (m *MqttFactory) MakeObject(ctx context.Context) (*pool.PooledObject, error) {
	Debugf(logger, "-> Create mqtt pooled object for : %s, %s\n", m.url, m.config.ClientId)
	client := mqtt.NewAdapterWithAuth(m.url, m.clientId, m.username, m.password)
	if m.useSSL {
		client.SetUseSSL(true)
	}
	client.SetAutoReconnect(true)
	client.SetCleanSession(true)
	client.SetQoS(m.config.QoS)
	err := client.Connect()
	if err != nil {
		return nil, err
	}
	return pool.NewPooledObject(
			&MqttPoolObject{
				Client: client,
			}),
		nil
}

func (f *MqttFactory) DestroyObject(ctx context.Context, object *pool.PooledObject) error {
	mqtt := object.Object.(*MqttPoolObject)
	mqtt.Client.Disconnect()
	return nil
}

func (f *MqttFactory) ValidateObject(ctx context.Context, object *pool.PooledObject) bool {
	// do validate
	return true
}

func (f *MqttFactory) ActivateObject(ctx context.Context, object *pool.PooledObject) error {
	// do activate
	return nil
}

func (f *MqttFactory) PassivateObject(ctx context.Context, object *pool.PooledObject) error {
	// do passivate
	return nil
}

func (f *MqttFactory) BorrowObject(ctx context.Context) (interface{}, error) {
	return f.objectPool.BorrowObject(ctx)
}

func (f *MqttFactory) ReturnObject(ctx context.Context, object interface{}) error {
	return f.objectPool.ReturnObject(ctx, object)
}

func (f *MqttFactory) StoppedChannel() chan struct{} {
	return f.stoppedChan
}

func (f *MqttFactory) Hostname() string {
	return f.config.Hostname
}

func registerListener(factory *MqttFactory) {
	go func() {
		<-factory.stoppedChan
		releasePool(factory)
	}()
}

func releasePool(factory *MqttFactory) {
	Debugf(logger, "Release mongo connections from pool.")
	ctx := context.Background()
	if factory.objectPool != nil {
		factory.objectPool.Close(ctx)
	}
}
