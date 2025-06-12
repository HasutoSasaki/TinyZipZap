// Package rle implements Run-Length Encoding compression algorithm.
// RLE（ランレングス符号化）は、同じ文字が連続して現れる場合に
// 「文字 + 出現回数」の形式で圧縮する最も基本的な圧縮アルゴリズムです。
package rle

import (
	"bytes"
	"fmt"

	"github.com/sasakihasuto/tinyzipzap/pkg/common"
)

// RLECompressor はRun-Length Encoding圧縮を実装します
type RLECompressor struct{}

// NewRLECompressor は新しいRLECompressorを作成します
func NewRLECompressor() *RLECompressor {
	return &RLECompressor{}
}

// Name はアルゴリズム名を返します
func (r *RLECompressor) Name() string {
	return "Run-Length Encoding (RLE)"
}

// Compress はRLEアルゴリズムでデータを圧縮します
// 形式: [文字][カウント][文字][カウント]...
// カウントは1-255の範囲で、255を超える場合は分割します
func (r *RLECompressor) Compress(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return []byte{}, nil
	}
	
	var compressed bytes.Buffer
	
	currentByte := data[0]
	count := 1
	
	for i := 1; i < len(data); i++ {
		if data[i] == currentByte && count < 255 {
			count++
		} else {
			// 現在の文字とカウントを出力
			compressed.WriteByte(currentByte)
			compressed.WriteByte(byte(count))
			
			// 次の文字に移行
			currentByte = data[i]
			count = 1
		}
	}
	
	// 最後の文字とカウントを出力
	compressed.WriteByte(currentByte)
	compressed.WriteByte(byte(count))
	
	return compressed.Bytes(), nil
}

// Decompress はRLE圧縮されたデータを展開します
func (r *RLECompressor) Decompress(data []byte) ([]byte, error) {
	if len(data)%2 != 0 {
		return nil, fmt.Errorf("RLE: 圧縮データのサイズが不正です（奇数バイト）")
	}
	
	var decompressed bytes.Buffer
	
	for i := 0; i < len(data); i += 2 {
		if i+1 >= len(data) {
			return nil, fmt.Errorf("RLE: データが不完全です")
		}
		
		char := data[i]
		count := int(data[i+1])
		
		if count == 0 {
			return nil, fmt.Errorf("RLE: カウントが0です")
		}
		
		// 指定された回数だけ文字を繰り返し
		for j := 0; j < count; j++ {
			decompressed.WriteByte(char)
		}
	}
	
	return decompressed.Bytes(), nil
}

// AnalyzeData はRLE圧縮に適したデータかどうかを分析します
func AnalyzeData(data []byte) {
	if len(data) == 0 {
		fmt.Println("データが空です")
		return
	}
	
	// 連続する文字の長さを分析
	runs := make(map[int]int) // 連続長 -> 出現回数
	currentChar := data[0]
	currentLength := 1
	totalRuns := 0
	
	for i := 1; i < len(data); i++ {
		if data[i] == currentChar {
			currentLength++
		} else {
			runs[currentLength]++
			totalRuns++
			currentChar = data[i]
			currentLength = 1
		}
	}
	runs[currentLength]++
	totalRuns++
	
	fmt.Printf("=== RLE分析結果 ===\n")
	fmt.Printf("総ラン数: %d\n", totalRuns)
	fmt.Printf("平均ラン長: %.2f\n", float64(len(data))/float64(totalRuns))
	
	// 長いランの統計
	longRuns := 0
	for length, count := range runs {
		if length > 3 {
			longRuns += count
		}
	}
	
	fmt.Printf("長いラン (4文字以上): %d (%.1f%%)\n", 
		longRuns, float64(longRuns)/float64(totalRuns)*100)
	
	// RLE圧縮効果の予測
	originalSize := len(data)
	estimatedCompressed := totalRuns * 2 // 各ランは文字+カウントの2バイト
	
	fmt.Printf("予想圧縮サイズ: %d bytes\n", estimatedCompressed)
	fmt.Printf("予想圧縮率: %.2f%%\n", 
		float64(estimatedCompressed)/float64(originalSize)*100)
}

// CompressWithStats は圧縮と統計計算を同時に行います
func (r *RLECompressor) CompressWithStats(data []byte) ([]byte, common.CompressionStats, error) {
	compressed, err := r.Compress(data)
	if err != nil {
		return nil, common.CompressionStats{}, err
	}
	
	stats := common.CompressionStats{
		OriginalSize:   int64(len(data)),
		CompressedSize: int64(len(compressed)),
		Algorithm:      r.Name(),
	}
	stats.CalculateRatio()
	
	return compressed, stats, nil
}
