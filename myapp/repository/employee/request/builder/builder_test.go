package builder

import "testing"

func TestBuild(t *testing.T) {
	actual := NewBuilder().Build("sales", "10", "20", "active")

	if actual.Keyword != "sales" ||
		actual.DepartmentID != "10" ||
		actual.PositionID != "20" ||
		actual.EmploymentStatus != "active" {
		t.Fatalf("Build() = %#v", actual)
	}
}
