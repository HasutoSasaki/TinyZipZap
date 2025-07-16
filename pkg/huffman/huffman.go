// Package huffman implements Huffman Coding compression algorithm.
// Huffman符号化は、文字の出現頻度に基づいてより短い符号を割り当てる圧縮アルゴリズムです。
package huffman

import (
	"container/heap"
	"fmt"
	"sort"

	"github.com/sasakihasuto/tinyzipzap/pkg/common"
)

// Compressor はHuffman Coding圧縮を実装します
type Compressor struct{}

// NewCompressor は新しいCompressorを作成します
func NewCompressor() *Compressor {
	return &Compressor{}
}

// Name はアルゴリズム名を返します
func (h *Compressor) Name() string {
	return "Huffman Coding"
}

// Node はHuffman木のノードを表します
type Node struct {
	Char  byte  // 文字（リーフノードの場合）
	Freq  int   // 頻度
	Left  *Node // 左の子ノード
	Right *Node // 右の子ノード
}

// IsLeaf はリーフノードかどうかを判定します
func (n *Node) IsLeaf() bool {
	return n.Left == nil && n.Right == nil
}

// NodeHeap はヒープ操作のための構造体
type NodeHeap []*Node

func (h NodeHeap) Len() int           { return len(h) }
func (h NodeHeap) Less(i, j int) bool { return h[i].Freq < h[j].Freq }
func (h NodeHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *NodeHeap) Push(x interface{}) {
	*h = append(*h, x.(*Node))
}

func (h *NodeHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// buildFrequencyTable は文字の出現頻度テーブルを構築します
func buildFrequencyTable(data []byte) map[byte]int {
	freq := make(map[byte]int)
	for _, b := range data {
		freq[b]++
	}
	return freq
}

// buildTree はHuffman木を構築します
func buildTree(freq map[byte]int) *Node {
	if len(freq) == 0 {
		return nil
	}

	// 単一文字の場合
	if len(freq) == 1 {
		for char, f := range freq {
			return &Node{Char: char, Freq: f}
		}
	}

	// ヒープを初期化
	h := &NodeHeap{}
	heap.Init(h)

	// 各文字をソートしてからノードとしてヒープに追加（安定性を保証）
	var chars []byte
	for char := range freq {
		chars = append(chars, char)
	}
	sort.Slice(chars, func(i, j int) bool { return chars[i] < chars[j] })

	for _, char := range chars {
		heap.Push(h, &Node{Char: char, Freq: freq[char]})
	}

	// Huffman木を構築
	for h.Len() > 1 {
		left := heap.Pop(h).(*Node)
		right := heap.Pop(h).(*Node)

		merged := &Node{
			Freq:  left.Freq + right.Freq,
			Left:  left,
			Right: right,
		}
		heap.Push(h, merged)
	}

	return heap.Pop(h).(*Node)
}

// buildCodeTable はHuffman符号テーブルを構築します
func buildCodeTable(root *Node) map[byte]string {
	if root == nil {
		return make(map[byte]string)
	}

	// 単一文字の場合
	if root.IsLeaf() {
		return map[byte]string{root.Char: "0"}
	}

	codes := make(map[byte]string)
	var buildCodes func(*Node, string)
	buildCodes = func(node *Node, code string) {
		if node == nil {
			return
		}
		if node.IsLeaf() {
			codes[node.Char] = code
			return
		}
		buildCodes(node.Left, code+"0")
		buildCodes(node.Right, code+"1")
	}

	buildCodes(root, "")
	return codes
}

// Compress はHuffmanアルゴリズムでデータを圧縮します
func (h *Compressor) Compress(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return []byte{}, nil
	}

	// 頻度テーブルを構築
	freq := buildFrequencyTable(data)

	// Huffman木を構築
	root := buildTree(freq)
	if root == nil {
		return nil, fmt.Errorf("failed to build Huffman tree")
	}

	// 符号テーブルを構築
	codes := buildCodeTable(root)

	// 圧縮データを構築
	var compressed []byte

	// ヘッダー: 文字数 + 頻度テーブル
	compressed = append(compressed, byte(len(freq)))

	// 頻度テーブルをソートして保存
	var chars []byte
	for char := range freq {
		chars = append(chars, char)
	}
	sort.Slice(chars, func(i, j int) bool { return chars[i] < chars[j] })

	for _, char := range chars {
		compressed = append(compressed, char)
		f := freq[char]
		// 頻度を4バイトで保存
		compressed = append(compressed,
			byte(f>>24), byte(f>>16), byte(f>>8), byte(f))
	}

	// データを符号化
	var bits []byte
	bitCount := 0
	currentByte := byte(0)

	for _, b := range data {
		code := codes[b]
		for _, bit := range code {
			if bit == '1' {
				currentByte |= (1 << (7 - bitCount))
			}
			bitCount++
			if bitCount == 8 {
				bits = append(bits, currentByte)
				currentByte = 0
				bitCount = 0
			}
		}
	}

	// 最後のバイトを処理
	if bitCount > 0 {
		bits = append(bits, currentByte)
	}

	// データ長（ビット数）を保存
	dataLen := len(data)
	compressed = append(compressed,
		byte(dataLen>>24), byte(dataLen>>16), byte(dataLen>>8), byte(dataLen))

	// 余分なビット数を保存
	paddingBits := byte(8 - bitCount)
	if bitCount == 0 {
		paddingBits = 0
	}
	compressed = append(compressed, paddingBits)

	// 符号化されたデータを追加
	compressed = append(compressed, bits...)

	return compressed, nil
}

// Decompress はHuffman圧縮されたデータを展開します
func (h *Compressor) Decompress(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return []byte{}, nil
	}

	offset := 0

	// 文字数を読み取り
	if offset >= len(data) {
		return nil, fmt.Errorf("invalid compressed data: missing character count")
	}
	charCount := int(data[offset])
	offset++

	// 頻度テーブルを再構築
	freq := make(map[byte]int)
	for i := 0; i < charCount; i++ {
		if offset+4 >= len(data) {
			return nil, fmt.Errorf("invalid compressed data: incomplete frequency table")
		}
		char := data[offset]
		offset++
		f := int(data[offset])<<24 | int(data[offset+1])<<16 |
			int(data[offset+2])<<8 | int(data[offset+3])
		offset += 4
		freq[char] = f
	}

	// Huffman木を再構築
	root := buildTree(freq)
	if root == nil {
		return nil, fmt.Errorf("failed to rebuild Huffman tree")
	}

	// データ長を読み取り
	if offset+4 >= len(data) {
		return nil, fmt.Errorf("invalid compressed data: missing data length")
	}
	dataLen := int(data[offset])<<24 | int(data[offset+1])<<16 |
		int(data[offset+2])<<8 | int(data[offset+3])
	offset += 4

	// 余分なビット数を読み取り
	if offset >= len(data) {
		return nil, fmt.Errorf("invalid compressed data: missing padding bits")
	}
	paddingBits := int(data[offset])
	offset++

	// 符号化されたデータを展開
	var result []byte
	current := root

	// 単一文字の場合
	if root.IsLeaf() {
		for i := 0; i < dataLen; i++ {
			result = append(result, root.Char)
		}
		return result, nil
	}

	for offset < len(data) && len(result) < dataLen {
		b := data[offset]

		for i := 7; i >= 0 && len(result) < dataLen; i-- {
			// 最後のバイトの余分なビットをスキップ
			if offset == len(data)-1 && i < paddingBits {
				break
			}

			bit := (b >> i) & 1
			if bit == 1 {
				current = current.Right
			} else {
				current = current.Left
			}

			if current.IsLeaf() {
				result = append(result, current.Char)
				current = root
			}
		}
		offset++
	}

	return result, nil
}

var _ common.Compressor = (*Compressor)(nil)
