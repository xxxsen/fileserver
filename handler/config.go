package handler

type config struct {
	users             map[string]string
	maxUploadThread   int
	maxDownloadThread int
}

type Option func(c *config)

func WithUser(u, p string) Option {
	return func(c *config) {
		c.users[u] = p
	}
}

func WithUsers(us map[string]string) Option {
	return func(c *config) {
		for u, p := range us {
			c.users[u] = p
		}
	}
}

func WithMaxUploadThread(cnt int) Option {
	return func(c *config) {
		c.maxUploadThread = cnt
	}
}

func WithMaxDownloadThread(cnt int) Option {
	return func(c *config) {
		c.maxDownloadThread = cnt
	}
}
