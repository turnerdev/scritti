package core

import (
	"strings"
	"testing"
)

func TestParseElement(t *testing.T) {
	t.Run("Test parse indent, style", func(t *testing.T) {
		want := ComponentSourceLine{
			indent: 2,
			style:  "test",
		}
		got := parseElement("\t\ttest")

		if want != got {
			t.Errorf("got %q want %q", got, want)
		}
	})
	t.Run("Test parse indent, tag, style", func(t *testing.T) {
		want := ComponentSourceLine{
			indent: 1,
			tag:    "span",
			style:  "test",
		}
		got := parseElement("\tspan.test")

		if want != got {
			t.Errorf("got %q want %q", got, want)
		}
	})
	t.Run("NewStyle - test basic tree parsing", func(t *testing.T) {
		want := ComponentSourceLine{
			indent: 0,
			tag:    "button",
			style:  "primary",
			text:   "Save",
		}
		got := parseElement(`button.primary "Save"`)

		if want != got {
			t.Errorf("got %q want %q", got, want)
		}
	})
	t.Run("Test parse text", func(t *testing.T) {
		want := ComponentSourceLine{
			indent: 3,
			text:   `hello "world"`, // TODO
		}
		got := parseElement(`   "hello \"world\""`)

		if want != got {
			t.Errorf("got %q want %q", got, want)
		}
	})
}

func TestNewStyle(t *testing.T) {
	t.Run("NewStyle - test basic tree parsing", func(t *testing.T) {
		source := strings.Join([]string{
			"class1",
			"class2",
			"class3",
		}, "\n")

		style, err := NewStyle(source)
		if err != nil {
			t.Error(err)
		}
		if len(style.classes) != 3 {
			t.Errorf("Found %d classes, expected 3", len(style.classes))
		}
	})
}

func TestNewComponent(t *testing.T) {
	t.Run("NewComponent - test basic tree parsing", func(t *testing.T) {
		source := strings.Join([]string{
			"root",
			"\tchild1",
			"\t\tgrandchild1",
			"\t\tgrandchild2",
			"\t\tgrandchild3",
			"\tchild2",
		}, "\n")

		component, err := NewComponent(source)

		if err != nil {

		}
		if len(component.children) != 2 {
			t.Errorf("Found %d children, expected 2", len(component.children))
		}
		if len(component.children[0].children) != 3 {
			t.Errorf("Found %d grandchildren under child1, expected 3", len(component.children))
		}
		if len(component.children[1].children) != 0 {
			t.Errorf("Found %d grandchildren under child2, expected 0", len(component.children))
		}
	})
}

func TestDependencies(t *testing.T) {

	t.Run("Test Component dependencies", func(t *testing.T) {
		component, err := NewComponent("root\n\tchild\n\t\tgrandchild")
		if err != nil {
			t.Error(err)
		}
		want := [2]AssetKey{
			{StyleType, "root"},
			{StyleType, "child"},
		}
		var got [2]AssetKey
		copy(got[:], getDependencyKeys(component))

		if want != got {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("Test Component dependencies", func(t *testing.T) {
		component, err := NewComponent("root\n\tchild\n\t\tgrandchild")
		if err != nil {
			t.Error(err)
		}
		want := [2]AssetKey{
			{StyleType, "child"},
			{StyleType, "grandchild"},
		}
		var got [2]AssetKey
		copy(got[:], getDependencyKeys(component.children[0]))

		if want != got {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
