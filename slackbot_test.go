// Add more tests
// ¯\_(ツ)_/¯

package slackbot_test

import(
  "testing"

  "github.com/premshree/lib-slackbot"
)

const FAKE_TOKEN = "token"

func TestNew(t *testing.T) {
  if slackbot.New(FAKE_TOKEN) == nil {
    t.Error("Expected new slackbot object.")
  }
}
