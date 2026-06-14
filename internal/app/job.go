package app

import (
	"context"
	"errors"
	"reflect"
	"runtime"
	"time"

	"github.com/mbranch/safe-go"
	"github.com/rs/zerolog/log"
)

const (
	ctxTimeout   = 30 * time.Second
	timerTimeout = 5 * time.Second
)

func (pp *PostProcessor) Workers(ctx context.Context) {
	go pp.runJob(ctx, timerTimeout, pp.ProcessPosts)
	// go pp.runJob(ctx, timerTimeout, pp.pingCLient)

	<-pp.stopCh
}

func (pp *PostProcessor) runJob(
	ctx context.Context,
	timerDuration time.Duration,
	job func(context.Context) error,
) {
	timer := time.NewTimer(timerDuration)

	method := runtime.FuncForPC(reflect.ValueOf(job).Pointer())
	if method != nil {
		ctx = context.WithValue(ctx, "method", method.Name())
	}

mainloop:
	for {
		select {
		case <-pp.stopCh:
			break mainloop
		case <-ctx.Done():
			log.Err(ctx.Err()).Ctx(ctx).Str("method", method.Name()).Msg("context done")
			break mainloop
		case <-timer.C:
			log.Info().Ctx(ctx).Str("method", method.Name()).Msg("start job")
			func() {
				ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
				defer cancel()

				err := safe.Do(func() (err error) {
					return job(ctx)
				})
				if err != nil {
					var panicErr safe.PanicError
					if errors.As(err, &panicErr) {
						log.Err(err).Ctx(ctx).
							Str("method", method.Name()).
							Interface("panic_value", panicErr.Panic()).
							Interface("stack", panicErr.StackTrace()).
							Msg("job panicked")
					}
					log.Err(err).Ctx(ctx).Str("method", method.Name()).Msg("job failed")
				}
			}()

			timer.Reset(timerDuration)
		}
	}
}
