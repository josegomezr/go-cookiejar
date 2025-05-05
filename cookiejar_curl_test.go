// _ gotcha! this one I really wrote, look at the mess, of course it's mine.

package cookiejar

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

var timeRef = time.Unix(1746471451, 0)

func TestEntryToCurlLine(t *testing.T) {
	refEntry := entry{
		Domain:   "example.com",
		Expires:  timeRef,
		HostOnly: true,
		HttpOnly: false,
		Name:     "cookiename",
		Path:     "/",
		Secure:   false,
		Value:    "cookievalue",
	}

	t.Run("Host-only cookies", func(t *testing.T) {
		var entry entry
		entry = refEntry

		expected := "example.com	false	/	false	1746471451	cookiename	cookievalue"
		got := entry.toCurlLine()

		if got != expected {
			t.Fatalf("Serialization failed.\nExpected=%s\nGot     =%s", expected, got)
		}
	})

	t.Run("Non-host-only (includes subdomains) cookies", func(t *testing.T) {
		var entry entry
		entry = refEntry
		entry.HostOnly = false

		expected := ".example.com	true	/	false	1746471451	cookiename	cookievalue"
		got := entry.toCurlLine()

		if got != expected {
			t.Fatalf("Serialization failed.\nExpected=%s\nGot     =%s", expected, got)
		}
	})

	t.Run("Secure cookies", func(t *testing.T) {
		var entry entry
		entry = refEntry
		entry.Secure = true

		expected := "example.com	false	/	true	1746471451	cookiename	cookievalue"
		got := entry.toCurlLine()

		if got != expected {
			t.Fatalf("Serialization failed.\nExpected=%s\nGot     =%s", expected, got)
		}
	})

	t.Run("HTTP-only cookies", func(t *testing.T) {
		var entry entry
		entry = refEntry
		entry.HttpOnly = true

		expected := "#HttpOnly_example.com	false	/	false	1746471451	cookiename	cookievalue"
		got := entry.toCurlLine()

		if got != expected {
			t.Fatalf("Serialization failed.\nExpected=%s\nGot     =%s", expected, got)
		}
	})

	t.Run("Secure with subdomains HTTP-only cookies", func(t *testing.T) {
		var entry entry
		entry = refEntry
		entry.HttpOnly = true
		entry.HostOnly = false
		entry.Secure = true

		expected := "#HttpOnly_.example.com	true	/	true	1746471451	cookiename	cookievalue"
		got := entry.toCurlLine()

		if got != expected {
			t.Fatalf("Serialization failed.\nExpected=%s\nGot     =%s", expected, got)
		}
	})
}

func TestJarWriteAsCurl(t *testing.T) {
	buff := strings.Builder{}

	expires := tNow.Add(time.Duration(86400) * time.Second)
	jar := newTestJar()
	jarTest{
		description: "WriteToCurl.",
		fromURL:     "https://www.host.test",
		setCookies: []string{
			"name=value; max-age=86400; secure; httponly",
		},
		content: "name=value",
		queries: []query{},
	}.run(t, jar)

	jar.WriteAsCurl(&buff)
	got := buff.String()[len(curlFileBanner)+1:]
	expected := fmt.Sprintf("#HttpOnly_www.host.test	false	/	true	%d	name	value\n", expires.Unix())

	if got != expected {
		t.Fatalf("Cookiejar serialization failed.\nexpected=%q\ngot     =%q", expected, got)
	}
}

func TestJarReadFromCurl(t *testing.T) {
	expires := time.Now().Add(time.Duration(86400) * time.Second)
	input := fmt.Sprintf("#HttpOnly_www.host.test	false	/	true	%d	name	value\n", expires.Unix())

	jar := newTestJar()
	err := jar.ReadFromCurl(strings.NewReader(input))
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	jarTest{
		description: "WriteToCurl.",
		fromURL:     "https://www.host.test",
		setCookies:  []string{},
		content:     "name=value",
		queries:     []query{},
	}.run(t, jar)
}
