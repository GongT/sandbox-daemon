//go:build !test_env

package test_handler

func NewTestForwarder() *testForwarder {
	return &testForwarder{}
}
