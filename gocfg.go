package gocfg

import (
	"fmt"
	"strings"
	"sync"
	"io/ioutil"
	"github.com/truexf/goutil"
	"strconv"
)

type GoConfig struct {
	cfgFile string
	cfgData map[string]map[string]string
	sync.Mutex
}

func NewGoConfig(fn string) (ret *GoConfig, e error) {
	ret = new(GoConfig)	
	if fn != "" {
		if err := ret.ReadConfig(fn); err != nil {
			return nil,fmt.Errorf("read config file %s fail,%s", fn,err.Error())
		}
	}
	return ret, nil
}

func (m *GoConfig) ReadConfig(fn string) error {
	m.Lock()
	defer m.Unlock()
	data,err := ioutil.ReadFile(fn)
	if err != nil {
		return err
	}
	m.cfgData = make(map[string]map[string]string)
	dataS := string(data)
	dataSlice := goutil.SplitByLine(dataS)
	var section map[string]string = nil
	for _, v := range dataSlice {
		s := strings.TrimSpace(v)
		if len(s) == 0 {
			continue
		}		
		if s[:2] == "//" {
			continue
		}
		if s[0] == '#' {
			continue
		}
		if s[:1] == "[" {
			if s[len(s)-1:] != "]" {
				continue
			}
			section = make(map[string]string)
			sec := s[1:len(s)-1]
			m.cfgData[sec] = section			
		} else {
			if (section == nil) {
				continue
			}
			l,r := goutil.SplitLR(s,"=")
			l = strings.TrimSpace(l)
			r = strings.TrimSpace(r)
			if l != "" {
				section[l] = r
			}
			continue
		}
	}
	m.cfgFile = fn
	return nil
}

func (m *GoConfig) Reload() {
	if m.cfgFile != "" {
		m.ReadConfig(m.cfgFile)
	}
}

func (m *GoConfig) SaveToFile(fn string) error {
	m.Lock()
	defer m.Unlock()
	ret := make([]string,0,100)
	for k,v := range m.cfgData {
		ret = append(ret,"")
		ret = append(ret,fmt.Sprintf("[%s]", k))
		for kC,vC := range v {
			ret = append(ret, fmt.Sprintf("%s = %s",kC,vC))
		}
	}
	return ioutil.WriteFile(fn,[]byte(strings.Join(ret,"\n")),0666)
}

func (m *GoConfig) Get(sec,ident,dft string) string {
	m.Lock()
	defer m.Unlock()
	if secMap,ok := m.cfgData[sec]; ok {
		if v,ok := secMap[ident]; ok {
			return v
		} else {
			return dft
		}
	} else {
		return dft
	}
}

func (m *GoConfig) Set(sec,k,v string) {
	if sec == "" || k == "" {
		return
	}
	m.Lock()
	defer m.Unlock()
	if secMap,ok := m.cfgData[sec]; ok {
		secMap[k] = v
	} else {
		sec := make(map[string]string)
		sec[k] = v
	}
	return 
}



func (m *GoConfig) GetIntDefault(sec,k string, dft int) int {
	ret := m.Get(sec,k,"")
	if ret == "" {
		return dft
	} else {
		if retInt,err := strconv.Atoi(ret); err != nil {
			return dft
		} else {
			return retInt
		}
	}
}

func (m *GoConfig) GetFloatDefault(sec,k string, dft float64) float64 {
	ret := m.Get(sec,k,"")
	if ret == "" {
		return dft
	} else {
		if retFloat,err := strconv.ParseFloat(ret,64); err != nil {
			return dft
		} else {
			return retFloat
		}
	}
}

func (m *GoConfig) GetBoolDefault(sec,k string, dft bool) bool {
	ret := m.Get(sec,k,"")
	if ret == "" {
		return dft
	} else {
		return ret == "1"
	}
}