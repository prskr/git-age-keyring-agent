package testx

import "github.com/stretchr/testify/assert"

type ValueAssertionFunc[T any] func(t assert.TestingT, val T, vals ...T) bool
