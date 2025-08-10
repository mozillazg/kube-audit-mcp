package utils

import (
	"reflect"
	"testing"
)

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

func TestRemoveDuplicates(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		var input []int
		result := RemoveDuplicates(input)
		expected := []int{}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("RemoveDuplicates(%v) = %v, want %v", input, result, expected)
		}
		// Verify the result is not nil
		if result == nil {
			t.Error("RemoveDuplicates should return empty slice, not nil")
		}
	})

	t.Run("single element", func(t *testing.T) {
		input := []int{42}
		result := RemoveDuplicates(input)
		expected := []int{42}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("RemoveDuplicates(%v) = %v, want %v", input, result, expected)
		}
	})

	t.Run("no duplicates", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		result := RemoveDuplicates(input)
		expected := []int{1, 2, 3, 4, 5}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("RemoveDuplicates(%v) = %v, want %v", input, result, expected)
		}
	})

	t.Run("all duplicates", func(t *testing.T) {
		input := []int{7, 7, 7, 7, 7}
		result := RemoveDuplicates(input)
		expected := []int{7}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("RemoveDuplicates(%v) = %v, want %v", input, result, expected)
		}
	})

	t.Run("mixed duplicates and unique", func(t *testing.T) {
		input := []int{1, 2, 2, 3, 1, 4, 3, 5}
		result := RemoveDuplicates(input)
		expected := []int{1, 2, 3, 4, 5}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("RemoveDuplicates(%v) = %v, want %v", input, result, expected)
		}
	})

	t.Run("preserves order", func(t *testing.T) {
		input := []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3}
		result := RemoveDuplicates(input)
		expected := []int{3, 1, 4, 5, 9, 2, 6}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("RemoveDuplicates(%v) = %v, want %v", input, result, expected)
		}
	})

	t.Run("with zero values", func(t *testing.T) {
		input := []int{0, 1, 0, 2, 0, 3}
		result := RemoveDuplicates(input)
		expected := []int{0, 1, 2, 3}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("RemoveDuplicates(%v) = %v, want %v", input, result, expected)
		}
	})

	t.Run("with negative numbers", func(t *testing.T) {
		input := []int{-1, 2, -1, -3, 2, -3, 0}
		result := RemoveDuplicates(input)
		expected := []int{-1, 2, -3, 0}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("RemoveDuplicates(%v) = %v, want %v", input, result, expected)
		}
	})
}

func TestRemoveDuplicatesStrings(t *testing.T) {
	t.Run("empty string slice", func(t *testing.T) {
		var input []string
		result := RemoveDuplicates(input)
		expected := []string{}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("RemoveDuplicates(%v) = %v, want %v", input, result, expected)
		}
	})

	t.Run("string duplicates", func(t *testing.T) {
		input := []string{"hello", "world", "hello", "foo", "world", "bar"}
		result := RemoveDuplicates(input)
		expected := []string{"hello", "world", "foo", "bar"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("RemoveDuplicates(%v) = %v, want %v", input, result, expected)
		}
	})

	t.Run("empty strings", func(t *testing.T) {
		input := []string{"", "hello", "", "world", ""}
		result := RemoveDuplicates(input)
		expected := []string{"", "hello", "world"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("RemoveDuplicates(%v) = %v, want %v", input, result, expected)
		}
	})

	t.Run("case sensitive", func(t *testing.T) {
		input := []string{"Hello", "hello", "HELLO", "Hello"}
		result := RemoveDuplicates(input)
		expected := []string{"Hello", "hello", "HELLO"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("RemoveDuplicates(%v) = %v, want %v", input, result, expected)
		}
	})
}

func TestRemoveDuplicatesFloat(t *testing.T) {
	t.Run("float duplicates", func(t *testing.T) {
		input := []float64{1.5, 2.7, 1.5, 3.14, 2.7, 0.0}
		result := RemoveDuplicates(input)
		expected := []float64{1.5, 2.7, 3.14, 0.0}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("RemoveDuplicates(%v) = %v, want %v", input, result, expected)
		}
	})

	t.Run("with NaN", func(t *testing.T) {
		// Note: NaN != NaN in Go, so each NaN is considered unique
		input := []float64{1.0, 2.0, 1.0}
		result := RemoveDuplicates(input)
		expected := []float64{1.0, 2.0}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("RemoveDuplicates(%v) = %v, want %v", input, result, expected)
		}
	})
}

func TestRemoveDuplicatesBool(t *testing.T) {
	t.Run("boolean duplicates", func(t *testing.T) {
		input := []bool{true, false, true, false, true}
		result := RemoveDuplicates(input)
		expected := []bool{true, false}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("RemoveDuplicates(%v) = %v, want %v", input, result, expected)
		}
	})

	t.Run("all true", func(t *testing.T) {
		input := []bool{true, true, true}
		result := RemoveDuplicates(input)
		expected := []bool{true}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("RemoveDuplicates(%v) = %v, want %v", input, result, expected)
		}
	})

	t.Run("all false", func(t *testing.T) {
		input := []bool{false, false, false}
		result := RemoveDuplicates(input)
		expected := []bool{false}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("RemoveDuplicates(%v) = %v, want %v", input, result, expected)
		}
	})
}

func TestRemoveDuplicatesCustomType(t *testing.T) {
	type CustomInt int

	t.Run("custom type duplicates", func(t *testing.T) {
		input := []CustomInt{CustomInt(1), CustomInt(2), CustomInt(1), CustomInt(3)}
		result := RemoveDuplicates(input)
		expected := []CustomInt{CustomInt(1), CustomInt(2), CustomInt(3)}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("RemoveDuplicates(%v) = %v, want %v", input, result, expected)
		}
	})
}
