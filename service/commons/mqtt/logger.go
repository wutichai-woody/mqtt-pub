package mqtt

import (
	"fmt"
	"techberry-go/common/v2/facade/pdk"
)

func Debugf(_logger pdk.Logger, format string, v ...interface{}) {
	if _logger != nil {
		_logger.Info().Msgf(format, v...)
	} else {
		fmt.Printf(format, v...)
	}
}

func Errorf(_logger pdk.Logger, format string, v ...interface{}) {
	if _logger != nil {
		_logger.Error().Msgf(format, v...)
	} else {
		fmt.Printf(format, v...)
	}
}
