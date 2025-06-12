package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/sasakihasuto/tinyzipzap/pkg/common"
	"github.com/sasakihasuto/tinyzipzap/pkg/rle"
)

const version = "1.0.0"

func main() {
	var (
		algorithm = flag.String("algo", "rle", "圧縮アルゴリズム (rle, huffman)")
		compress  = flag.Bool("c", false, "圧縮モード")
		decompress = flag.Bool("d", false, "展開モード") 
		analyze   = flag.Bool("a", false, "分析モード")
		input     = flag.String("i", "", "入力ファイル")
		output    = flag.String("o", "", "出力ファイル")
		verbose   = flag.Bool("v", false, "詳細出力")
		showVersion = flag.Bool("version", false, "バージョン表示")
	)
	
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "TinyZipZap - 圧縮アルゴリズム学習ツール v%s\n\n", version)
		fmt.Fprintf(os.Stderr, "使用方法:\n")
		fmt.Fprintf(os.Stderr, "  %s [オプション]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "オプション:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n例:\n")
		fmt.Fprintf(os.Stderr, "  # ファイルをRLEで圧縮\n")
		fmt.Fprintf(os.Stderr, "  %s -c -algo rle -i sample.txt -o sample.rle\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  # 圧縮ファイルを展開\n")
		fmt.Fprintf(os.Stderr, "  %s -d -algo rle -i sample.rle -o output.txt\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  # ファイルを分析\n")
		fmt.Fprintf(os.Stderr, "  %s -a -algo rle -i sample.txt\n\n", os.Args[0])
	}
	
	flag.Parse()
	
	if *showVersion {
		fmt.Printf("TinyZipZap v%s\n", version)
		return
	}
	
	// 基本的な引数チェック
	if *input == "" {
		fmt.Fprintf(os.Stderr, "エラー: 入力ファイルが指定されていません\n\n")
		flag.Usage()
		os.Exit(1)
	}
	
	// モードの確認
	modeCount := 0
	if *compress { modeCount++ }
	if *decompress { modeCount++ }
	if *analyze { modeCount++ }
	
	if modeCount == 0 {
		fmt.Fprintf(os.Stderr, "エラー: モード(-c, -d, -a)を指定してください\n\n")
		flag.Usage()
		os.Exit(1)
	}
	
	if modeCount > 1 {
		fmt.Fprintf(os.Stderr, "エラー: 複数のモードは同時に指定できません\n\n")
		flag.Usage()
		os.Exit(1)
	}
	
	// ファイルの読み込み
	data, err := ioutil.ReadFile(*input)
	if err != nil {
		log.Fatalf("ファイル読み込みエラー: %v", err)
	}
	
	if *verbose {
		fmt.Printf("入力ファイル: %s (%s)\n", *input, common.FormatBytes(int64(len(data))))
		fmt.Printf("アルゴリズム: %s\n", strings.ToUpper(*algorithm))
		fmt.Printf("データサイズ: %d bytes\n", len(data))
		if len(data) > 0 {
			fmt.Printf("エントロピー: %.3f bits/byte\n", common.CalculateEntropy(data))
		}
		fmt.Println()
	}
	
	// アルゴリズムの選択
	var compressor common.Compressor
	switch strings.ToLower(*algorithm) {
	case "rle":
		compressor = rle.NewRLECompressor()
	default:
		log.Fatalf("未対応のアルゴリズム: %s", *algorithm)
	}
	
	// モードに応じた処理
	switch {
	case *analyze:
		handleAnalyze(compressor, data, *verbose)
	case *compress:
		handleCompress(compressor, data, *input, *output, *verbose)
	case *decompress:
		handleDecompress(compressor, data, *input, *output, *verbose)
	}
}

func handleAnalyze(compressor common.Compressor, data []byte, verbose bool) {
	fmt.Printf("=== データ分析結果 ===\n")
	fmt.Printf("アルゴリズム: %s\n", compressor.Name())
	fmt.Printf("データサイズ: %s (%d bytes)\n", common.FormatBytes(int64(len(data))), len(data))
	
	if len(data) > 0 {
		entropy := common.CalculateEntropy(data)
		fmt.Printf("エントロピー: %.3f bits/byte\n", entropy)
		fmt.Printf("理論的最小サイズ: %.1f bytes\n", entropy * float64(len(data)) / 8)
	}
	
	fmt.Println()
	
	// アルゴリズム固有の分析
	if _, ok := compressor.(*rle.RLECompressor); ok {
		rle.AnalyzeData(data)
		fmt.Println()
	}
	
	// 実際に圧縮してみる
	fmt.Println("=== 圧縮テスト ===")
	switch comp := compressor.(type) {
	case *rle.RLECompressor:
		_, stats, err := comp.CompressWithStats(data)
		if err != nil {
			log.Fatalf("圧縮テストエラー: %v", err)
		}
		common.PrintCompressionStats(stats)
	default:
		compressed, err := compressor.Compress(data)
		if err != nil {
			log.Fatalf("圧縮テストエラー: %v", err)
		}
		
		stats := common.CompressionStats{
			OriginalSize:   int64(len(data)),
			CompressedSize: int64(len(compressed)),
			Algorithm:      compressor.Name(),
		}
		stats.CalculateRatio()
		common.PrintCompressionStats(stats)
	}
}

func handleCompress(compressor common.Compressor, data []byte, inputFile, outputFile string, verbose bool) {
	if outputFile == "" {
		outputFile = inputFile + ".compressed"
	}
	
	var compressed []byte
	var stats common.CompressionStats
	var err error
	
	// 統計付き圧縮があれば使用
	if rleComp, ok := compressor.(*rle.RLECompressor); ok {
		compressed, stats, err = rleComp.CompressWithStats(data)
	} else {
		compressed, err = compressor.Compress(data)
		if err == nil {
			stats = common.CompressionStats{
				OriginalSize:   int64(len(data)),
				CompressedSize: int64(len(compressed)),
				Algorithm:      compressor.Name(),
			}
			stats.CalculateRatio()
		}
	}
	
	if err != nil {
		log.Fatalf("圧縮エラー: %v", err)
	}
	
	// ファイルに書き込み
	err = ioutil.WriteFile(outputFile, compressed, 0644)
	if err != nil {
		log.Fatalf("ファイル書き込みエラー: %v", err)
	}
	
	fmt.Printf("✅ 圧縮完了: %s -> %s\n", inputFile, outputFile)
	
	if verbose {
		fmt.Println()
		common.PrintCompressionStats(stats)
	} else {
		fmt.Printf("圧縮率: %.2f%% (%s -> %s)\n", 
			stats.Ratio*100,
			common.FormatBytes(stats.OriginalSize),
			common.FormatBytes(stats.CompressedSize))
	}
}

func handleDecompress(compressor common.Compressor, data []byte, inputFile, outputFile string, verbose bool) {
	if outputFile == "" {
		ext := filepath.Ext(inputFile)
		if ext == ".compressed" {
			outputFile = strings.TrimSuffix(inputFile, ext)
		} else {
			outputFile = inputFile + ".decompressed"
		}
	}
	
	decompressed, err := compressor.Decompress(data)
	if err != nil {
		log.Fatalf("展開エラー: %v", err)
	}
	
	// ファイルに書き込み
	err = ioutil.WriteFile(outputFile, decompressed, 0644)
	if err != nil {
		log.Fatalf("ファイル書き込みエラー: %v", err)
	}
	
	fmt.Printf("✅ 展開完了: %s -> %s\n", inputFile, outputFile)
	
	if verbose {
		fmt.Printf("圧縮サイズ: %s (%d bytes)\n", 
			common.FormatBytes(int64(len(data))), len(data))
		fmt.Printf("展開サイズ: %s (%d bytes)\n", 
			common.FormatBytes(int64(len(decompressed))), len(decompressed))
	}
}
