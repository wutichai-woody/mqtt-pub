package controllers

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"techberry-go/micronode/service/commons/mqtt"
	"time"
)

var poolManager *mqtt.PoolManager
var once sync.Once

func init() {
	once.Do(func() {
		poolManager = mqtt.New()
		defer poolManager.Close()
	})

}

func (c *ServiceController) Reload(input any) (any, error) {
	if poolManager != nil {
		poolManager.Close()
	} else {
		poolManager = mqtt.New()
	}
	return input, nil
}

func (c *ServiceController) Echo(input any) (any, error) {
	fmt.Printf("pool manager : %p\n", poolManager)
	return input, nil
}

func (c *ServiceController) Execute(input any) (any, error) {
	m := input.(map[string]any)
	pool_cfg, err := c.getPoolConfig(input)
	if err != nil {
		return make(map[string]any), err
	}
	var topics []string
	var message string
	mapHandler := c.Handler.Map(false)
	t, ok := mapHandler.GetByPath("topic", m)
	if ok {
		switch v := t.(type) {
		case []string:
			topics = v
		case string:
			topics = append(topics, v)
		}
	} else {
		return make(map[string]any), errors.New("No topic.")
	}
	t, ok = mapHandler.GetByPath("message", m)
	if ok {
		switch v := t.(type) {
		case string:
			message = v
		}
	} else {
		return make(map[string]any), errors.New("No message.")
	}
	if len(topics) > 0 && message != "" {
		pool := poolManager.GetMqttPool(c.Logger, pool_cfg)
		ctx := context.Background()
		obj1, err := pool.BorrowObject(ctx)
		if err != nil {
			return make(map[string]any), err
		}
		o := obj1.(*mqtt.MqttPoolObject)
		for _, topic := range topics {
			if pool_cfg.Retained {
				o.Client.SetQoS(pool_cfg.QoS)
				o.Client.PublishAndRetain(topic, []byte(message))
			} else {
				o.Client.PublishWithQOS(topic, pool_cfg.QoS, []byte(message))
			}
		}
		err = pool.ReturnObject(ctx, obj1)
		if err != nil {
			return make(map[string]any), err
		}
		result_map := map[string]any{
			"status": true,
		}
		return result_map, nil
	}
	return input, nil
}

func (c *ServiceController) getPoolConfig(input interface{}) (*mqtt.PoolConfig, error) {
	m := input.(map[string]any)
	mapHandler := c.Handler.Map(false)
	url := mapHandler.String(m, "credential.url", "")
	username := mapHandler.String(m, "credential.username", "")
	password := mapHandler.String(m, "credential.password", "")
	retained := mapHandler.Bool(m, "credential.retained", false)
	qos := mapHandler.Int(m, "credential.qos", 0)

	if url != "" && username != "" && password != "" {
		return nil, errors.New("No credential.")
	}
	config := mqtt.PoolConfig{
		Url:      url,
		Username: username,
		Password: password,
		ClientId: getClientId("pubsub"),
		Retained: retained,
		QoS:      qos,
	}

	return &config, nil
}

func errorStrMessage(status_code int, error_code string, err error) string {
	return fmt.Sprintf("{\"error\": {\"code\": \"%s\", \"status_code\": %d, \"message\": \"%s\"} }", error_code, status_code, err.Error())
}

func getClientId(clientPrefix string) string {
	rgen := rand.New(rand.NewSource(time.Now().UnixNano()))
	// Generate random client id if one wasn't supplied.
	b := make([]byte, 16)
	rgen.Read(b)
	return fmt.Sprintf("%s-%x", clientPrefix, b)
}
