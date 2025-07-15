package lz77

import (
	"bytes"
	"testing"
)

func TestLZ77Compressor_Name(t *testing.T) {
	compressor := NewCompressor()
	expected := "LZ77"
	if compressor.Name() != expected {
		t.Errorf("Expected %s, got %s", expected, compressor.Name())
	}
}

func TestLZ77Compressor_EmptyData(t *testing.T) {
	compressor := NewCompressor()

	// 空のデータをテスト
	compressed, err := compressor.Compress([]byte{})
	if err != nil {
		t.Fatalf("Compress failed: %v", err)
	}

	decompressed, err := compressor.Decompress(compressed)
	if err != nil {
		t.Fatalf("Decompress failed: %v", err)
	}

	if len(decompressed) != 0 {
		t.Errorf("Expected empty result, got %v", decompressed)
	}
}

func TestLZ77Compressor_SingleCharacter(t *testing.T) {
	compressor := NewCompressor()
	original := []byte("a")

	compressed, err := compressor.Compress(original)
	if err != nil {
		t.Fatalf("Compress failed: %v", err)
	}

	decompressed, err := compressor.Decompress(compressed)
	if err != nil {
		t.Fatalf("Decompress failed: %v", err)
	}

	if !bytes.Equal(original, decompressed) {
		t.Errorf("Original: %s, Decompressed: %s", original, decompressed)
	}
}

func TestLZ77Compressor_RepeatingPattern(t *testing.T) {
	compressor := NewCompressor()
	original := []byte("aaaaaaaaaa") // 10個の'a'

	compressed, err := compressor.Compress(original)
	if err != nil {
		t.Fatalf("Compress failed: %v", err)
	}

	decompressed, err := compressor.Decompress(compressed)
	if err != nil {
		t.Fatalf("Decompress failed: %v", err)
	}

	if !bytes.Equal(original, decompressed) {
		t.Errorf("Original: %s, Decompressed: %s", original, decompressed)
	}
}

func TestLZ77Compressor_BasicText(t *testing.T) {
	compressor := NewCompressor()
	testCases := [][]byte{
		[]byte("hello world"),
		[]byte("abcdefghijklmnopqrstuvwxyz"),
		[]byte("The quick brown fox jumps over the lazy dog"),
		[]byte("aaaaabbbbccccdddd"),
	}

	for i, original := range testCases {
		compressed, err := compressor.Compress(original)
		if err != nil {
			t.Fatalf("Test case %d: Compress failed: %v", i, err)
		}

		decompressed, err := compressor.Decompress(compressed)
		if err != nil {
			t.Fatalf("Test case %d: Decompress failed: %v", i, err)
		}

		if !bytes.Equal(original, decompressed) {
			t.Errorf("Test case %d: Original: %s, Decompressed: %s", i, original, decompressed)
		}
	}
}

func TestLZ77Compressor_LongRepeatingText(t *testing.T) {
	compressor := NewCompressor()

	// 長い反復パターンをテスト
	original := []byte("abcdabcdabcdabcdabcdabcdabcdabcd")

	compressed, err := compressor.Compress(original)
	if err != nil {
		t.Fatalf("Compress failed: %v", err)
	}

	decompressed, err := compressor.Decompress(compressed)
	if err != nil {
		t.Fatalf("Decompress failed: %v", err)
	}

	if !bytes.Equal(original, decompressed) {
		t.Errorf("Original: %s, Decompressed: %s", original, decompressed)
	}

	// LZ77では反復パターンが圧縮されるはずなので、圧縮後のサイズをチェック
	if len(compressed) >= len(original) {
		t.Logf("Warning: Compressed size (%d) is not smaller than original (%d)", len(compressed), len(original))
	}
}

func TestLZ77Compressor_WindowSizeEdgeCase(t *testing.T) {
	compressor := NewCompressor()

	// ウィンドウサイズを超える長いデータをテスト
	data := make([]byte, 5000) // WindowSize (4096) より大きい
	for i := range data {
		data[i] = byte(i%26 + 'a') // a-z を繰り返し
	}

	compressed, err := compressor.Compress(data)
	if err != nil {
		t.Fatalf("Compress failed: %v", err)
	}

	decompressed, err := compressor.Decompress(compressed)
	if err != nil {
		t.Fatalf("Decompress failed: %v", err)
	}

	if !bytes.Equal(data, decompressed) {
		t.Errorf("Data mismatch for large input")
	}
}

func TestToken_IsLiteral(t *testing.T) {
	// リテラルトークンのテスト
	literalToken := Token{Distance: 0, Length: 0, Literal: 'a'}
	if !literalToken.IsLiteral() {
		t.Error("Expected token to be literal")
	}

	// マッチトークンのテスト
	matchToken := Token{Distance: 5, Length: 3, Literal: 'b'}
	if matchToken.IsLiteral() {
		t.Error("Expected token to not be literal")
	}
}

func TestLZ77Compressor_FindLongestMatch(t *testing.T) {
	compressor := NewCompressor()

	// テストデータ: "abcabcabc"
	data := []byte("abcabcabc")

	// 位置0では一致なし
	distance, length := compressor.FindLongestMatch(data, 0)
	if distance != 0 || length != 0 {
		t.Errorf("Expected no match at position 0, got distance=%d, length=%d", distance, length)
	}

	// 位置3では "abc" が一致するはず
	distance, length = compressor.FindLongestMatch(data, 3)
	if distance != 3 || length != 3 {
		t.Errorf("Expected match at position 3: distance=3, length=3, got distance=%d, length=%d", distance, length)
	}

	// 位置6では "abc" が一致するはず（より近い方を参照）
	distance, length = compressor.FindLongestMatch(data, 6)
	if distance != 3 || length != 3 {
		t.Errorf("Expected match at position 6: distance=3, length=3, got distance=%d, length=%d", distance, length)
	}
}

func TestLZ77Compressor_MinimumMatchLength(t *testing.T) {
	compressor := NewCompressor()

	// 長さ2の一致は無視されるはず（最小マッチ長は3）
	data := []byte("abcdefab") // "ab" が2文字だけ一致

	distance, length := compressor.FindLongestMatch(data, 6)
	if distance != 0 || length != 0 {
		t.Errorf("Expected no match for length 2, got distance=%d, length=%d", distance, length)
	}

	// 長さ3の一致は検出されるはず
	data = []byte("abcdefabc") // "abc" が3文字一致

	distance, length = compressor.FindLongestMatch(data, 6)
	if distance != 6 || length != 3 {
		t.Errorf("Expected match for length 3: distance=6, length=3, got distance=%d, length=%d", distance, length)
	}
}

func TestLZ77Compressor_BinaryData(t *testing.T) {
	compressor := NewCompressor()

	// バイナリデータのテスト
	original := []byte{0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03, 0xFF, 0xFE}

	compressed, err := compressor.Compress(original)
	if err != nil {
		t.Fatalf("Compress failed: %v", err)
	}

	decompressed, err := compressor.Decompress(compressed)
	if err != nil {
		t.Fatalf("Decompress failed: %v", err)
	}

	if !bytes.Equal(original, decompressed) {
		t.Errorf("Original: %v, Decompressed: %v", original, decompressed)
	}
}

func TestLZ77Compressor_InvalidData(t *testing.T) {
	compressor := NewCompressor()

	// 不正な圧縮データのテスト
	invalidData := []byte{1, 0, 5} // 不完全なマッチトークン

	_, err := compressor.Decompress(invalidData)
	if err == nil {
		t.Error("Expected error for invalid compressed data")
	}

	// 不正な距離のテスト
	invalidDistanceData := []byte{1, 0, 10, 3, 'a'} // 距離10だが結果が空

	_, err = compressor.Decompress(invalidDistanceData)
	if err == nil {
		t.Error("Expected error for invalid distance")
	}
}

// ベンチマークテスト
func BenchmarkLZ77Compress(b *testing.B) {
	compressor := NewCompressor()
	data := []byte("The quick brown fox jumps over the lazy dog. " +
		"The quick brown fox jumps over the lazy dog. " +
		"The quick brown fox jumps over the lazy dog.")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := compressor.Compress(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLZ77Decompress(b *testing.B) {
	compressor := NewCompressor()
	data := []byte("The quick brown fox jumps over the lazy dog. " +
		"The quick brown fox jumps over the lazy dog. " +
		"The quick brown fox jumps over the lazy dog.")

	compressed, err := compressor.Compress(data)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := compressor.Decompress(compressed)
		if err != nil {
			b.Fatal(err)
		}
	}
}
