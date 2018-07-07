package acmconf

import (
	"reflect"
	"errors"
	"strings"
	"sync"
	"github.com/ilibs/json5"
	"github.com/verystar/goacm"
	"time"
	"fmt"
)

type acmItem struct {
	key      string
	dataId   string
	group    string
	multiple bool
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

func (c *Config) unmarshal(result string, v interface{}, item *acmItem) error {

	if item.multiple {
		result = fmt.Sprintf(`{"%s":%s}`, item.key, result)
	}
	buf := []byte(result)
	return json5.Unmarshal(buf, v)
}

func (c *Config) Listen(fn func(key string, v interface{})) {

	c.cache.Range(func(key, value interface{}) bool {
		item := key.(*acmItem)
		go func(item *acmItem, v interface{}) {
			for {
				ret, err := c.Client.Subscribe(item.dataId, item.group, "")
				if err == nil {
					c.unmarshal(ret, v, item)
					fn(item.key, v)
				}
				time.Sleep(1 * time.Second)
			}
		}(item, value)

		return true
	})
}

func (c *Config) getCacheKey(group, dataId string) string {
	return strings.Join([]string{c.Client.NameSpace, dataId, group}, "-")
}

func (c *Config) getTags(str string) map[string]*acmItem {
	m := make(map[string]*acmItem, 0)
	tags := strings.Split(str, ",")
	multiple := len(tags) > 1

	for _, v := range tags {
		tmp := strings.Split(v, ":")
		if len(tmp) != 2 {
			continue
		}
		m[v] = &acmItem{
			key:      v,
			group:    tmp[0],
			dataId:   tmp[1],
			multiple: multiple,
		}
	}
	return m
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
		tags := c.getTags(f)
		for _, item := range tags {
			result, err := c.Client.GetConfig(item.dataId, item.group)
			if err != nil {
				return err
			}

			err = c.unmarshal(result, val.Field(i).Addr().Interface(), item)
			if err != nil {
				return err
			}

			c.cache.Store(item, val.Field(i).Addr().Interface())
		}
	}

	return nil
}
