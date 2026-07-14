package main

import (
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"

	"engine/internal/version"
)

const crockford32 = "0123456789ABCDEFGHJKMNPQRSTVWZ"

func main() {
	args := filterDoubleDash(os.Args[1:])

	for _, a := range args {
		if a == "-v" || a == "--version" || a == "version" {
			fmt.Println(version.String())
			return
		}
	}

	if len(args) >= 1 && args[0] == "create" {
		title, err := parseCreateArgs(args[1:])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := writeCreate(title); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	fmt.Println("engine: ok")
}

func filterDoubleDash(in []string) []string {
	out := make([]string, 0, len(in))
	for _, a := range in {
		if a != "--" {
			out = append(out, a)
		}
	}
	return out
}

func parseCreateArgs(args []string) (title string, err error) {
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--title":
			if i+1 >= len(args) {
				return "", fmt.Errorf("create: --title requires a value")
			}
			title = args[i+1]
			i++
		default:
			return "", fmt.Errorf("create: unknown argument %q", args[i])
		}
	}
	if title == "" {
		return "", fmt.Errorf("create: --title is required")
	}
	return title, nil
}

func newID26() (string, error) {
	var rnd [26]byte
	if _, err := rand.Read(rnd[:]); err != nil {
		return "", err
	}
	b := make([]byte, 26)
	for i := range b {
		b[i] = crockford32[rnd[i]%32]
	}
	return string(b), nil
}

func writeCreate(title string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	id, err := newID26()
	if err != nil {
		return err
	}
	name := id + ".listello.md"
	body := fmt.Sprintf("---\nid: %s\ntitle: %s\n", id, title)
	return os.WriteFile(filepath.Join(wd, name), []byte(body), 0o644)
}
