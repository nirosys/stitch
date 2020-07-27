package eval

import (
	"github.com/nirosys/stitch/object"
)

// ObjectResolver /////////////////////////////////////////////////////////////
type ObjectResolver interface {
	Resolve(name string) (object.Object, error)
}
