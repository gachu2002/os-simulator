package sim

import "testing"

func TestFilesystemResolveAndInvariants(t *testing.T) {
	fs := NewFileSystem()
	inodeID, traversal, err := fs.Resolve("/docs/readme.txt")
	if err != nil {
		t.Fatalf("resolve failed: %v", err)
	}
	if inodeID == 0 || len(traversal) != 3 {
		t.Fatalf("unexpected traversal result inode=%d traversal=%v", inodeID, traversal)
	}

	if err := fs.Invariants(); err != nil {
		t.Fatalf("filesystem invariants failed: %v", err)
	}
}

func TestFilesystemBlockMapReadWrite(t *testing.T) {
	fs := NewFileSystem()
	inodeID, _, err := fs.Resolve("/docs/readme.txt")
	if err != nil {
		t.Fatalf("resolve failed: %v", err)
	}

	written, blocks, off := fs.WriteInode(inodeID, []byte("-extra-data"), 0)
	if written != len("-extra-data") || len(blocks) == 0 || off == 0 {
		t.Fatalf("write mapping unexpected written=%d blocks=%v off=%d", written, blocks, off)
	}

	data, readBlocks, _, err := fs.ReadInode(inodeID, 5, 0)
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}
	if len(data) != 5 || len(readBlocks) == 0 {
		t.Fatalf("read mapping unexpected data=%q blocks=%v", string(data), readBlocks)
	}
	if err := fs.Invariants(); err != nil {
		t.Fatalf("filesystem invariants failed after io: %v", err)
	}
}
