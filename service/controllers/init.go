package controllers

import (
	"techberry-go/common/v2/facade"
	"techberry-go/micronode/service/commons/mqtt"
)

var MqttPoolManager *mqtt.PoolManager
var RedisConn facade.CacheHandler
var MqttPool mqtt.MqttPool
