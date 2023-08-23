package config

import (
	"errors"
	"github.com/ethereum-optimism/superchain-registry/superchain"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	l1Networks      = []Chain{}
	defaultNetworks = []Network{}
)

type Chain struct {
	ChainId       uint64
	Name          string
	RPCUrl        string
	BlockExplorer string
}

type AddressType int

const ( // iota is reset to 0
	EOA AddressType = iota
	Implementation
	Proxy
)

type L1Addresses struct {
	AddressManager                    common.Address
	L1CrossDomainMessengerProxy       common.Address
	L1ERC721BridgeProxy               common.Address
	L1StandardBridgeProxy             common.Address
	L2OutputOracleProxy               common.Address
	OptimismMintableERC20FactoryProxy common.Address
	OptimismPortalProxy               common.Address
	ProxyAdmin                        common.Address
	SystemConfig                      common.Address
	BatchInbox                        common.Address
}

type OpStackChain struct {
	Chain
	L1ChainId    uint64
	SequencerUrl string
	L1Addresses  L1Addresses
}

type Network struct {
	Name   string
	Chains []OpStackChain
}

func GetNetworks() []Network {
	return defaultNetworks
}

func GetChains() []OpStackChain {
	result := []OpStackChain{}
	for _, network := range GetNetworks() {
		for _, chain := range network.Chains {
			result = append(result, chain)
		}
	}
	return result
}

func GetChainById(chainId uint64) (OpStackChain, error) {
	for _, chain := range GetChains() {
		if chain.ChainId == chainId {
			return chain, nil
		}
	}

	return OpStackChain{}, errors.New("unknown chain id")
}

func init() {
	for _, schain := range superchain.Superchains {
		l1 := schain.Config.L1
		l1Networks = append(l1Networks, Chain{
			ChainId:       l1.ChainID,
			Name:          schain.Superchain,
			RPCUrl:        l1.PublicRPC,
			BlockExplorer: l1.Explorer,
		})

		network := Network{
			Name:   schain.Superchain,
			Chains: []OpStackChain{},
		}

		for _, chainId := range schain.ChainIDs {
			rollup := superchain.OPChains[chainId]
			addresses := superchain.Addresses[chainId]

			network.Chains = append(network.Chains, OpStackChain{
				Chain: Chain{
					ChainId:       rollup.ChainID,
					Name:          rollup.Name,
					RPCUrl:        rollup.PublicRPC,
					BlockExplorer: rollup.Explorer,
				},
				SequencerUrl: rollup.SequencerRPC,
				L1ChainId:    l1.ChainID,
				L1Addresses: L1Addresses{
					AddressManager:                    common.Address(addresses.AddressManager),
					L1CrossDomainMessengerProxy:       common.Address(addresses.L1CrossDomainMessengerProxy),
					L1ERC721BridgeProxy:               common.Address(addresses.L1ERC721BridgeProxy),
					L1StandardBridgeProxy:             common.Address(addresses.L1StandardBridgeProxy),
					L2OutputOracleProxy:               common.Address(addresses.L2OutputOracleProxy),
					OptimismMintableERC20FactoryProxy: common.Address(addresses.OptimismMintableERC20FactoryProxy),
					OptimismPortalProxy:               common.Address(addresses.OptimismPortalProxy),
					ProxyAdmin:                        common.Address(addresses.ProxyAdmin),
					SystemConfig:                      common.Address(rollup.SystemConfigAddr),
					BatchInbox:                        common.Address(rollup.BatchInboxAddr),
				},
			})
		}

		defaultNetworks = append(defaultNetworks, network)
	}
}

func GetRPC(chain OpStackChain) (*ethclient.Client, error) {
	return ethclient.Dial(chain.RPCUrl)
}

func GetL1RPC(chain OpStackChain) (*ethclient.Client, error) {
	for _, eth := range l1Networks {
		if eth.ChainId == chain.L1ChainId {
			return ethclient.Dial(eth.RPCUrl)
		}
	}
	return nil, errors.New("unable to resolve L1 RPC")
}
