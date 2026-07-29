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
	GoPath  ConfigPath
}

func errPathF(path twoPath, msg string, args ...any) error {
	err := errors.Errorf(msg, args...)

	return &errorWithPath{
		Err:     errors.WithStack(err),
		CfgPath: path.tags,
		GoPath:  path.golang,
	}
}

func errPathW(path twoPath, err error, msg string, args ...any) error {
	if msg == "" && len(args) == 0 {
		return errPath(path, err)
	}

	if e, dup := err.(*errorWithPath); dup {
		e.Err = errors.WithMessagef(e.Err, msg, args...)
		return e
	}

	err = errors.WithMessagef(err, msg, args...)

	return &errorWithPath{
		Err:     errors.WithStack(err),
		CfgPath: path.tags,
		GoPath:  path.golang,
	}
}

func errPath(path twoPath, err error) error {
	if err == nil {
		return nil
	}
	if _, dup := err.(*errorWithPath); dup {
		return err
	}

	return &errorWithPath{
		Err:     errors.WithStack(err),
		CfgPath: path.tags,
		GoPath:  path.golang,
	}
}

func (e *errorWithPath) Error() string {
	return e.Unwrap().Error()
}

func (e *errorWithPath) Unwrap() error {
	return errors.WithMessagef(e.Err, "配置路径\"%v\"映射为\"%v\"", s(e.CfgPath), s(e.GoPath))
}
