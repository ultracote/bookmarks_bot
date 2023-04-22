package router

import (
	"github.com/mymmrac/telego"
)

type Sender struct {
	bot             *telego.Bot
	previousHandler Handler
}

func NewSender(bot *telego.Bot, previousHandler Handler) *Sender {
	return &Sender{
		bot:             bot,
		previousHandler: previousHandler,
	}
}

func (s *Sender) Handle(request *Request) (*Response, error) {
	previousResponse, err := s.previousHandler.Handle(request)
	if err != nil {
		return nil, err
	}

	for _, message := range previousResponse.GetMessages() {
		// set message chat id
		message.ChatID = telego.ChatID{
			ID: request.GetChatID(),
		}

		// TODO: handle errors
		_, err := s.bot.SendMessage(message)
		if err != nil {
			return nil, err
		}
	}

	// clean sent messages
	previousResponse.SetMessages(nil)

	return previousResponse, nil
}
