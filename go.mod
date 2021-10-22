module github.com/go-spring/starter-go-redis

go 1.14

require (
	github.com/go-redis/redis/v8 v8.11.4
	github.com/go-spring/spring-base v1.1.0-beta.0.20211022224302-dea9f41f5f6d
	github.com/go-spring/spring-core v1.0.6-0.20211022224649-f0f6fffd8bc2
	github.com/go-spring/spring-go-redis v0.0.0-20211022225754-689e8d2dd56d
	github.com/go-spring/starter-core v1.1.0-beta.0.20211022230035-68dc9bcad473
)

//replace (
//	github.com/go-spring/spring-core => ../../spring/spring-core
//	github.com/go-spring/starter-core => ../starter-core
//)
