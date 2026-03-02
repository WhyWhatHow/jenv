<div align="center">
<img src="assets/jenv-logo.png" width="200" height="200" alt="JEnv ロゴ">

# Jenv: Java 環境マネージャー

![GitHub release](https://img.shields.io/github/v/release/WhyWhatHow/jenv)
![Build Status](https://img.shields.io/github/actions/workflow/status/WhyWhatHow/jenv/release.yml?branch=main)
![Version](https://img.shields.io/badge/version-v0.6.7-blue)

[English](README.md) | [中文](README_zh.md) | 日本語

**🚀 [ランディングページでクイックスタート](https://jenv-win.vercel.app)** - プラットフォーム自動検出（Windows/Linux/macOS）によるワンクリックダウンロード • 多言語対応（EN/中文/日本語） • 隔週更新

</div>

## 最新のアップデート (v0.6.7)

### 🚀 パフォーマンスの向上
- **超高速 JDK スキャン**: スキャン時間を 3 秒から 300ms に短縮（90% の向上）。[詳細はこちら](doc/PERFORMANCE_jp.md)
- **並列処理**: Goroutine を使用した Dispatcher-Worker モデルの実装
- **スマートフィルタリング**: 不要なディレクトリをスキップする積極的なプリフィルタリング
- **進捗追跡**: リアルタイムのスキャン進捗と詳細な統計表示

### ✅ クロスプラットフォーム対応
- **macOS 対応完了**: macOS（Intel/Apple Silicon）への完全対応
- **Linux 対応完了**: マルチシェル対応による完全な Linux 互換性
- **Windows 最適化**: パス検証の強化と互換性の修正

### 🔧 技術的な強化
- **Java パス検証**: Windows での JDK 検出の信頼性を向上
- **環境管理**: クロスプラットフォームでの環境変数処理を最適化
- **設定のクリーンアップ**: 未使用オプションの削除とコードの保守性向上

---

## 概要

`Jenv` は、システム上の複数の Java バージョンを管理するためのコマンドラインツールです。異なる Java バージョンを簡単に切り替え、新しい Java インストールを追加し、Java 環境を効率的に管理できます。

## 特徴

### 効率的な Java バージョン管理

- **シンボリックリンクベースのアーキテクチャ**
    - シンボリックリンクによる高速なバージョン切り替え
    - 1 回限りのシステム PATH 設定
    - システム再起動後も変更が持続
    - すべてのコンソールウィンドウで即座に反映

### クロスプラットフォーム対応

- **Windows 対応**
    - レジストリベースの環境変数管理（Windows 標準）
    - 管理者権限の自動処理
    - 最小権限の原則による UAC プロンプトの最小化
    - Windows 10/11 システムでの優れたパフォーマンス

- **Linux/Unix 対応**
    - シェル設定ファイルベースの環境管理
    - ユーザーレベルおよびシステムレベルの設定オプション
    - マルチシェル環境（bash, zsh, fish）のサポート
    - インテリジェントな権限処理

- **macOS 対応**
    - Intel および Apple Silicon アーキテクチャのサポート
    - macOS 標準の JDK 配置場所に対応

### モダンな CLI 体験

- **直感的なインターフェース**
    - わかりやすいコマンド構造
    - ライト/ダークテーマのサポート
    - 読みやすさを向上させるカラー出力
    - 詳細なヘルプドキュメント

### 高度な機能

- **スマートな JDK 管理**
    - システム全体での JDK スキャン
    - **超高速スキャン（3s → 300ms）**: Dispatcher-Worker モデルによる並列処理
    - エイリアスベースの JDK 管理
    - 現在の JDK ステータス追跡
    - 簡単な JDK 追加と削除

---

## インストール

### リリースから
[Releases ページ](https://github.com/WhyWhatHow/jenv/releases)から最新リリースをダウンロードしてください。

### ソースからビルド

#### 前提条件

- Go 1.21 以上
- Git
- **Windows**: システムシンボリックリンクの作成に管理者権限が必要

#### ビルド手順

1. リポジトリをクローン:
```bash
git clone https://github.com/WhyWhatHow/jenv.git
cd jenv
```

2. プロジェクトをビルド:

```bash
cd src
# Windows (PowerShell)
go build -o jenv.exe
# Linux/macOS
go build -o jenv
```

## 使い方

### 初期設定

```bash
# jenv の初期化（初回使用時に必要）
jenv init

# jenv をシステム PATH に追加
jenv add-to-path
```

### JDK の追加と削除

```bash
# エイリアス名で新しい JDK を追加
jenv add jdk8 "C:\Program Files\Java\jdk1.8.0_291"
# JDK の削除
jenv remove jdk8
```

### JDK の切り替え

```bash
# インストール済み JDK の一覧表示
jenv list

# 特定の JDK バージョンに切り替え
jenv use jdk8

# 現在使用中の JDK を表示
jenv current
```

### システムのスキャン

```bash
# インストール済み JDK を自動検出
jenv scan c:\  # Windows
jenv scan /usr/lib/jvm  # Linux
```

---

## なぜこのプロジェクトが作られたのか？

Linux や macOS のユーザーには `sdkman` やオリジナルの `jenv`（bash ベース）のような成熟した Java バージョン管理ツールがありますが、Windows ユーザーにはこれまで最適とは言えない選択肢しかありませんでした。既存の [JEnv-for-Windows](https://github.com/FelixSelter/JEnv-for-Windows) などのソリューションは、現代の Windows システムでパフォーマンスのボトルネックに直面することがあります。

このプロジェクトは、主に 2 つの動機から生まれました：

1.  **Windows のギャップを埋める**: Windows 開発者に、Unix 系システムと同等（またはそれ以上）の堅牢で高性能な Java バージョン管理ツールを提供すること。
2.  **パフォーマンス重視**: システムの規模や複雑に関わらず、JDK のスキャンと切り替えをほぼ瞬時に行うこと。

私たちの目標は、**Windows における事実上の標準 Java 環境マネージャー**になると同時に、Windows、Linux、macOS をまたいで作業する開発者にシームレスで統一された体験を提供することです。

## ライセンス

このプロジェクトは Apache License 2.0 の下でライセンスされています。詳細は [LICENSE](LICENSE) ファイルをご覧ください。
