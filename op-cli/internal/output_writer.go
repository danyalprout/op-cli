package internal

import (
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"op-cli/op-cli/config"
	"syscall"
)

type Context struct {
	cmd *cobra.Command
	// Env vars
}

func NewContext(command *cobra.Command) *Context {
	return &Context{
		cmd: command,
	}
}

const jsonFormat = "json"
const tableFormat = "table"

var outputFormats = []string{jsonFormat, tableFormat}

func AddChainIdFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().Uint64P("chain-id", "c", 0, "Add chain id")
}

func (c *Context) GetRollupOrFail() config.OpStackChain {
	if c.HasChainId() {
		rollup, err := config.GetChainById(c.RequireChainId())
		if err != nil {
			c.cmd.PrintErrf("unknown rollup id")
			syscall.Exit(1)
		}
		return rollup
	}

	c.cmd.PrintErrf("must select a rollup")
	syscall.Exit(1)
	return config.OpStackChain{}
}

func (c *Context) GetSelectedRollupsOrAll() []config.OpStackChain {
	if c.HasChainId() {
		rollup, err := config.GetChainById(c.RequireChainId())
		if err != nil {
			c.cmd.PrintErrf("unknown rollup id")
			syscall.Exit(1)
		}
		return []config.OpStackChain{rollup}
	} else {
		result := []config.OpStackChain{}
		networks := config.GetNetworks()
		for _, network := range networks {
			for _, chain := range network.Chains {
				result = append(result, chain)
			}
		}
		return result
	}

}

func (c *Context) RequireChainId() uint64 {
	chainId, err := c.cmd.Flags().GetUint64("chain-id")

	if err != nil || chainId == 0 {
		c.cmd.PrintErrf("must provide a chain-id via --chain-id or -c")
		syscall.Exit(1)
	}

	return chainId
}

func (c *Context) WriteOutput(jsonPrinter func() ([]byte, error), tablePrinter func(table *tablewriter.Table)) {
	requestedFormat := c.cmd.Flag("fmt").Value.String()

	if requestedFormat == jsonFormat {
		body, err := jsonPrinter()
		if err != nil {
			panic(err)
		}
		_, _ = c.cmd.OutOrStdout().Write(body)
	} else if requestedFormat == tableFormat {
		table := tablewriter.NewWriter(c.cmd.OutOrStdout())
		tablePrinter(table)
		table.Render()
	} else {
		c.cmd.PrintErrf("unsupported format %s, options are %v\n", requestedFormat, outputFormats)
		syscall.Exit(1)
	}
}

func (c *Context) FatalError(err error) {
	_, _ = c.cmd.ErrOrStderr().Write([]byte(err.Error()))
	syscall.Exit(1)
}

func (c *Context) HasChainId() bool {
	chainId, err := c.cmd.Flags().GetUint64("chain-id")
	return err == nil && chainId != 0
}
