package main

import (
	"docker/cgroups"
	"docker/cgroups/subsystems"
	"docker/container"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

func Run(tty bool, comArray []string, res *subsystems.ResourceConfig, volume string) {
	parent, writePipe := container.NewParentProcess(tty, volume)
	if parent == nil {
		log.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	// 创建cgroup manager，并通过调用set和apply设置资源限制并使限制在容器上生效
	cgroupManager := cgroups.NewCgroupManager("mydocker-cgroup")
	defer cgroupManager.Destroy()
	// 设置资源限制
	cgroupManager.Set(res)

	// 将容器进程加入到各个subsystem挂载对应的cgroup中
	cgroupManager.Apply(parent.Process.Pid)

	// 对容器设置完限制之后，初始化容器
	sendInitCommand(comArray, writePipe)

	parent.Wait()

	// 删除
	mntURL := "/root/mnt"
	rootURL := "/root"
	container.DeleteWorkSpace(rootURL, mntURL, volume)
	os.Exit(0)
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	log.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}
