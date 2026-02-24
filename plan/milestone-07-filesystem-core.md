# Milestone 07: Filesystem Core

## Goal

Provide a persistence model suitable for OS teaching flows.

## Scope

- Inode model
- Directory lookup and path resolution
- Block mapping

## Deliverables

- Path traversal animation
- Basic file operations integrated with syscall path

## Exit Criteria

- Filesystem invariants pass consistently

## Key Risks

- Scope creep into journaling and advanced FS features

## Suggested First Tasks

1. Define inode and directory entry structures.
2. Implement path resolver and open file table hooks.
3. Add block mapping read/write flow.
4. Add invariant tests for inode and directory consistency.
