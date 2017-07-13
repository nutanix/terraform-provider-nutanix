package cfg

/*
Uses Viper to read values from ENV variable or commandline flags transparently.

a flag with name "x-y" can be set in CLI as <binary> --x-y
if the same flag has to be set in ENV, it has tobe set as X_Y
*/

import (
	"flag"
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"strings"
)

// ConfigPath is path for virtualmachineconfig package
var ConfigPath string
// SchemaPath is path for virtualmachineschema package
var SchemaPath string
// StructPath is path for virtualmachineconfig package
var StructPath string

func init() {
	flag.String("schema-path", "", "path for creating the schema file")
	flag.String("config-path", "", "path for creating the config file")
	flag.String("sdk-path", "", "path for sdk repo")

	//pflag configuration
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	viper.BindPFlags(pflag.CommandLine)
	pflag.Parse()
	pflag.Visit(func(f *pflag.Flag) {
		fmt.Printf("GOPATH %s overridden: %s -> %s\n", f.Name, f.DefValue, f.Value)
	})

	//Env configuration
	viper.AutomaticEnv()
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)

	//Config Init
	SchemaPath = viper.GetString("schema-path")
	ConfigPath = viper.GetString("config-path")
	StructPath = viper.GetString("sdk-path")
}
