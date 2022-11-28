package commons

func (c *ServiceCommon) GetMessageReadSignal(msg_type string, message any) map[string]any {
	if msg_type == "MESSAGE_READ" {
		m := map[string]any{
			"id": c.Handler.String(false).GenerateUuid(false),
			"message": map[string]any{
				"id": message,
			},
			"type": msg_type,
		}
		return m
	}
	return make(map[string]any)
}

func (c *ServiceCommon) GetMessageSignal(msg_type string, message any, target string) map[string]any {
	m := map[string]any{
		"id":      c.Handler.String(false).GenerateUuid(false),
		"message": message,
		"target":  target,
		"type":    msg_type,
	}
	return m
}
