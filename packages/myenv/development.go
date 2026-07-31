//go:build !release

package myenv

import "testing"

const IsDebug = true
const IsRelease = false

var IsTesting = testing.Testing()
