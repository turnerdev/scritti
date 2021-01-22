package core

import "testing"

func TestComponent(t *testing.T) {

	t.Run("Test component Append", func(t *testing.T) {
		root := Component{name: "root"}
		child1 := root.Append("child1")
		child2 := root.Append("child2").Append("child2.1")

		if len(root.children) != 2 {
			t.Errorf("Expected 2 children, found %d", len(root.children))
		}

		if child1.parent != &root {
			t.Error("Root parent incorrectly linked")
		}

		if child2 != root.children[1].children[0] || child2.parent.parent != &root {
			t.Error("Grandchild incorrectly linked")
		}

	})

}

func TestParseComponent(t *testing.T) {
	source := "root\n" +
		"  child1\n" +
		"    child1.1\n" +
		"  child2\n"
	result := ParseComponent(source)

	if len(result.children) != 2 {
		t.Errorf("Expected 2 children, found %d", len(result.children))
	}

	if result.name != "root" {
		t.Errorf("Expected root, found %q", result.name)
	}

	t.Run("Test component depth", func(t *testing.T) {
		if len(result.children[0].children) != 1 {
			t.Errorf("Expected 1st child to have no children")
		}
		if len(result.children[1].children) != 0 {
			t.Errorf("Expected 1st child to have no children")
		}
		if len(result.children[0].children[0].children) != 0 {
			t.Errorf("Expected no great-grandchild")
		}
	})
}
