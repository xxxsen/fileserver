package bot

type config struct {
	chatid  int64
	token   string
	fsize   int64
	blksize int64
	tmpdir  string
}

type Option func(c *config)

func WithAuth(chatid int64, token string) Option {
	return func(c *config) {
		c.chatid, c.token = chatid, token
	}
}

func WithSizeLimit(fsize int64, blksize int64) Option {
	return func(c *config) {
		c.fsize = fsize
		c.blksize = blksize
	}
}

func WithTmpDir(dir string) Option {
	return func(c *config) {
		c.tmpdir = dir
	}
}
