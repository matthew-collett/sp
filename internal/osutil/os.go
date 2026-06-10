package osutil

import (
	"encoding/json"
	"os"
	"os/exec"
	"runtime"
)

func OpenBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func OpenApplication(appName string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
		args = []string{"-a", appName}
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", appName + ":"}
	default: // linux
		cmd = appName
	}

	return exec.Command(cmd, args...).Start()
}

func CloseApplication(appName string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "darwin":
		cmd = "osascript"
		args = []string{"-e", "quit app \"" + appName + "\""}
	case "windows":
		cmd = "taskkill"
		args = []string{"/IM", appName + ".exe", "/F"}
	default: // linux
		cmd = "pkill"
		args = []string{"-x", appName}
	}

	return exec.Command(cmd, args...).Start()
}

func WriteFile(data any, path string, perms os.FileMode) error {
	contents, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, contents, perms); err != nil {
		return err
	}

	return nil
}

func ReadFileSilent(path string, res any) {
	data, _ := os.ReadFile(path)
	_ = json.Unmarshal(data, res)
}

func ReadFileSilentString(path string) string {
	data, _ := os.ReadFile(path)
	return string(data)
}
