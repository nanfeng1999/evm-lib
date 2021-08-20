package kernel

type ETHChainHandler struct{}

func (ethChainHandler *ETHChainHandler) GetBlockHeaderHash(uint64) Hash {
	//just return a fake value
	return HexToHash("this is a demo")
}
