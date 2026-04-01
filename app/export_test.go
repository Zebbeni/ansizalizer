package app

import (
	"os"
	"path/filepath"
	"testing"
)

func createTestFile(t *testing.T, dir, name string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestBuildExportQueueImageFiles(t *testing.T) {
	src := t.TempDir()
	dest := t.TempDir()

	createTestFile(t, src, "photo.png")
	createTestFile(t, src, "image.jpg")
	createTestFile(t, src, "readme.txt") // should be skipped

	jobs, err := buildExportQueue(src, dest, false)
	if err != nil {
		t.Fatalf("buildExportQueue: %v", err)
	}
	if len(jobs) != 2 {
		t.Errorf("got %d jobs, want 2 (png + jpg)", len(jobs))
	}

	// Check destination paths end with .ansi
	for _, job := range jobs {
		if ext := filepath.Ext(job.destinationPath); ext != ".ansi" {
			t.Errorf("destination %q should have .ansi extension, got %q", job.destinationPath, ext)
		}
	}
}

func TestBuildExportQueueSkipsNonImages(t *testing.T) {
	src := t.TempDir()
	dest := t.TempDir()

	createTestFile(t, src, "readme.txt")
	createTestFile(t, src, "data.csv")
	createTestFile(t, src, "script.go")

	jobs, err := buildExportQueue(src, dest, false)
	if err != nil {
		t.Fatalf("buildExportQueue: %v", err)
	}
	if len(jobs) != 0 {
		t.Errorf("got %d jobs, want 0 (no image files)", len(jobs))
	}
}

func TestBuildExportQueueNoSubDirs(t *testing.T) {
	src := t.TempDir()
	dest := t.TempDir()
	sub := filepath.Join(src, "subdir")
	os.Mkdir(sub, 0755)

	createTestFile(t, src, "top.png")
	createTestFile(t, sub, "nested.png")

	jobs, err := buildExportQueue(src, dest, false)
	if err != nil {
		t.Fatalf("buildExportQueue: %v", err)
	}
	if len(jobs) != 1 {
		t.Errorf("got %d jobs, want 1 (only top-level)", len(jobs))
	}
}

func TestBuildExportQueueWithSubDirs(t *testing.T) {
	src := t.TempDir()
	dest := t.TempDir()
	sub := filepath.Join(src, "subdir")
	os.Mkdir(sub, 0755)

	createTestFile(t, src, "top.png")
	createTestFile(t, sub, "nested.png")
	createTestFile(t, sub, "nested.jpg")

	jobs, err := buildExportQueue(src, dest, true)
	if err != nil {
		t.Fatalf("buildExportQueue: %v", err)
	}
	if len(jobs) != 3 {
		t.Errorf("got %d jobs, want 3 (1 top + 2 nested)", len(jobs))
	}
}

func TestBuildExportQueueEmptyDir(t *testing.T) {
	src := t.TempDir()
	dest := t.TempDir()

	jobs, err := buildExportQueue(src, dest, false)
	if err != nil {
		t.Fatalf("buildExportQueue: %v", err)
	}
	if len(jobs) != 0 {
		t.Errorf("got %d jobs, want 0", len(jobs))
	}
}

func TestBuildExportQueueAllExtensions(t *testing.T) {
	src := t.TempDir()
	dest := t.TempDir()

	for _, ext := range []string{".png", ".jpg", ".jpeg", ".gif"} {
		createTestFile(t, src, "file"+ext)
	}

	jobs, err := buildExportQueue(src, dest, false)
	if err != nil {
		t.Fatalf("buildExportQueue: %v", err)
	}
	if len(jobs) != 4 {
		t.Errorf("got %d jobs, want 4 (all image extensions)", len(jobs))
	}
}

func TestBuildExportQueueDestinationPaths(t *testing.T) {
	src := t.TempDir()
	dest := t.TempDir()

	createTestFile(t, src, "myimage.png")

	jobs, err := buildExportQueue(src, dest, false)
	if err != nil {
		t.Fatalf("buildExportQueue: %v", err)
	}
	if len(jobs) != 1 {
		t.Fatalf("got %d jobs, want 1", len(jobs))
	}

	expected := filepath.Join(dest, "myimage.ansi")
	if jobs[0].destinationPath != expected {
		t.Errorf("destinationPath = %q, want %q", jobs[0].destinationPath, expected)
	}
}
