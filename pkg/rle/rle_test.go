package rle

import (
	"bytes"
	"testing"
)

func TestRLEBasic(t *testing.T) {
	compressor := NewRLECompressor()
	
	tests := []struct {
		name     string
		input    string
		expected string // 期待される圧縮結果（可視化のため文字列で表現）
	}{
		{
			name:     "シンプルな繰り返し",
			input:    "aaabbb",
			expected: "a3b3", // a:3回, b:3回
		},
		{
			name:     "単一文字",
			input:    "a",
			expected: "a1",
		},
		{
			name:     "繰り返しなし",
			input:    "abcd",
			expected: "a1b1c1d1",
		},
		{
			name:     "長い繰り返し",
			input:    string(bytes.Repeat([]byte("x"), 100)),
			expected: "x100",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 圧縮
			compressed, err := compressor.Compress([]byte(tt.input))
			if err != nil {
				t.Fatalf("圧縮エラー: %v", err)
			}
			
			// 展開
			decompressed, err := compressor.Decompress(compressed)
			if err != nil {
				t.Fatalf("展開エラー: %v", err)
			}
			
			// 元のデータと一致するかチェック
			if !bytes.Equal([]byte(tt.input), decompressed) {
				t.Errorf("展開結果が一致しません\n入力: %s\n展開結果: %s", tt.input, string(decompressed))
			}
			
			// 圧縮結果の確認（バイトレベル）
			t.Logf("入力: %s (%d bytes)", tt.input, len(tt.input))
			t.Logf("圧縮結果: %v (%d bytes)", compressed, len(compressed))
			t.Logf("展開結果: %s (%d bytes)", string(decompressed), len(decompressed))
		})
	}
}

func TestRLEEdgeCases(t *testing.T) {
	compressor := NewRLECompressor()
	
	t.Run("空データ", func(t *testing.T) {
		compressed, err := compressor.Compress([]byte{})
		if err != nil {
			t.Fatalf("空データの圧縮でエラー: %v", err)
		}
		
		decompressed, err := compressor.Decompress(compressed)
		if err != nil {
			t.Fatalf("空データの展開でエラー: %v", err)
		}
		
		if len(decompressed) != 0 {
			t.Errorf("空データの処理結果が空でない: %v", decompressed)
		}
	})
	
	t.Run("255文字以上の繰り返し", func(t *testing.T) {
		// 300文字の'a'を作成
		input := bytes.Repeat([]byte("a"), 300)
		
		compressed, err := compressor.Compress(input)
		if err != nil {
			t.Fatalf("長い繰り返しの圧縮でエラー: %v", err)
		}
		
		decompressed, err := compressor.Decompress(compressed)
		if err != nil {
			t.Fatalf("長い繰り返しの展開でエラー: %v", err)
		}
		
		if !bytes.Equal(input, decompressed) {
			t.Errorf("長い繰り返しの処理が正しくありません\n期待: %d文字の'a'\n結果: %d文字", 
				len(input), len(decompressed))
		}
		
		// 圧縮結果は'a',255,'a',45 のような形になるはず
		t.Logf("300文字のaを圧縮: %v", compressed)
	})
	
	t.Run("不正な圧縮データ", func(t *testing.T) {
		// 奇数バイトの圧縮データ（不正）
		invalidData := []byte{0x41} // 1バイトのみ
		
		_, err := compressor.Decompress(invalidData)
		if err == nil {
			t.Error("不正データでエラーが発生しませんでした")
		}
		
		// カウントが0の圧縮データ（不正）
		invalidData2 := []byte{0x41, 0x00} // 'A', count=0
		
		_, err = compressor.Decompress(invalidData2)
		if err == nil {
			t.Error("カウント0でエラーが発生しませんでした")
		}
	})
}

func TestRLEWithStats(t *testing.T) {
	compressor := NewRLECompressor()
	
	// 圧縮に適したデータ（多くの繰り返し）
	goodData := []byte("aaaaabbbbcccccddddd")
	compressed, stats, err := compressor.CompressWithStats(goodData)
	if err != nil {
		t.Fatalf("統計付き圧縮でエラー: %v", err)
	}
	
	t.Logf("良いデータの圧縮統計:")
	t.Logf("  元のサイズ: %d", stats.OriginalSize)
	t.Logf("  圧縮後: %d", stats.CompressedSize)
	t.Logf("  圧縮率: %.2f%%", stats.Ratio*100)
	
	// 圧縮に不適なデータ（繰り返しなし）
	badData := []byte("abcdefghijklmnopqrstuvwxyz")
	compressed2, stats2, err := compressor.CompressWithStats(badData)
	if err != nil {
		t.Fatalf("統計付き圧縮でエラー: %v", err)
	}
	
	t.Logf("悪いデータの圧縮統計:")
	t.Logf("  元のサイズ: %d", stats2.OriginalSize)
	t.Logf("  圧縮後: %d", stats2.CompressedSize)
	t.Logf("  圧縮率: %.2f%%", stats2.Ratio*100)
	
	// 圧縮効果の確認
	if stats.Ratio >= 1.0 {
		t.Error("繰り返しの多いデータで圧縮効果がありません")
	}
	
	if stats2.Ratio <= 1.0 {
		t.Log("注意: 繰り返しのないデータでも圧縮されました（予想外）")
	}
	
	// 正しく展開できるか確認
	decompressed, _ := compressor.Decompress(compressed)
	decompressed2, _ := compressor.Decompress(compressed2)
	
	if !bytes.Equal(goodData, decompressed) {
		t.Error("良いデータの展開結果が一致しません")
	}
	
	if !bytes.Equal(badData, decompressed2) {
		t.Error("悪いデータの展開結果が一致しません")
	}
}

// ベンチマークテスト
func BenchmarkRLECompress(b *testing.B) {
	compressor := NewRLECompressor()
	data := bytes.Repeat([]byte("Hello World! "), 1000) // 繰り返しのあるデータ
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := compressor.Compress(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRLEDecompress(b *testing.B) {
	compressor := NewRLECompressor()
	data := bytes.Repeat([]byte("Hello World! "), 1000)
	compressed, _ := compressor.Compress(data)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := compressor.Decompress(compressed)
		if err != nil {
			b.Fatal(err)
		}
	}
}
