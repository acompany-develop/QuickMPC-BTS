## grpcurlでのdebug
※ portはよしなに変更してください
```bash
grpcurl -d '{"job_id": 1, "amount": 10}' beaver_triple_service:54100 enginetobts.EngineToBts/GetTriples
```

## grpcサーバのヘルスチェック
```bash
grpc_health_probe -addr=localhost:54100
```

## JWT token の生成

YAML ファイルを入力に JWT token を生成します

```console
root@container:/QuickMPC-BTS# go run Cmd/JWTGenerator/main.go     # generate from sample.yml
root@container:/QuickMPC-BTS# go run Cmd/JWTGenerator/main.go \
>                                 -file /path/to/config.yml   \
>                                 -o ./output/envs                # use own configuration
root@container:/QuickMPC-BTS# go run Cmd/JWTGenerator/main.go -h  # show help
```

クライアントとサーバ向けにそれぞれ `.env` ファイル形式の設定ファイルが書き込まれます
