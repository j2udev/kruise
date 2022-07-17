package main

import (
	"github.com/j2udevelopment/kruise/cmd"
	"github.com/j2udevelopment/kruise/pkg/kruise"
	"github.com/spf13/cobra"
)

func main() {
	kruise.Initialize()
	cobra.CheckErr(cmd.NewKruiseKmd().Execute())
}
