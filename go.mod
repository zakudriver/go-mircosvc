module github.com/Zhan9Yunhua/blog-svr

go 1.12

require (
	github.com/Zhan9Yunhua/logger v0.0.0-20190429041551-fbc2f63be669
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.3.0
	github.com/golang/protobuf v1.3.2
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/micro/go-api v0.7.0
	github.com/micro/go-micro v1.8.0
	github.com/nats-io/nats-server/v2 v2.0.2 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/yaml.v2 v2.2.2
)

replace github.com/hashicorp/consul => github.com/hashicorp/consul v1.5.1
