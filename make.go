package log

import (
	"io"
	"os"
	"errors"
)

var (
	_foreground = false
	_level = LevelInfo
	_out io.Writer = os.Stdout
	_cache []*Logger = make([]*Logger,0,10)

	errInvalidFilename = errors.New("Invalid filename when foreground logger is disabled.")
)

// each package can assign this to a global var
func MakeLog(prefix string) *Logger {
	logger := New(_out, "[" + prefix + "] " , LstdFlags ,_level)
	_cache = append(_cache, logger)
	return logger
}


// after checkout configuration file, using this to set each logger that created using MakeLog
// if foreground is true, filename, maxsize, maxcount will be ignored
// if foreground is false, filename, maxsize, maxcount will be used to rotate
func InitLog(foreground bool, level int, filename string, maxSize,maxCount int) error {

	_foreground = foreground
	_level = level
	for _,l := range _cache {
		l.SetLevel(_level)
	}

	if foreground {
		return nil
	}

	if filename == "" {
		return errInvalidFilename
	}

	_maxSize := 10*1024*1024 //default 10M
	_maxCount := 5

	if maxSize > 0 {
		_maxSize = maxSize
	}

	if maxCount > 0 {
		_maxCount = maxCount
	}

	io, err := Open(filename,0,int64(_maxSize),int64(_maxCount))

	if err != nil {
		return err
	}

	_out = io

	for _,l := range _cache {
		l.SetOutput(_out)
	}

	return nil
}