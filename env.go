package env

import (
	"context"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/imdario/mergo"
	"github.com/unistack-org/micro/v3/config"
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

func (c *envConfig) Init(opts ...config.Option) error {
	for _, o := range opts {
		o(&c.opts)
	}
	return nil
}

func (c *envConfig) Load(ctx context.Context) error {
	for _, fn := range c.opts.BeforeLoad {
		if err := fn(ctx, c); err != nil && !c.opts.AllowFail {
			return err
		}
	}

	src, err := config.Zero(c.opts.Struct)
	if err == nil {
		err = c.fillValues(ctx, reflect.ValueOf(src))
	}
	if err != nil && !c.opts.AllowFail {
		return err
	}

	if err = mergo.Merge(c.opts.Struct, src, mergo.WithOverride, mergo.WithTypeCheck, mergo.WithAppendSlice); err != nil && !c.opts.AllowFail {
		return err
	}

	for _, fn := range c.opts.AfterLoad {
		if err := fn(ctx, c); err != nil && !c.opts.AllowFail {
			return err
		}
	}

	return nil
}

func (c *envConfig) fillValue(ctx context.Context, value reflect.Value, val string) error {
	switch value.Kind() {
	case reflect.Map:
		t := value.Type()
		nvals := strings.FieldsFunc(val, func(c rune) bool { return c == ',' || c == ';' })
		if value.IsNil() {
			value.Set(reflect.MakeMapWithSize(t, len(nvals)))
		}
		kt := t.Key()
		et := t.Elem()
		for _, nval := range nvals {
			kv := strings.FieldsFunc(nval, func(c rune) bool { return c == '=' })
			mkey := reflect.Indirect(reflect.New(kt))
			mval := reflect.Indirect(reflect.New(et))
			if err := c.fillValue(ctx, mkey, kv[0]); err != nil {
				return err
			}
			if err := c.fillValue(ctx, mval, kv[1]); err != nil {
				return err
			}
			value.SetMapIndex(mkey, mval)
		}
	case reflect.Slice, reflect.Array:
		nvals := strings.FieldsFunc(val, func(c rune) bool { return c == ',' || c == ';' })
		value.Set(reflect.MakeSlice(reflect.SliceOf(value.Type().Elem()), len(nvals), len(nvals)))
		for idx, nval := range nvals {
			nvalue := reflect.Indirect(reflect.New(value.Type().Elem()))
			if err := c.fillValue(ctx, nvalue, nval); err != nil {
				return err
			}
			value.Index(idx).Set(nvalue)
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
		return config.ErrInvalidStruct
	}

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
		switch value.Kind() {
		case reflect.Struct:
			value.Set(reflect.Indirect(reflect.New(value.Type())))
			if err := c.fillValues(ctx, value); err != nil {
				return err
			}
			continue
		case reflect.Ptr:
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
			continue
		}

		if err := c.fillValue(ctx, value, val); err != nil {
			return err
		}
	}

	return nil
}

func (c *envConfig) Save(ctx context.Context) error {
	for _, fn := range c.opts.BeforeSave {
		if err := fn(ctx, c); err != nil && !c.opts.AllowFail {
			return err
		}
	}

	for _, fn := range c.opts.AfterSave {
		if err := fn(ctx, c); err != nil && !c.opts.AllowFail {
			return err
		}
	}

	return nil
}

func (c *envConfig) String() string {
	return "env"
}

func NewConfig(opts ...config.Option) config.Config {
	options := config.NewOptions(opts...)
	if len(options.StructTag) == 0 {
		options.StructTag = DefaultStructTag
	}
	return &envConfig{opts: options}
}
