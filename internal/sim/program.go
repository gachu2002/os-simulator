package sim

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseProgram(spec string) ([]Instruction, error) {
	if strings.TrimSpace(spec) == "" {
		return nil, fmt.Errorf("program is empty")
	}

	parts := strings.Split(spec, ";")
	prog := make([]Instruction, 0, len(parts))
	for _, part := range parts {
		line := strings.TrimSpace(part)
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		op := strings.ToUpper(fields[0])
		switch op {
		case "COMPUTE", "BLOCK":
			if len(fields) != 2 {
				return nil, fmt.Errorf("%s expects one argument", op)
			}
			n, err := strconv.Atoi(fields[1])
			if err != nil || n <= 0 {
				return nil, fmt.Errorf("%s arg must be positive int", op)
			}
			prog = append(prog, Instruction{Op: op, Arg: n})
		case "EXIT":
			if len(fields) != 1 {
				return nil, fmt.Errorf("EXIT takes no args")
			}
			prog = append(prog, Instruction{Op: op})
		case "SYSCALL":
			if len(fields) < 2 || len(fields) > 3 {
				return nil, fmt.Errorf("SYSCALL expects 1 or 2 arguments")
			}
			sys := strings.ToLower(fields[1])
			if sys != SysOpen && sys != SysRead && sys != SysWrite && sys != SysSleep && sys != SysExit {
				return nil, fmt.Errorf("unknown syscall %q", sys)
			}
			arg := 0
			argText := ""
			if len(fields) == 3 {
				n, err := strconv.Atoi(fields[2])
				if err == nil {
					arg = n
				} else {
					argText = fields[2]
				}
			}
			prog = append(prog, Instruction{Op: op, Syscall: sys, Arg: arg, ArgText: argText})
		case "ACCESS":
			if len(fields) != 3 {
				return nil, fmt.Errorf("ACCESS expects 2 arguments")
			}
			addr, err := strconv.ParseUint(fields[1], 0, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid access address %q", fields[1])
			}
			mode := strings.ToLower(fields[2])
			if mode != string(AccessRead) && mode != string(AccessWrite) {
				return nil, fmt.Errorf("ACCESS mode must be r or w")
			}
			prog = append(prog, Instruction{Op: op, Addr: addr, Access: AccessType(mode)})
		default:
			return nil, fmt.Errorf("unknown op %q", op)
		}
	}

	if len(prog) == 0 {
		return nil, fmt.Errorf("program has no instructions")
	}

	return prog, nil
}
