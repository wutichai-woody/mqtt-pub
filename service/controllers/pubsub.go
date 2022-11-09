package controllers

import (
	"context"
	"errors"
	"fmt"
	"techberry-go/micronode/service/commons/firebase"
	"techberry-go/micronode/service/commons/mqtt"
	"time"
)

func (c *ServiceController) Echo(input any) (any, error) {
	return input, nil
}

func (c *ServiceController) Notify(input any) (any, error) {
	m := input.(map[string]any)
	output, err := firebase.SendFirebaseNotification(m)
	if err != nil {
		b := []byte(fmt.Sprintf("{\"error\": \"%s\"}", err.Error()))
		return b, nil
	}
	return output, nil
}

func (c *ServiceController) SyncCache(input any) (any, error) {
	m := input.(map[string]any)
	if _, ok := m["data"]; ok {
		data := m["data"].([]map[string]interface{})
		key := ""
		value := ""
		expire := -1
		for _, val := range data {
			for k, v := range val {
				if k == "key" {
					key = v.(string)
				} else if k == "value" {
					value = v.(string)
				} else if k == "expire" {
					expire = v.(int)
				}
			}
			if key != "" && value != "" {
				if expire == -1 {
					expire = 10 * 365 * 24 * 60
				}
				RedisConn.Set(key, value, time.Duration(expire))
			}
		}
	}
	return make(map[string]any), nil
}

func (c *ServiceController) PublishMessageRead(input any) (any, error) {
	c.Logger.Info().Msgf("pubsub.publishMessageRead() called.")
	m := input.(map[string]any)
	c.Logger.Info().Msgf("pubsub.publishMessageRead() : %v", m)
	for k, v := range m {
		c.Logger.Info().Msgf("pubsub.publishMessageRead() key : %s, v : %v", k, v)
		topicToken, err := RedisConn.Get(k, false)
		if err == nil && topicToken != "" {
			c.Logger.Info().Msgf("found topic token : %s (%s)", topicToken, k)
			msg := map[string]any{
				"topics":  []string{topicToken},
				"message": c.ServiceCommon.GetMessageSignal("MESSAGE_READ", v),
			}
			c.Broadcast(msg)
		} else {
			c.Logger.Info().Msgf("pubsub.publishMessageRead() topic token not found.")
		}
	}
	return make(map[string]any), nil
}

func (c *ServiceController) Broadcast(input any) (any, error) {
	ctx := context.Background()
	obj1, err := MqttPool.BorrowObject(ctx)
	if err != nil {
		return make(map[string]any), err
	}

	m := input.(map[string]any)
	c.Logger.Info().Msgf("pubsub.broadcast() : %v", m)
	var topics []string
	mapHandler := c.Handler.Map(false)
	var t []any = mapHandler.GetArray(m, "topics")
	topics = make([]string, len(t))
	for i, v := range t {
		topics[i] = fmt.Sprint(v)
	}

	message := mapHandler.String(m, "message", "")
	if c.Handler.String(false).IsEmptyString(message) {
		return make(map[string]any), errors.New("No message.")
	}
	var o *mqtt.MqttPoolObject
	if len(topics) > 0 {
		o = obj1.(*mqtt.MqttPoolObject)
		for _, topic := range topics {
			fmt.Printf("topic : %s, message : %s\n", topic, message)
			o.Client.Publish(topic, []byte(message))
		}
		err := MqttPool.ReturnObject(ctx, obj1)
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

/*
localhost:8124/megw/apis/node/mqtt_pub/v1.0/publishMessageRead
{
	"topics": [
		"topicToken1",
		"topicToken2",
		....
	],
	"message": "xxxxxxxxxxxx"
}

{
	"redis_key1": [
		msg_id1,
		msg_id2,
		....
	],
	"redis_key2": [
		msg_id3,
		msg_id4,
		....
	]
}
*/

/*
func (c *ServiceController) Reload(input any) (any, error) {
	if poolManager != nil {
		poolManager.Close()
	} else {
		poolManager = mqtt.New()
	}
	return input, nil
}
*/

/*
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
*/

func errorStrMessage(status_code int, error_code string, err error) string {
	return fmt.Sprintf("{\"error\": {\"code\": \"%s\", \"status_code\": %d, \"message\": \"%s\"} }", error_code, status_code, err.Error())
}