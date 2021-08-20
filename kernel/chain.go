package kernel

type ChainHandler interface {
	GetBlockHeaderHash(uint64) Hash
}
