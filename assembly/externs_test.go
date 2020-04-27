package assembly

import (
	"os"
	"testing"
)

func TestParseExterns(t *testing.T) {
	f, err := os.Open("./udon_funcs_data.txt")
	if err != nil {
		t.Errorf("open file: %s", err)
	}
	defer f.Close()
	_, err = ParseExterns(f)
	if err != nil {
		t.Errorf("%v", err)
	}
}
