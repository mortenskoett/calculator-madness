package page_test

import (
	"bytes"
	"testing"
	"viewer/pkg/page"
)

func Test_StatusPage_not_empty(t *testing.T) {

	param := page.StatusParams{
		Title: "Calculation status",
	}

    var b bytes.Buffer

	err := page.Status(&b, param)
	if err != nil {
		t.Fatalf("no error should happen here: %+v", err)
	}

	if b.Len() == 0 {
		t.Errorf("length should be > 0")
	}
}
