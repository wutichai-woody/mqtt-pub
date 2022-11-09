package commons

import (
	"encoding/json"
	"techberry-go/micronode/service/models"
)

func (c *ServiceCommon) GetMessageSignal(msg_type string, message any) string {
	var model models.OutputChatMessageSignal
	if msg_type == "MESSAGE_READ" {
		model = models.OutputChatMessageSignal{
			Id: c.Handler.String(false).GenerateUuid(false),
			Message: map[string]any{
				"id": message,
			},
			Type: msg_type,
		}
	}
	b, err := json.Marshal(model)
	if err == nil {
		return string(b)
	} else {
		return ""
	}
}
