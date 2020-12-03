package env

import (
	"github.com/unistack-org/micro/v3/config"
	"context"
	"reflect"
	"os"
	"strconv"
)

var (
	DefaultStructTag = "env"
)

type envConfig struct {
	opts config.Options
}

func (c *envConfig) Options() config.Options {
	return c.opts
}


func (c *envConfig) Init(opts...config.Option) error {
	for _, o := range opts {
		o(&c.opts)
	}
	return nil
}

func (c *envConfig) Load(ctx context.Context) error {
	fields := reflect.TypeOf(c.opts.Struct).Elem()
	values := reflect.ValueOf(c.opts.Struct).Elem()

	for idx := 0; idx < fields.NumField(); idx++ {
		field := fields.Field(idx)
		value := values.Field(idx)
		if !value.CanSet() {
			continue
		}
		tag, ok := field.Tag.Lookup(c.opts.StructTag)
		if !ok {
			continue
		}
		val, ok := os.LookupEnv(tag)
		if !ok {
			continue
		}
		switch value.Kind() {
		case reflect.Bool:
			v, err := strconv.ParseBool(val)
			if err != nil {
				return err
			}
			value.Set(reflect.ValueOf(v))
		case reflect.String:
			value.Set(reflect.ValueOf(val))
		case reflect.Float32:
			v, err := strconv.ParseFloat(val, 32)
			if err != nil {
				return err
			}
			value.Set(reflect.ValueOf(v))
		case reflect.Float64:
			v, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return err
			}
			value.Set(reflect.ValueOf(v))
		case reflect.Int:
			v, err := strconv.ParseInt(val, 10, 0)
			if err != nil {
				return err
			}
			value.Set(reflect.ValueOf(v))
		case reflect.Int8:
			v, err := strconv.ParseInt(val, 10, 8)
			if err != nil {
				return err
			}
			value.Set(reflect.ValueOf(v))
		case reflect.Int16:
			v, err := strconv.ParseInt(val, 10, 16)
			if err != nil {
				return err
			}
			value.Set(reflect.ValueOf(v))
		case reflect.Int32:
			v, err := strconv.ParseInt(val, 10, 32)
			if err != nil {
				return err
			}
			value.Set(reflect.ValueOf(v))
		case reflect.Int64:
			v, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return err
			}
			value.Set(reflect.ValueOf(v))
		case reflect.Uint:
			v, err := strconv.ParseUint(val, 10, 0)
			if err != nil {
				return err
			}
			value.Set(reflect.ValueOf(v))
		case reflect.Uint8:
			v, err := strconv.ParseUint(val, 10, 8)
			if err != nil {
				return err
			}
			value.Set(reflect.ValueOf(v))
		case reflect.Uint16:
			v, err := strconv.ParseUint(val, 10, 16)
			if err != nil {
				return err
			}
			value.Set(reflect.ValueOf(v))
		case reflect.Uint32:
			v, err := strconv.ParseUint(val, 10, 32)
			if err != nil {
				return err
			}
			value.Set(reflect.ValueOf(v))
		case reflect.Uint64:
			v, err := strconv.ParseUint(val, 10, 64)
			if err != nil {
				return err
			}
			value.Set(reflect.ValueOf(v))
		}
	}

	return nil
}

func (c *envConfig) Save(ctx context.Context) error {
	return nil
}

func (c *envConfig) String() string {
	return "env"
}

func NewConfig(opts...config.Option) config.Config {
	options := config.NewOptions(opts...)
	if len(options.StructTag) == 0 {
		options.StructTag = DefaultStructTag
	}
	return &envConfig{opts:options}
}
