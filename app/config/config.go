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

type FilePathOption []string

func (o *FilePathOption) String() string {
	return fmt.Sprintf("%v", *o)
}

func (o *FilePathOption) Set(s string) error {
	*o = append(*o, s)
	return nil
}

type Config struct {
	FilePaths []string
	Cmd       string
}

func NewConfig() (*Config, error) {
	var paths FilePathOption
	flag.Var(&paths, "f", "specify file path")
	flag.Parse()

	if len(paths) == 0 {
		return nil, ConfigParseError{"file path not specified"}
	}
	if len(flag.Args()) != 1 {
		return nil, ConfigParseError{"too many argument or not specified"}
	}
	cmd := flag.Args()[0]
	return &Config{
		FilePaths: paths,
		Cmd:       cmd,
	}, nil
}
