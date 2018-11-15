package main

import "soloos/tinyiron"

func (p *Env) CtrNetBlockWrite(ir *tinyiron.Request) {
	ir.ApiOutputSuccess("hi")
}
