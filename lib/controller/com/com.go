package com

import "github.com/SmartEnergyPlatform/device-manager/lib/config"

type Com struct {
	config config.Config
}

func New(config config.Config) *Com {
	return &Com{config: config}
}
