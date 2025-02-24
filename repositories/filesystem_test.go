package repositories

import (
	"os"
	"testing"
)

const TEST_DIR string = "testing"

// / Util function for removing a directory and file
func removeFile(dir string, t *testing.T) {
	if err := os.Remove(dir); err != nil {
		t.Fatalf("Failed to remove directory after test: %v", err.Error())
	}
}

func TestHasDirectoryReturnsTrueIfDirectoryExists(t *testing.T) {
	// Arrange
	if err := os.Mkdir(TEST_DIR, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err.Error())
	}
	service := NewFilesystemRepository()

	// Act
	hasDir, err := service.HasDirectory(TEST_DIR)
	if err != nil {
		t.Fatalf("Error occured: %v", err.Error())
	}

	// Assert
	removeFile(TEST_DIR, t)
	if !hasDir {
		t.Fatal("Failed to detect directory")
	}
}

func TestHasDirectoryReturnsFalseIfDirectoryDoesntExists(t *testing.T) {
	// Arrange
	service := NewFilesystemRepository()

	// Act
	hasDir, err := service.HasDirectory("non-existent")
	if err != nil {
		t.Fatalf("Error occured: %v", err.Error())
	}

	// Assert
	if hasDir {
		t.Fatal("Detected directory that was not present")
	}
}

func TestCreateDirectoryWorks(t *testing.T) {
	// Arrange
	service := NewFilesystemRepository()

	// Act
	if err := service.CreateDirectory(TEST_DIR); err != nil {
		t.Fatalf("Failed to create directory: %v", err.Error())
	}

	// Assert
	if hasDir, err := service.HasDirectory(TEST_DIR); err != nil || !hasDir {
		t.Fatalf("Failed to create directory")
	}
	removeFile(TEST_DIR, t)
}

func TestCreateDirectoryReturnsExpectedErrorIfFileExists(t *testing.T) {
	// Arrange
	service := NewFilesystemRepository()

	// Act
	if err := service.CreateDirectory(TEST_DIR); err != nil {
		t.Fatalf("Failed to create directory: %v", err.Error())
	}
	err := service.CreateDirectory(TEST_DIR)

	// Assert
	removeFile(TEST_DIR, t)
	if err.Error() != "file with that name already exists" {
		t.Fatalf("error did not return with expected content: %v", err.Error())
	}
}
