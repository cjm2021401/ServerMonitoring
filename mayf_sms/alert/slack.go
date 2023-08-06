package alert

import (
	"os"
	"strconv"

	"custom_sms/agent/model"
	"custom_sms/config"

	"github.com/slack-go/slack"
)

func CPUAlert(cpustats []model.CPUStat) error {
	if cpustats[0].Used > config.Env.Alert.CpuUsageLimit {

		hostname, _ := os.Hostname()
		attachment := slack.Attachment{
			Pretext: hostname + "'s CPU Used > " + strconv.FormatFloat(float64(config.Env.Alert.CpuUsageLimit), 'f', -1, 32),
			Text:    "server :" + hostname + "\nlimit : " + strconv.FormatFloat(float64(config.Env.Alert.CpuUsageLimit), 'f', -1, 32) + "\n current : " + strconv.FormatFloat(float64(cpustats[0].Used), 'f', -1, 32),
		}
		api := slack.New(config.Env.Slack.Token)
		_, _, err := api.PostMessage(
			config.Env.Slack.Channel,
			slack.MsgOptionText("", false),
			slack.MsgOptionAttachments(attachment),
			slack.MsgOptionAsUser(false),
		)
		return err

	}
	return nil

}

func MemoryAlert(memoryStat model.MemoryStat) error {
	if memoryStat.MemUsed > config.Env.Alert.MemoryUsageLimit {

		hostname, _ := os.Hostname()
		attachment := slack.Attachment{
			Pretext: hostname + "'s Memory Used > " + strconv.FormatFloat(float64(config.Env.Alert.MemoryUsageLimit), 'f', -1, 32),
			Text:    "server :" + hostname + "\nlimit : " + strconv.FormatFloat(float64(config.Env.Alert.MemoryUsageLimit), 'f', -1, 32) + "\n current : " + strconv.FormatFloat(float64(memoryStat.MemUsed), 'f', -1, 32),
		}
		api := slack.New(config.Env.Slack.Token)
		_, _, err := api.PostMessage(
			config.Env.Slack.Channel,
			slack.MsgOptionText("", false),
			slack.MsgOptionAttachments(attachment),
			slack.MsgOptionAsUser(false),
		)
		return err

	}
	return nil
}
