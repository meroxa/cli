package cmd

import "testing"

func TestListResourcesCmd(t *testing.T) {
	cmd := ListResourcesCmd()
	cmd.Execute()
}
