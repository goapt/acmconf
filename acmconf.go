package acmconf

import (
	"reflect"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/ilibs/json5"
	"github.com/verystar/goacm"
)

type acmItem struct {
	dataId string
	group  string
	v      interface{}
}

type Config struct {
	Client *goacm.Client
	Tag    string
	cache  sync.Map
}

func NewConfig(option func(c *goacm.Client)) (*Config, error) {
	client, err := goacm.NewClient(option)

	if err != nil {
		return nil, err
	}

	return &Config{
		Client: client,
		Tag:    "acmconf",
		cache:  sync.Map{},
	}, nil
}

func (c *Config) Get(dataId, group string) (string, error) {
	return c.Client.GetConfig(dataId, group)
}

//Unmarshal is json5 unmarshal to struct and support xpath
func (c *Config) Unmarshal(dataId, group string, v interface{}) error {
	result, err := c.Get(dataId, group)

	if err != nil {
		return err
	}

	buf := []byte(result)
	err = json5.Unmarshal(buf, v)
	if err == nil {
		c.cache.Store(c.getCacheKey(dataId, group), &acmItem{
			dataId: dataId,
			group:  group,
			v:      v,
		})
	}
	return err
}

func (c *Config) Listen(dataId, group string, v interface{}, fn func()) {
	go func(dataId, grpup string, v interface{}) {
		for {
			_, err := c.Client.Subscribe(dataId, group, "")
			if err == nil {
				c.Unmarshal(dataId, group, v)
				fn()
			}
			time.Sleep(1 * time.Second)
		}
	}(dataId, group, v)
}

func (c *Config) getCacheKey(dataId, group string) string {
	return strings.Join([]string{c.Client.NameSpace, dataId, group}, "-")
}

func (c *Config) Load(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("config:Load(non-pointer)")
	}
	val := rv.Elem()
	t := reflect.TypeOf(v).Elem()

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i).Tag.Get(c.Tag)
		if f == "-" || f == "" {
			continue
		}
		tmp := strings.Split(f, ",")

		if len(tmp) != 2 {
			continue
		}

		err := c.Unmarshal(tmp[0], tmp[1], val.Field(i).Addr().Interface())
		if err != nil {
			return err
		}
	}

	return nil
}
