package infomaniakai

import (
	"os"
	"strconv"
)

const (
	envProductID = "INFOMANIAK_PRODUCT_ID"
	envApiKey    = "INFOMANIAK_API_KEY"
)

type config struct {
	productId int
	key       string
}

type Option interface {
	apply(*config)
}

type optionFunc func(c *config)

func (o optionFunc) apply(c *config) {
	o(c)
}

func WithApiToken(key string) Option {
	return optionFunc(func(c *config) {
		c.key = key
	})
}

func WithProductID(id int) Option {
	return optionFunc(func(c *config) {
		if id >= 0 {
			c.productId = id
		}
	})
}

func defaultConfig() *config {
	productID, _ := strconv.Atoi(os.Getenv(envProductID))
	return &config{
		productId: productID,
		key:       os.Getenv(envApiKey),
	}
}
