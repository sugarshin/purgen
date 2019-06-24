package app

import (
	"strings"
	"testing"
)

func TestIsURI(t *testing.T) {
	if result := isURI("https://sugarshin.net"); result != true {
		t.Errorf("isURI(\"https://sugarshin.net\") => %T", result)
	}
}

func TestGetImageSourcesFromReader(t *testing.T) {
	html := `<html>
		<head>
			<meta charset="utf-8">
			<title>Example</title>
		<body>
			<img src="https://sugarshin.net/images/s.png">
		</body>
	</html>
	`
	r := strings.NewReader(html)
	results, err := getImageSourcesFromReader(r)
	if err != nil {
		t.Error("getImageSourcesFromReader(html)", err)
	}

	if l := len(results); l != 1 {
		t.Errorf("getImageSourcesFromReader(html) length => %d", l)
	}
}
