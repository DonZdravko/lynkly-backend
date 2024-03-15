package logging

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"lynkly-backend/internal/common"
	"net/http"
)

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Panic(args ...interface{})
	WithFields(fields logrus.Fields) *logrus.Entry
	WithRequest(r *http.Request) *logrus.Entry
	SetLevel(level logrus.Level)
}

const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel logrus.Level = iota
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel
)

type lynklyLogger struct {
	*logrus.Entry
}

func NewLogger(packageName string) Logger {
	log := logrus.New()
	// Configure logger (optional)
	return &lynklyLogger{log.WithFields(logrus.Fields{"package": packageName})}
}

func (l *lynklyLogger) Debug(args ...interface{}) {
	l.Entry.Debug(args...)
}

func (l *lynklyLogger) Info(args ...interface{}) {
	l.Entry.Info(args...)
}

func (l *lynklyLogger) Warn(args ...interface{}) {
	l.Entry.Warn(args...)
}

func (l *lynklyLogger) Error(args ...interface{}) {
	l.Entry.Error(args...)
}

func (l *lynklyLogger) Panic(args ...interface{}) {
	l.Entry.Panic(args...)
}

func (l *lynklyLogger) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.Entry.WithFields(fields)
}

func (l *lynklyLogger) WithRequest(r *http.Request) *logrus.Entry {
	return l.WithFields(requestFields(r))
}

func (l *lynklyLogger) SetLevel(level logrus.Level) {
	l.Entry.Logger.SetLevel(level)
}

const (
	AccountIDName = "accountID"
	ClientIDName  = "clientID"
	SubKey        = "sub"
	UserKey       = "user"
)

func requestFields(r *http.Request) logrus.Fields {
	accountID := "N/A"
	clientID := "N/A"
	currentUser := common.ContextGet(r, UserKey)
	if currentUser != nil {
		claims := currentUser.(*jwt.Token).Claims.(jwt.MapClaims)
		clientID = claims[SubKey].(string)
	}

	return logrus.Fields{
		AccountIDName: accountID,
		ClientIDName:  clientID}
}
