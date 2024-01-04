package handler

type config struct {
	users              map[string]string
	maxUploadThread    int
	maxDownloadThread  int
	enableFakeS3       bool
	fakeS3Buckets      []string
	enableRefererCheck bool
	referers           []string
	enableWebUI        bool
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

func WithEnableFakeS3(v bool) Option {
	return func(c *config) {
		c.enableFakeS3 = v
	}
}

func WithFakeS3BucketList(s []string) Option {
	return func(c *config) {
		c.fakeS3Buckets = s
	}
}

func WithEnableRefererCheck(v bool) Option {
	return func(c *config) {
		c.enableRefererCheck = v
	}
}

func WithRefererList(v []string) Option {
	return func(c *config) {
		c.referers = v
	}
}

func WithEnableWebUI(v bool) Option {
	return func(c *config) {
		c.enableWebUI = v
	}
}
