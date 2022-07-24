package kruise

func Trace(err error) {
	if err != nil {
		Logger.Trace(err)
	}
}

func Tracef(err error, msg string, args ...interface{}) {
	if err != nil {
		Logger.Tracef(msg, args...)
		Logger.Trace(err)
	}
}

func Debug(err error) {
	if err != nil {
		Logger.Debug(err)
	}
}

func Debugf(err error, msg string, args ...interface{}) {
	if err != nil {
		Logger.Debugf(msg, args...)
		Logger.Debug(err)
	}
}

func Info(err error) {
	if err != nil {
		Logger.Info(err)
	}
}

func Infof(err error, msg string, args ...interface{}) {
	if err != nil {
		Logger.Infof(msg, args...)
		Logger.Info(err)
	}
}

func Warn(err error) {
	if err != nil {
		Logger.Warn(err)
	}
}

func Warnf(err error, msg string, args ...interface{}) {
	if err != nil {
		Logger.Warnf(msg, args...)
		Logger.Warn(err)
	}
}

func Error(err error) {
	if err != nil {
		Logger.Error(err)
	}
}

func Errorf(err error, msg string, args ...interface{}) {
	if err != nil {
		Logger.Errorf(msg, args...)
		Logger.Error(err)
	}
}

func Fatal(err error) {
	if err != nil {
		Logger.Fatal(err)
	}
}

func Fatalf(err error, msg string, args ...interface{}) {
	if err != nil {
		Logger.Errorf(msg, args...)
		Logger.Fatal(err)
	}
}

func Panic(err error) {
	if err != nil {
		Logger.Panic(err)
	}
}

func Panicf(err error, msg string, args ...interface{}) {
	if err != nil {
		Logger.Errorf(msg, args...)
		Logger.Panic(err)
	}
}
