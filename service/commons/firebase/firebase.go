package firebase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	fb "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

type firebaseNotiInput struct {
	Config       map[string]interface{} `structs:"config" json:"config" bson:"config"`
	Tokens       []string               `structs:"tokens" json:"tokens" bson:"tokens"`
	Data         map[string]string      `structs:"data" json:"data" bson:"data"`
	Notification struct {
		Title    string `structs:"title" json:"title" bson:"title"`
		Body     string `structs:"body" json:"body" bson:"body"`
		ImageURL string `structs:"imageURL" json:"imageURL" bson:"imageURL"`
	} `structs:"notification" json:"notification" bson:"notification"`
}

func SendFirebaseNotification(rootMap map[string]interface{}) (result map[string]interface{}, err error) {
	defer fmt.Println("end SendMessage")

	// Prepare data
	// Convert config to byte[]
	rootBytes, err := json.Marshal(rootMap)
	if err != nil {
		return nil, err
	}
	// Convert to model
	rootModel := firebaseNotiInput{}
	err = json.Unmarshal(rootBytes, &rootModel)
	if err != nil {
		return nil, err
	}

	// Read Config
	config := rootModel.Config
	tokens := rootModel.Tokens
	notiObj := rootModel.Notification
	dataMap := rootModel.Data

	// Validate Data
	if len(config) <= 0 {
		return nil, errors.New("config is empty")
	}
	if notiObj.Title == "" {
		return nil, errors.New("notification is empty")
	}
	title := notiObj.Title
	body := notiObj.Body
	imageURL := notiObj.ImageURL
	if dataMap == nil {
		dataMap = make(map[string]string)
	}

	// Convert config to byte[]
	bConfig, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	// Init-App
	opt := option.WithCredentialsJSON(bConfig)
	app, err := fb.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}

	// Init-Message
	message := &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title:    title,
			Body:     body,
			ImageURL: imageURL,
		},
		Data:   dataMap,
		Tokens: tokens,
	}

	// Init Client
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		return nil, err
	}

	// Call SendMulticast Notification
	br, err := client.SendMulticast(context.Background(), message)
	if err != nil {
		return nil, err
	}

	// Make Result
	result = make(map[string]interface{})
	result["successCount"] = br.SuccessCount
	result["failureCount"] = br.FailureCount

	// Make Result List
	tkLen := len(tokens)
	list := make([]interface{}, 0)
	response := br.Responses
	for i, rp := range response {
		token := ""
		if i < tkLen {
			token = tokens[i]
		}
		msg := make(map[string]interface{})
		msg["token"] = token
		msg["success"] = rp.Success
		msg["messageId"] = rp.MessageID
		if rp.Error != nil {
			msg["error"] = rp.Error.Error()
		}
		list = append(list, msg)
	}
	result["responses"] = list

	return result, nil
}
