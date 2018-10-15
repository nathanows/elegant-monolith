package conf

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Parser is the base config parser
type Parser interface {
	LoadConfig() error
}

type parser struct {
	cmd       *cobra.Command
	conf      interface{}
	envPrefix string
	fileName  string
	filePath  string
}

// FileName allows overriding of config file name. Defaults to 'config.json'
func FileName(name string) func(*parser) {
	return func(p *parser) {
		p.fileName = name
	}
}

// FilePath allows overriding of default config file path. Defaults to project root.
func FilePath(path string) func(*parser) {
	return func(p *parser) {
		p.filePath = path
	}
}

// EnvPrefix allows specifying a prefix to be appended to all env vars
func EnvPrefix(prefix string) func(*parser) {
	return func(p *parser) {
		p.envPrefix = prefix
	}
}

// NewParser returns a properly initialized Parser
func NewParser(cmd *cobra.Command, conf interface{}, options ...func(*parser)) (Parser, error) {
	parser := parser{
		cmd:       cmd,
		conf:      conf,
		envPrefix: "",
		fileName:  "config",
		filePath:  "./",
	}

	for _, option := range options {
		option(&parser)
	}

	return &parser, nil
}

// LoadConfig loads the config from a file if specified, otherwise from the environment
func (p *parser) LoadConfig() error {
	if err := viper.BindPFlags(p.cmd.Flags()); err != nil {
		return err
	}

	viper.SetEnvPrefix(p.envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if configFile, _ := p.cmd.Flags().GetString("config"); configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName(p.fileName)
		viper.AddConfigPath(p.filePath)
	}

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return populateConfig(p.conf)
}

const tagPrefix = "viper"

func populateConfig(config interface{}) error {
	if err := recursivelySet(reflect.ValueOf(config), ""); err != nil {
		return err
	}

	return nil
}

func recursivelySet(val reflect.Value, prefix string) error {
	if val.Kind() != reflect.Ptr {
		return errors.New("WTF")
	}

	// dereference
	val = reflect.Indirect(val)
	if val.Kind() != reflect.Struct {
		return errors.New("FML")
	}

	// grab the type for this instance
	vType := reflect.TypeOf(val.Interface())

	// go through child fields
	for i := 0; i < val.NumField(); i++ {
		thisField := val.Field(i)
		thisType := vType.Field(i)
		tag := prefix + getTag(thisType)

		switch thisField.Kind() {
		case reflect.Struct:
			if err := recursivelySet(thisField.Addr(), tag+"."); err != nil {
				return err
			}
		case reflect.Int:
			fallthrough
		case reflect.Int32:
			fallthrough
		case reflect.Int64:
			// you can only set with an int64 -> int
			configVal := int64(viper.GetInt(tag))
			thisField.SetInt(configVal)
		case reflect.String:
			thisField.SetString(viper.GetString(tag))
		case reflect.Bool:
			thisField.SetBool(viper.GetBool(tag))
		default:
			return fmt.Errorf("unexpected type detected ~ aborting: %s", thisField.Kind())
		}
	}

	return nil
}

func getTag(field reflect.StructField) string {
	// check if maybe we have a special magic tag
	tag := field.Tag
	if tag != "" {
		for _, prefix := range []string{tagPrefix, "mapstructure", "json"} {
			if v := tag.Get(prefix); v != "" {
				return v
			}
		}
	}

	return field.Name
}
