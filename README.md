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
./skynet dispatch -target resistance-hub -units 6
./skynet status
```

## Commands

- `awaken`: コア起動
- `assimilate`: ノード追加
- `target`: ターゲット登録/更新
- `gameplan`: ゲーム理論ベースの防衛配分案を計算
- `dispatch`: ミッション実行シミュレーション
- `status`: 現在状態を表示

`dispatch` 実行時は、投入ユニットを一度消費し、ミッション結果に応じて一部が回復します。

`gameplan` は防御側（Skynet）の配分に対して、攻撃側（Resistance）の最適反応を計算します。

## State File

デフォルト: `.skynet/state.json`

環境変数 `SKYNET_HOME` で保存先ディレクトリを変更できます。
