package accesslog

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func CleanupString(s string) string {
	if s == "" {
		return ""
	}

	switch s[0] {
	case '"':
		n, err := strconv.Unquote(s)
		if err == nil {
			return n
		}
		return s
	case '[':
		n, err := time.Parse("[02/Jan/2006:15:04:05 -0700]", s)
		if err == nil {
			return n.Format(time.RFC3339)
		}
		return s
	}
	return s
}

type Request struct {
	Method string
	URL    *url.URL
	Proto  string
}

func ParseRequest(raw string) (Request, error) {
	s := strings.SplitN(CleanupString(raw), " ", 3)
	if len(s) == 1 || len(s) > 3 {
		return Request{}, fmt.Errorf("url format error %q", raw)
	}

	u, err := url.Parse(s[1])
	if err != nil {
		return Request{}, err
	}

	r := Request{
		Method: s[0],
		URL:    u,
	}

	if len(s) >= 3 {
		r.Proto = s[2]
	}

	return r, nil
}

func (r Request) String() string {
	return fmt.Sprintf("%s %s %s", r.Method, r.URL.String(), r.Proto)
}
