package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// describePortOwner returns a human-readable description of the process
// listening on the given port, or empty string if it can't be determined.
func describePortOwner(port int) string {
	pid := findListenerPID(port)
	if pid == 0 {
		return ""
	}

	name := processName(pid)
	if name == "" {
		return ""
	}

	if strings.HasPrefix(name, "viewmd") {
		dir := processCwd(pid)
		if dir != "" {
			return fmt.Sprintf("A viewmd instance (PID %d) is serving %s", pid, dir)
		}
		return fmt.Sprintf("A viewmd instance (PID %d) is using that port.", pid)
	}

	return fmt.Sprintf("Process %q (PID %d) is using that port.", name, pid)
}

// findListenerPID returns the PID of the process listening on the given port.
func findListenerPID(port int) int {
	out, err := exec.Command("lsof", "-i", fmt.Sprintf(":%d", port), "-sTCP:LISTEN", "-t").Output()
	if err != nil {
		return 0
	}
	line := strings.TrimSpace(strings.Split(string(out), "\n")[0])
	pid, err := strconv.Atoi(line)
	if err != nil {
		return 0
	}
	return pid
}

// processName returns the name of the process with the given PID.
func processName(pid int) string {
	out, err := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "comm=").Output()
	if err != nil {
		return ""
	}
	name := strings.TrimSpace(string(out))
	// ps may return the full path; extract the basename
	if i := strings.LastIndex(name, "/"); i >= 0 {
		name = name[i+1:]
	}
	return name
}

// processCwd returns the working directory of the process with the given PID.
func processCwd(pid int) string {
	out, err := exec.Command("lsof", "-a", "-p", strconv.Itoa(pid), "-d", "cwd", "-Fn").Output()
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(out), "\n") {
		if strings.HasPrefix(line, "n") {
			return line[1:]
		}
	}
	return ""
}
