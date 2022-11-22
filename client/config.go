package client

type config struct {
	addr string
	ak   string
	sk   string
}

type Option func(c *config)

func WithHost(addr string) Option {
	return func(c *config) {
		c.addr = addr
	}
}

func WithKey(ak, sk string) Option {
	return func(c *config) {
		c.ak, c.sk = ak, sk
	}
}
