package tracing

import (
	"telegram-bot/internal/clients/tg"
	"telegram-bot/internal/logger"
	"telegram-bot/internal/model/messages"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

func init() {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
	}

	_, err := cfg.InitGlobalTracer("tg")
	if err != nil {
		logger.Fatal("Cannot init tracing", "err", err)
	}
}

// TracingMiddleware Функция трейсинга.
func TracingMiddleware(next tg.HandlerFunc) tg.HandlerFunc {
	handler := tg.HandlerFunc(
		func(tgUpdate tgbotapi.Update, c *tg.Client, msgModel *messages.Model) {
			span, ctx := opentracing.StartSpanFromContext(msgModel.GetCtx(), "ProcessingMessages")
			defer span.Finish()
			if spanContext, ok := span.Context().(jaeger.SpanContext); ok {
				logger.Info("start span trace", "traceId", spanContext.TraceID().String())
			}
			// Выполнение процесса обработки сообщения.
			msgModel.SetCtx(ctx)
			next.RunFunc(tgUpdate, c, msgModel)

			ext.SpanKindRPCClient.Set(span)
		},
	)

	return handler
}
