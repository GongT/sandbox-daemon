package config

import (
	"gitlab.com/tozd/go/errors"
)

func s(paths ConfigPath) string {
	return paths.JoinWithDot("<root>")
}

type errorWithPath struct {
	Err     error
	CfgPath ConfigPath
}

func errPathF(path ConfigPath, msg string, args ...any) error {
	err := errors.Errorf(msg, args...)

	return &errorWithPath{
		Err:     errors.WithStack(err),
		CfgPath: path,
	}
}

func errPathW(path ConfigPath, err error, msg string, args ...any) error {
	err = errors.WithMessagef(err, msg, args...)

	return &errorWithPath{
		Err:     errors.WithStack(err),
		CfgPath: path,
	}
}

func errPath(path ConfigPath, err error) error {
	if err == nil {
		return nil
	}
	if _, dup := err.(*errorWithPath); dup {
		return err
	}

	return &errorWithPath{
		Err:     errors.WithStack(err),
		CfgPath: path,
	}
}

func (e *errorWithPath) Error() string {
	return e.Unwrap().Error()
}

func (e *errorWithPath) Unwrap() error {
	return errors.WithMessagef(e.Err, "配置路径\"%v\"", s(e.CfgPath))
}
