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

/*
const DEFAULT_CONTEXT_TIMEOUT = 10 * time.Second

var client_chan chan adapter.MqttClient

var (
	MqttPools   map[string]MqttPool
	once        sync.Once
	stoppedChan chan struct{}
	logger      log.Modular
)

type MqttConfig struct {
	Url      string
	ClientId string
	Username string
	Password string
}

type MqttPool struct {
	client      adapter.MqttClient
	errors      chan error
	connections int
	config      *MqttConfig
	poolSize    int // Default pool size is 90. However the default pool size from official mongo driver is 100 for golang.
	mutex       *sync.Mutex
}

func init() {
	once.Do(func() {
		client_chan = make(chan adapter.MqttClient, 100)
		stoppedChan = make(chan struct{})
		registerListener()
	})
}

func NewMqttPool(_logger log.Modular) chan struct{} {
	logger = _logger
	return stoppedChan
}

func registerListener() {
	go func() {
		<-stoppedChan
		releasePool()
	}()
}

func getContextTimeOut() (context.Context, context.CancelFunc) {
	ctx, cancelFn := context.WithTimeout(context.Background(), DEFAULT_CONTEXT_TIMEOUT)
	return ctx, cancelFn
}

func (m *MqttPool) Disconnect() error {
	if err := m.client.Disconnect(); err != nil {
		if logger != nil {
			logger.Errorf("Close the mqtt client failed, err=%v", err)
		}
		return err
	}
	return nil
}

func (m *MqttPool) ReturnConnection(conn adapter.MqttClient) error {
	select {
	case client_chan <- conn:
		m.mutex.Lock()
		m.connections--
		m.mutex.Unlock()
		return nil
	}
}

func (m *MqttPool) createChannel() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	conn := mqtt.NewAdapterWithAuth(m.config.Url, m.config.ClientId, m.config.Username, m.config.Password)

	err := conn.Connect()
	if err == nil {
		m.client = conn
		client_chan <- conn
		m.connections++
	} else {
		m.errors <- errors.New("can't create connection.")
	}
}

func (m *MqttPool) GetConnection() (adapter.MqttClient, error) {
	for {
		select {
		case client := <-client_chan:
			if client.IsConnected() && client.IsConnectionOpen() {
				return client, nil
			} else {
				return nil, errors.New("connection is close or disconnected.")
			}
		case err := <-m.errors:
			return nil, err
		default:
			if m.connections < m.poolSize {
				m.createChannel()
			}
		}
	}
}

func releasePool() {
	if logger != nil {
		logger.Infof("Release mqtt connections from pool.")
	}
	for _, v := range MqttPools {
		if err := v.client.Disconnect(); err != nil {
			logger.Errorf("Close the mqtt client failed, err=%v", err)
		}
	}
}
*/

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
		Debugf(logger, "mqtt pool key : %s exists.", key)
		return v
	} else {
		Debugf(logger, "create mqtt pool key : %s.", key)
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
	Debugf(logger, "create mqtt pooled object : %s, %s", m.url, m.config.ClientId)
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
