package settings

import (
	"fmt"

	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/patrickmn/go-cache"
	"gopkg.in/ini.v1"
	"strings"
	"time"
)

var (
	Cfg         *ini.File
	Interface   string
	BlockTime   time.Duration
	WhiteIPlist map[string]bool

	SshLog string
	Cache  map[string]*cache.Cache
)

func init() {
	log.SetPrefix("[xsec-ssh-firewall] ")
	var err error
	source := "conf/app.ini"
	Cfg, err = ini.Load(source)
	//log.Println(Cfg, err)
	if err != nil {
		log.Panicln(err)
	}

	Interface = Cfg.Section("").Key("INTERFACE").MustString("eth0")
	blockTime := Cfg.Section("").Key("BLOCKTIME").MustInt64(60)
	BlockTime = time.Duration(blockTime)
	whiteIps := Cfg.Section("").Key("WHITE_IPLIST").MustString("127.0.0.1")
	WhiteIPlist = make(map[string]bool)
	for _, ip := range strings.Split(whiteIps, ",") {
		WhiteIPlist[ip] = true
	}

	curDir, _ := filepath.Abs(".")
	secLogs := Cfg.Section("logs")
	SshLog = secLogs.Key("sshd_log").MustString(fmt.Sprintf("%v/logs/openssh/", curDir))

	if _, err := os.Stat(SshLog); !os.IsExist(err) {
		logPath := path.Dir(SshLog)
		os.MkdirAll(logPath, 0755)
		os.Create(path.Join(SshLog, "auth.log"))
		os.Create(path.Join(SshLog, "syslog"))
	}
	Cache = make(map[string]*cache.Cache)
	log.Printf("SSHLogPath: %v, Interface: %v, BlockTime: %v, WhiteIplist: %v\n", SshLog, Interface,
		BlockTime*time.Minute, WhiteIPlist)
}
