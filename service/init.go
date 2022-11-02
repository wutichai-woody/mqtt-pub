package service

import (
	"techberry-go/common/v2/facade/adapter"
	"techberry-go/common/v2/facade/pdk"

	"github.com/spf13/viper"
)

func LoadOnce() func(adapter.LogEvent, pdk.Connector, *viper.Viper) {
	return func(logger adapter.LogEvent, connector pdk.Connector, config *viper.Viper) {
		host := config.GetString("services.redis.host")
		enable := config.GetBool("services.redis.enable")
		port := config.GetInt("services.redis.port")
		dbnum := config.GetInt("services.redis.dbnum")
		poolSize := config.GetInt("services.redis.poolSize")
		if enable {
			connector.GetRedisConnection(host, port, dbnum, poolSize)
		}
	}
}
