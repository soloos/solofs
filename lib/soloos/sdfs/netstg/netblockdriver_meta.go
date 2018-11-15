package netstg

func (p *NetBlockDriver) AllocBlock(inodePath string) {
	p.snetClientDriver.Write(p.nameNodePeer)
}
