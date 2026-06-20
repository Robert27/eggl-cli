//go:build linux

package kill

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const tcpListenState = "0A"

type linuxFinder struct{}

func (linuxFinder) FindListeners(ctx context.Context, port int) ([]Process, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	inodes, err := listenerInodes(port)
	if err != nil {
		return nil, err
	}
	if len(inodes) == 0 {
		return nil, nil
	}

	byPID := make(map[int]Process)
	for _, inode := range inodes {
		pids, err := pidsForInode(inode)
		if err != nil {
			return nil, err
		}
		for _, pid := range pids {
			if _, ok := byPID[pid]; ok {
				continue
			}
			byPID[pid] = Process{
				PID:  pid,
				Name: processName(pid),
			}
		}
	}

	processes := make([]Process, 0, len(byPID))
	for _, proc := range byPID {
		processes = append(processes, proc)
	}
	return processes, nil
}

func listenerInodes(port int) ([]uint64, error) {
	var inodes []uint64
	for _, path := range []string{"/proc/net/tcp", "/proc/net/tcp6"} {
		data, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, fmt.Errorf("read %s: %w", path, err)
		}
		found, err := parseProcNetTCP(string(data), port)
		if err != nil {
			return nil, err
		}
		inodes = append(inodes, found...)
	}
	return inodes, nil
}

func parseProcNetTCP(data string, port int) ([]uint64, error) {
	targetPort := fmt.Sprintf("%04X", port)
	lines := strings.Split(data, "\n")
	var inodes []uint64

	for i, line := range lines {
		if i == 0 || strings.TrimSpace(line) == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 10 {
			continue
		}
		if fields[3] != tcpListenState {
			continue
		}

		local := fields[1]
		colon := strings.LastIndex(local, ":")
		if colon < 0 {
			continue
		}
		if !strings.EqualFold(local[colon+1:], targetPort) {
			continue
		}

		inode, err := strconv.ParseUint(fields[9], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse inode from %q: %w", line, err)
		}
		inodes = append(inodes, inode)
	}

	return inodes, nil
}

func pidsForInode(inode uint64) ([]int, error) {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil, fmt.Errorf("read /proc: %w", err)
	}

	target := fmt.Sprintf("socket:[%d]", inode)
	var pids []int
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}

		fdDir := filepath.Join("/proc", entry.Name(), "fd")
		fds, err := os.ReadDir(fdDir)
		if err != nil {
			continue
		}

		for _, fd := range fds {
			link, err := os.Readlink(filepath.Join(fdDir, fd.Name()))
			if err != nil {
				continue
			}
			if link == target {
				pids = append(pids, pid)
				break
			}
		}
	}

	return pids, nil
}
