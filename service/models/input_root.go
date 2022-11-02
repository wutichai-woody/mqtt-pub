package models

import (
	"techberry-go/common/v2/facade/helper"
)

type InputRoot struct {
	helper.DataContext
}

type CacheInputRoot struct {
	Data []struct {
		Key    string `json:"key"`
		Value  string `json:"value"`
		Expire int    `json:"expire"`
	} `json:"data"`
}
