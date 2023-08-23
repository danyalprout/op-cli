package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum-optimism/optimism/op-bindings/bindings"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	cfg "op-cli/op-cli/config"
	"op-cli/op-cli/internal"
)

func init() {
	internal.AddChainIdFlag(owners)
	rootCmd.AddCommand(owners)
}

type OwnerInfo struct {
	Name  string
	Owner string
	Note  string
}

func eip1967Owner(ctx *internal.Context, rpc *ethclient.Client, address common.Address) common.Address {
	eip1967Slot := common.HexToHash("0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103")
	data, err := rpc.StorageAt(context.TODO(), address, eip1967Slot, nil)
	if err != nil {
		ctx.FatalError(fmt.Errorf("failed to fetch data for %s", address.String()))
		return [20]byte{}
	}
	return common.Address(data[len(data)-20:])
}

func ownerMethod(ctx *internal.Context, rpc *ethclient.Client, address common.Address) common.Address {
	contract, err := bindings.NewAddressManager(address, rpc)
	if err != nil {
		ctx.FatalError(err)
		return common.Address{}
	}

	owner, err := contract.Owner(&bind.CallOpts{
		BlockNumber: nil,
	})
	if err != nil {
		ctx.FatalError(fmt.Errorf("failed to fetch data for %s", address.String()))
		return [20]byte{}
	}

	return owner
}

var owners = &cobra.Command{
	Use:   "owners",
	Short: "For each address print any ownership info",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := internal.NewContext(cmd)
		rollup := ctx.GetRollupOrFail()
		rpc, err := cfg.GetL1RPC(rollup)

		if err != nil {
			ctx.FatalError(err)
		}

		owners := []OwnerInfo{
			{
				"Address Manager",
				ownerMethod(ctx, rpc, rollup.L1Addresses.AddressManager).String(),
				"Ownable",
			},
			{
				"Proxy Admin",
				ownerMethod(ctx, rpc, rollup.L1Addresses.ProxyAdmin).String(),
				"Ownable",
			},
			{
				"L1 ERC721 Bridge Proxy",
				eip1967Owner(ctx, rpc, rollup.L1Addresses.L1ERC721BridgeProxy).String(),
				"EIP-1967",
			},
			{
				"Batch Inbox",
				"N/A",
				"EOA",
			},
		}

		ctx.WriteOutput(func() ([]byte, error) {
			return json.Marshal(owners)
		}, func(table *tablewriter.Table) {
			table.SetHeader([]string{"Name", "Owner", "Type"})
			for _, oi := range owners {
				table.Append([]string{oi.Name, oi.Owner, oi.Note})
			}
		})
	},
}
