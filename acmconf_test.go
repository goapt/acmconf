package acmconf

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/verystar/goacm"
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

type App struct {
	Test map[string]string `acmconf:"test,test"`
}

var beforJsonConf = `
{
  "name": "lili",
  "mail": "test@test.com",
  "hello": "world",
}
`

var afterJsonConf = `
{
  "name": "didi",
  "mail": "tmp@test.com",
  "hello": "verystar",
}
`

func TestConfig_Listen(t *testing.T) {
	conf := getConfig()
	_, err := conf.Client.Publish("test", "test", beforJsonConf)

	if err != nil {
		t.Error(err)
	}

	//Sleep 3 secend, Ensure configuration effective
	time.Sleep(3 * time.Second)

	app := &App{}
	err = conf.Load(app)

	if err != nil {
		t.Error(err)
	}

	conf.Listen("test", "test", &app.Test, func() {
		if app.Test["name"] != "didi" {
			t.Error("app not update")
		}
	})

	time.Sleep(5 * time.Second)
	_, err = conf.Client.Publish("test", "test", afterJsonConf)

	if err != nil {
		t.Error(err)
	}
	time.Sleep(5 * time.Second)

	if app.Test["name"] != "didi" {
		t.Error("app not update")
	}
}
