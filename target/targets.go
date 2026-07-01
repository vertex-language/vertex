package target

import "runtime"

func DefaultTarget() string {
	goos := runtime.GOOS
	goarch := runtime.GOARCH
	switch goarch {
	case "amd64", "arm64", "riscv64":
	default:
		goarch = "amd64"
	}
	switch goos {
	case "linux", "darwin", "windows":
		return goos + "-" + goarch
	default:
		return "linux-amd64"
	}
}

func SupportedTargets() []string {
	return []string{
		"linux-amd64",
		"linux-arm64",
		"linux-riscv64",
		"darwin-amd64",
		"darwin-arm64",
		"windows-amd64",
		"windows-arm64",
		"freestanding-amd64",
		"freestanding-arm64",
		"freestanding-riscv64",
	}
}