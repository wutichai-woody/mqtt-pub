package accessors

import (
	"techberry-go/common/v2/facade"
	"techberry-go/common/v2/facade/pdk"
	"techberry-go/micronode/service/commons"

	"github.com/spf13/viper"
)

type ServiceAccessor struct {
	Credential    *pdk.RequestCredential
	TraceId       string
	Logger        pdk.Logger
	Context       pdk.Context
	Connector     pdk.Connector
	Config        *viper.Viper
	Handler       facade.Handler
	Redis         facade.CacheHandler
	Mqtt          adapter.MqttPool
	Version       string
	ServiceCommon *commons.ServiceCommon
}
