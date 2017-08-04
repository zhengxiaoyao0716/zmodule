package zmodule

import (
	"fmt"
	"testing"
)

func TestParseCmd(t *testing.T) {
	fmt.Println("\n================================")
	parseCmd("test", []string{"-h"}, Cmds)

	fmt.Println("\n================================")
	parseCmd("test", []string{"version"}, Cmds)

	fmt.Println("\n================================")
	parseCmd("test", []string{"service", "-h"}, Cmds)

	fmt.Println("\n================================")
}
