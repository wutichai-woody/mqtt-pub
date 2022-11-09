package models

import (
	"techberry-go/common/v2/facade"
)

type OutputRoot struct {
	facade.JsonResponse
}

type OutputChatMessageSignal struct {
	Id      string `structs:"id" json:"id" bson:"id" mapstructure:"id"`
	Message any    `structs:"message" json:"message" bson:"message" mapstructure:"message"`
	Type    string `structs:"type" json:"type" bson:"type" mapstructure:"type"`
}
