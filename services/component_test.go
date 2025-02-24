package services

import "testing"

func TestToSentanceCase(t *testing.T) {
	// Arrange
	word := "cRazyStYlEs"
	expected := "Crazystyles"

	// Act
	result := toSentenceCase(word)

	// Assert
	if result != expected {
		t.Fatalf("Wanted %v, got %v", expected, result)
	}
}
func TestToSentanceCaseWithOneLetter(t *testing.T) {
	// Arrange
	word := "c"
	expected := "C"

	// Act
	result := toSentenceCase(word)

	// Assert
	if result != expected {
		t.Fatalf("Wanted %v, got %v", expected, result)
	}
}
