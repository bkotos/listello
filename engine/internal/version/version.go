package version

import "runtime/debug"

// String returns a human-friendly build/version identifier.
func String() string {
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return "dev"
}

