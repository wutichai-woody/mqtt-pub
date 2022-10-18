package controllers

import (
	"techberry-go/common/v2/facade"
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
	ServiceNode     pdk.ServiceNode
	Config          *viper.Viper
	Handler         facade.Handler
	Version         string
	ServiceAccessor *accessors.ServiceAccessor
	ServiceCommon   *commons.ServiceCommon
}
