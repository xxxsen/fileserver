package utils

func CalcFileBlockCount(sz uint64, blksz uint64) int {
	return int((sz + blksz - 1) / blksz)
}

func CalcBlockSize(sz uint64, blksz uint64, blkid int) uint64 {
	blkcnt := CalcFileBlockCount(sz, blksz)
	if blkid >= blkcnt || blkid < 0 {
		return 0
	}
	if blkid < blkcnt-1 {
		return blksz
	}
	return sz - blksz*(uint64(blkcnt)-1)
}
