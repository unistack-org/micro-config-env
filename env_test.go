package env_test

import (
	"github.com/unistack-org/micro/v3/config"
	env	"github.com/unistack-org/micro-config-env"
	"testing"
	"context"
	"os"
)

type Config struct {
	StringValue string `env:"STRING_VALUE"`
	BoolValue bool `env:"BOOL_VALUE"`
	StringSlice []string `env:"STRING_SLICE"`
	IntSlice []int `env:"INT_SLICE"`
	MapStringValue map[string]string `env:"MAP_STRING"`
	MapIntValue map[string]int `env:"MAP_INT"`
}

func TestEnv(t *testing.T) {
	ctx := context.Background()
	conf := &Config{StringValue: "before_load"}
	cfg := env.NewConfig(config.Struct(conf))

	if err := cfg.Init(); err != nil {
		t.Fatal(err)
	}

	if err := cfg.Load(ctx); err !=nil {
		t.Fatal(err)
	}

	if conf.StringValue != "before_load" {
		t.Fatalf("something wrong with env config: %v", conf)
	}


	os.Setenv("STRING_VALUE","STRING_VALUE")
	os.Setenv("BOOL_VALUE","true")
	os.Setenv("STRING_SLICE", "STRING_SLICE1,STRING_SLICE2;STRING_SLICE3")
	os.Setenv("INT_SLICE", "1,2,3,4,5")
	os.Setenv("MAP_STRING", "key1=val1,key2=val2")
	os.Setenv("MAP_INT", "key1=1,key2=2")

	if err := cfg.Load(ctx); err !=nil {
		t.Fatal(err)
	}
	if conf.StringValue != "STRING_VALUE" {
		t.Fatalf("something wrong with env config: %v", conf)
	}

	if !conf.BoolValue {
		t.Fatalf("something wrong with env config: %v", conf)
	}

	if len(conf.StringSlice) != 3 {
		t.Fatalf("something wrong with env config: %v", conf)
	}

	if len(conf.MapStringValue) != 2 {
		t.Fatalf("something wrong with env config: %v", conf)
	}

	if len(conf.MapIntValue) != 2 {
		t.Fatalf("something wrong with env config: %v", conf)
	}

}
