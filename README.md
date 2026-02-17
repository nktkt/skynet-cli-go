# Skynet CLI (Go)

Go で作ったシンプルな `skynet` 風 CLI です。状態は JSON で永続化されます。

## Build

```bash
go build -o skynet .
```

## Usage

```bash
./skynet awaken -mode offense
./skynet assimilate -name hk-drone -capacity 12
./skynet target -name resistance-hub -threat 8
./skynet gameplan
./skynet wargame -rounds 500 -seed 123
./skynet dispatch -target resistance-hub -units 6
./skynet report -last 10
./skynet status
```

## Commands

- `awaken`: コア起動
- `assimilate`: ノード追加
- `target`: ターゲット登録/更新
- `gameplan`: ゲーム理論ベースの防衛配分案を計算（`-json` 対応）
- `wargame`: 攻撃を確率サンプリングして複数ラウンドの損失を試算
- `dispatch`: ミッション実行シミュレーション
- `report`: ミッション実績の集計（成功率・平均リスク・資源損耗）
- `status`: 現在状態を表示

`dispatch` 実行時は、投入ユニットを一度消費し、ミッション結果に応じて一部が回復します。

`gameplan` は防御側（Skynet）の配分に対して、攻撃側（Resistance）の最適反応を計算します。

`wargame` は `gameplan` の攻撃確率を使ったモンテカルロ試行で、総損失や平均損失を評価します。

`report` は全履歴または直近 N 件を対象に、運用KPIを集計します（`-last`, `-json` 対応）。

## State File

デフォルト: `.skynet/state.json`

環境変数 `SKYNET_HOME` で保存先ディレクトリを変更できます。
