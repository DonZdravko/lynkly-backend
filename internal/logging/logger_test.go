package logging

import "github.com/sirupsen/logrus"

// MockLogger defines a mockable interface for testing purposes
type MockLogger interface {
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	WithFields(fields logrus.Fields) *MockLogger
}
