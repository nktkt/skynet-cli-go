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
- `dispatch`: ミッション実行シミュレーション
- `gameplan`: ゲーム理論ベースの防衛配分案を計算（`-json` 対応）
- `wargame`: 攻撃を確率サンプリングして複数ラウンドの損失を試算
- `report`: ミッション実績の集計（成功率・平均リスク・資源損耗）
- `status`: 現在状態を表示

## State File

デフォルト: `.skynet/state.json`

環境変数 `SKYNET_HOME` で保存先ディレクトリを変更できます。

## Development

```bash
go test ./...
go build ./...
```
