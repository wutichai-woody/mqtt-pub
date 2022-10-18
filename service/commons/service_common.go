package commons

import (
	"techberry-go/common/v2/facade"
	"techberry-go/common/v2/facade/pdk"

	"github.com/spf13/viper"
)

type ServiceCommon struct {
	Credential  *pdk.RequestCredential
	TraceId     string
	Logger      pdk.Logger
	Context     pdk.Context
	Connector   pdk.Connector
	ServiceNode pdk.ServiceNode
	Config      *viper.Viper
	Handler     facade.Handler
	Version     string
}
