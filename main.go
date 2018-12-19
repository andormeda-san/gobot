package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/nlopes/slack"
)

// Slackparams is slack parameters
type Slackparams struct {
	tokenID   string
	botID     string
	channelID string
	rtm       *slack.RTM
}

func main() {
	params := Slackparams{
		tokenID:   "xxxxx-xxxxx-xxxxx-xxxxx",
		botID:     "<@xxxxx>",
		channelID: "xxxxx",
	}

	api := slack.New(params.tokenID)

	params.rtm = api.NewRTM()
	go params.rtm.ManageConnection()

	for msg := range params.rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			if err := params.ValidateMessageEvent(ev); err != nil {
				log.Printf("[ERROR] Failed to handle message: %s", err)
			}
		}
	}
}

// ValidateMessageEvent is Validate Message Event
func (s *Slackparams) ValidateMessageEvent(ev *slack.MessageEvent) error {
	// Only response in specific channel. Ignore else.
	if ev.Channel != s.channelID {
		log.Printf("%s %s", ev.Channel, ev.Msg.Text)
		return nil
	}

	// Only response mention to bot. Ignore else.
	if !strings.HasPrefix(ev.Msg.Text, s.botID) {
		log.Printf("%s %s", ev.Channel, ev.Msg.Text)
		return nil
	}

	// Parse message start
	m := strings.Split(strings.TrimSpace(ev.Msg.Text), " ")[1:]
	if len(m) == 0 {
		return fmt.Errorf("invalid message")
	}

	if m[0] == "ホリネズミってなに?" {
		_, lead := wiki()
		s.rtm.SendMessage(s.rtm.NewOutgoingMessage(lead, ev.Channel))
		return nil
	}
	s.rtm.SendMessage(s.rtm.NewOutgoingMessage("メンション付いてるな。呼んだか？", ev.Channel))
	return nil
}

func wiki() (title string, lead string) {
	url := "https://ja.wikipedia.org/wiki/%E3%83%9B%E3%83%AA%E3%83%8D%E3%82%BA%E3%83%9F"
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Println("Wikipedia scraping failed.")
		os.Exit(1)
	}
	title = doc.Find("#firstHeading").Text()
	lead = doc.Find("#mw-content-text p").First().Text()
	fmt.Println(title)
	fmt.Println(lead)
	return title, lead
}
