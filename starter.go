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

package StarterGoRedis

import (
	"github.com/go-spring/spring-core/conf"
	"github.com/go-spring/spring-core/gs"
	"github.com/redis/go-redis/v9"
)

// Config defines Redis connection configuration.
type Config struct {
	Addr     string `value:"${addr}"`
	Password string `value:"${password:=}"`
}

// Factory defines an interface for creating Redis clients.
type Factory interface {
	CreateClient(c Config) (*redis.Client, error)
}

type DefaultFactory struct{}

// CreateClient creates a new Redis client based on the provided configuration.
func (DefaultFactory) CreateClient(c Config) (*redis.Client, error) {
	return redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		Password: c.Password,
	}), nil
}

func init() {
	const key = "spring.go-redis"

	// Register a module that initializes Redis clients
	gs.Module([]gs.ConditionOnProperty{
		gs.OnProperty(key),
	}, func(p conf.Properties) error {

		// Bind configuration into a map of name -> Config
		var m map[string]Config
		if err := p.Bind(&m, "${"+key+"}"); err != nil {
			return err
		}

		// Register DefaultFactory as a bean implementing Factory,
		// but only if no other Factory bean has been provided.
		gs.Object(&DefaultFactory{}).
			Condition(gs.OnMissingBean[Factory]()).
			Export(gs.As[Factory]())

		// For each Redis configuration entry,
		// create and register a Redis client bean.
		for name, c := range m {
			gs.Provide(func(factory Factory) (*redis.Client, error) { // create
				return factory.CreateClient(c)
			}).Destroy(func(client *redis.Client) error { // destroy
				return client.Close()
			}).Name(name)
		}
		return nil
	})
}
