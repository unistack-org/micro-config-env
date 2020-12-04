package env

import (
	"github.com/unistack-org/micro/v3/config"
	"context"
	"reflect"
	"strings"
	"os"
	"strconv"
	"errors"
)

var (
	DefaultStructTag = "env"
	ErrInvalidStruct = errors.New("invalid struct specified")
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
	valueOf := reflect.ValueOf(c.opts.Struct)

	if err := c.fillValues(ctx, valueOf); err != nil {
		return err
	}

	return nil
}

func (c *envConfig) fillValue(ctx context.Context, value reflect.Value, val string) error {
	switch value.Kind() {
	case reflect.Slice, reflect.Array:
		nvals := strings.FieldsFunc(val, func(c rune) bool { return c == ',' || c == ';' })
	//	value = value.Elem()
		value.Set(reflect.MakeSlice(reflect.SliceOf(value.Type().Elem()), len(nvals), len(nvals)))
		for idx, nval := range nvals {
			nvalue := reflect.Indirect(reflect.New(value.Type().Elem()))
			if err := c.fillValue(ctx, nvalue, nval); err != nil {
				return err
			}
			value.Index(idx).Set(nvalue)
			//value.Set(reflect.Append(value, nvalue))
		}
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
		value.Set(reflect.ValueOf(float32(v)))
	case reflect.Float64:
		v, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		}
		value.Set(reflect.ValueOf(float64(v)))
	case reflect.Int:
		v, err := strconv.ParseInt(val, 10, 0)
		if err != nil {
			return err
		}
		value.Set(reflect.ValueOf(int(v)))
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
		value.Set(reflect.ValueOf(int16(v)))
	case reflect.Int32:
		v, err := strconv.ParseInt(val, 10, 32)
		if err != nil {
			return err
		}
		value.Set(reflect.ValueOf(int32(v)))
	case reflect.Int64:
		v, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}
		value.Set(reflect.ValueOf(int64(v)))
	case reflect.Uint:
		v, err := strconv.ParseUint(val, 10, 0)
		if err != nil {
			return err
		}
		value.Set(reflect.ValueOf(uint(v)))
	case reflect.Uint8:
		v, err := strconv.ParseUint(val, 10, 8)
		if err != nil {
			return err
		}
		value.Set(reflect.ValueOf(uint8(v)))
	case reflect.Uint16:
		v, err := strconv.ParseUint(val, 10, 16)
		if err != nil {
			return err
		}
		value.Set(reflect.ValueOf(uint16(v)))
	case reflect.Uint32:
		v, err := strconv.ParseUint(val, 10, 32)
		if err != nil {
			return err
		}
		value.Set(reflect.ValueOf(uint32(v)))
	case reflect.Uint64:
		v, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return err
		}
		value.Set(reflect.ValueOf(uint64(v)))
	}
	return nil
}

func (c *envConfig) fillValues(ctx context.Context, valueOf reflect.Value) error {
	var values reflect.Value

	if valueOf.Kind() == reflect.Ptr {
		values = valueOf.Elem()
	} else {
	values = valueOf
	}

	if values.Kind() == reflect.Invalid {
		return ErrInvalidStruct
	}

	//	if values.Kind() != reflect.Struct {
	//		return c.fillValue(ctx, values)
	//	}

	fields := values.Type()

	for idx := 0; idx < fields.NumField(); idx++ {
		field := fields.Field(idx)
		value := values.Field(idx)
		if !value.CanSet() {
			continue
		}
		if len(field.PkgPath) != 0 {
			continue
		}
		if value.Kind() == reflect.Ptr {
			if value.IsNil() {
				if value.Type().Elem().Kind() != reflect.Struct {
					// nil pointer to a non-struct: leave it alone
					break
				}
				// nil pointer to struct: create a zero instance
				value.Set(reflect.New(value.Type().Elem()))
			}
			value = value.Elem()
			if err := c.fillValues(ctx, value); err != nil {
				return err
			}
			continue
		}
		tag, ok := field.Tag.Lookup(c.opts.StructTag)
		if !ok {
			return nil
		}
		val, ok := os.LookupEnv(tag)
		if !ok {
			return nil
		}

		if err := c.fillValue(ctx, value, val); err != nil {
			return err
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
