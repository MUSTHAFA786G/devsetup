package cloner

import "testing"

func TestNormalizeURL(t *testing.T) {
	cases := []struct{ in, want string }{
		{"https://github.com/user/repo", "https://github.com/user/repo.git"},
		{"https://github.com/user/repo.git", "https://github.com/user/repo.git"},
		{"https://github.com/user/repo/", "https://github.com/user/repo.git"},
		{"git@github.com:user/repo.git", "https://github.com/user/repo.git"},
	}
	for _, c := range cases {
		got := normalizeURL(c.in)
		if got != c.want {
			t.Errorf("normalizeURL(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestExtractRepoName(t *testing.T) {
	cases := []struct{ url, want string }{
		{"https://github.com/user/myproject.git", "myproject"},
		{"https://github.com/org/my-repo.git", "my-repo"},
	}
	for _, c := range cases {
		got := extractRepoName(c.url)
		if got != c.want {
			t.Errorf("extractRepoName(%q) = %q, want %q", c.url, got, c.want)
		}
	}
}
