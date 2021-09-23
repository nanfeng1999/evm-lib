package kernel

import "evm/common"

type ChainHandler interface {
	GetBlockHeaderHash(uint64) common.Hash
}
