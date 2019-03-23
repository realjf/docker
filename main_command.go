package main

import (
	"docker/cgroups/subsystems"
	"docker/container"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var runCommand = cli.Command{
	Name: "run",
	Usage: `Create a container with namespace and cgroups limit docker run -ti [command]`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name: "ti",
			Usage: "enable tty",
		},
		cli.BoolFlag{
			Name: "d",
			Usage: "detach container",
		},
		cli.StringFlag{
			Name: "m",
			Usage: "memory limit",
		},
		cli.StringFlag{
			Name: "cpushare",
			Usage: "cpushare limit",
		},
		cli.StringFlag{
			Name: "cpuset",
			Usage: "cpuset limit",
		},
		cli.StringFlag{
			Name: "name",
			Usage: "container name",
		},
		cli.StringFlag{
			Name: "v",
			Usage: "volume",
		},
		cli.StringSliceFlag{
			Name: "e",
			Usage: "set environment",
		},
		cli.StringFlag{
			Name: "net",
			Usage: "container network",
		},
		cli.StringSliceFlag{
			Name: "p",
			Usage: "port mapping",
		},
	},

	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("Missing container command")
		}
		var cmdArray []string
		for _, arg := range context.Args() {
			cmdArray = append(cmdArray, arg)
		}

		cmdArray = cmdArray[0:]

		createTty := context.Bool("ti")
		detach := context.Bool("d")


		if createTty && detach {
			return fmt.Errorf("ti and d paramter can not both provided")
		}

		resConf := &subsystems.ResourceConfig{
			MemoryLimit: context.String("m"),
			CpuSet: context.String("cpuset"),
			CpuShare: context.String("cpushare"),
		}
		log.Infof("createTty %v", createTty)

		// 把volume参数传给Run函数
		volume := context.String("v")

		Run(createTty, cmdArray, resConf, volume)
		return nil
	},
}

// 初始化命令行
var initCommand = cli.Command{
	Name: "init",
	Usage: "Init container process run user's process in container. Do not call it outside",
	Action: func(context *cli.Context) error {
		log.Infof("init come on")
		err := container.RunContainerInitProcess()
		return err
	},
}

var listCommand = cli.Command{
	Name: "ps",
	Usage: "list all the containers",
	Action: func(context *cli.Context) error {
		ListContainers()
		return nil
	},
}

var logCommand = cli.Command{
	Name: "logs",
	Usage:"print logs of a container",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("Please input your container name")
		}
		containerName := context.Args().Get(0)
		logContainer(containerName)
		return nil
	},
}

var commitCommand = cli.Command{
	Name: "commit",
	Usage: "commit a container into image",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 2 {
			return fmt.Errorf("Missing container name and image name")
		}
		containerName := context.Args().Get(0)
		imageName := context.Args().Get(1)
		commitContainer(containerName, imageName)
		return nil
	},
}
