package mqtt

import (
	"context"
	"fmt"
	"testing"
)

func TestConnection(t *testing.T) {
	config := PoolConfig{
		Url:      "mqtt://3.0.31.237:7883",
		ClientId: "abc",
		Username: "mqttuser",
		Password: "Admin#1$0",
	}
	poolManager := New()
	pool := poolManager.GetMqttPool(nil, &config)

	ctx := context.Background()
	obj1, err := pool.BorrowObject(ctx)
	if err != nil {
		panic(err)
	}
	o1 := obj1.(*MqttPoolObject)
	fmt.Printf("obj1 : %p\n", o1.Client)

	obj2, err := pool.BorrowObject(ctx)
	if err != nil {
		panic(err)
	}
	o2 := obj2.(*MqttPoolObject)
	fmt.Printf("obj2 : %p\n", o2.Client)

	err = o2.Client.PublishWithQOS("test_topic", 0, []byte("11111abcdefgh1111"))
	fmt.Printf("error : %v\n", err)
	err = pool.ReturnObject(ctx, obj1)
	if err != nil {
		panic(err)
	}
	err = pool.ReturnObject(ctx, obj2)
	if err != nil {
		panic(err)
	}
	close(pool.StoppedChannel())
}
