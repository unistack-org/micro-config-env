package env

import (
	"context"
	"os"
	"testing"
	"time"

	"go.unistack.org/micro/v4/config"
	rutil "go.unistack.org/micro/v4/util/reflect"
)

type Config struct {
	StringValue    string            `env:"STRING_VALUE,STRING_VALUE2"`
	BoolValue      bool              `env:"BOOL_VALUE"`
	StringSlice    []string          `env:"STRING_SLICE"`
	IntSlice       []int             `env:"INT_SLICE"`
	MapStringValue map[string]string `env:"MAP_STRING"`
	MapIntValue    map[string]int    `env:"MAP_INT"`
}

func TestMerge(t *testing.T) {
	defer func() {
		for _, v := range []string{"STRING_VALUE", "BOOL_VALUE", "STRING_SLICE", "INT_SLICE", "MAP_STRING", "MAP_INT"} {
			if err := os.Unsetenv(v); err != nil {
				t.Fatal(err)
			}
		}
	}()

	ctx := context.Background()
	type Nested struct {
		Name string `env:"NAME_VALUE"`
	}
	type Cfg struct {
		Name   string `env:"NAME_VALUE"`
		Nested Nested
	}

	conf := &Cfg{}

	cfg := NewConfig(config.Struct(conf))

	if err := cfg.Init(); err != nil {
		t.Fatal(err)
	}

	if err := cfg.Load(ctx, config.LoadOverride(true), config.LoadAppend(true)); err != nil {
		t.Fatal(err)
	}

	w, err := cfg.Watch(ctx, config.WatchInterval(50*time.Millisecond, 500*time.Millisecond))
	defer func() {
		if err := w.Stop(); err != nil {
			t.Fatal(err)
		}
	}()
	if err != nil {
		t.Fatal(err)
	}

	os.Setenv("NAME_VALUE", "after")
	changes, err := w.Next()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range changes {
		if err := rutil.SetFieldByPath(conf, v, k); err != nil {
			t.Fatal(err)
		}
	}

	if conf.Name != "after" || conf.Nested.Name != "after" {
		t.Fatalf("changes %#+v not applied to %#+v", changes, conf)
	}
}

func TestLoad(t *testing.T) {
	defer func() {
		for _, v := range []string{"STRING_VALUE", "STRING_VALUE2", "BOOL_VALUE", "STRING_SLICE", "INT_SLICE", "MAP_STRING", "MAP_INT"} {
			if err := os.Unsetenv(v); err != nil {
				t.Fatal(err)
			}
		}
	}()

	ctx := context.Background()
	conf := &Config{StringValue: "before_load"}
	cfg := NewConfig(config.Struct(conf))

	if err := cfg.Init(); err != nil {
		t.Fatal(err)
	}

	if err := cfg.Load(ctx, config.LoadOverride(true), config.LoadAppend(true)); err != nil {
		t.Fatal(err)
	}

	if conf.StringValue != "before_load" {
		t.Fatalf("something wrong with env config: %#+v", conf)
	}

	os.Setenv("STRING_VALUE", "STRING_VALUE")
	os.Setenv("BOOL_VALUE", "true")
	os.Setenv("STRING_SLICE", "STRING_SLICE1,STRING_SLICE2;STRING_SLICE3")
	os.Setenv("INT_SLICE", "1,2,3,4,5")
	os.Setenv("MAP_STRING", "key1=val1,key2=val2")
	os.Setenv("MAP_INT", "key1=1,key2=2")

	if err := cfg.Load(ctx, config.LoadOverride(true), config.LoadAppend(true)); err != nil {
		t.Fatal(err)
	}
	if conf.StringValue != "STRING_VALUE" {
		t.Fatalf("something wrong with env config: %#+v", conf)
	}

	if !conf.BoolValue {
		t.Fatalf("something wrong with env config: %#+v", conf)
	}

	if len(conf.StringSlice) != 3 {
		t.Fatalf("something wrong with env config: %#+v", conf.StringSlice)
	}

	if len(conf.MapStringValue) != 2 {
		t.Fatalf("something wrong with env config: %#+v", conf.MapStringValue)
	}

	if len(conf.MapIntValue) != 2 {
		t.Fatalf("something wrong with env config: %#+v", conf.MapIntValue)
	}

	for _, v := range []string{"STRING_VALUE", "STRING_VALUE2", "BOOL_VALUE", "STRING_SLICE", "INT_SLICE", "MAP_STRING", "MAP_INT"} {
		if err := os.Unsetenv(v); err != nil {
			t.Fatal(err)
		}
	}
}

func TestSave(t *testing.T) {
	defer func() {
		for _, v := range []string{"STRING_VALUE", "STRING_VALUE2", "BOOL_VALUE", "STRING_SLICE", "INT_SLICE", "MAP_STRING", "MAP_INT"} {
			if err := os.Unsetenv(v); err != nil {
				t.Fatal(err)
			}
		}
	}()

	ctx := context.Background()
	conf := &Config{StringValue: "MICRO_CONFIG_ENV"}
	cfg := NewConfig(config.Struct(conf))

	if err := cfg.Init(); err != nil {
		t.Fatal(err)
	}

	if _, ok := os.LookupEnv("STRING_VALUE"); ok {
		if err := os.Unsetenv("STRING_VALUE"); err != nil {
			t.Fatal(err)
		}
	}

	if err := os.Unsetenv("STRING_VALUE"); err != nil {
		t.Fatal(err)
	}

	if err := cfg.Save(ctx); err != nil {
		t.Fatal(err)
	}

	if _, ok := os.LookupEnv("STRING_VALUE"); !ok {
		t.Fatal("env value STRING_VALUE=MICRO_CONFIG_ENV not set")
	}

	if err := os.Unsetenv("STRING_VALUE"); err != nil {
		t.Fatal(err)
	}

	for _, tv := range []string{"STRING_VALUE", "BOOL_VALUE", "STRING_SLICE", "INT_SLICE", "MAP_STRING", "MAP_INT"} {
		if v, ok := os.LookupEnv("STRING_VALUE"); ok {
			t.Fatalf("env value %s=%s set", tv, v)
		}
	}

	for _, v := range []string{"STRING_VALUE", "STRING_VALUE2", "BOOL_VALUE", "STRING_SLICE", "INT_SLICE", "MAP_STRING", "MAP_INT"} {
		if err := os.Unsetenv(v); err != nil {
			t.Fatal(err)
		}
	}
}

func TestLoadMultiple(t *testing.T) {
	defer func() {
		for _, v := range []string{"STRING_VALUE", "STRING_VALUE2", "BOOL_VALUE", "STRING_SLICE", "INT_SLICE", "MAP_STRING", "MAP_INT"} {
			if err := os.Unsetenv(v); err != nil {
				t.Fatal(err)
			}
		}
	}()

	ctx := context.Background()
	conf := &Config{StringValue: "before_load"}
	cfg := NewConfig(config.Struct(conf))

	if err := cfg.Init(); err != nil {
		t.Fatal(err)
	}

	if err := cfg.Load(ctx, config.LoadOverride(true), config.LoadAppend(true)); err != nil {
		t.Fatal(err)
	}

	if conf.StringValue != "before_load" {
		t.Fatalf("something wrong with env config: %#+v", conf)
	}

	os.Setenv("STRING_VALUE", "STRING_VALUE1")
	os.Setenv("STRING_VALUE2", "STRING_VALUE2")
	defer func() {
		for _, v := range []string{"STRING_VALUE", "STRING_VALUE2"} {
			if err := os.Unsetenv(v); err != nil {
				t.Fatal(err)
			}
		}
	}()

	if err := cfg.Load(ctx, config.LoadOverride(true), config.LoadAppend(true)); err != nil {
		t.Fatal(err)
	}
	if conf.StringValue != "STRING_VALUE2" {
		t.Fatalf("something wrong with env config: %#+v", conf)
	}

	for _, v := range []string{"STRING_VALUE", "STRING_VALUE2", "BOOL_VALUE", "STRING_SLICE", "INT_SLICE", "MAP_STRING", "MAP_INT"} {
		if err := os.Unsetenv(v); err != nil {
			t.Fatal(err)
		}
	}
}
