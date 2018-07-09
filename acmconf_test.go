package acmconf

import (
	"log"
	"os"
	"testing"
	"github.com/verystar/goacm"
	"fmt"
)

func getConfig() *Config {
	conf, err := NewConfig(func(c *goacm.Client) {
		c.AccessKey = os.Getenv("AccessKey")
		c.SecretKey = os.Getenv("SecretKey")
		c.EndPoint = "acm.aliyun.com"
		c.NameSpace = os.Getenv("NameSpace")
		c.TimeOut = 10
	})

	if err != nil {
		log.Fatal(err)
	}

	return conf
}

type tmpConf struct {
	Enable       bool   `json:"enable"`
	Driver       string `json:"driver"`
	Dsn          string `json:"dsn"`
	MaxOpenConns int    `toml:"max_open_conns" json:"max_open_conns"`
	MaxIdleConns int    `toml:"max_idle_conns" json:"max_idle_conns"`
	ShowSql      bool   `toml:"show_sql" json:"show_sql"`
}

type App struct {
	DB   map[string]*tmpConf `acmconf:"[verypay:database.pay]"`
}

func TestConfig_Load(t *testing.T) {
	conf := getConfig()
	app := &App{}
	err := conf.Load(app)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(app.DB["verypay:database.pay"])
}

//func TestConfig_Listen(t *testing.T) {
//	conf := getConfig()
//	//Sleep 3 secend, Ensure configuration effective
//	time.Sleep(3 * time.Second)
//
//	app := &App{}
//	err := conf.Load(app)
//	fmt.Println(app.DB["verypay:database.pay"])
//
//	if err != nil {
//		t.Error(err)
//	}
//
//	conf.Listen(func(key string, v interface{}) {
//		fmt.Println("update", key)
//	})
//
//	time.Sleep(8 * time.Second)
//	publish("test2",t)
//
//
//	if err != nil {
//		t.Error(err)
//	}
//	time.Sleep(5 * time.Second)
//
//	fmt.Println(app.DB["verypay:database.pay"])
//	if app.DB["verypay:database.pay"].Dsn != "test2:test@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Asia%2FShanghai" {
//		t.Error("app not update")
//	}
//}
