package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"os"
	"os/exec"
	"os/signal"
	"strings"

	"github.com/denisbrodbeck/machineid"
	"github.com/joho/godotenv"
	"github.com/pterm/pterm"
	"github.com/qeesung/image2ascii/convert"
)

func SH(c string) (string, error) {
	cmd := exec.Command("/bin/sh", "-c", c)
	cmd.Env = os.Environ()
	o, err := cmd.CombinedOutput()
	return string(o), err
}

func SHInDir(c, dir string, envs ...string) (string, error) {
	cmd := exec.Command("/bin/sh", "-c", c)
	cmd.Env = append(os.Environ(), envs...)
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

func K3sBin() string {
	for _, p := range []string{"/usr/bin/k3s", "/usr/local/bin/k3s"} {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}

	return ""
}

func WriteEnv(envFile string, config map[string]string) error {
	content, _ := os.ReadFile(envFile)
	env, _ := godotenv.Unmarshal(string(content))

	for key, val := range config {
		env[key] = val
	}

	return godotenv.Write(env, envFile)
}

func Flavor() string {
	release, _ := godotenv.Read("/etc/os-release")
	v := release["NAME"]
	return strings.ReplaceAll(v, "kairos-", "")
}

func IsOpenRCBased() bool {
	f := Flavor()
	return strings.Contains(f, "alpine")
}

func ShellSTDIN(s, c string) (string, error) {
	cmd := exec.Command("/bin/sh", "-c", c)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = bytes.NewBuffer([]byte(s))
	o, err := cmd.CombinedOutput()
	return string(o), err
}

func SetEnv(env []string) {
	for _, e := range env {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) >= 2 {
			os.Setenv(pair[0], pair[1])
		}
	}
}

func OnSignal(fn func(), sig ...os.Signal) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, sig...)
	go func() {
		<-sigs
		fn()
	}()
}

func Shell() *exec.Cmd {
	cmd := exec.Command("/bin/sh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd
}

func Prompt(t string) (string, error) {
	if t != "" {
		pterm.Info.Println(t)
	}
	answer, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(answer), nil
}

func PrintBanner(d []byte) {
	img, _, _ := image.Decode(bytes.NewReader(d))

	convertOptions := convert.DefaultOptions
	convertOptions.FixedWidth = 100
	convertOptions.FixedHeight = 40

	// Create the image converter
	converter := convert.NewImageConverter()
	fmt.Print(converter.Image2ASCIIString(img, &convertOptions))
}
