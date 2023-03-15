package utils

import (
	"fmt"
	"github.com/denisbrodbeck/machineid"
	"github.com/joho/godotenv"
	"os"
	"os/exec"
)

func SH(c string) (string, error) {
	cmd := exec.Command("/bin/sh", "-c", c)
	cmd.Env = os.Environ()
	o, err := cmd.CombinedOutput()
	return string(o), err
}

func SHInDir(c, dir string) (string, error) {
	cmd := exec.Command("/bin/sh", "-c", c)
	cmd.Env = os.Environ()
	cmd.Dir = dir
	o, err := cmd.CombinedOutput()
	return string(o), err
}

func Exists(path string) bool {
	_, err := os.Stat(path)

	return !os.IsNotExist(err)
}

// UUID TODO: move this into a machine submodule
func UUID() string {
	if os.Getenv("UUID") != "" {
		return os.Getenv("UUID")
	}
	id, _ := machineid.ID()
	hostname, _ := os.Hostname()
	return fmt.Sprintf("%s-%s", id, hostname)
}

func OSRelease(key string) (string, error) {
	release, err := godotenv.Read("/etc/os-release")
	if err != nil {
		return "", err
	}
	v, exists := release[key]
	if !exists {
		return "", fmt.Errorf("key not found")
	}
	return v, nil
}

func FindCommand(def string, options []string) string {
	for _, p := range options {
		path, err := exec.LookPath(p)
		if err == nil {
			return path
		}
	}

	// Otherwise return default
	return def
}
