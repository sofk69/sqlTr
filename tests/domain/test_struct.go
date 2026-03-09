package domain

import "testing"

type Test struct {
	name          string
	performAction func(*Test, string)
	verifyResult  func(*testing.T, *Test, string)
}
