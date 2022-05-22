package test

import (
	"os"
	"testing"
)

func TestOs(t *testing.T) {
	dirs, err := os.ReadDir("../")
	if err != nil {
		return
	}
	for _, dir := range dirs {
		if dir.IsDir() {
		
		} else {

		}
	}
}
