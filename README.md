# fmc

fmc (FrontMatter Compiler) は、Frontmatter を含むテキストファイルを JSON ファイルへコンパイルする CLI ツールです。

HTTP サーバー機能は持たず、生成した JSON を Nginx・S3・CloudFront 等の静的配信基盤で配信することを前提としています。

## インストール

```bash
go install github.com/add20/fmc/cmd/fmc@latest
```

## はじめかた

```bash
# プロジェクトの初期化
fmc init

# ビルド
fmc build

# 監視（変更時に自動ビルド）
fmc watch

# 出力ディレクトリの削除
fmc clean
```

`fmc init` を実行すると以下の構成が生成されます。

```
.
├── contents/        # ソースファイルを置くディレクトリ
├── dist/            # JSON の出力先
└── settings/
    └── config.toml  # 設定ファイル
```

## 対応フォーマット

拡張子による制限はありません。YAML・TOML の Frontmatter をサポートします。

### YAML

```
---
title: Hello
tags:
  - go
---

本文
```

### TOML

```
+++
title = "Hello"
tags = ["go"]
+++

本文
```

## 出力形式

各ファイルは `dist/` 以下に `元ファイル名.json` として出力されます。

```json
{
  "slug": "README",
  "srcPath": "2026/06/README.md",
  "frontMatter": {
    "title": "Hello"
  },
  "content": "本文"
}
```

`dist/index.json` には全ドキュメントの一覧が生成されます。

```json
[
  {
    "slug": "README",
    "path": "2026/06/README.md.json",
    "title": "Hello"
  }
]
```

## 設定

`settings/config.toml` で入出力ディレクトリを変更できます。

```toml
[contents]
dir = "contents"

[output]
dir = "dist"
```

## slug

contents からの相対パスにおいて、ファイル名から全ての拡張子を除去した値です。REST API のパスとして利用できます。

```
README.md               → README
archive.tar.gz          → archive
2026/06/README.md       → 2026/06/README
2026/06/archive.tar.gz  → 2026/06/archive
```

異なるディレクトリの同名ファイルは slug が異なるため重複しません。同一ディレクトリ内で拡張子のみ異なるファイルが存在する場合はビルドエラーになります。

## 開発

```bash
go build ./...
go test ./...
```
