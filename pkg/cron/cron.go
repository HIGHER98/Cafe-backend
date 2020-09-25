package cron

import (
	"cafe/models"
	"cafe/pkg/logging"
	"cafe/pkg/setting"

	"gopkg.in/robfig/cron.v1"
)

func Setup() {
	c := cron.New()
	c.AddFunc("@daily", Anonymise)
	c.Start()
}

//Anonymises data `days` old
func Anonymise() {
	if err := models.Anonymise(setting.CronSetting.Days); err != nil {
		logging.Error("Failed to anonymise: ", err)
		return
	}
	logging.Info("Successfully anonymised data")
}
