package config

import (
	"flag"
	"fmt"
)

const USAGE = `
usage: monitor [options] command
  options:
    -f (required) secify file path
`

func (e ConfigParseError) Error() string {
	return fmt.Sprintf("%s\n%s", e.message, USAGE)
}

type ConfigParseError struct {
	message string
}

type Config struct {
	FilePath string
	Cmd      string
}

func NewConfig() (*Config, error) {
	f := flag.String("f", "", "specify file path")
	flag.Parse()

	if len(*f) == 0 {
		return nil, ConfigParseError{"file path not specified"}
	}
	if len(flag.Args()) != 1 {
		return nil, ConfigParseError{"too many argument or not specified"}
	}
	cmd := flag.Args()[0]
	return &Config{
		FilePath: *f,
		Cmd:      cmd,
	}, nil
}
