## grpcurlでのdebug
※ portはよしなに変更してください
```bash
grpcurl -d '{"job_id": 1, "amount": 10}' beaver_triple_service:54100 enginetobts.EngineToBts/GetTriples
```

## grpcサーバのヘルスチェック
```bash
grpc_health_probe -addr=localhost:54100
```