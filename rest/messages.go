package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"net/url"
	"strconv"

	"github.com/facundomedica/Rocket.Chat.Go.SDK/models"
)

type MessagesResponse struct {
	Status
	Messages []models.Message `json:"messages"`
}

type MessageResponse struct {
	Status
	Message models.Message `json:"message"`
}

// Sends a message to a channel. The name of the channel has to be not nil.
// The message will be html escaped.
//
// https://rocket.chat/docs/developer-guides/rest-api/chat/postmessage
func (c *Client) Send(channel *models.Channel, msg string) error {
	body := fmt.Sprintf(`{ "channel": "%s", "text": "%s"}`, channel.Name, html.EscapeString(msg))
	return c.Post("chat.postMessage", bytes.NewBufferString(body), new(MessageResponse))
}

// PostMessage send a message to a channel. The channel or roomId has to be not nil.
// The message will be json encode.
//
// https://rocket.chat/docs/developer-guides/rest-api/chat/postmessage
func (c *Client) PostMessage(msg *models.PostMessage) (*MessageResponse, error) {
	body, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	response := new(MessageResponse)
	err = c.Post("chat.postMessage", bytes.NewBuffer(body), response)
	return response, err
}

// Get messages from a channel. The channel id has to be not nil. Optionally a
// count can be specified to limit the size of the returned messages.
//
// https://rocket.chat/docs/developer-guides/rest-api/channels/history
func (c *Client) GetMessages(channel *models.Channel, page *models.Pagination) ([]models.Message, error) {
	params := url.Values{
		"roomId": []string{channel.ID},
	}

	if page != nil {
		params.Add("count", strconv.Itoa(page.Count))
	}

	response := new(MessagesResponse)
	if err := c.Get("channels.history", params, response); err != nil {
		return nil, err
	}

	return response.Messages, nil
}

// SetReaction Toggles the authenticated user’s reaction to the provided message.
//
// https://rocket.chat/docs/developer-guides/rest-api/chat/react/
func (c *Client) SetReaction(messageID, emoji string, shouldReact bool) error {

	payload := map[string]interface{}{}
	payload["messageId"] = messageID
	payload["emoji"] = emoji
	payload["shouldReact"] = shouldReact

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	response := new(Status)
	err = c.Post("chat.react", bytes.NewBuffer(body), response)
	return err
}

// UpdateMessage updates the text of a message
//
// https://rocket.chat/docs/developer-guides/rest-api/chat/update/
func (c *Client) UpdateMessage(messageID, roomID string, newText string) (*MessageResponse, error) {

	payload := map[string]interface{}{}
	payload["messageId"] = messageID
	payload["roomId"] = roomID
	payload["text"] = newText

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	response := new(MessageResponse)
	err = c.Post("chat.update", bytes.NewBuffer(body), response)
	return response, err
}

// GetMessage gets a single message by its ID
//
// https://rocket.chat/docs/developer-guides/rest-api/chat/getmessage/
func (c *Client) GetMessage(messageID string) (*MessageResponse, error) {

	params := url.Values{
		"messageId": []string{messageID},
	}

	response := new(MessageResponse)
	err := c.Get("chat.getMessage", params, response)
	return response, err
}
