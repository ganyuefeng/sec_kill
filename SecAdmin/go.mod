module SecAdmin

go 1.13

require (
	github.com/astaxie/beego v1.12.3
	github.com/coreos/etcd v3.3.25+incompatible
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/go-sql-driver/mysql v1.5.0
	github.com/google/uuid v1.1.2 // indirect
	github.com/jmoiron/sqlx v1.2.0
	go.uber.org/zap v1.16.0 // indirect
	google.golang.org/genproto v0.0.0-20201211151036-40ec1c210f7a // indirect
	google.golang.org/grpc v1.34.0 // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
