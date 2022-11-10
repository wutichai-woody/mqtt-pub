package service

import (
	"techberry-go/common/v2/facade"
	"techberry-go/common/v2/facade/pdk"
	"techberry-go/micronode/service/accessors"
	"techberry-go/micronode/service/commons"
	"techberry-go/micronode/service/controllers"

	"github.com/spf13/viper"
)

type ServicePlugin struct {
	Credential        *pdk.RequestCredential
	TraceId           string
	Logger            pdk.Logger
	Context           pdk.Context
	Connector         pdk.Connector
	Config            *viper.Viper
	Handler           facade.Handler
	Redis             facade.CacheHandler
	Mqtt              adapter.MqttPool
	Version           string
	ServiceController *controllers.ServiceController
	ServiceAccessor   *accessors.ServiceAccessor
	ServiceCommon     *commons.ServiceCommon
}
