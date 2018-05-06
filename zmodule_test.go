package zmodule

import (
	"fmt"
	"testing"
)

func TestParseCmd(t *testing.T) {
	fmt.Println("\n================================")
	ParseCmd(Cmds)("test", []string{"-h"})

	fmt.Println("\n================================")
	ParseCmd(Cmds)("test", []string{"version"})

	fmt.Println("\n================================")
	ParseCmd(Cmds)("test", []string{"service", "-h"})

	fmt.Println("\n================================")
}
