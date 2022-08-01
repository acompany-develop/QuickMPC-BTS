module github.com/acompany-develop/QuickMPC-BTS/src/BeaverTripleService

go 1.14

replace github.com/acompany-develop/QuickMPC-BTS/src/Proto/EngineToBts => ./../Proto/EngineToBts

require (
	github.com/acompany-develop/QuickMPC-BTS/src/Proto/EngineToBts v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.21.0
	golang.org/x/tools/gopls v0.7.3 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/grpc v1.42.0
	google.golang.org/protobuf v1.27.1 // indirect
)
