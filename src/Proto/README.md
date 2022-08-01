## protocからの生成コマンド
```
cd /Proto
protoc --go_out=/Proto --go_opt=paths=source_relative \
    --go-grpc_out=/Proto --go-grpc_opt=paths=source_relative \
    ./EngineToBts/engine_to_bts.proto
```