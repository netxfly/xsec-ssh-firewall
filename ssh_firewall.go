package main

import (
	"xsec-ssh-firewall/settings"
	"xsec-ssh-firewall/util"
)

func main() {
	go util.MonitorLog(settings.SshLog)

	util.Schedule(30)

}
