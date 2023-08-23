package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	cfg "op-cli/op-cli/config"
	"op-cli/op-cli/internal"
	"strconv"
)

func init() {
	config.AddCommand(networks)
	config.AddCommand(rollups)

	internal.AddChainIdFlag(addresses)
	config.AddCommand(addresses)

	rootCmd.AddCommand(config)
}

var config = &cobra.Command{
	Use: "config",
}

var networks = &cobra.Command{
	Use:   "networks",
	Short: "Print network configuration",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := internal.NewContext(cmd)

		networks := cfg.GetNetworks()

		ctx.WriteOutput(func() ([]byte, error) {
			return json.Marshal(networks)
		}, func(table *tablewriter.Table) {
			table.SetHeader([]string{"Name", "Chains"})
			for _, network := range networks {
				var chainSummary string
				for _, chain := range network.Chains {
					if chainSummary != "" {
						chainSummary += ", "
					}
					description := fmt.Sprintf("%s (%d)", chain.Name, chain.ChainId)
					chainSummary += description
				}
				table.Append([]string{network.Name, chainSummary})
			}
		})
	},
}

type ChainConfig struct {
	cfg.OpStackChain
	Network string
}

var rollups = &cobra.Command{
	Use:   "rollups",
	Short: "Print chain configuration",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := internal.NewContext(cmd)

		networks := cfg.GetNetworks()

		result := []ChainConfig{}

		for _, network := range networks {
			for _, chain := range network.Chains {
				result = append(result, ChainConfig{
					OpStackChain: chain,
					Network:      network.Name,
				})
			}
		}

		ctx.WriteOutput(func() ([]byte, error) {
			return json.Marshal(result)
		}, func(table *tablewriter.Table) {
			table.SetHeader([]string{"Network", "Name", "ID", "RPC"})
			for _, cfg := range result {
				table.Append([]string{cfg.Network, cfg.Name, strconv.Itoa(int(cfg.ChainId)), cfg.RPCUrl})
			}
		})
	},
}

type AddressInfo struct {
	Rollup      string
	L1Addresses cfg.L1Addresses
}

var addresses = &cobra.Command{
	Use:   "addresses",
	Short: "Print chain system addresses",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := internal.NewContext(cmd)
		rollups := ctx.GetSelectedRollupsOrAll()

		result := []AddressInfo{}

		for _, chain := range rollups {
			result = append(result, AddressInfo{
				Rollup:      chain.Name,
				L1Addresses: chain.L1Addresses,
			})
		}

		ctx.WriteOutput(func() ([]byte, error) {
			return json.Marshal(result)
		}, func(table *tablewriter.Table) {
			table.SetHeader([]string{"Network", "Address", "value"})

			for _, cfg := range result {
				table.Append([]string{
					cfg.Rollup,
					"Address Manager",
					cfg.L1Addresses.AddressManager.Hex(),
				})

				table.Append([]string{
					cfg.Rollup,
					"Cross Domain Messenger Proxy",
					cfg.L1Addresses.L1CrossDomainMessengerProxy.Hex(),
				})

				table.Append([]string{
					cfg.Rollup,
					"L1 ERC721 Bridge Proxy",
					cfg.L1Addresses.L1ERC721BridgeProxy.Hex(),
				})

				table.Append([]string{
					cfg.Rollup,
					"L1 Standard Bridge Proxy",
					cfg.L1Addresses.L1StandardBridgeProxy.Hex(),
				})

				table.Append([]string{
					cfg.Rollup,
					"L2 Output Oracle Proxy",
					cfg.L1Addresses.L2OutputOracleProxy.Hex(),
				})

				table.Append([]string{
					cfg.Rollup,
					"Mintable ERC20 Factory Proxy",
					cfg.L1Addresses.OptimismMintableERC20FactoryProxy.Hex(),
				})

				table.Append([]string{
					cfg.Rollup,
					"Portal Proxy",
					cfg.L1Addresses.OptimismPortalProxy.Hex(),
				})

				table.Append([]string{
					cfg.Rollup,
					"Proxy Admin",
					cfg.L1Addresses.ProxyAdmin.Hex(),
				})

				table.Append([]string{
					cfg.Rollup,
					"System Config",
					cfg.L1Addresses.SystemConfig.Hex(),
				})

				table.Append([]string{
					cfg.Rollup,
					"Batch Inbox",
					cfg.L1Addresses.BatchInbox.Hex(),
				})
			}
		})
	},
}
