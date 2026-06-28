package builder

import "testing"

func TestBuild(t *testing.T) {
	actual := NewBuilder().Build("name", "desc")

	if actual.OrderBy != "name" || actual.OrderDirection != "desc" {
		t.Fatalf("Build() = %#v", actual)
	}
}
