package util

import (
	"io/ioutil"

	"github.com/nlopes/slack"
	"github.com/spf13/viper"
)

// Alert sends a Slack Alert
func Alert(message string) error {
	buffer, err := ioutil.ReadFile(viper.GetString("slackToken"))
	if err != nil {
		return err
	}
	api := slack.New(string(buffer))
	// If you set debugging, it will log all requests to the console
	// Useful when encountering issues
	// api.SetDebug(true)
	//bleh, err := api.JoinChannel("test98331")
	api.PostMessage(
		viper.GetString("slackChannel"),
		message,
		slack.PostMessageParameters{})

	return nil
}
