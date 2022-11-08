package commons

import "techberry-go/micronode/service/models"

func (c *ServiceCommon) GetMessageSignal(msg_type string, message any) map[string]any {
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

	return c.Handler.Struct(false).ToMap(model)
}
