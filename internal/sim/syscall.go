package sim

import "fmt"

const (
	SysOpen  = "open"
	SysRead  = "read"
	SysWrite = "write"
	SysSleep = "sleep"
	SysExit  = "exit"
)

type SyscallResult struct {
	ReturnValue uint64
	Blocked     bool
	SleepTicks  Tick
	Exit        bool
	AsyncDevice string
	AsyncOp     string
	AsyncBytes  int
	Path        string
	Traversal   []int
	Blocks      []int
	FD          int
}

type SyscallDispatcher interface {
	Handle(proc *Process, name string, arg int, argText string) (SyscallResult, error)
}

type KernelDispatcher struct {
	fs *FileSystem
}

func NewKernelDispatcher(fs *FileSystem) *KernelDispatcher {
	return &KernelDispatcher{fs: fs}
}

func (d *KernelDispatcher) Handle(proc *Process, name string, arg int, argText string) (SyscallResult, error) {
	switch name {
	case SysOpen:
		path := argText
		if path == "" {
			path = "/docs/readme.txt"
		}
		inodeID, traversal, err := d.fs.Resolve(path)
		if err != nil {
			return SyscallResult{}, err
		}
		fd := proc.NextFD
		proc.NextFD++
		proc.OpenFiles[fd] = OpenFile{Path: path, InodeID: inodeID, Offset: 0}
		return SyscallResult{ReturnValue: uint64(fd), FD: fd, Path: path, Traversal: traversal}, nil
	case SysRead:
		if arg < 0 {
			return SyscallResult{}, fmt.Errorf("read size must be non-negative")
		}
		fd, ok := firstOpenFD(proc)
		if !ok {
			return SyscallResult{}, fmt.Errorf("read requires open file")
		}
		return SyscallResult{Blocked: true, AsyncDevice: DeviceDisk, AsyncOp: SysRead, AsyncBytes: arg, FD: fd}, nil
	case SysWrite:
		if arg < 0 {
			return SyscallResult{}, fmt.Errorf("write size must be non-negative")
		}
		fd, ok := firstOpenFD(proc)
		if !ok {
			return SyscallResult{}, fmt.Errorf("write requires open file")
		}
		return SyscallResult{Blocked: true, AsyncDevice: DeviceTerminal, AsyncOp: SysWrite, AsyncBytes: arg, FD: fd}, nil
	case SysSleep:
		if arg <= 0 {
			return SyscallResult{}, fmt.Errorf("sleep ticks must be positive")
		}
		return SyscallResult{Blocked: true, SleepTicks: Tick(arg)}, nil
	case SysExit:
		return SyscallResult{Exit: true}, nil
	default:
		return SyscallResult{}, fmt.Errorf("unknown syscall %q", name)
	}
}

func firstOpenFD(proc *Process) (int, bool) {
	fd := 0
	for k := range proc.OpenFiles {
		if fd == 0 || k < fd {
			fd = k
		}
	}
	if fd == 0 {
		return 0, false
	}
	return fd, true
}
