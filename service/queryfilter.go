package service

import (
	. "github.com/ccsdsmo/malgo/mal"
)

type QueryFilter interface {
	Composite
	QueryFilter() QueryFilter
}
