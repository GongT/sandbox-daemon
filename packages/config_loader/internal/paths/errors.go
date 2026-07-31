package paths

import (
	"gitlab.com/tozd/go/errors"
)

func s(paths ConfigPath) string {
	return paths.JoinWithDot("<root>")
}

type ErrorWithPath struct {
	Err     error
	CfgPath ConfigPath
	GoPath  ConfigPath
}

func ErrPathF(path TwoPath, msg string, args ...any) error {
	err := errors.Errorf(msg, args...)

	return &ErrorWithPath{
		Err:     errors.WithStack(err),
		CfgPath: path.Tags,
		GoPath:  path.Golang,
	}
}

func ErrPathW(path TwoPath, err error, msg string, args ...any) error {
	if msg == "" && len(args) == 0 {
		return ErrPath(path, err)
	}

	if e, dup := err.(*ErrorWithPath); dup {
		e.Err = errors.WithMessagef(e.Err, msg, args...)
		return e
	}

	err = errors.WithMessagef(err, msg, args...)

	return &ErrorWithPath{
		Err:     errors.WithStack(err),
		CfgPath: path.Tags,
		GoPath:  path.Golang,
	}
}

func ErrPath(path TwoPath, err error) error {
	if err == nil {
		return nil
	}
	if _, dup := err.(*ErrorWithPath); dup {
		return err
	}

	return &ErrorWithPath{
		Err:     errors.WithStack(err),
		CfgPath: path.Tags,
		GoPath:  path.Golang,
	}
}

func (e *ErrorWithPath) Error() string {
	return e.Unwrap().Error()
}

func (e *ErrorWithPath) Unwrap() error {
	return errors.WithMessagef(e.Err, "配置路径\"%v\"映射为\"%v\"", s(e.CfgPath), s(e.GoPath))
}
