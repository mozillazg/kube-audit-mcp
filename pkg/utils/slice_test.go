package utils

import "testing"

func TestContains(t *testing.T) {
	t.Run("string slice - item exists", func(t *testing.T) {
		slice := []string{"apple", "banana", "cherry"}
		if !Contains(slice, "banana") {
			t.Error("Expected Contains to return true for existing item 'banana'")
		}
	})

	t.Run("string slice - item does not exist", func(t *testing.T) {
		slice := []string{"apple", "banana", "cherry"}
		if Contains(slice, "grape") {
			t.Error("Expected Contains to return false for non-existing item 'grape'")
		}
	})

	t.Run("string slice - empty slice", func(t *testing.T) {
		var slice []string
		if Contains(slice, "anything") {
			t.Error("Expected Contains to return false for empty slice")
		}
	})

	t.Run("string slice - nil slice", func(t *testing.T) {
		var slice []string = nil
		if Contains(slice, "anything") {
			t.Error("Expected Contains to return false for nil slice")
		}
	})

	t.Run("string slice - empty string search", func(t *testing.T) {
		slice := []string{"", "apple", "banana"}
		if !Contains(slice, "") {
			t.Error("Expected Contains to return true for empty string in slice")
		}
	})

	t.Run("int slice - item exists", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5}
		if !Contains(slice, 3) {
			t.Error("Expected Contains to return true for existing item 3")
		}
	})

	t.Run("int slice - item does not exist", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5}
		if Contains(slice, 6) {
			t.Error("Expected Contains to return false for non-existing item 6")
		}
	})

	t.Run("int slice - negative numbers", func(t *testing.T) {
		slice := []int{-5, -3, 0, 2, 10}
		if !Contains(slice, -3) {
			t.Error("Expected Contains to return true for existing negative item -3")
		}
		if Contains(slice, -1) {
			t.Error("Expected Contains to return false for non-existing negative item -1")
		}
	})

	t.Run("int slice - zero value", func(t *testing.T) {
		slice := []int{1, 0, 3}
		if !Contains(slice, 0) {
			t.Error("Expected Contains to return true for zero value")
		}
	})

	t.Run("float64 slice - item exists", func(t *testing.T) {
		slice := []float64{1.1, 2.2, 3.3}
		if !Contains(slice, 2.2) {
			t.Error("Expected Contains to return true for existing float 2.2")
		}
	})

	t.Run("float64 slice - item does not exist", func(t *testing.T) {
		slice := []float64{1.1, 2.2, 3.3}
		if Contains(slice, 4.4) {
			t.Error("Expected Contains to return false for non-existing float 4.4")
		}
	})

	t.Run("bool slice - true exists", func(t *testing.T) {
		slice := []bool{false, true, false}
		if !Contains(slice, true) {
			t.Error("Expected Contains to return true for existing bool true")
		}
	})

	t.Run("bool slice - false exists", func(t *testing.T) {
		slice := []bool{true, true, true}
		if Contains(slice, false) {
			t.Error("Expected Contains to return false for non-existing bool false")
		}
	})

	t.Run("single item slice - match", func(t *testing.T) {
		slice := []string{"only"}
		if !Contains(slice, "only") {
			t.Error("Expected Contains to return true for single matching item")
		}
	})

	t.Run("single item slice - no match", func(t *testing.T) {
		slice := []string{"only"}
		if Contains(slice, "other") {
			t.Error("Expected Contains to return false for single non-matching item")
		}
	})

	t.Run("duplicate items - first occurrence", func(t *testing.T) {
		slice := []int{1, 2, 2, 3, 2}
		if !Contains(slice, 2) {
			t.Error("Expected Contains to return true for duplicate item")
		}
	})

	t.Run("large slice performance", func(t *testing.T) {
		// Create a large slice for performance testing
		slice := make([]int, 10000)
		for i := 0; i < 10000; i++ {
			slice[i] = i
		}

		// Test item at the beginning
		if !Contains(slice, 0) {
			t.Error("Expected Contains to return true for item at beginning of large slice")
		}

		// Test item at the end
		if !Contains(slice, 9999) {
			t.Error("Expected Contains to return true for item at end of large slice")
		}

		// Test non-existing item
		if Contains(slice, 10000) {
			t.Error("Expected Contains to return false for non-existing item in large slice")
		}
	})
}

// Test with custom struct type to ensure generics work properly
type Person struct {
	Name string
	Age  int
}

func TestContainsCustomType(t *testing.T) {
	t.Run("custom struct - item exists", func(t *testing.T) {
		people := []Person{
			{Name: "Alice", Age: 30},
			{Name: "Bob", Age: 25},
			{Name: "Charlie", Age: 35},
		}
		target := Person{Name: "Bob", Age: 25}
		if !Contains(people, target) {
			t.Error("Expected Contains to return true for existing custom struct")
		}
	})

	t.Run("custom struct - item does not exist", func(t *testing.T) {
		people := []Person{
			{Name: "Alice", Age: 30},
			{Name: "Bob", Age: 25},
		}
		target := Person{Name: "David", Age: 40}
		if Contains(people, target) {
			t.Error("Expected Contains to return false for non-existing custom struct")
		}
	})

	t.Run("custom struct - partial match should not work", func(t *testing.T) {
		people := []Person{
			{Name: "Alice", Age: 30},
			{Name: "Bob", Age: 25},
		}
		// Same name but different age
		target := Person{Name: "Bob", Age: 30}
		if Contains(people, target) {
			t.Error("Expected Contains to return false for partial match of custom struct")
		}
	})
}
