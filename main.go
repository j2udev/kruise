package main

import (
	"github.com/j2udev/kruise/cmd"
	"github.com/j2udev/kruise/internal/kruise"
	"github.com/spf13/cobra"
)

func main() {
	kruise.Initialize()
	cobra.CheckErr(cmd.NewKruiseCmd().Execute())
}
