package commons

import "techberry-go/micronode/service/models"

func (c *ServiceCommon) GetMessageSignal(msg_type string, message any) map[string]any {
	model := models.OutputChatMessageSignal{
		Id:      c.Handler.String(false).GenerateUuid(false),
		Message: message,
		Type:    msg_type,
	}
	return c.Handler.Struct(false).ToMap(model)
}
