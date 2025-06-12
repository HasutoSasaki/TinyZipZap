package common

import "io"

// Compressor は圧縮アルゴリズムの共通インターフェース
type Compressor interface {
	// Compress はデータを圧縮します
	Compress(data []byte) ([]byte, error)
	
	// Decompress は圧縮されたデータを展開します
	Decompress(data []byte) ([]byte, error)
	
	// Name はアルゴリズム名を返します
	Name() string
}

// StreamCompressor はストリーミング圧縮のインターフェース
type StreamCompressor interface {
	// CompressStream はストリームを圧縮します
	CompressStream(src io.Reader, dst io.Writer) error
	
	// DecompressStream は圧縮ストリームを展開します
	DecompressStream(src io.Reader, dst io.Writer) error
}

// CompressionStats は圧縮統計情報
type CompressionStats struct {
	OriginalSize   int64   // 元のサイズ
	CompressedSize int64   // 圧縮後のサイズ
	Ratio          float64 // 圧縮率
	Algorithm      string  // 使用アルゴリズム
}

// CalculateRatio は圧縮率を計算します
func (s *CompressionStats) CalculateRatio() {
	if s.OriginalSize > 0 {
		s.Ratio = float64(s.CompressedSize) / float64(s.OriginalSize)
	}
}
