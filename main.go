package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"

	"github.com/wrfly/container-web-tty/config"
)

func appBefore(c *cli.Context) error {
	// set log-level
	if c.Bool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
	}
	return nil
}

func envVars(e string) []string {
	e = strings.ToUpper(e)
	return []string{"WEB_TTY_" + strings.Replace(e, "-", "_", -1)}
}

func main() {
	conf := config.Config{
		Backend: config.BackendConfig{},
	}
	appFlags := []cli.Flag{
		&cli.IntFlag{
			Name:        "port",
			Aliases:     []string{"p"},
			EnvVars:     envVars("port"),
			Usage:       "server port",
			Value:       8080,
			Destination: &conf.Port,
		},
		&cli.BoolFlag{
			Name:        "debug",
			Aliases:     []string{"d"},
			EnvVars:     envVars("debug"),
			Usage:       "debug log level",
			Destination: &conf.Debug,
		},
		&cli.StringFlag{
			Name:        "backend",
			Aliases:     []string{"b"},
			EnvVars:     envVars("backend"),
			Value:       "docker",
			Usage:       "backend type, docker or kubectl for now",
			Destination: &conf.Backend.Type,
		},
		&cli.StringFlag{
			Name:        "docker-path",
			EnvVars:     envVars("docker-path"),
			Value:       "/usr/bin/docker",
			Usage:       "docker cli path",
			Destination: &conf.Backend.DockerPath,
		},
		&cli.StringFlag{
			Name:        "kubectl-path",
			EnvVars:     envVars("kubectl-path"),
			Value:       "/usr/bin/kubectl",
			Usage:       "kubectl cli path",
			Destination: &conf.Backend.KubectlPath,
		},
		&cli.StringSliceFlag{
			Name:    "extra-args",
			EnvVars: envVars("extra-args"),
			Usage:   "extra args for your backend",
		},
		&cli.StringSliceFlag{
			Name:    "servers",
			EnvVars: envVars("servers"),
			Usage:   "upstream servers, for proxy mode",
		},
	}

	app := &cli.App{
		Name:   "reg",
		Usage:  "docker registry.v2 cli",
		Before: appBefore,
		Flags:  appFlags,
		EnableShellCompletion: true,
		Authors: []*cli.Author{
			&cli.Author{Name: "wrfly", Email: "mr.wrfly@gmail.com"},
		},
		Version: fmt.Sprintf("v%s-%s @%s", Version, CommitID, BuildAt),
		Action: func(c *cli.Context) error {
			conf.Backend.ExtraArgs = c.StringSlice("extra-args")
			conf.Servers = c.StringSlice("servers")
			logrus.Debugf("got config: %+v", conf)
			run(c, conf)
			return nil
		},
	}

	app.Run(os.Args)
}