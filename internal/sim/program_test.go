package sim

import "testing"

func TestParseProgram(t *testing.T) {
	prog, err := ParseProgram("COMPUTE 3; ACCESS 0x1000 r; SYSCALL read 2; EXIT")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	if len(prog) != 4 {
		t.Fatalf("expected 4 instructions, got %d", len(prog))
	}

	if prog[0].Op != "COMPUTE" || prog[0].Arg != 3 {
		t.Fatalf("unexpected first instruction: %+v", prog[0])
	}
	if prog[1].Op != "ACCESS" || prog[1].Addr != 0x1000 || prog[1].Access != AccessRead {
		t.Fatalf("unexpected second instruction: %+v", prog[1])
	}
	if prog[2].Op != "SYSCALL" || prog[2].Syscall != SysRead || prog[2].Arg != 2 {
		t.Fatalf("unexpected third instruction: %+v", prog[2])
	}
}
