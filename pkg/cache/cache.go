package cache

import (
	"fmt"
)

// Type Key defines a valid key which can be fetched from Redis
// It is composed of the Org Name and the Repo Name
type Key struct {
	OrgName  string
	RepoName string
}

// Gets the Redis path for the key
func (k *Key) getKeyPath() string {
	// Hard coded github as this is the only source we support today
	return fmt.Sprintf("github/%v/%v", k.OrgName, k.RepoName)
}

// Interface for KV store
type Cache interface {
	GetKey(key Key) (string, error)
	UpdateKey(key Key, value string) error
	GetCookstyleVersion(key Key) (string, error)
	UpdateCookstyleVersion(key Key) (string, error)
}

// var ctx = context.Background()

// func ExampleClient() {
// 	rdb := redis.NewClient(&redis.Options{
// 		Addr:     "localhost:6379",
// 		Password: os.Getenv("REDIS_PASSWORD"),
// 		DB:       0, // use default DB
// 	})

// 	err := rdb.Set(ctx, "key", "value", 0).Err()
// 	if err != nil {
// 		t
// 		panic(err)
// 	}

// 	val, err := rdb.Get(ctx, "key").Result()
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("key", val)

// 	val2, err := rdb.Get(ctx, "key2").Result()
// 	if err == redis.Nil {
// 		fmt.Println("key2 does not exist")
// 	} else if err != nil {
// 		panic(err)
// 	} else {
// 		fmt.Println("key2", val2)
// 	}
// 	// Output: key value
// 	// key2 does not exist
// }
