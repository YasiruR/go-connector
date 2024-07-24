package errors

import "fmt"

// Pkg internal errors should not be defined as they may subject to variations
// depending on the underlying implementation

const (
	urnFailed = "URN error (context: %s, function: %s) - %s"
	pkgFailed = "package error (name: %s, function: %s) - %s"
)

func URNFailed(ctx, function string, err error) error {
	return fmt.Errorf(urnFailed, ctx, function, err)
}

func PkgFailed(name, function string, err error) error {
	return fmt.Errorf(pkgFailed, name, function, err)
}
