package sim

import (
	"fmt"
	"sort"
	"strings"
)

type InodeType string

const (
	InodeDir  InodeType = "dir"
	InodeFile InodeType = "file"
)

type OpenFile struct {
	Path    string
	InodeID int
	Offset  int
}

type Inode struct {
	ID           int
	Type         InodeType
	Name         string
	Parent       int
	Entries      map[string]int
	DirectBlocks []int
	Size         int
}

type FileSystem struct {
	blockSize int
	nextInode int
	nextBlock int
	inodes    map[int]*Inode
	blocks    map[int][]byte
	rootInode int
}

func NewFileSystem() *FileSystem {
	fs := &FileSystem{
		blockSize: 16,
		nextInode: 1,
		nextBlock: 0,
		inodes:    map[int]*Inode{},
		blocks:    map[int][]byte{},
	}

	root := fs.newDir("", 0)
	root.Parent = root.ID
	fs.rootInode = root.ID

	docs := fs.newDir("docs", root.ID)
	root.Entries["docs"] = docs.ID

	readme := fs.newFile("readme.txt", docs.ID)
	docs.Entries["readme.txt"] = readme.ID
	_, _, _ = fs.WriteInode(readme.ID, []byte("ostep-sim"), 0)

	return fs
}

func (fs *FileSystem) newDir(name string, parent int) *Inode {
	inode := &Inode{ID: fs.nextInode, Type: InodeDir, Name: name, Parent: parent, Entries: map[string]int{}}
	fs.nextInode++
	fs.inodes[inode.ID] = inode
	return inode
}

func (fs *FileSystem) newFile(name string, parent int) *Inode {
	inode := &Inode{ID: fs.nextInode, Type: InodeFile, Name: name, Parent: parent}
	fs.nextInode++
	fs.inodes[inode.ID] = inode
	return inode
}

func (fs *FileSystem) Resolve(path string) (int, []int, error) {
	if path == "" || path[0] != '/' {
		return 0, nil, fmt.Errorf("path must be absolute")
	}
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 1 && parts[0] == "" {
		return fs.rootInode, []int{fs.rootInode}, nil
	}

	current := fs.rootInode
	traversal := []int{current}
	for _, part := range parts {
		inode := fs.inodes[current]
		if inode.Type != InodeDir {
			return 0, nil, fmt.Errorf("non-directory in path")
		}
		next, ok := inode.Entries[part]
		if !ok {
			return 0, nil, fmt.Errorf("path not found: %s", path)
		}
		current = next
		traversal = append(traversal, current)
	}

	return current, traversal, nil
}

func (fs *FileSystem) ReadInode(inodeID int, n int, offset int) ([]byte, []int, int, error) {
	inode, ok := fs.inodes[inodeID]
	if !ok || inode.Type != InodeFile {
		return nil, nil, offset, fmt.Errorf("inode %d not readable file", inodeID)
	}
	if n < 0 {
		return nil, nil, offset, fmt.Errorf("read length must be non-negative")
	}
	if offset > inode.Size {
		offset = inode.Size
	}
	remaining := inode.Size - offset
	if n > remaining {
		n = remaining
	}
	out := make([]byte, 0, n)
	blocks := []int{}
	for len(out) < n {
		filePos := offset + len(out)
		blockIdx := filePos / fs.blockSize
		blockOff := filePos % fs.blockSize
		if blockIdx >= len(inode.DirectBlocks) {
			break
		}
		blk := inode.DirectBlocks[blockIdx]
		blocks = append(blocks, blk)
		data := fs.blocks[blk]
		need := n - len(out)
		take := fs.blockSize - blockOff
		if take > need {
			take = need
		}
		out = append(out, data[blockOff:blockOff+take]...)
	}
	return out, uniqInts(blocks), offset + len(out), nil
}

func (fs *FileSystem) WriteInode(inodeID int, data []byte, offset int) (int, []int, int) {
	inode := fs.inodes[inodeID]
	if offset > inode.Size {
		offset = inode.Size
	}
	blocks := []int{}
	written := 0
	for written < len(data) {
		filePos := offset + written
		blockIdx := filePos / fs.blockSize
		blockOff := filePos % fs.blockSize
		for blockIdx >= len(inode.DirectBlocks) {
			blk := fs.nextBlock
			fs.nextBlock++
			inode.DirectBlocks = append(inode.DirectBlocks, blk)
			fs.blocks[blk] = make([]byte, fs.blockSize)
		}
		blk := inode.DirectBlocks[blockIdx]
		blocks = append(blocks, blk)
		chunk := fs.blockSize - blockOff
		if chunk > len(data)-written {
			chunk = len(data) - written
		}
		copy(fs.blocks[blk][blockOff:blockOff+chunk], data[written:written+chunk])
		written += chunk
	}
	if offset+written > inode.Size {
		inode.Size = offset + written
	}
	return written, uniqInts(blocks), offset + written
}

func (fs *FileSystem) Invariants() error {
	root, ok := fs.inodes[fs.rootInode]
	if !ok || root.Type != InodeDir {
		return fmt.Errorf("root inode invalid")
	}
	for id, inode := range fs.inodes {
		if inode.Type == InodeDir {
			for name, childID := range inode.Entries {
				if name == "" {
					return fmt.Errorf("empty directory entry name in inode %d", id)
				}
				if _, ok := fs.inodes[childID]; !ok {
					return fmt.Errorf("dangling dir entry inode=%d name=%s child=%d", id, name, childID)
				}
			}
		}
		if inode.Type == InodeFile {
			for _, blk := range inode.DirectBlocks {
				if _, ok := fs.blocks[blk]; !ok {
					return fmt.Errorf("inode %d references missing block %d", id, blk)
				}
			}
		}
	}
	return nil
}

func uniqInts(vals []int) []int {
	if len(vals) == 0 {
		return nil
	}
	m := map[int]bool{}
	out := make([]int, 0, len(vals))
	for _, v := range vals {
		if m[v] {
			continue
		}
		m[v] = true
		out = append(out, v)
	}
	sort.Ints(out)
	return out
}
