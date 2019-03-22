package subsystems

type ResourceConfig struct {
	MemoryLimit string
	CpuShare string
	CpuSet string
}

type Subsystem interface {
	Name() string // 返回subsystem的名字
	Set(path string, res *ResourceConfig) error // 设置某个cgroup在这个subsystem中的资源限制
	Apply(path string, pid int) error // 将进程添加到某个cgroup
	Remove(path string) error // 移除某个cgroup
}

var (
	SubsystemsIns = []Subsystem{
		&CpusetSubSystem{},
		&MemorySubSystem{},
		&CpuSubSystem{},
	}
)


