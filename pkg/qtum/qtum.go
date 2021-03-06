package qtum

import (
	"github.com/pkg/errors"
	"github.com/qtumproject/janus/pkg/utils"
)

type Qtum struct {
	*Client
	*Method
	chain string
}

const (
	ChainMain    = "main"
	ChainTest    = "test"
	ChainRegTest = "regtest"
)

var AllChains = []string{ChainMain, ChainRegTest, ChainTest}

func New(c *Client, chain string) (*Qtum, error) {
	if !utils.InStrSlice(AllChains, chain) {
		return nil, errors.New("invalid qtum chain")
	}

	return &Qtum{
		Client: c,
		Method: &Method{Client: c},
		chain:  chain,
	}, nil
}

func (c *Qtum) Chain() string {
	return c.chain
}
