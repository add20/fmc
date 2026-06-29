# fmc v1.0 実装方針

## 概要

fmc (FrontMatter Compiler) は Frontmatter を含むテキストファイル群を JSON ファイルへコンパイルする CLI ツールである。

HTTPサーバー機能は持たず、生成された JSON は Nginx、S3、CloudFront 等の静的配信基盤で配信することを前提とする。

---

# ディレクトリ構成

```text
fmc/
├── cmd/
│   └── fmc/
│       └── main.go
│
├── internal/
│   ├── cli/
│   ├── config/
│   ├── compiler/
│   ├── frontmatter/
│   └── watcher/
│
├── testdata/
│
├── go.mod
├── go.sum
└── README.md
```

watch機能は後から追加する。

まずは build を完成させる。

---

# Goバージョン

```text
Go 1.26
```

---

# 利用ライブラリ

## CLI

Cobra

## TOML

pelletier/go-toml

## YAML

go-yaml/yaml

## watch

fsnotify

## テスト

testify

---

# FrontMatter型

```go
type FrontMatter map[string]any
```

---

# JSON出力

```go
json.MarshalIndent
```

を利用する。

整形済みJSONを出力する。

---

# エラー設計

独自エラー型を利用する。

例

```go
type ErrorCode string

const (
    ErrDuplicateSlug    ErrorCode = "DUPLICATE_SLUG"
    ErrFrontMatterParse ErrorCode = "FRONTMATTER_PARSE"
    ErrConfigLoad       ErrorCode = "CONFIG_LOAD"
    ErrWriteFile        ErrorCode = "WRITE_FILE"
)

type FMCError struct {
    Code    ErrorCode
    Message string
    Cause   error
}
```

---

# Config

```go
type Config struct {
    Contents struct {
        Dir string `toml:"dir"`
    } `toml:"contents"`

    Output struct {
        Dir string `toml:"dir"`
    } `toml:"output"`
}
```

設定ファイル

```text
settings/config.toml
```

---

# Compiler API

v1.0では状態を持たない。

```go
func Build(cfg config.Config) error
```

を採用する。

Compiler構造体は作らない。

---

# Watch API

```go
func Watch(cfg config.Config) error
```

内部で fsnotify を利用する。

---

# index.json

並び順は

```text
srcPath 昇順
```

とする。

---

# title / draft

IndexEntry は以下とする。

```go
type IndexEntry struct {
    Slug  string  `json:"slug"`
    Path  string  `json:"path"`
    Title *string `json:"title"`
    Draft bool    `json:"draft"`
}
```

---

# frontMatter.title

title が存在しない場合

```json
{
  "title": null
}
```

を出力する。

title が文字列でない場合も

```json
{
  "title": null
}
```

を出力する。

エラーにはしない。

---

# frontMatter.draft

draft が存在しない場合、または bool でない場合は

```json
{
  "draft": false
}
```

を出力する。

エラーにはしない。

---

# watch時の挙動

## ファイル作成

全再ビルド

```text
Build()
```

を実行する。

---

## ファイル更新

全再ビルド

```text
Build()
```

を実行する。

差分ビルドは行わない。

---

## ファイル削除

全再ビルド

```text
Build()
```

を実行する。

差分削除は行わない。

---

# watch時のslug重複

例

```text
README.md
README.txt
```

↓

```text
README
```

重複

対応

```text
エラー表示
index.json削除
全出力停止
```

watchプロセス自体は継続する。

---

# watch時のFrontMatterパース失敗

対応

```text
エラー表示
該当JSON削除
watch継続
```

---

# init

既存ファイルは上書きしない。

例

```text
contents/
```

が存在する場合

```text
作成をスキップ
```

する。

エラーにはしない。

---

# clean

```text
dist/
```

が存在しなくても成功扱い。

---

# テスト方針

testify を利用する。

構成例

```text
internal/
├── compiler/
│   ├── slug_test.go
│   ├── validator_test.go
│   └── builder_test.go
│
├── frontmatter/
│   └── parser_test.go
│
├── config/
│   └── loader_test.go
```

---

# testdata

```text
testdata/
├── contents/
├── expected/
└── configs/
```

を用意する。

---

# 実装順序

## Phase 1

build機能完成

対象

```text
config
frontmatter
compiler
cli(build)
cli(init)
cli(clean)
```

完成条件

```bash
go test ./...
go run ./cmd/fmc build
```

が動作する。

---

## Phase 2

watch追加

対象

```text
internal/watcher
cli/watch
```

---

# 一括生成方針

watchを除いた build 機能一式は一括生成する。

対象

```text
go.mod
cmd/
internal/
testdata/
README.md
テストコード
```

その後 watch を追加する。
