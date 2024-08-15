package browse

import (
	"context"
	"encoding/json"
	"testing"
	"time"
)

const str = ``

func Test_New(t *testing.T) {
	req := &BrowseRequest{}
	err := json.Unmarshal([]byte(str), &req.Cookies)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	page, err := New(req, ctx)
	if err != nil {
		t.Fatal(err)
	}

	defer page.Close()
	defer page.Browser().Close()

	resp, err := req.Hijack(page)
	if err != nil {
		t.Fatal(err)
	}

	err = page.Navigate(req.PageURL)
	if err != nil {
		t.Fatal(err)
	}
	page.WaitLoad()

	for _, bri := range req.Instructions {
		_, err = bri.Act(page)
		if err != nil && bri.Fatal {
			t.Fatal(err)
		}
	}
	time.Sleep(30 * time.Second)

	b, err := json.Marshal(resp)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", b)

}
