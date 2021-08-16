package gocfg

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestGoConfig(t *testing.T) {
	cfgData := `[DEFAULT]
bind = 5050

[QUEUE_LOG]
//this is comment
broker = 192.168.13.86:9092,192.168.13.87:9092,192.168.13.88:9092,192.168.13.89:9092,192.168.13.90:9092,192.168.13.91:9092,192.168.13.92:9092
#this is comment too
filelog.dir = /home/wwwroot/service/ae/queuelog //this is not comment

[STAT]
second_limit = 60
request_limit = 1000000
write_path = /home/wwwroot/service/ae/statis

[DSP]
data-dir = /home/wwwroot/g.ggxt.net/engine/data/core/adx/dsp
id = 100008,100009,100010,100011,100016,100017,100018,100019,100020,100021,100022,100024,100025,100026,100027,100028,100029,100030,100031,100032,100033,100034,100035,100036,100040,100041,100042,100043,100044,100045,100046,100051
timeout = 50

[AD]`
	fn := "/tmp/testgoconfig"
	ioutil.WriteFile(fn, []byte(cfgData), 0666)
	cfg, e := NewGoConfig(fn)
	if e != nil {
		t.Errorf(e.Error())
		return
	}
	fmt.Printf("dsp.data-dir: %s\n", cfg.Get("DSP", "data-dir", ""))
	fmt.Printf(cfg.Get("AD", "abc", "124"))
	fmt.Printf("\n%d\n", cfg.GetIntDefault("STAT", "second_limit", -100))
	fmt.Printf("%f\n", cfg.GetFloatDefault("DSP", "timeout", -100.0))
	fmt.Printf("%t\n", cfg.GetBoolDefault("DSP", "id", false))
	fmt.Printf("%t\n", cfg.GetBoolDefault("DSP", "id", true))
	cfg.SaveToFile("/tmp/testgoconfig1")
	if d, e := ioutil.ReadFile(fn + "1"); e == nil {
		fmt.Println(string(d))
	} else {
		fmt.Println(e.Error())
	}

}
