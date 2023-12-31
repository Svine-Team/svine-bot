//go:build mage

package main

import (
	"github.com/magefile/mage/sh"
)

func Build() error {
	return sh.Run("go", "build", "-o", "tmp", "./...")
}

func BuildLive() error {
	return sh.RunV("air", "*> ./tmp/run-verbose.log")
}
