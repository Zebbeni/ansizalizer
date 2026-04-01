package modal

import (
	"os"
	"path/filepath"
	"testing"
)

func createFile(t *testing.T, dir, name string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestCountImagesBasic(t *testing.T) {
	dir := t.TempDir()
	createFile(t, dir, "a.png")
	createFile(t, dir, "b.jpg")
	createFile(t, dir, "c.txt") // not an image

	if got := countImages(dir, false); got != 2 {
		t.Errorf("countImages() = %d, want 2", got)
	}
}

func TestCountImagesEmpty(t *testing.T) {
	dir := t.TempDir()
	if got := countImages(dir, false); got != 0 {
		t.Errorf("countImages() = %d, want 0", got)
	}
}

func TestCountImagesNoSubDirs(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "sub")
	os.Mkdir(sub, 0755)

	createFile(t, dir, "top.png")
	createFile(t, sub, "nested.png")

	if got := countImages(dir, false); got != 1 {
		t.Errorf("countImages(useSubDirs=false) = %d, want 1", got)
	}
}

func TestCountImagesWithSubDirs(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "sub")
	os.Mkdir(sub, 0755)

	createFile(t, dir, "top.png")
	createFile(t, sub, "nested.png")
	createFile(t, sub, "nested.gif")

	if got := countImages(dir, true); got != 3 {
		t.Errorf("countImages(useSubDirs=true) = %d, want 3", got)
	}
}

func TestCountImagesAllExtensions(t *testing.T) {
	dir := t.TempDir()
	for _, ext := range []string{".png", ".jpg", ".jpeg", ".gif"} {
		createFile(t, dir, "file"+ext)
	}

	if got := countImages(dir, false); got != 4 {
		t.Errorf("countImages() = %d, want 4", got)
	}
}

func TestCountImagesNonexistentDir(t *testing.T) {
	if got := countImages("/nonexistent/path", false); got != 0 {
		t.Errorf("countImages(nonexistent) = %d, want 0", got)
	}
}
