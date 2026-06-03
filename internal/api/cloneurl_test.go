package api

import "testing"

func Test__cloneURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		publicURL string
		repoID    string
		want      string
	}{
		{
			publicURL: "http://localhost:8080/git",
			repoID:    "orgs/111/canvases/222",
			want:      "http://localhost:8080/git/orgs/111/canvases/222.git",
		},
		{
			publicURL: "http://localhost:8080/git/",
			repoID:    "acme/widgets",
			want:      "http://localhost:8080/git/acme/widgets.git",
		},
		{
			publicURL: "",
			repoID:    "acme/widgets",
			want:      "",
		},
	}

	for _, tc := range tests {
		if got := cloneURL(tc.publicURL, tc.repoID); got != tc.want {
			t.Fatalf("cloneURL(%q, %q) = %q, want %q", tc.publicURL, tc.repoID, got, tc.want)
		}
	}
}
