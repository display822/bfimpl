package log

import (
	"encoding/json"
	"fmt"
	beelog "github.com/astaxie/beego/logs"
	"io"
	"net"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	LevelEmergency = iota
	LevelAlert
	LevelCritical
	LevelError
	LevelWarning
	LevelNotice
	LevelInformational
	LevelDebug
)

var levelName [8]string

//go|toolname|serviceip|servicename|time|level|funcion|content
const msgpattern = "go|%s|%s|%s|%s|%s|%s|%s\n"
const timeformat = "2006-01-02 15:04:05"

var ethip string

type dlogsWriter struct {
	sync.Mutex
	innerWriter    io.WriteCloser
	ToolName       string `json:"toolname"`
	ServiceName    string `json:"servicename"`
	ServiceIP      string `json:"serviceip"`
	TLogIPPort     string `json:"tlogipport"`     // 127.0.0.1:6666 with udp
	NetType        string `json:"nettype"`        // tcp or udp
	Level          int    `json:"level"`          // min level
	Reconnect      bool   `json:"reconnect"`      // if config changed at runtime
	ReconnectOnMsg bool   `json:"reconnectOnMsg"` // reconnect on every msg
}

func (d *dlogsWriter) Init(config string) error {
	return json.Unmarshal([]byte(config), d)
}

func (d *dlogsWriter) WriteMsg(when time.Time, msg string, level int) error {
	if level > d.Level {
		return nil
	}
	if d.needToConnectOnMsg() {
		err := d.connect()
		if err != nil {
			return err
		}
	}

	if d.ReconnectOnMsg {
		defer d.innerWriter.Close()
	}
	i := strings.IndexByte(msg, ']')
	if i > 0 && i+2 < len(msg) {
		msg = msg[i+2:]
	}
	var levelname string
	if level > 0 && level < LevelDebug {
		levelname = levelName[level]
	} else {
		levelname = "DEBUG"
	}
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	_, filename := path.Split(file)
	funname := fmt.Sprintf("%s.%d", filename, line)
	timefmt := when.Format(timeformat)
	data := fmt.Sprintf(msgpattern,
		d.ToolName, d.ServiceIP, d.ServiceName,
		timefmt, levelname, funname,
		msg)

	d.Lock()
	_, _ = d.innerWriter.Write([]byte(data))
	d.Unlock()
	return nil
}

func (d *dlogsWriter) Destroy() {
	if d.innerWriter != nil {
		_ = d.innerWriter.Close()
	}
}

func (d *dlogsWriter) connect() error {
	if d.innerWriter != nil {
		_ = d.innerWriter.Close()
		d.innerWriter = nil
	}

	conn, err := net.Dial(d.NetType, d.TLogIPPort)
	if err != nil {
		return err
	}

	if tcpConn, ok := conn.(*net.TCPConn); ok {
		_ = tcpConn.SetKeepAlive(true)
	}

	d.innerWriter = conn
	return nil
}

func (d *dlogsWriter) needToConnectOnMsg() bool {
	if d.Reconnect {
		d.Reconnect = false
		return true
	}

	if d.innerWriter == nil {
		return true
	}

	return d.ReconnectOnMsg
}

// empty
func (d *dlogsWriter) Flush() {
}

func NewDlogsWriter() beelog.Logger {
	w := new(dlogsWriter)
	w.Level = LevelInformational
	w.ReconnectOnMsg = false
	w.ServiceIP = ethip
	return w
}

func getEthIP(name string) string {
	var ip string
	in, err := net.InterfaceByName(name)
	if err == nil {
		addr, err := in.Addrs()
		if err == nil && len(addr) > 0 {
			ip = addr[0].String()
			if strings.IndexByte(ip, '/') > 0 {
				ip = ip[0:strings.IndexByte(ip, '/')]
			}
		}
	}
	return ip
}

func init() {
	levelName[LevelEmergency] = "FATAL"
	levelName[LevelAlert] = "ALERT"
	levelName[LevelCritical] = "CRITICAL"
	levelName[LevelError] = "ERROR"
	levelName[LevelWarning] = "WARN"
	levelName[LevelNotice] = "NOTICE"
	levelName[LevelInformational] = "INFO"
	levelName[LevelDebug] = "DEBUG"
	serverip := os.Getenv("serverip")
	ethip = "127.0.0.1"
	if serverip != "" {
		ethip = serverip
	} else if serverip = getEthIP("eth1"); serverip != "" {
		ethip = serverip
	} else if serverip = getEthIP("eth0"); serverip != "" {
		ethip = serverip
	}
}
