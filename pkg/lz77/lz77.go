// Package lz77 implements LZ77 compression algorithm.
// LZ77は辞書ベースの圧縮アルゴリズムで、過去に出現した文字列を参照することで圧縮を行います。
package lz77

import (
	"github.com/sasakihasuto/tinyzipzap/pkg/common"
)

// Compressor はLZ77圧縮を実装します
type Compressor struct {
	encoder *Encoder
	decoder *Decoder
}

// NewCompressor は新しいCompressorを作成します
func NewCompressor() *Compressor {
	const (
		defaultWindowSize = 4096 // 4KB
		defaultBufferSize = 18   // 最大マッチ長
	)

	return &Compressor{
		encoder: NewEncoder(defaultWindowSize, defaultBufferSize),
		decoder: NewDecoder(),
	}
}

// Name はアルゴリズム名を返します
func (l *Compressor) Name() string {
	return "LZ77"
}

// Compress はLZ77アルゴリズムでデータを圧縮します
func (l *Compressor) Compress(data []byte) ([]byte, error) {
	tokens := l.encoder.Encode(data)
	return TokensToBytes(tokens), nil
}

// Decompress はLZ77圧縮されたデータを展開します
func (l *Compressor) Decompress(data []byte) ([]byte, error) {
	tokens, err := l.decoder.Decode(data)
	if err != nil {
		return nil, err
	}

	return l.decoder.TokensToData(tokens)
}

// FindLongestMatch は最長一致を検索します（テスト用の公開メソッド）
func (l *Compressor) FindLongestMatch(data []byte, pos int) (distance int, length int) {
	match := l.encoder.matcher.FindLongestMatch(data, pos)
	return match.Distance, match.Length
}

// コンパイル時にインターフェースの実装を確認
var _ common.Compressor = (*Compressor)(nil)
