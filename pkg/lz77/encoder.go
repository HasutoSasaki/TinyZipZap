package lz77

import (
	"encoding/binary"
)

// Encoder はLZ77のエンコード処理を担当します
type Encoder struct {
	matcher *Matcher
}

// NewEncoder は新しいEncoderを作成します
func NewEncoder(windowSize, bufferSize int) *Encoder {
	return &Encoder{
		matcher: NewMatcher(windowSize, bufferSize),
	}
}

// Encode はデータをLZ77トークンの配列にエンコードします
func (e *Encoder) Encode(data []byte) []Token {
	if len(data) == 0 {
		return []Token{}
	}

	var tokens []Token
	pos := 0

	for pos < len(data) {
		match := e.matcher.FindLongestMatch(data, pos)

		if match.Length > 0 {
			// マッチが見つかった場合
			nextChar := e.getNextChar(data, pos+match.Length)

			tokens = append(tokens, NewMatchToken(
				uint16(match.Distance),
				uint8(match.Length),
				nextChar,
			))

			pos += match.Length + 1 // マッチ長 + 次の文字
		} else {
			// マッチが見つからない場合はリテラル
			tokens = append(tokens, NewLiteralToken(data[pos]))
			pos++
		}
	}

	return tokens
}

// getNextChar は指定位置の次の文字を取得します（範囲外の場合は0を返す）
func (e *Encoder) getNextChar(data []byte, pos int) byte {
	if pos < len(data) {
		return data[pos]
	}
	return 0
}

// TokensToBytes はトークン配列をバイナリ形式にシリアライズします
func TokensToBytes(tokens []Token) []byte {
	var result []byte

	for _, token := range tokens {
		if token.IsLiteral() {
			// リテラル: フラグ(0) + 文字
			result = append(result, 0, token.Literal)
		} else {
			// マッチ: フラグ(1) + 距離(2バイト) + 長さ(1バイト) + リテラル(1バイト)
			result = append(result, 1)
			distanceBytes := make([]byte, 2)
			binary.BigEndian.PutUint16(distanceBytes, token.Distance)
			result = append(result, distanceBytes...)
			result = append(result, token.Length, token.Literal)
		}
	}

	return result
}
