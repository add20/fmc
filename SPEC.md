# fmc v1.0 Specification

## Overview

fmc (FrontMatter Compiler) は、Frontmatter を含むテキストファイル群を JSON ファイルへコンパイルするコマンドラインツールである。

fmc 自身は HTTP サーバー機能を持たない。

生成された JSON ファイルは Nginx、Apache、S3、CloudFront 等の静的ファイルサーバーで配信することを想定する。

---

## Goals

* シンプルであること
* DB を必要としないこと
* 常駐プロセスを必要としないこと
* Git 管理しやすいこと
* 静的配信と相性が良いこと
* Frontmatter と本文を JSON に変換すること

---

## Commands

### fmc init

カレントディレクトリに初期ディレクトリ構成を生成する。

```bash
fmc init
```

生成結果

```text
.
├── contents/
├── dist/
└── settings/
    └── config.toml
```

---

### fmc build

contents ディレクトリ以下を走査し、JSON ファイルを生成する。

```bash
fmc build
```

---

### fmc watch

contents ディレクトリ以下を監視する。

```bash
fmc watch
```

監視対象の変更

* ファイル作成
* ファイル更新
* ファイル削除

変更時は対応する JSON ファイルを更新する。

---

### fmc clean

dist ディレクトリを削除する。

```bash
fmc clean
```

---

## Configuration

設定ファイル

```text
settings/config.toml
```

初期内容

```toml
[contents]
dir = "contents"

[output]
dir = "dist"
```

### [index]

index.json の各エントリに含める追加の Frontmatter キーを指定する。任意。

```toml
[index]
fields = ["category", "tags"]
```

指定したキーは各エントリの `frontMatter` に入れ子で出力される。未指定の場合は出力しない。

---

## Supported Source Files

contents 以下の全ファイルを対象とする。

拡張子による制限は設けない。

例

```text
contents/
├── README.md
├── hello.txt
├── index.adoc
└── memo
```

すべてコンパイル対象。

---

## Frontmatter Formats

以下の形式をサポートする。

### YAML Frontmatter

```text
---
title: Hello
tags:
  - go
  - cms
---
```

### TOML Frontmatter

```text
+++
title = "Hello"

tags = ["go", "cms"]
+++
```

---

## Document Model

コンパイル結果の JSON は以下の構造を持つ。

```json
{
  "slug": "README",
  "srcPath": "2026/06/README.md",
  "frontMatter": {
    "title": "最初に読むファイル"
  },
  "content": "# README\n\n本文"
}
```

### slug

contents ディレクトリからの相対パスにおいて、ファイル名から全ての拡張子を除去した値。REST API のパスとして利用できる識別子とする。

例

```text
README.md               → README
README.md.txt           → README
archive.tar.gz          → archive
2026/06/README.md       → 2026/06/README
2026/06/archive.tar.gz  → 2026/06/archive
```

---

### srcPath

contents ディレクトリからの相対パス。

例

```text
contents/2026/06/README.md
```

↓

```json
{
  "srcPath": "2026/06/README.md"
}
```

---

### frontMatter

Frontmatter をパースした結果。

任意のキーを許可する。

---

### content

Frontmatter を除去した本文。

---

## Output Directory Structure

出力ディレクトリは source の構造を保持する。

入力

```text
contents/
└── 2026
    └── 06
        └── README.md
```

出力

```text
dist/
└── 2026
    └── 06
        └── README.md.json
```

---

## Output File Naming

出力ファイル名は

```text
元ファイル名 + ".json"
```

とする。

例

```text
README.md
```

↓

```text
README.md.json
```

例

```text
README.md.txt
```

↓

```text
README.md.txt.json
```

---

## index.json

fmc build は dist/index.json を生成する。

index.json は全ドキュメントの一覧である。

例

```json
[
  {
    "slug": "README",
    "path": "2026/06/README.md.json",
    "title": "最初に読むファイル",
    "draft": false,
    "frontMatter": {
      "category": "blog",
      "tags": ["go", "cms"]
    }
  }
]
```

### path

dist からの相対パス。

### title

frontMatter.title の値。

title が存在しない場合は null とする。

### draft

frontMatter.draft の値。

draft が存在しない場合、または bool でない場合は false とする。

### frontMatter

設定 `[index] fields` で指定したキーのうち、そのドキュメントの Frontmatter に存在するものを入れ子で出力する。

指定キーがそのドキュメントに存在しない場合は、そのキーを省略する。

`fields` 未指定、または該当キーが一つもない場合は `frontMatter` 自体を出力しない。

---

## Duplicate Slug Detection

slug が重複した場合はビルドエラーとする。

slug はディレクトリパスを含むため、異なるディレクトリに同名ファイルがあっても重複しない。
同一ディレクトリ内で拡張子のみ異なるファイルが存在する場合に重複となる。

例

```text
contents/
├── README.md
└── README.txt
```

両方とも

```text
README
```

となるためエラー。

```text
contents/
├── 2026/06/README.md
└── 2026/07/README.md
```

それぞれ

```text
2026/06/README
2026/07/README
```

となるため重複しない。

例（エラーメッセージ）

```text
ERROR: duplicate slug detected

slug: README

files:
- contents/README.md
- contents/README.txt
```

ビルドは失敗し、出力を行わない。

---

## Watch Behavior

### File Created

対応する JSON を生成する。

### File Modified

対応する JSON を再生成する。

### File Deleted

対応する JSON を削除する。

例

```text
contents/2026/06/README.md
```

削除

↓

```text
dist/2026/06/README.md.json
```

削除

さらに

```text
dist/index.json
```

も再生成する。

---

## Build Failure

以下の場合はビルド失敗とする。

* Frontmatter パース失敗
* YAML 構文エラー
* TOML 構文エラー
* slug 重複
* 出力ファイル書き込み失敗

終了コードは 1 を返す。

---

## Non Goals

v1.0 では以下を提供しない。

* HTTP サーバー
* REST API
* GraphQL
* データベース
* テンプレートエンジン
* HTML 生成
* RSS 生成
* 差分ビルドキャッシュ
* プラグイン機構

---

## Philosophy

fmc は Frontmatter 付きテキストファイルを JSON に変換することだけに責務を限定する。

配信、検索、キャッシュ、CDN、認証などは利用者側の責務とする。

```text
contents
    ↓
fmc build
    ↓
dist/*.json
    ↓
Nginx / S3 / CloudFront
```

以上を fmc v1.0 の仕様とする。
