package logger

import (
	"context"
	"fmt"
	"strings"

	"github.com/golang/glog"
	"github.com/openshift-online/rh-trex/pkg/util"
)

type OCMLogger interface {
	V(level int32) OCMLogger
	Infof(format string, args ...interface{})
	Extra(key string, value interface{}) OCMLogger
	Info(message string)
	Warning(message string)
	Error(message string)
	Fatal(message string)
}

var _ OCMLogger = &logger{}

type extra map[string]interface{}

type logger struct {
	context   context.Context
	level     int32
	accountID string
	// TODO username is unused, should we be logging it? Could be pii
	username string
	extra    extra
}

// NewOCMLogger creates a new logger instance with a default verbosity of 1
func NewOCMLogger(ctx context.Context) OCMLogger {
	logger := &logger{
		context:   ctx,
		level:     1,
		extra:     make(extra),
		accountID: util.GetAccountIDFromContext(ctx),
	}
	return logger
}

func (l *logger) prepareLogPrefix(message string, extra extra) string {
	prefix := " "

	if txid, ok := l.context.Value("txid").(int64); ok {
		prefix = fmt.Sprintf("[tx_id=%d]%s", txid, prefix)
	}

	if l.accountID != "" {
		prefix = fmt.Sprintf("[accountID=%s]%s", l.accountID, prefix)
	}

	if opid, ok := l.context.Value(OpIDKey).(string); ok {
		prefix = fmt.Sprintf("[opid=%s]%s", opid, prefix)
	}

	var args []string
	for k, v := range extra {
		args = append(args, fmt.Sprintf("%s=%v", k, v))
	}

	return fmt.Sprintf("%s %s %s", prefix, message, strings.Join(args, " "))
}

func (l *logger) prepareLogPrefixf(format string, args ...interface{}) string {
	orig := fmt.Sprintf(format, args...)
	prefix := " "

	if txid, ok := l.context.Value("txid").(int64); ok {
		prefix = fmt.Sprintf("[tx_id=%d]%s", txid, prefix)
	}

	if l.accountID != "" {
		prefix = fmt.Sprintf("[accountID=%s]%s", l.accountID, prefix)
	}

	if opid, ok := l.context.Value(OpIDKey).(string); ok {
		prefix = fmt.Sprintf("[opid=%s]%s", opid, prefix)
	}

	return fmt.Sprintf("%s%s", prefix, orig)
}

func (l *logger) V(level int32) OCMLogger {
	return &logger{
		context:   l.context,
		accountID: l.accountID,
		username:  l.username,
		level:     level,
	}
}

func (l *logger) Infof(format string, args ...interface{}) {
	prefixed := l.prepareLogPrefixf(format, args...)
	glog.V(glog.Level(l.level)).Infof("%s", prefixed)
}

func (l *logger) Extra(key string, value interface{}) OCMLogger {
	l.extra[key] = value
	return l
}

func (l *logger) Info(message string) {
	l.log(message, glog.V(glog.Level(l.level)).Infoln)
}

func (l *logger) Warning(message string) {
	l.log(message, glog.Warningln)
}

func (l *logger) Error(message string) {
	l.log(message, glog.Errorln)
}

func (l *logger) Fatal(message string) {
	l.log(message, glog.Fatalln)
}

func (l *logger) log(message string, glogFunc func(args ...interface{})) {
	prefixed := l.prepareLogPrefix(message, l.extra)
	glogFunc(prefixed)
}
