# TinyZipZap 🗜️

Go で学ぶ圧縮アルゴリズム実装プロジェクト

## 📖 概要

TinyZipZap は、圧縮アルゴリズムを学習するための Go プロジェクトです。基本的な圧縮アルゴリズムをシンプルに実装し、その仕組みを理解することを目的としています。

## 🎯 実装済みアルゴリズム

### ✅ Run-Length Encoding (RLE)

- 最もシンプルな圧縮アルゴリズム
- 同じ文字の連続を「文字+回数」で表現
- 繰り返しの多いデータに効果的

### 🚧 予定しているアルゴリズム

- [ ] Huffman Coding (頻度ベースの圧縮)
- [ ] LZ77 (辞書ベースの圧縮)
- [ ] 簡易 Deflate (LZ77 + Huffman)

## 🚀 使用方法

### ビルド

```bash
go build -o tinyzipzap cmd/tinyzipzap/main.go
```

### 基本的な使用例

#### ファイルの分析

```bash
./tinyzipzap -a -algo rle -i examples/sample.txt
```

#### ファイルの圧縮

```bash
./tinyzipzap -c -algo rle -i examples/sample.txt -o sample.rle
```

#### ファイルの展開

```bash
./tinyzipzap -d -algo rle -i sample.rle -o restored.txt
```

#### 詳細出力付き

```bash
./tinyzipzap -c -algo rle -i examples/sample.txt -v
```

## 📁 プロジェクト構造

```
TinyZipZap/
├── README.md                    # このファイル
├── go.mod                       # Goモジュール設定
├── cmd/
│   └── tinyzipzap/
│       └── main.go             # CLIツール
├── pkg/
│   ├── common/
│   │   ├── types.go            # 共通インターフェース
│   │   └── utils.go            # ユーティリティ関数
│   └── rle/
│       ├── encoder.go          # RLE実装
│       └── rle_test.go         # RLEテスト
├── examples/
│   └── sample.txt              # テスト用サンプル
└── docs/                       # ドキュメント（予定）
```

## 🧪 テスト実行

```bash
# 全テスト実行
go test ./...

# RLEのテストのみ
go test ./pkg/rle/

# ベンチマーク実行
go test -bench=. ./pkg/rle/

# 詳細出力付きテスト
go test -v ./pkg/rle/
```

## 📚 学習ポイント

### Run-Length Encoding (RLE)

- **仕組み**: 連続する同じ文字を「文字+回数」で表現
- **適用場面**: ロゴ画像、単色領域の多い画像など
- **特徴**:
  - 実装が簡単
  - 繰り返しのないデータでは逆に大きくなる
  - リアルタイム処理に適している

### データ分析機能

- エントロピー計算（情報理論）
- 圧縮率の予測
- ランレングス分析

## 🔧 開発

### 新しいアルゴリズムの追加

1. `pkg/` 以下に新しいパッケージを作成
2. `common.Compressor` インターフェースを実装
3. テストファイルを作成
4. `main.go` に追加

### 設計原則

- シンプルで理解しやすい実装
- 充実したテストとベンチマーク
- 詳細な分析機能
- 学習に役立つコメント

## 📊 実行例

```bash
$ ./tinyzipzap -a -algo rle -i examples/sample.txt

=== データ分析結果 ===
アルゴリズム: Run-Length Encoding (RLE)
データサイズ: 1.1 KB (1137 bytes)
エントロピー: 4.693 bits/byte
理論的最小サイズ: 667.4 bytes

=== RLE分析結果 ===
総ラン数: 387
平均ラン長: 2.94
長いラン (4文字以上): 19 (4.9%)
予想圧縮サイズ: 774 bytes
予想圧縮率: 68.07%

=== 圧縮テスト ===
=== 圧縮統計 ===
アルゴリズム: Run-Length Encoding (RLE)
元のサイズ:   1.1 KB (1137 bytes)
圧縮後サイズ: 756 bytes (756 bytes)
圧縮率:       66.49% (0.665)
削減率:       33.51%
```

## 🎓 学習リソース

- [情報理論とエントロピー](https://ja.wikipedia.org/wiki/情報エントロピー)
- [データ圧縮の基礎](https://ja.wikipedia.org/wiki/データ圧縮)
- [ランレングス符号化](https://ja.wikipedia.org/wiki/ランレングス符号)

## 📄 ライセンス

MIT License
