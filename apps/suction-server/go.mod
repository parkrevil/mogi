module suction-server

go 1.25.0

require (
	github.com/quic-go/quic-go v0.54.0
	go-shared/common v0.0.0
	go-shared/pb v0.0.0
	go.uber.org/fx v1.24.0
	go.uber.org/zap v1.27.0
	google.golang.org/protobuf v1.33.0
)

require (
	github.com/joho/godotenv v1.5.1 // indirect
	go.uber.org/dig v1.19.0 // indirect
	go.uber.org/mock v0.5.2 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.41.0 // indirect
	golang.org/x/mod v0.27.0 // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/tools v0.36.0 // indirect
)

replace go-shared/common => ../../packages/go-shared/common

replace go-shared/pb => ../../packages/go-shared/pb
