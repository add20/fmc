# CLAUDE.md

このファイルは、Claude Code (claude.ai/code) がこのリポジトリで作業する際のガイダンスを提供します。

## プロジェクト概要

`fmc` (FrontMatter Compiler) は、Frontmatter を含むテキストファイルを JSON ファイルへコンパイルする Go 製 CLI ツールである。HTTP サーバー機能は持たず、JSON の生成のみを行う。生成した JSON は Nginx、S3、CloudFront 等の静的配信基盤で配信することを前提とする。

仕様の詳細は `SPEC.md`、実装方針は `IMPLEMENTATION_POLICY.md` を参照。

## コマンド

```bash
# ビルド
go build ./...

# ツールの実行
go run ./cmd/fmc <command>

# テスト
go test ./...
go test ./internal/compiler/...   # パッケージ単体
go test -run TestBuild ./...      # テスト単体

# Lint（golangci-lint が使える場合）
golangci-lint run
```

## アーキテクチャ

```
cmd/fmc/main.go          # Cobra ルートコマンド。各サブコマンドを登録する
internal/
  cli/                   # Cobra サブコマンド: build, init, clean, watch
  config/                # 設定ローダー — settings/config.toml を読み込む
  compiler/              # ビルドロジック: slug 生成、重複検出、JSON 出力
  frontmatter/           # YAML / TOML frontmatter パーサー
  watcher/               # fsnotify を使った watch ループ (Phase 2)
testdata/
  contents/              # テスト用ソースファイル
  expected/              # 期待する JSON 出力
  configs/               # テスト用 config.toml バリアント
```

## 設計上の重要事項

**Go バージョン**: 1.26

**依存ライブラリ**: Cobra (CLI)、`pelletier/go-toml` (TOML)、`go-yaml/yaml` (YAML)、`fsnotify` (watch)、`testify` (テスト)

**Compiler 構造体は作らない** — v1.0 では状態を持たない関数 `Build(cfg config.Config) error` を採用する。`Watch(cfg config.Config) error` も同様。

**slug**: ファイル名から全ての拡張子を除去した値 (`README.md.txt` → `README`)。slug が重複した場合はビルドエラーとし、出力は一切行わない。

**出力ファイル名**: 元のファイル名に `.json` を付加する (`README.md` → `README.md.json`)。ディレクトリ構造は `dist/` 以下に保持される。

**index.json**: `dist/index.json` に生成。`srcPath` の昇順でソートする。各エントリは `IndexEntry{Slug, Path, Title *string}`。`frontMatter.title` が存在しない、または文字列でない場合、`title` は `null` とする（エラーにはしない）。

**FrontMatter 型**: `map[string]any`。YAML (`---`) と TOML (`+++`) の両デリミタをサポートする。ファイル拡張子による制限はない。

**エラー型**: `FMCError{Code ErrorCode, Message string, Cause error}` の独自型を使う。エラーコードは `DUPLICATE_SLUG`、`FRONTMATTER_PARSE`、`CONFIG_LOAD`、`WRITE_FILE`。

**設定ファイルのパス**: `settings/config.toml`（カレントディレクトリからの相対パス）。

**`fmc init`**: 既存のファイル・ディレクトリは上書きしない。既存の場合は作成をスキップし、エラーにはしない。

**`fmc clean`**: `dist/` が存在しなくても成功扱いとする。

**watch の挙動 (Phase 2)**: 変更（作成・更新・削除）が発生した場合は常に `Build()` を全実行する（差分ビルドなし）。watch 中に slug 重複が発生した場合はエラーを表示し `index.json` を削除して出力を停止するが、watch プロセス自体は継続する。watch 中に Frontmatter のパースに失敗した場合はエラーを表示し該当 JSON を削除して watch を継続する。

## 実装フェーズ

- **Phase 1**（現在の対象）: `config`、`frontmatter`、`compiler`、`cli/build`、`cli/init`、`cli/clean` — `go test ./...` と `go run ./cmd/fmc build` が動作すれば完了。
- **Phase 2**: `internal/watcher` + `cli/watch`。
