package repositories

import (
	"fmt"
	"os"
	"slices"
	"testing"
)

const TEST_DIR string = "testing"

// Util function for removing a directory and file
func delete(file string, t *testing.T) {
	if err := os.RemoveAll(file); err != nil {
		t.Fatalf("Failed to remove directory/file: %v", err.Error())
	}
}

// Util function for creating a directory and file
func mkdir(dir string, t *testing.T) {
	if err := os.Mkdir(dir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err.Error())
	}
}

// Util function for writing a string to a file
func createFile(filename, content string, t *testing.T) {
	file, err := os.Create(filename)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err.Error())
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}
}

func TestHasDirectoryOrFileReturnsTrueIfDirectoryExists(t *testing.T) {
	// Arrange
	mkdir(TEST_DIR, t)
	defer delete(TEST_DIR, t)
	filesystem, _ := NewFilesystemRepository()

	// Act
	hasDir, err := filesystem.HasDirectoryOrFile(TEST_DIR)
	if err != nil {
		t.Fatalf("Error occured: %v", err.Error())
	}

	// Assert
	if !hasDir {
		t.Fatal("Failed to detect directory")
	}
}

func TestHasDirectoryOrFileReturnsFalseIfDirectoryDoesntExists(t *testing.T) {
	// Arrange
	filesystem, _ := NewFilesystemRepository()

	// Act
	hasDir, err := filesystem.HasDirectoryOrFile("non-existent")
	if err != nil {
		t.Fatalf("Error occured: %v", err.Error())
	}

	// Assert
	if hasDir {
		t.Fatal("Detected directory that was not present")
	}
}

func TestReadDirRecursive(t *testing.T) {
	// Arrange
	mkdir("testdir", t)
	defer delete("testdir", t)
	createFile("testdir/a.txt", "mock-content", t)
	createFile("testdir/test.txt", "mock-content", t)
	mkdir("testdir/nested", t)
	createFile("testdir/nested/test.txt", "mock-content", t)
	mkdir("testdir/nested/verynested", t)
	createFile("testdir/nested/verynested/test.txt", "mock-content", t)
	filesystem, _ := NewFilesystemRepository()

	// Act
	files, err := filesystem.ReadDirRecursive("testdir")
	if err != nil {
		t.Fatalf("Failed to read dir: %v", err.Error())
	}

	// Assert
	expected := []string{
		"testdir/a.txt",
		"testdir/nested/test.txt",
		"testdir/nested/verynested/test.txt",
		"testdir/test.txt",
	}

	if !slices.Equal(expected, files) {
		t.Fatalf("Expected %v, got %v", expected, files)
	}
}

func TestReadFileWorks(t *testing.T) {
	// Arrange
	filename := "test.txt"
	filecontent := "mock-content"
	createFile(filename, filecontent, t)
	defer delete(filename, t)
	filesystem, _ := NewFilesystemRepository()

	// Act
	content, err := filesystem.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	// Assert
	if content != filecontent {
		t.Fatalf("File had incorrect contents: %v", content)
	}
}

func TestCreateDirectoryWorks(t *testing.T) {
	// Arrange
	filesystem, _ := NewFilesystemRepository()

	// Act
	if err := filesystem.CreateDirectory(TEST_DIR); err != nil {
		t.Fatalf("Failed to create directory: %v", err.Error())
	}
	defer delete(TEST_DIR, t)

	// Assert
	if hasDir, err := filesystem.HasDirectoryOrFile(TEST_DIR); err != nil || !hasDir {
		t.Fatalf("Failed to create directory")
	}
}

func TestCreateDirCreatesRequiredParentDirectories(t *testing.T) {
	// Arrange
	SECOND_DIR := fmt.Sprintf("%v/test", TEST_DIR)
	filesystem, _ := NewFilesystemRepository()

	// Act
	if err := filesystem.CreateDirectory(SECOND_DIR); err != nil {
		t.Fatalf("Failed to create directory: %v", err.Error())
	}
	defer delete(TEST_DIR, t)

	// Assert
	if hasDir, err := filesystem.HasDirectoryOrFile(SECOND_DIR); err != nil || !hasDir {
		t.Fatalf("Failed to create directory")
	}
}

func TestCreateDirectoryReturnsExpectedErrorIfFileExists(t *testing.T) {
	// Arrange
	filesystem, _ := NewFilesystemRepository()

	// Act
	if err := filesystem.CreateDirectory(TEST_DIR); err != nil {
		t.Fatalf("Failed to create directory: %v", err.Error())
	}
	err := filesystem.CreateDirectory(TEST_DIR)
	defer delete(TEST_DIR, t)

	// Assert
	if err.Error() != "directory with that name already exists" {
		t.Fatalf("error did not return with expected content: %v", err.Error())
	}
}
