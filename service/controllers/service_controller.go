package controllers

import (
	"techberry-go/common/v2/facade"
	"techberry-go/common/v2/facade/adapter"
	"techberry-go/common/v2/facade/pdk"
	"techberry-go/micronode/service/accessors"
	"techberry-go/micronode/service/commons"

	"github.com/spf13/viper"
)

type ServiceController struct {
	Credential      *pdk.RequestCredential
	TraceId         string
	Logger          pdk.Logger
	Context         pdk.Context
	Connector       pdk.Connector
	Config          *viper.Viper
	Handler         facade.Handler
	Redis           facade.CacheHandler
	MqttPool        adapter.MqttPool
	Version         string
	ServiceAccessor *accessors.ServiceAccessor
	ServiceCommon   *commons.ServiceCommon
}
