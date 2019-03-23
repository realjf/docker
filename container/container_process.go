package container

import (
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"syscall"
)

var (
	RUNNING             string = "running"
	STOP                string = "stopped"
	Exit                string = "exited"
	DefaultInfoLocation string = "/var/run/mydocker/%s/"
	ConfigName          string = "config.json"
	ContainerLogFile    string = "container.log"
	RootUrl             string = "/root"
	MntUrl              string = "/root/mnt/%s"
	WriteLayerUrl       string = "/root/writeLayer/%s"
)

type ContainerInfo struct {
	Pid         string   `json:"pid"`         // 容器的init进程在宿主机上的pid
	Id          string   `json:"id"`          // 容器id
	Name        string   `json:"name"`        // 容器名
	Command     string   `json:"command"`     // 容器内init运行命令
	CreatedTime string   `json:"createTime"`  // 创建时间
	Status      string   `json:"status"`      // 容器状态
	Volume      string   `json:"volume"`      // 容器的数据卷
	PortMapping []string `json:"portmapping"` // 端口映射
}

func NewParentProcess(tty bool, volume string) (*exec.Cmd, *os.File) {
	readPipe, writePipe, err := NewPipe()
	if err != nil {
		log.Errorf("New pipe error %v", err)
		return nil, nil
	}

	//initCmd, err := os.Readlink("/proc/self/exe")
	//if err != nil {
	//	log.Errorf("get init process error %v", err)
	//	return nil, nil
	//}

	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	cmd.ExtraFiles = []*os.File{readPipe}
	mntURL := "/root/mnt/"
	rootURL := "/root/"
	NewWorkSpace(rootURL, mntURL, volume)
	cmd.Dir = mntURL
	return cmd, writePipe
}

func NewPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	return read, write, nil
}
