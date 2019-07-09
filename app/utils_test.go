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

func TestGetResourcesFromReader(t *testing.T) {
	html := `<html>
		<head>
			<meta charset="utf-8">
			<title>Example</title>
			<script src="https://sugarshin.net/app.js"></script>
			<script src="https://sugarshin.net/vendor/foo.js"></script>
			<link ref="stylesheet" href="/app.css">
			<link ref="stylesheet" href="/vendor.css">
		<body>
			<img src="https://sugarshin.net/images/s.png" srcset="/example-img@x2.jpg 2x, /example-img@x3.jpg 3x">
			<script src="bar.js"></script>
		</body>
	</html>
	`
	r := strings.NewReader(html)
	results, err := getResourcesFromReader(r)
	if err != nil {
		t.Error("getResourcesFromReader(html)", err)
	}

	if l := len(results); l != 8 {
		t.Errorf("getResourcesFromReader(html) length => %d", l)
	}
}
