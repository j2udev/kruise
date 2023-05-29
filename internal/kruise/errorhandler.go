package kruise

// Debug is used to catch non-nil errors and Log them at a debug level
func Debug(err error) {
	if err != nil {
		Logger.Debug(err)
	}
}

// Debugf is used to catch non-nil errors and Log them at a debug level with a
// custom message
func Debugf(err error, msg string, args ...interface{}) {
	if err != nil {
		Logger.Debugf(msg, args...)
		Logger.Debug(err)
	}
}

// Info is used to catch non-nil errors and Log them at a info level
func Info(err error) {
	if err != nil {
		Logger.Info(err)
	}
}

// Infof is used to catch non-nil errors and Log them at a info level with a
// custom message
func Infof(err error, msg string, args ...interface{}) {
	if err != nil {
		Logger.Infof(msg, args...)
		Logger.Info(err)
	}
}

// Warn is used to catch non-nil errors and Log them at a warn level
func Warn(err error) {
	if err != nil {
		Logger.Warn(err)
	}
}

// Warnf is used to catch non-nil errors and Log them at a warn level with a
// custom message
func Warnf(err error, msg string, args ...interface{}) {
	if err != nil {
		Logger.Warnf(msg, args...)
		Logger.Warn(err)
	}
}

// Error is used to catch non-nil errors and Log them at a error level
func Error(err error) {
	if err != nil {
		Logger.Error(err)
	}
}

// Errorf is used to catch non-nil errors and Log them at a error level with a
// custom message
func Errorf(err error, msg string, args ...interface{}) {
	if err != nil {
		Logger.Errorf(msg, args...)
		Logger.Error(err)
	}
}

// Fatal is used to catch non-nil errors and Log them at a fatal level
//
// This bails out of execution
func Fatal(err error) {
	if err != nil {
		Logger.Fatal(err)
	}
}

// Fatalf is used to catch non-nil errors and Log them at a fatal level with a
// custom message
//
// This bails out of execution
func Fatalf(err error, msg string, args ...interface{}) {
	if err != nil {
		Logger.Errorf(msg, args...)
		Logger.Fatal(err)
	}
}
