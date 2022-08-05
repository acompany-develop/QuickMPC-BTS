# QuickMPC-BTS
[![Build BTS and Run Build and Test in Container](https://github.com/acompany-develop/QuickMPC-BTS/actions/workflows/main.yml/badge.svg)](https://github.com/acompany-develop/QuickMPC-BTS/actions/workflows/main.yml)

[QuickMPC](https://github.com/acompany-develop/QuickMPC-BTS)で使われるTripleを生成するサービス

## ローカルでの起動方法
`src/BeaverTripleService/`で以下のコマンドを実行
```sh
make run
```
これにより, 同一ホストネットワーク内であれば `127.0.0.1:64101`で接続可能
以下の様に`grpcurl`でCLIからリクエストを送ることも可能
```sh
$ grpcurl -plaintext -d '{"job_id": 1, "amount": 5}' 127.0.0.1:64101 enginetobts.EngineToBts/GetTriples
{
  "triples": [
    {
      "a": "8218791",
      "b": "9915790",
      "c": "9079842"
    },
    {
      "a": "4645217",
      "b": "-5036032",
      "c": "8470037"
    },
    {
      "a": "-426163",
      "b": "495462",
      "c": "9562034"
    },
    {
      "a": "-2297368",
      "b": "-2005170",
      "c": "-150099"
    },
    {
      "a": "5474059",
      "b": "-3305195",
      "c": "7914345"
    }
  ]
}
```

## テスト方法
`Test/`で以下のコマンドを実行
```sh
make test
```
特定のtestを指定して実行したい場合は以下のようにする
```sh
make test t=./BeaverTripleService/TripleGenerator
# Test/BeaverTripleService/TripleGenerator/ 直下のみのテストを実行したい場合
make test p=unit # `uint*test.sh`を実行したい場合
make test m=build # `*test.sh`のbuild処理のみ実行したい場合
make test m=run # `*test.sh`のrun処理のみ実行したい場合
```

## 開発方法
`src/BeaverTripleService/`で以下のコマンドを実行
```sh
make up-build
make upd-build # バックグラウンドで起動したい場合はこちら
```

その後, VSCodeの左タブから`Remote Explorer` > 上のトグルから`Containers`を選択 > `beavertripleservice`にカーソルを合わせる > 新規フォルダアイコンを選択 > 開く場所を選択してsrc_btsコンテナの中で開発が行えるようになる.

![image](https://user-images.githubusercontent.com/33140349/142567126-52b8e392-a81c-4630-bf6c-6f801653770a.png)

## Container Image

GitHub Packages のコンテナレジストリにイメージを用意している

| tag             | description                                                    |
|-----------------|----------------------------------------------------------------|
| stable          | 最新安定版のイメージ                                           |
| s${date}        | 日時: `${date}` に作成された安定版のイメージ                   |
| stable-alpine   | `stable` の軽量イメージ                                        |
| s${date}-alpine | `s${date}` の軽量イメージ                                      |
| ${id}           | GitHub Actions の実行 ID: `${id}` で作成される開発版のイメージ |
| ${id}-alpine    | `${id}` の軽量イメージ                                         |

Dockerfile で使用される各 build stage については以下のリンクを参照

[QuickMPC/Test/README.md#how-to-develop-docker-composeyml](https://github.com/acompany-develop/QuickMPC/blob/579ba7332caf75162a1ac6c425fb83c04655c095/Test/README.md#how-to-develop-docker-composeyml)
