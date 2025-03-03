package cron

import (
	"context"
	"runtime/debug"

	"github.com/robfig/cron/v3"
	"github.com/xxxsen/common/logutil"
	"go.uber.org/zap"
)

type ICronJob interface {
	Name() string
	Expression() string
	Run(ctx context.Context) error
}

type Cron struct {
	cr *cron.Cron
}

func New() *Cron {
	cr := cron.New(cron.WithSeconds())
	return &Cron{
		cr: cr,
	}
}

func (c *Cron) Start() {
	c.cr.Start()
}

func (c *Cron) AddJob(impl ICronJob) error {
	expr := impl.Expression()
	job := c.jobWrap(impl)
	if _, err := c.cr.AddJob(expr, job); err != nil {
		return err
	}
	return nil
}

func (c *Cron) jobWrap(impl ICronJob) cron.FuncJob {
	ctx := context.Background()
	return func() {
		defer func() {
			if e := recover(); e != nil {
				logutil.GetLogger(ctx).Error("run job panic",
					zap.String("job", impl.Name()), zap.Any("panic", e), zap.String("stack", string(debug.Stack())))
			}
		}()
		if err := impl.Run(ctx); err != nil {
			logutil.GetLogger(ctx).Error("run job failed", zap.String("job", impl.Name()), zap.Error(err))
		}
	}
}
