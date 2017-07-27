package slackbot

import(
  "log"
  "strings"

  "github.com/nlopes/slack"
)

const HELP = "help"

type Bot struct {
  api *slack.Client
  commands map[string]command
}

type command struct {
  Name string
  Description string
  Callback fn
}

type fn func(*Bot, string, ...string)

// Initializes a new slackbot
func New(slackToken string) *Bot {
  return &Bot{
    api: slack.New(slackToken),
    commands: map[string]command{ },
  }
}

// AddCommand lets you add a command that your slack bot can respond to. It passes back
// the bot (*slackbot.Bot), a channel(string), and a variadic args to the callback.
func (b *Bot) AddCommand(message, description string, callback fn, args ...string) {
  b.commands[message] = command{
    Name: message,
    Description: description,
    Callback: callback,
  };
}

// Once you add commands to your bot, you need to call Run() so your bot can start
// listening to commands
func (b *Bot) Run() {
  rtm := b.api.NewRTM()
  go rtm.ManageConnection()

  for msg := range rtm.IncomingEvents {
    switch ev := msg.Data.(type) {
    case *slack.MessageEvent:
      go b.handleMessage(ev.Msg)
    case *slack.RTMError:
      log.Printf("Error: %s\n", ev.Error())
    default:
    }
  }
}

// A handy function you can use within your AddCommand callbacks so the bot
// can reply to commands
func (b *Bot) Reply(channel string, reply string) {
  _, _, err := b.api.PostMessage(channel, reply, slack.PostMessageParameters{})
  if err != nil {
    log.Fatal(err)
  }
}

func (b *Bot) handleMessage(msg slack.Msg) {
  messageSlice := strings.Split(msg.Text, " ")
  command := messageSlice[0]
  channel := msg.Channel
  var args []string
  if len(messageSlice) > 1 {
    args = messageSlice[1:]
  }
  if _, ok := b.commands[command]; ok {
    if args != nil && args[0] == HELP {
      b.Reply(channel, b.commands[command].Description)
    } else {
      b.commands[command].Callback(b, channel, args...)
    }
  }
}
