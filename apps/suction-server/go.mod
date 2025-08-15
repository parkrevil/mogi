module suction-server

go 1.25.0

require go-shared/common v0.0.0

require (
	go.uber.org/dig v1.19.0 // indirect
	go.uber.org/fx v1.24.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
)

replace go-shared/common => ../../packages/go-shared/common
