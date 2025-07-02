package utils

import (
	"context"
	"reflect"
	"testing"
)

func TestCoupler_Couple(t *testing.T) {
	t.Run("正常系: 基本的なカップリング", func(t *testing.T) {
		// Arrange
		coupler := NewCoupler()
		ctx := context.Background()
		input := [][]string{
			{"Alice", "Bob"},
			{"X", "Y", "Z"},
		}

		// Act
		result, err := coupler.Couple(ctx, input)

		// Assert
		if err != nil {
			t.Errorf("期待していないエラーが発生: %v", err)
		}

		// カップリングは最大の集合サイズまで実行される（空文字列含む）
		expectedLength := 3 // max(2, 3) = 3
		if len(result) != expectedLength {
			t.Errorf("結果の長さが期待値と異なります: got %v, want %v", len(result), expectedLength)
		}

		// 各カップリング結果の検証
		validFirstElements := map[string]bool{"Alice": true, "Bob": true, "": true}
		validSecondElements := map[string]bool{"X": true, "Y": true, "Z": true, "": true}

		for i, couple := range result {
			if len(couple) != 2 {
				t.Errorf("カップル %d の要素数が期待値と異なります: got %v, want %v", i, len(couple), 2)
			}

			if !validFirstElements[couple[0]] {
				t.Errorf("カップル %d の第1要素が無効です: got %v", i, couple[0])
			}
			if !validSecondElements[couple[1]] {
				t.Errorf("カップル %d の第2要素が無効です: got %v", i, couple[1])
			}
		}
	})

	t.Run("正常系: 異なる長さの集合", func(t *testing.T) {
		// Arrange
		coupler := NewCoupler()
		ctx := context.Background()
		input := [][]string{
			{"A"},
			{"1", "2", "3", "4"},
		}

		// Act
		result, err := coupler.Couple(ctx, input)

		// Assert
		if err != nil {
			t.Errorf("期待していないエラーが発生: %v", err)
		}

		// 最大の集合サイズまでカップリングが行われる
		expectedLength := 4 // max(1, 4) = 4
		if len(result) != expectedLength {
			t.Errorf("結果の長さが期待値と異なります: got %v, want %v", len(result), expectedLength)
		}

		validFirstElements := map[string]bool{"A": true, "": true}
		validSecondElements := map[string]bool{"1": true, "2": true, "3": true, "4": true, "": true}

		for i, couple := range result {
			if len(couple) != 2 {
				t.Errorf("カップル %d の要素数が期待値と異なります: got %v, want %v", i, len(couple), 2)
			}

			if !validFirstElements[couple[0]] {
				t.Errorf("カップル %d の第1要素が無効です: got %v", i, couple[0])
			}

			if !validSecondElements[couple[1]] {
				t.Errorf("カップル %d の第2要素が無効です: got %v", i, couple[1])
			}
		}
	})

	t.Run("正常系: 3つの集合", func(t *testing.T) {
		// Arrange
		coupler := NewCoupler()
		ctx := context.Background()
		input := [][]string{
			{"Red", "Blue"},
			{"Circle", "Square"},
			{"Big", "Small"},
		}

		// Act
		result, err := coupler.Couple(ctx, input)

		// Assert
		if err != nil {
			t.Errorf("期待していないエラーが発生: %v", err)
		}

		if len(result) != 2 { // 各集合に2要素ずつあるので2回のカップリング
			t.Errorf("結果の長さが期待値と異なります: got %v, want %v", len(result), 2)
		}

		for i, couple := range result {
			if len(couple) != 3 {
				t.Errorf("カップル %d の要素数が期待値と異なります: got %v, want %v", i, len(couple), 3)
			}

			// 各要素が有効な値であることを確認
			validColors := map[string]bool{"Red": true, "Blue": true}
			validShapes := map[string]bool{"Circle": true, "Square": true}
			validSizes := map[string]bool{"Big": true, "Small": true}

			if !validColors[couple[0]] {
				t.Errorf("カップル %d の色が無効です: got %v", i, couple[0])
			}
			if !validShapes[couple[1]] {
				t.Errorf("カップル %d の形が無効です: got %v", i, couple[1])
			}
			if !validSizes[couple[2]] {
				t.Errorf("カップル %d のサイズが無効です: got %v", i, couple[2])
			}
		}
	})

	t.Run("異常系: 空の集合配列", func(t *testing.T) {
		// Arrange
		coupler := NewCoupler()
		ctx := context.Background()
		input := [][]string{}

		// Act
		result, err := coupler.Couple(ctx, input)

		// Assert
		if err == nil {
			t.Error("エラーが期待されていましたが、nilが返されました")
		}
		if result != nil {
			t.Errorf("nilが期待されていましたが、%vが返されました", result)
		}
	})

	t.Run("正常系: 元の配列が変更されないことを確認", func(t *testing.T) {
		// Arrange
		coupler := NewCoupler()
		ctx := context.Background()
		input := [][]string{
			{"A", "B"},
			{"1", "2"},
		}
		originalInput := make([][]string, len(input))
		for i, set := range input {
			originalInput[i] = make([]string, len(set))
			copy(originalInput[i], set)
		}

		// Act
		_, err := coupler.Couple(ctx, input)

		// Assert
		if err != nil {
			t.Errorf("期待していないエラーが発生: %v", err)
		}

		if !reflect.DeepEqual(input, originalInput) {
			t.Error("元の配列が変更されてしまいました")
		}
	})
}

func TestParseItemSets(t *testing.T) {
	t.Run("正常系: 基本的なパース", func(t *testing.T) {
		// Arrange
		input := []string{
			"apple,banana,cherry",
			"red,green,blue",
		}
		separator := ","

		// Act
		result := ParseItemSets(input, separator)

		// Assert
		expected := [][]string{
			{"apple", "banana", "cherry"},
			{"red", "green", "blue"},
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("結果が期待値と異なります: got %v, want %v", result, expected)
		}
	})

	t.Run("正常系: 空白の除去", func(t *testing.T) {
		// Arrange
		input := []string{
			" apple , banana , cherry ",
			"red,  green,blue  ",
		}
		separator := ","

		// Act
		result := ParseItemSets(input, separator)

		// Assert
		expected := [][]string{
			{"apple", "banana", "cherry"},
			{"red", "green", "blue"},
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("結果が期待値と異なります: got %v, want %v", result, expected)
		}
	})

	t.Run("正常系: 空行の除去", func(t *testing.T) {
		// Arrange
		input := []string{
			"apple,banana",
			"",
			"   ",
			"red,green",
		}
		separator := ","

		// Act
		result := ParseItemSets(input, separator)

		// Assert
		expected := [][]string{
			{"apple", "banana"},
			{"red", "green"},
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("結果が期待値と異なります: got %v, want %v", result, expected)
		}
	})

	t.Run("正常系: 異なる区切り文字", func(t *testing.T) {
		// Arrange
		input := []string{
			"apple|banana|cherry",
			"red|green",
		}
		separator := "\\|" // 正規表現でエスケープが必要

		// Act
		result := ParseItemSets(input, separator)

		// Assert
		expected := [][]string{
			{"apple", "banana", "cherry"},
			{"red", "green"},
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("結果が期待値と異なります: got %v, want %v", result, expected)
		}
	})
}

func TestCoupling_BackwardCompatibility(t *testing.T) {
	t.Run("正常系: 後方互換性の確認", func(t *testing.T) {
		// Arrange
		input := [][]string{
			{"A", "B"},
			{"1", "2"},
		}
		coupling := [][]string{}

		// Act
		result := Coupling(input, coupling)

		// Assert
		if len(result) != 2 {
			t.Errorf("結果の長さが期待値と異なります: got %v, want %v", len(result), 2)
		}

		for i, couple := range result {
			if len(couple) != 2 {
				t.Errorf("カップル %d の要素数が期待値と異なります: got %v, want %v", i, len(couple), 2)
			}
		}
	})
}

func TestGetItemSets_BackwardCompatibility(t *testing.T) {
	t.Run("正常系: 後方互換性の確認", func(t *testing.T) {
		// Arrange
		input := []string{"a,b,c", "1,2,3"}
		splitter := ","

		// Act
		result := GetItemSets(input, splitter)

		// Assert
		expected := [][]string{
			{"a", "b", "c"},
			{"1", "2", "3"},
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("結果が期待値と異なります: got %v, want %v", result, expected)
		}
	})
}