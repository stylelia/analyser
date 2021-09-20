package redis

// TODO: Keys will be a composite - read Redis docs!
// Redis setup for now
type Redis struct{}

func (r *Redis) GetKey(key string) (string, error) {
	return key, nil
}

func (r *Redis) UpdateKey(key, value string) error {
	return nil
}
