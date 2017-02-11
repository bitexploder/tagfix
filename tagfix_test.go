package main

import (
	"os/user"
	"testing"
)

func TestExpandHomedir(t *testing.T) {
	usr, _ := user.Current()
	dir := usr.HomeDir + "/"

	paths := []string{"~/", "~/Downloads/afile.txt", "", "~", "____~/", "/home/someuser/a/regular/path"}

	for _, p := range paths {
		switch p {
		case "~/":
			expanded := expandHome(p)
			if expanded != dir {
				t.Fatalf("~/ expansion error")
			}
		case "~/Downloads/afile.txt":
			expanded := expandHome(p)
			expected := dir + "Downloads/afile.txt"
			if expanded != expected {
				t.Fatalf("expansion error: %s != %s", expanded, expected)
			}
		case "":
			expanded := expandHome(p)
			expected := ""
			if expanded != expected {
				t.Fatalf("expansion error: expected empty string, got %s", expanded)
			}
		case "~":
			expanded := expandHome(p)
			expected := "~"
			if expanded != expected {
				t.Fatalf("expansion error: don't expand naked tilde, got %s", expanded)
			}
		case "____~/":
			expanded := expandHome(p)
			expected := "____~/"
			if expanded != expected {
				t.Fatalf("expansion error: don't expand in the middle of a string", expanded)
			}
		case "/home/someuser/a/regular/path":
			expanded := expandHome(p)
			expected := "/home/someuser/a/regular/path"
			if expanded != expected {
				t.Fatalf("expansion error: don't expand without ~/ at start")
			}
		}
	}
}
