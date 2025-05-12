package setting

import (
	"service/config"
)

type setting struct {
	Name         string
	Env          string
	QueueProgram string
}

var Setting setting

func (s *setting) IsProduction() bool {
	if s.Env == "production" {
		return true
	}
	return false
}

func NewSetting(config *config.Config) {
	Setting = setting{
		Name:         config.App.Name,
		Env:          config.App.Environment,
		QueueProgram: config.Setting.QueueProgram,
	}
}
