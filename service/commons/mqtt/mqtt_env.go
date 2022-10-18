package mqtt

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

func GetHyperMqttConfig(clientPrefix string) *PoolConfig {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "default_host"
	}
	mqttUrl := os.Getenv("HYPER_MQTT_URL")
	mqttUsername := os.Getenv("HYPER_MQTT_USERNAME")
	mqttPassword := os.Getenv("HYPER_MQTT_PASSWORD")

	if mqttUrl != "" && mqttUsername != "" && mqttPassword != "" {
		rgen := rand.New(rand.NewSource(time.Now().UnixNano()))
		// Generate random client id if one wasn't supplied.
		b := make([]byte, 16)
		rgen.Read(b)
		mqttClientId := fmt.Sprintf("%s-%x", clientPrefix, b)
		config := &PoolConfig{
			Url:      mqttUrl,
			ClientId: mqttClientId,
			Username: mqttUsername,
			Password: mqttPassword,
			Hostname: hostname,
		}
		return config
	}
	return nil
}
