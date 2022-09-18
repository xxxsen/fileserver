package utils

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

type pair struct {
	fsz    uint64
	bsz    uint64
	bid    int
	realsz uint64
}

func TestCalcBlockSize(t *testing.T) {
	lst := []*pair{
		{fsz: 436567521, bsz: 200 * 1024 * 1024, bid: 2, realsz: 17137121},
		{fsz: 436567521, bsz: 200 * 1024 * 1024, bid: 1, realsz: 200 * 1024 * 1024},
		{fsz: 436567521, bsz: 200 * 1024 * 1024, bid: 0, realsz: 200 * 1024 * 1024},
		{fsz: 436567521, bsz: 200 * 1024 * 1024, bid: 3, realsz: 0},
		{fsz: 85770472, bsz: 20 * 1024 * 1024, bid: 0, realsz: 20 * 1024 * 1024},
		{fsz: 85770472, bsz: 20 * 1024 * 1024, bid: 1, realsz: 20 * 1024 * 1024},
		{fsz: 85770472, bsz: 20 * 1024 * 1024, bid: 2, realsz: 20 * 1024 * 1024},
		{fsz: 85770472, bsz: 20 * 1024 * 1024, bid: 3, realsz: 20 * 1024 * 1024},
		{fsz: 85770472, bsz: 20 * 1024 * 1024, bid: 4, realsz: 1884392},
	}

	for _, item := range lst {
		calcsz := CalcBlockSize(item.fsz, item.bsz, item.bid)
		assert.Equal(t, item.realsz, calcsz)
	}
}
