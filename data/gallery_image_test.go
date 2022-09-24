package data

import "testing"

func Test_guessIsImage(t *testing.T) {
	if isImage := guessIsImage("someFile.pdf", 1000); isImage {
		t.Errorf("someFile.pdf shouldn't be considered as an image")
	}

	if isImage := guessIsImage("someFile.jpg", 1000); !isImage {
		t.Errorf("someFile.jpg should be considered as an image")
	}

	if isImage := guessIsImage("oversized.jpg", 31*1048576); isImage {
		t.Errorf("oversized.jpg shouldn't be considered as an image")
	}
}
