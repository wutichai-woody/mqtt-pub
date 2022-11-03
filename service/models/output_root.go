package models

import (
	"techberry-go/common/v2/facade"
)

type OutputRoot struct {
	facade.JsonResponse
}

type OutputChatMessageSignal struct {
	Id      string `json:"id" bson:"id"`
	Message any    `json:"message" bson:"message"`
	Type    string `json:"type" bson:"type"`
}
