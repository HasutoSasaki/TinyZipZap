package huffman

import (
	"bytes"
	"testing"
)

func TestCompressor_Name(t *testing.T) {
	compressor := NewCompressor()
	expected := "Huffman Coding"
	if compressor.Name() != expected {
		t.Errorf("Expected %s, got %s", expected, compressor.Name())
	}
}

func TestCompressor_EmptyData(t *testing.T) {
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

func TestCompressor_SingleCharacter(t *testing.T) {
	compressor := NewCompressor()
	original := []byte("aaaa")

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

func TestCompressor_BasicText(t *testing.T) {
	compressor := NewCompressor()
	testCases := [][]byte{
		[]byte("hello world"),
		[]byte("aaabbbccc"),
		[]byte("The quick brown fox jumps over the lazy dog"),
		[]byte("11111000001111100000"),
	}

	for _, original := range testCases {
		t.Run(string(original), func(t *testing.T) {
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
		})
	}
}

func TestCompressor_LargeData(t *testing.T) {
	compressor := NewCompressor()

	// 大きなデータを作成
	var original []byte
	for i := 0; i < 1000; i++ {
		original = append(original, byte(i%10+'0'))
	}

	compressed, err := compressor.Compress(original)
	if err != nil {
		t.Fatalf("Compress failed: %v", err)
	}

	decompressed, err := compressor.Decompress(compressed)
	if err != nil {
		t.Fatalf("Decompress failed: %v", err)
	}

	if !bytes.Equal(original, decompressed) {
		t.Errorf("Data mismatch")
	}
}

func BenchmarkCompress(b *testing.B) {
	compressor := NewCompressor()
	data := []byte("The quick brown fox jumps over the lazy dog. " +
		"This is a sample text for benchmarking compression algorithms.")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := compressor.Compress(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecompress(b *testing.B) {
	compressor := NewCompressor()
	data := []byte("The quick brown fox jumps over the lazy dog. " +
		"This is a sample text for benchmarking compression algorithms.")

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
