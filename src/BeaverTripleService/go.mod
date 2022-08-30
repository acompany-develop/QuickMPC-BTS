module github.com/acompany-develop/QuickMPC-BTS/src/BeaverTripleService

go 1.14

replace github.com/acompany-develop/QuickMPC-BTS/src/Proto/EngineToBts => ./../Proto/EngineToBts

require (
	github.com/acompany-develop/QuickMPC-BTS/src/Proto/EngineToBts v0.0.0-00010101000000-000000000000
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	go.uber.org/zap v1.21.0
	golang.org/x/tools v0.1.8-0.20211014194737-fc98fb2abd48 // indirect
	google.golang.org/grpc v1.42.0
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)
