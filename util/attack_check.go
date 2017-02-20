package util

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hpcloud/tail"
	"github.com/patrickmn/go-cache"

	"xsec-ssh-firewall/settings"
)

func MonitorLog(SShLogPath string) {
	logName := fmt.Sprintf("%v/auth.log", SShLogPath)
	t, err := tail.TailFile(logName, tail.Config{Follow: true})
	// fmt.Println(t, err)
	if err == nil {
		for line := range t.Lines {
			//fmt.Printf("%v, %v, %v\n", line.Time.Format("2006-01-02 15:04:05"), line.Text, line.Err)
			CheckSSH(line)
		}
	}
}

func CheckSSH(logContent *tail.Line) {
	content := logContent.Text
	re, _ := regexp.Compile(`.* Failed password for (.+?) from (.+?) port .*`)
	ret := re.FindStringSubmatch(content)
	if len(ret) > 0 {
		ip := ret[2]
		user := ret[1]
		if _, ok := settings.WhiteIPlist[ip]; ok {
			log.Printf("%v in White list, ignore\n", ip)

		} else {
			log.Printf("%v, [ %v ] try to crack %v 's password\n", logContent.Time.Format("2006-01-02 15:04:05"), ip, user)
			if c, ok := settings.Cache[ip]; ok {
				if _, ok := c.Get("times"); ok {
					c.IncrementInt("times", 1)
					times, _ := c.Get("times")
					log.Printf("%v crack times:%v\n", ip, times)
					t := times.(int)
					if t >= 3 {
						// Refresh the traffic transfer policy
						RefreshPolicy()
					}
				} else {
					c.Set("times", 1, cache.DefaultExpiration)
				}
			} else {
				c := cache.New(settings.BlockTime*time.Minute, 30*time.Second)

				c.Set("times", 1, cache.DefaultExpiration)
				settings.Cache[ip] = c
			}
		}

	}
}
