package accesslog_test

import (
	"testing"

	"github.com/wzshiming/accesslog"
	"github.com/wzshiming/accesslog/nginx"
)

func TestEntry(t *testing.T) {
	raw := `192.168.0.1 - - [01/Jan/2000:00:00:00 +0800] "GET /example.jpg HTTP/1.0" 200 9999999 "http://192.168.0.1/" "curl/7.15.5" "-"`
	item, err := accesslog.ParseEntry[nginx.DefaultAccessLog]([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	if item.String() != raw {
		t.Errorf("item.String() = %s, want %s", item.String(), raw)
	}
	entry := item.Entry()

	if entry.RemoteAddr != "192.168.0.1" {
		t.Errorf("ip.AsString() = %s, want %s", entry.RemoteAddr, "192.168.0.1")
	}

	if entry.TimeLocal != "[01/Jan/2000:00:00:00 +0800]" {
		t.Errorf("time.AsTime() = %s, want %s", entry.TimeLocal, "[01/Jan/2000:00:00:00 +0800]")
	}

	req, err := accesslog.ParseRequest(entry.Request)
	if err != nil {
		t.Fatal(err)
	}
	if req.Method != "GET" {
		t.Errorf("req.Method = %s, want GET", req.Method)
	}

	t.Logf("%#v", item)
}
