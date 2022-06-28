package kernel

import "github.com/nanfeng1999/evm-lib/common"

type ChainHandler interface {
	GetBlockHeaderHash(uint64) common.Hash
}
