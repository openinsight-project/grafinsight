package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitEmails(t *testing.T) {
	testcases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "",
			expected: []string{},
		},
		{
			input:    "open@grafinsight.org",
			expected: []string{"open@grafinsight.org"},
		},
		{
			input:    "open@grafinsight.org;dev@grafinsight.org",
			expected: []string{"open@grafinsight.org", "dev@grafinsight.org"},
		},
		{
			input:    "open@grafinsight.org;dev@grafinsight.org,",
			expected: []string{"open@grafinsight.org", "dev@grafinsight.org"},
		},
		{
			input:    "dev@grafinsight.org,open@grafinsight.org",
			expected: []string{"dev@grafinsight.org", "open@grafinsight.org"},
		},
		{
			input:    "dev@grafinsight.org,open@grafinsight.org,",
			expected: []string{"dev@grafinsight.org", "open@grafinsight.org"},
		},
		{
			input:    "dev@grafinsight.org\nopen@grafinsight.org",
			expected: []string{"dev@grafinsight.org", "open@grafinsight.org"},
		},
		{
			input:    "dev@grafinsight.org\nopen@grafinsight.org\n",
			expected: []string{"dev@grafinsight.org", "open@grafinsight.org"},
		},
	}

	for _, tt := range testcases {
		emails := SplitEmails(tt.input)
		assert.Equal(t, tt.expected, emails)
	}
}
