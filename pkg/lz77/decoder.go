package lz77

import (
	"encoding/binary"
	"fmt"
)

// Decoder はLZ77のデコード処理を担当します
type Decoder struct{}

// NewDecoder は新しいDecoderを作成します
func NewDecoder() *Decoder {
	return &Decoder{}
}

// Decode はバイナリデータをLZ77トークンの配列にパースします
func (d *Decoder) Decode(data []byte) ([]Token, error) {
	if len(data) == 0 {
		return []Token{}, nil
	}

	var tokens []Token
	pos := 0

	for pos < len(data) {
		if pos >= len(data) {
			break
		}

		flag := data[pos]
		pos++

		if flag == 0 {
			// リテラル
			if pos >= len(data) {
				return nil, fmt.Errorf("invalid compressed data: missing literal")
			}
			tokens = append(tokens, NewLiteralToken(data[pos]))
			pos++
		} else {
			// マッチ
			if pos+4 > len(data) {
				return nil, fmt.Errorf("invalid compressed data: incomplete match token")
			}

			distance := binary.BigEndian.Uint16(data[pos : pos+2])
			pos += 2
			length := data[pos]
			pos++
			literal := data[pos]
			pos++

			tokens = append(tokens, NewMatchToken(distance, length, literal))
		}
	}

	return tokens, nil
}

// TokensToData はトークン配列を元のデータに復元します
func (d *Decoder) TokensToData(tokens []Token) ([]byte, error) {
	var result []byte

	for _, token := range tokens {
		if token.IsLiteral() {
			result = append(result, token.Literal)
		} else {
			// 距離チェック
			if int(token.Distance) > len(result) {
				return nil, fmt.Errorf("invalid distance: %d, result length: %d", token.Distance, len(result))
			}

			// マッチした文字列をコピー
			if err := d.copyMatch(&result, int(token.Distance), int(token.Length)); err != nil {
				return nil, err
			}

			// 次のリテラル文字を追加
			if token.Literal != 0 || len(result) < cap(result) {
				result = append(result, token.Literal)
			}
		}
	}

	return result, nil
}

// copyMatch はマッチした文字列を結果にコピーします
func (d *Decoder) copyMatch(result *[]byte, distance, length int) error {
	start := len(*result) - distance

	for i := 0; i < length; i++ {
		if start+i >= len(*result) {
			return fmt.Errorf("invalid match: trying to copy from future position")
		}
		*result = append(*result, (*result)[start+i])
	}

	return nil
}
