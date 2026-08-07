A.install swagger gen tool
  go install github.com/swaggo/swag/cmd/swag@latest

  gen swagger docs
  swag init -o ./docs -pdl 3

B.to get all command please run
go run main.go

C.to start/stop/restart the server please run
go run main.go start
go run main.go stop (not available in windows)
go run main.go restart (not available in windows)

D.to gen model and dal files please run
go run main.go gen