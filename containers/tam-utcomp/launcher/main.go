package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"
	"gopkg.in/yaml.v3"
)

var launchConfigFile string
var serverIniFile string

var uccPath = "/opt/ut2004/System/ucc-bin"
var uccExe string
var systemDir string

var defaultMap = "DM-Antalus.ut2"

func init() {
	flag.StringVar(&launchConfigFile, "launch", "launch.yml", "launch config")
	flag.StringVar(&serverIniFile, "ini", "/opt/ut2004/System/UT2004.ini", "UT2004 server ini")
}

func main() {
	flag.Parse()

	// Set default map if provided
	if flag.NArg() > 0 {
		defaultMap = flag.Arg(0)
	}

	config := parseConfig(launchConfigFile)

	ini, err := os.Open(serverIniFile)
	if err != nil {
		panic(err)
	}
	defer ini.Close()

	// Allow overwriting UCC executable (e.g. use 64-bit build instead)
	p, ok := os.LookupEnv("UCC")
	if ok {
		uccPath = p
	}

	systemDir, uccExe = filepath.Split(uccPath)
	tmpIniFile, err := os.CreateTemp("", "UT2004.ini.*")
	if err != nil {
		panic(err)
	}

	defer func() {
		tmpIniFile.Close()
		os.Remove(tmpIniFile.Name())
	}()

	// Create new INI file
	err = config.Transform(tmpIniFile, ini)
	if err != nil {
		panic(err)
	}

	tmpIniFile.Seek(0, io.SeekStart)

	// Replace UT2004.ini in the system directory
	updateSystemINI(tmpIniFile)

	playMap, err := config.EnrichMap(defaultMap)
	if err != nil {
		panic(err)
	}

	cmd := "./" + uccExe
	args := []string{uccExe, "server", playMap, "-nohomedir"}

	fmt.Fprintf(os.Stderr, "Running %s %v\n", cmd, args[1:])

	// Execute UCC relative to system directory
	err = os.Chdir(systemDir)
	if err != nil {
		panic(err)
	}

	err = unix.Exec(cmd, args, nil)
	if err != nil {
		panic(err)
	}

	// Any code after this should be unreachable
}

func parseConfig(path string) Config {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)

	var config Config
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}

	return config
}

func updateSystemINI(r io.Reader) {
	f, err := os.Create(filepath.Join(systemDir, "UT2004.ini"))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = io.Copy(f, r)
	if err != nil {
		panic(err)
	}
}
