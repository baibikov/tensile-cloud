package utils

import "testing"

func TestFileType(t *testing.T) {
	tests := []struct {
		Name string
		Want string
	}{
		{
			Name: "file.txt",
			Want: "txt",
		},
		{
			Name: "file.txt.zip",
			Want: "zip",
		},
		{
			Name: "file",
			Want: "file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			actual := FileType(tt.Name)
			if actual != tt.Want {
				t.Errorf("name is not correct actual - %s, want - %s", actual, tt.Want)
			}
		})
	}
}
