package services

import "testing"

func TestIsValidSqliteType(t *testing.T) {
    tests := []struct {
        name     string
		input    string
        expected bool
    }{
		// General Formatting
        {
            name:     "accepted name in caps",
            input:    "INT",
            expected: true,
        },
        {
			// The function should be parsed formated types
            name:     "should reject name in lowercase",
            input:    "int",
            expected: false,
        },
        {
            name:     "rejects name with no valid part",
            input:    "asdfg",
            expected: false,
        },
        {
            name:     "rejects name with a valid part",
            input:    "DOUBLEf",
            expected: false,
        },
		// Testing Decimal Type
        {
            name:     "Decimal type accepts valid input",
            input:    "DECIMAL(123, 124)",
            expected: true,
        },
        {
            name:     "Decimal type doesn't accept only one input",
            input:    "DECIMAL(123)",
            expected: false,
        },
        {
            name:     "Decimal type doesn't accept only no inputs",
            input:    "DECIMAL()",
            expected: false,
        },
        {
            name:     "Decimal type doesn't accept blank inputs",
            input:    "DECIMAL(,)",
            expected: false,
        },
		// Testing character types
        {
            name:     "CHARACTER accepts a count input",
            input:    "CHARACTER(143)",
            expected: true,
        },
        {
            name:     "CHARACTER accepts only one input",
            input:    "CHARACTER(123, 123)",
            expected: false,
        },
        {
            name:     "CHARACTER rejects no input",
            input:    "CHARACTER()",
            expected: false,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            result := isValidSqliteType(tc.input)
            if result != tc.expected {
                t.Errorf("expected %v, got %v", tc.expected, result)
            }
		})
	}
}

