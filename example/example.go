/*
 * Copyright 2025 The Go-Spring Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/go-spring/spring-core/gs"
	"github.com/redis/go-redis/v9"

	StarterGoRedis "github.com/go-spring/starter-go-redis"
)

type Service struct {
	Redis *redis.Client `autowire:""`
}

// AnotherRedisFactory is a custom implementation of the Factory interface.
type AnotherRedisFactory struct{}

func (AnotherRedisFactory) CreateClient(c StarterGoRedis.Config) (*redis.Client, error) {
	return redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		Password: c.Password,
	}), nil
}

func main() {

	// Register a custom Factory bean to replace the default one.
	gs.Provide(func() StarterGoRedis.Factory {
		return &AnotherRedisFactory{}
	})

	// Here `s` is not referenced by any other object,
	// so we need to register it as a root object.
	s := &Service{}
	gs.Root(gs.Object(s))

	// Define a handler to GET a Redis key value.
	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		str, err := s.Redis.Get(r.Context(), "key").Result()
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		_, _ = w.Write([]byte(str))
	})

	// Define a handler to SET a Redis key value.
	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		str, err := s.Redis.Set(r.Context(), "key", "value", 0).Result()
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		_, _ = w.Write([]byte(str))
	})

	gs.Run()

	// Example usage:
	//
	// ~ curl http://127.0.0.1:9090/get
	// redis: nil%
	// ~ curl http://127.0.0.1:9090/set
	// OK%
	// ~ curl http://127.0.0.1:9090/get
	// value%
}

// ----------------------------------------------------------------------------
// Change working directory
// ----------------------------------------------------------------------------

// init sets the working directory of the application to the directory
// where this source file resides.
// This ensures that any relative file operations are based on the source file location,
// not the process launch path.
func init() {
	var execDir string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		execDir = filepath.Dir(filename)
	}
	err := os.Chdir(execDir)
	if err != nil {
		panic(err)
	}
	workDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Println(workDir)
}
