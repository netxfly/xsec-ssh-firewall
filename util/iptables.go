package util

import (
	"log"
	"os/exec"

	"xsec-ssh-firewall/settings"
)

func SetIptables() {

	for ip := range settings.Cache {
		if ip == "127.0.0.1" {
			log.Printf("Local ip: %v", ip)
			continue
		} else {
			log.Printf("Block ip: %v\n", ip)
			exec.Command("/sbin/iptables", "-t", "filter", "-A", "WHITELIST", "-i", settings.Interface, "-s", ip, "-j", "DROP").Output()

		}

	}
}

// Init iptables policy
func InitPolicy() {
	// set white list chain in filter table
	exec.Command("/sbin/iptables", "-t", "filter", "-N", "WHITELIST").Run()
	exec.Command("/sbin/iptables", "-t", "filter", "-F", "WHITELIST").Run()
	exec.Command("/sbin/iptables", "-t", "filter", "-A", "INPUT", "-j", "WHITELIST").Run()
}

// Delete Policy
func DeletePolicy() {
	// Flush rule
	exec.Command("/sbin/iptables", "-t", "filter", "-F").Run()
	// delete chain
	exec.Command("/sbin/iptables", "-t", "filter", "-X", "WHITELIST").Run()

}

func RefreshPolicy() {
	Stop()
	InitPolicy()
	SetIptables()
}

func Stop() {
	DeletePolicy()
}
