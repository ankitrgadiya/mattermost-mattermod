// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/mattermost/mattermost-mattermod/server"
	"gopkg.in/robfig/cron.v3"
)

func main() {
	var flagConfigFile string
	flag.StringVar(&flagConfigFile, "config", "config-mattermod.json", "")
	flag.Parse()

	server.LoadConfig(flagConfigFile)
	s := server.New()
	s.Start()
	defer s.Stop()

	//server.CleanOutdatedPRs()
	//server.CleanOutdatedIssues()

	c := cron.New()
	c.AddFunc("@daily", s.CheckPRActivity)
	c.AddFunc("@midnight", s.CleanOutdatedPRs)
	c.AddFunc("@every 2h", s.CheckSpinmintLifeTime)

	cronTicker := fmt.Sprintf("@every %dm", server.Config.TickRateMinutes)
	c.AddFunc(cronTicker, s.Tick)

	c.Start()
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}
