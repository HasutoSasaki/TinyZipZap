package lz77

// MatchResult はマッチング結果を表します
type MatchResult struct {
	Distance int
	Length   int
}

// Matcher はLZ77のマッチング処理を担当します
type Matcher struct {
	windowSize int
	bufferSize int
}

// NewMatcher は新しいMatcherを作成します
func NewMatcher(windowSize, bufferSize int) *Matcher {
	return &Matcher{
		windowSize: windowSize,
		bufferSize: bufferSize,
	}
}

// FindLongestMatch は最長一致を検索します
func (m *Matcher) FindLongestMatch(data []byte, pos int) MatchResult {
	if pos == 0 {
		return MatchResult{Distance: 0, Length: 0}
	}

	maxLength := 0
	bestDistance := 0

	// 検索開始位置を決定
	start := pos - m.windowSize
	if start < 0 {
		start = 0
	}

	// 先読みバッファの最大長を決定
	maxLookahead := len(data) - pos
	if maxLookahead > m.bufferSize {
		maxLookahead = m.bufferSize
	}

	// 検索ウィンドウ内で一致を探す
	for i := start; i < pos; i++ {
		matchLength := m.calculateMatchLength(data, i, pos, maxLookahead)

		// より長い一致が見つかった場合は更新（最小マッチ長は3）
		if matchLength > maxLength && matchLength >= 3 {
			maxLength = matchLength
			bestDistance = pos - i
		}
	}

	return MatchResult{Distance: bestDistance, Length: maxLength}
}

// calculateMatchLength は指定された位置からのマッチ長を計算します
func (m *Matcher) calculateMatchLength(data []byte, start, pos, maxLength int) int {
	length := 0
	for j := 0; j < maxLength && start+j < pos && data[start+j] == data[pos+j]; j++ {
		length++
	}
	return length
}
