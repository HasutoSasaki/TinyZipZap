package common

import (
	"fmt"
	"math"
)

// CountBytes はバイト配列内の各バイトの出現回数をカウントします
func CountBytes(data []byte) map[byte]int {
	counts := make(map[byte]int)
	for _, b := range data {
		counts[b]++
	}
	return counts
}

// CalculateEntropy はデータのエントロピーを計算します（情報理論）
func CalculateEntropy(data []byte) float64 {
	if len(data) == 0 {
		return 0
	}
	
	counts := CountBytes(data)
	total := float64(len(data))
	entropy := 0.0
	
	for _, count := range counts {
		if count > 0 {
			p := float64(count) / total
			entropy -= p * math.Log2(p)
		}
	}
	
	return entropy
}

// FormatBytes はバイト数を人間が読みやすい形式にフォーマットします
func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	
	units := []string{"KB", "MB", "GB", "TB"}
	return fmt.Sprintf("%.1f %s", float64(bytes)/float64(div), units[exp])
}

// PrintCompressionStats は圧縮統計を見やすく表示します
func PrintCompressionStats(stats CompressionStats) {
	fmt.Printf("=== 圧縮統計 ===\n")
	fmt.Printf("アルゴリズム: %s\n", stats.Algorithm)
	fmt.Printf("元のサイズ:   %s (%d bytes)\n", FormatBytes(stats.OriginalSize), stats.OriginalSize)
	fmt.Printf("圧縮後サイズ: %s (%d bytes)\n", FormatBytes(stats.CompressedSize), stats.CompressedSize)
	fmt.Printf("圧縮率:       %.2f%% (%.3f)\n", stats.Ratio*100, stats.Ratio)
	
	if stats.Ratio < 1.0 {
		reduction := (1.0 - stats.Ratio) * 100
		fmt.Printf("削減率:       %.2f%%\n", reduction)
	} else {
		increase := (stats.Ratio - 1.0) * 100
		fmt.Printf("サイズ増加:   %.2f%%\n", increase)
	}
}
