package lz77

// Token はLZ77のトークンを表します
type Token struct {
	Distance uint16 // 後方距離（0の場合はリテラル）
	Length   uint8  // マッチ長
	Literal  byte   // リテラル文字（Distance=0の場合に使用）
}

// IsLiteral はトークンがリテラルかどうかを判定します
func (t Token) IsLiteral() bool {
	return t.Distance == 0
}

// NewLiteralToken はリテラルトークンを作成します
func NewLiteralToken(literal byte) Token {
	return Token{
		Distance: 0,
		Length:   0,
		Literal:  literal,
	}
}

// NewMatchToken はマッチトークンを作成します
func NewMatchToken(distance uint16, length uint8, nextChar byte) Token {
	return Token{
		Distance: distance,
		Length:   length,
		Literal:  nextChar,
	}
}
