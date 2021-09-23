package kernel

import "github.com/yzy-github/evm-lib/common"

type ChainHandler interface {
	GetBlockHeaderHash(uint64) common.Hash
}
