package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type GoogleChatService struct {
	webhookURL           string
	lastNotificationTime time.Time
	mx                   *sync.Mutex
}

func NewGoogleChatService(webhookURL string) *GoogleChatService {
	return &GoogleChatService{
		webhookURL:           webhookURL,
		lastNotificationTime: time.Time{},
		mx:                   &sync.Mutex{},
	}
}

var _ Service = (*GoogleChatService)(nil)

func (s *GoogleChatService) PushCreateNotification(key string, description string, flagURL string) error {
	return s.pushNotification("✨", key, "Feature Flag Created", description, &flagURL)
}

func (s *GoogleChatService) PushDeleteNotification(key string, description string) error {
	return s.pushNotification("🗑️", key, "Feature Flag Removed", description, nil)
}

func (s *GoogleChatService) pushNotification(emoji string, key string, header string, description string, flagURL *string) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	// Google chat has a rate limit of 1 message per second.
	if time.Since(s.lastNotificationTime) < time.Second {
		time.Sleep(time.Second - time.Since(s.lastNotificationTime))
	}

	accessoryWidgets := []googleChatRequestModelAccessoryWidget{}
	if flagURL != nil {
		accessoryWidgets = append(accessoryWidgets, googleChatRequestModelAccessoryWidget{
			ButtonList: googleChatRequestModelButtonList{
				Buttons: []googleChatRequestModelButton{
					{
						Text: "View on PostHog",
						Icon: googleChatRequestModelButtonIcon{
							MaterialIcon: googleChatRequestModelButtonMaterialIcon{
								Name: "savings",
							},
						},
						OnClick: googleChatRequestModelButtonOnClick{
							OpenLink: googleChatRequestModelButtonOpenLink{
								Url: *flagURL,
							},
						},
					},
				},
			},
		})
	}

	requestModel := &googleChatRequestModel{
		CardContainers: []googleChatRequestModelCardContainer{
			{
				Card: googleChatRequestModelCard{
					Header: googleChatRequestModelHeader{
						Title: fmt.Sprintf("%s %s", emoji, header),
					},
					Sections: []googleChatRequestModelSection{
						{
							Header: "Name",
							Widgets: []googleChatRequestModelWidget{
								{
									TextParagraph: &googleChatRequestModelTextParagraph{
										Text: key,
									},
								},
							},
						},
						{
							Header: "Description",
							Widgets: []googleChatRequestModelWidget{
								{
									TextParagraph: &googleChatRequestModelTextParagraph{
										Text: description,
									},
								},
							},
						},
					},
				},
			},
		},
		AccessoryWidgets: accessoryWidgets,
	}

	rawRequestModel, err := json.Marshal(requestModel)
	if err != nil {
		return fmt.Errorf("failed to marshal google chat request model: %w", err)
	}

	req, err := http.NewRequest("POST", s.webhookURL, bytes.NewBuffer(rawRequestModel))
	if err != nil {
		return fmt.Errorf("failed to create google chat request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	s.lastNotificationTime = time.Now()

	if err != nil {
		return fmt.Errorf("failed to complete google chat http request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		rawResponseText, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("google chat request failed with status code: %d and failed to read response body: %w", res.StatusCode, err)
		}

		return fmt.Errorf("google chat request failed with status code: %d and response body: %s", res.StatusCode, string(rawResponseText))
	}

	return nil
}

type googleChatRequestModel struct {
	CardContainers   []googleChatRequestModelCardContainer   `json:"cardsV2"`
	AccessoryWidgets []googleChatRequestModelAccessoryWidget `json:"accessoryWidgets,omitempty"`
}

type googleChatRequestModelCardContainer struct {
	Card googleChatRequestModelCard `json:"card"`
}

type googleChatRequestModelCard struct {
	Header   googleChatRequestModelHeader    `json:"header"`
	Sections []googleChatRequestModelSection `json:"sections"`
}

type googleChatRequestModelHeader struct {
	Title string `json:"title"`
}

type googleChatRequestModelSection struct {
	Header  string                         `json:"header"`
	Widgets []googleChatRequestModelWidget `json:"widgets"`
}

type googleChatRequestModelWidget struct {
	TextParagraph *googleChatRequestModelTextParagraph `json:"textParagraph,omitempty"`
}

type googleChatRequestModelTextParagraph struct {
	Text string `json:"text"`
}

type googleChatRequestModelAccessoryWidget struct {
	ButtonList googleChatRequestModelButtonList `json:"buttonList"`
}

type googleChatRequestModelButtonList struct {
	Buttons []googleChatRequestModelButton `json:"buttons"`
}

type googleChatRequestModelButton struct {
	Text    string                              `json:"text"`
	Icon    googleChatRequestModelButtonIcon    `json:"icon"`
	OnClick googleChatRequestModelButtonOnClick `json:"onClick"`
}

type googleChatRequestModelButtonIcon struct {
	MaterialIcon googleChatRequestModelButtonMaterialIcon `json:"materialIcon"`
}

type googleChatRequestModelButtonMaterialIcon struct {
	Name string `json:"name"`
}

type googleChatRequestModelButtonOnClick struct {
	OpenLink googleChatRequestModelButtonOpenLink `json:"openLink"`
}

type googleChatRequestModelButtonOpenLink struct {
	Url string `json:"url"`
}
