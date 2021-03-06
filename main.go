package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"

	"github.com/b4b4r07/stein/pkg/logging"
	"github.com/mitchellh/cli"
)

const (
	// AppName is the application name
	AppName = "stein"
	// Version is the application version
	Version = "0.2.3"

	envEnvPrefix = "STEIN_"
)

// CLI represents the command-line interface
type CLI struct {
	Stdout io.Writer
	Stderr io.Writer
}

func main() {
	logWriter, err := logging.LogOutput()
	if err != nil {
		panic(err)
	}
	log.SetOutput(logWriter)

	log.Printf("[INFO] Stein version: %s", Version)
	log.Printf("[INFO] Go runtime version: %s", runtime.Version())
	log.Printf("[INFO] CLI args: %#v", os.Args)

	stein := CLI{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	app := cli.NewCLI(AppName, Version)
	app.Args = os.Args[1:]
	app.Commands = map[string]cli.CommandFactory{
		"apply": func() (cli.Command, error) {
			return &ApplyCommand{CLI: stein}, nil
		},
		"fmt": func() (cli.Command, error) {
			return &FmtCommand{CLI: stein}, nil
		},
	}
	exitStatus, err := app.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s: %v\n", AppName, err)
	}
	os.Exit(exitStatus)
}

func (c CLI) exit(msg interface{}) int {
	switch m := msg.(type) {
	case int:
		return m
	case nil:
		return 0
	case string:
		fmt.Fprintf(c.Stdout, "%s\n", m)
		return 0
	case error:
		fmt.Fprintf(c.Stderr, "[ERROR] %s: %s\n", AppName, m.Error())
		return 1
	default:
		panic(msg)
	}
}
