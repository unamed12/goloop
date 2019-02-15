package block

import (
	"bytes"
	"io"
	"testing"

	"github.com/icon-project/goloop/common/db"
	"github.com/icon-project/goloop/module"
	"github.com/stretchr/testify/assert"
)

func assertHasValidGenesisBlock(t *testing.T, bm module.BlockManager) {
	blk, err := bm.GetLastBlock()
	assert.Nil(t, err, "error")
	assert.Equal(t, gheight, blk.Height(), "height")
	id := blk.ID()

	blk, err = bm.GetBlockByHeight(gheight)
	assert.Nil(t, err, "error")
	assert.Equal(t, gheight, blk.Height(), "height")
	assert.Equal(t, id, blk.ID(), "ID")

	blk, err = bm.GetBlock(id)
	assert.Nil(t, err, "error")
	assert.Equal(t, gheight, blk.Height(), "height")
	assert.Equal(t, id, blk.ID(), "ID")
}

type blockGenerator struct {
	t  *testing.T
	sm *testServiceManager
	bm module.BlockManager
	n  int64
}

func newBlockGenerator(t *testing.T, c module.Chain) *blockGenerator {
	bg := &blockGenerator{}
	bg.t = t
	bg.sm = newTestServiceManager(c.Database())
	bg.bm = NewManager(c, bg.sm)
	return bg
}

func (bg *blockGenerator) getBlock(n int64) module.Block {
	bg.generateUntil(n)
	blk, err := bg.bm.GetBlockByHeight(n)
	assert.Nil(bg.t, err, "GetBlockByHeight")
	return blk
}

func (bg *blockGenerator) getReaderForBlock(n int64) io.Reader {
	buf := bytes.NewBuffer(nil)
	blk := bg.getBlock(n)
	blk.MarshalHeader(buf)
	blk.MarshalBody(buf)
	return buf
}

func (bg *blockGenerator) generateUntil(n int64) {
	blk, err := bg.bm.GetLastBlock()
	assert.Nil(bg.t, err, "GetLastBlock")
	for i := blk.Height(); i < n; i++ {
		pid := blk.ID()
		br := proposeSync(bg.bm, pid, newCommitVoteSet(true))
		blk = br.blk
		err := bg.bm.Finalize(blk)
		assert.Nil(bg.t, err, "Finalize")
	}
}

type blockManagerTestSetUp struct {
	gtx      *testTransaction
	database db.Database
	chain    *testChain
	sm       *testServiceManager
	bm       module.BlockManager
	bg       *blockGenerator
}

func newBlockManagerTestSetUp(t *testing.T) *blockManagerTestSetUp {
	s := &blockManagerTestSetUp{}
	s.database = newMapDB()
	s.gtx = newGenesisTX(defaultValidators)
	s.chain = newTestChain(s.database, s.gtx)
	s.sm = newServiceManager(s.chain)
	s.bm = NewManager(s.chain, s.sm)
	c := newTestChain(newMapDB(), s.gtx)
	s.bg = newBlockGenerator(t, c)
	return s
}

func getLastBlockID(t *testing.T, bm module.BlockManager) []byte {
	blk, err := bm.GetLastBlock()
	assert.Nil(t, err, "last block error")
	return blk.ID()
}

func getBadBlockID(t *testing.T, bm module.BlockManager) []byte {
	id := getLastBlockID(t, bm)
	pid := make([]byte, len(id))
	copy(pid, id)
	pid[0] = ^pid[0]
	return pid
}

type blockResult struct {
	blk      module.Block
	err      error
	cberr    error
	cbCalled bool
}

func (br *blockResult) assertOK(t *testing.T) {
	assert.NotNil(t, br.blk, "block")
	assert.Nil(t, br.err, "return error")
	assert.Nil(t, br.cberr, "cb error")
	assert.True(t, br.cbCalled, "cb called")
}

func (br *blockResult) assertError(t *testing.T) {
	assert.Nil(t, br.blk, "block")
	assert.NotNil(t, br.err, "return error")
	assert.Nil(t, br.cberr, "cb error")
	assert.False(t, br.cbCalled, "cb called")
}

func (br *blockResult) assertCBError(t *testing.T) {
	assert.Nil(t, br.blk, "block")
	assert.Nil(t, br.err, "return error")
	assert.NotNil(t, br.cberr, "cb error")
	assert.True(t, br.cbCalled, "cb called")
}

type cbResult struct {
	blk module.Block
	err error
}

func proposeSync(bm module.BlockManager, pid []byte, vs module.CommitVoteSet) *blockResult {
	ch := make(chan cbResult)
	_, err := bm.Propose(pid, vs, func(blk module.Block, err error) {
		ch <- cbResult{blk, err}
	})
	if err != nil {
		return &blockResult{nil, err, nil, false}
	}
	res := <-ch
	return &blockResult{res.blk, nil, res.err, true}
}

func importSync(bm module.BlockManager, r io.Reader) *blockResult {
	ch := make(chan cbResult)
	_, err := bm.Import(r, func(blk module.Block, err error) {
		ch <- cbResult{blk, err}
	})
	if err != nil {
		return &blockResult{nil, err, nil, false}
	}
	res := <-ch
	return &blockResult{res.blk, nil, res.err, true}
}

func TestBlockManager_New_HasValidGenesisBlock(t *testing.T) {
	s := newBlockManagerTestSetUp(t)
	bm := s.bm
	assertHasValidGenesisBlock(t, bm)
	blk, _ := bm.GetLastBlock()
	id := blk.ID()
	assert.Equal(t, s.gtx.Data.Effect.NextValidators.Bytes(), blk.NextValidators().Bytes())

	sm := newServiceManager(s.chain)
	bm = NewManager(s.chain, sm)
	assertHasValidGenesisBlock(t, bm)
	blk, _ = bm.GetLastBlock()
	assert.Equal(t, id, blk.ID(), "ID")
	assert.Equal(t, s.gtx.Data.Effect.NextValidators.Bytes(), blk.NextValidators().Bytes())
}

func TestBlockManager_Propose_ErrorOnBadParent(t *testing.T) {
	bm := newBlockManagerTestSetUp(t).bm
	pid := getBadBlockID(t, bm)
	br := proposeSync(bm, pid, newCommitVoteSet(false))
	br.assertError(t)
}

func TestBlockManager_Propose_ErrorOnInvalidCommitVoteSet(t *testing.T) {
	bm := newBlockManagerTestSetUp(t).bm
	pid := getLastBlockID(t, bm)

	cases := []struct {
		vs module.CommitVoteSet
		ok bool
	}{
		{newCommitVoteSet(false), false},
		{newCommitVoteSet(true), true},
	}
	// for height 1
	for _, c := range cases {
		br := proposeSync(bm, pid, c.vs)
		if c.ok {
			br.assertOK(t)
			err := bm.Finalize(br.blk)
			assert.Nil(t, err, "finalize error")
			pid = br.blk.ID()
		} else {
			br.assertError(t)
		}
	}
	// for height 2
	for _, c := range cases {
		br := proposeSync(bm, pid, c.vs)
		if c.ok {
			br.assertOK(t)
			err := bm.Finalize(br.blk)
			assert.Nil(t, err, "finalize error")
			pid = br.blk.ID()
		} else {
			br.assertError(t)
		}
	}
}

func TestBlockManager_Propose_ReturnsValidBlock(t *testing.T) {
	s := newBlockManagerTestSetUp(t)
	bm := s.bm
	sm := s.sm
	tx := newTestTransaction()
	tx.Data.Effect.NextValidators = newRandomTestValidatorList(2)
	sm.SendTransaction(tx)
	pid := getLastBlockID(t, bm)
	br := proposeSync(bm, pid, newCommitVoteSet(true))
	br.assertOK(t)
	blk := br.blk
	assert.Equal(t, gheight+1, blk.Height(), "height")
	assert.Equal(t, pid, blk.PrevID(), "prevID")
	assert.Equal(t, s.chain.Wallet().Address().Bytes(), blk.Proposer().Bytes())
	assert.Equal(t, s.gtx.Data.Effect.NextValidators.Bytes(), blk.NextValidators().Bytes())
	ntxs := blk.NormalTransactions()
	assert.NotNil(t, ntxs, "normal transactions")
	tx2, err := ntxs.Get(0)
	assert.Nil(t, err, "0th transaction")
	ttx2 := tx2.(*testTransaction)
	assert.NotNil(t, ttx2, "casting to testTransaction")
	assert.Equal(t, tx, ttx2, "transaction")
	br = proposeSync(bm, br.blk.ID(), newCommitVoteSet(true))
	br.assertOK(t)
	assert.Equal(t, tx.Data.Effect.NextValidators.Bytes(), br.blk.NextValidators().Bytes(), "validator list")
}

func TestBlockManager_Propose_Cancel(t *testing.T) {
	s := newBlockManagerTestSetUp(t)
	ec := make(chan struct{})
	s.sm.setTransitionExeChan(ec)
	pid := getLastBlockID(t, s.bm)
	br := proposeSync(s.bm, pid, newCommitVoteSet(true))
	blk := br.blk
	pid = blk.ID()

	canceler, err := s.bm.Propose(pid, newCommitVoteSet(true), func(blk module.Block, err error) {
		assert.Fail(t, "canceled proposal cb was called")
	})
	assert.Nil(t, err, "propose return error")
	res := canceler()
	assert.Equal(t, true, res, "canceler result")
}

func TestBlockManager_Import_ErrorOnBadParent(t *testing.T) {
	s := newBlockManagerTestSetUp(t)
	r := s.bg.getReaderForBlock(3)
	br := importSync(s.bm, r)
	br.assertError(t)
}

func TestBlockManager_Import_OK(t *testing.T) {
	s := newBlockManagerTestSetUp(t)
	for i := int64(1); i < 10; i++ {
		r := s.bg.getReaderForBlock(i)
		br := importSync(s.bm, r)
		br.assertOK(t)
		s.bm.Finalize(br.blk)
	}
}

func TestBlockManager_Import_Cancel(t *testing.T) {
	s := newBlockManagerTestSetUp(t)
	ec := make(chan struct{})
	r := s.bg.getReaderForBlock(1)
	s.sm.setTransitionExeChan(ec)
	canceler, err := s.bm.Import(r, func(blk module.Block, err error) {
		assert.Fail(t, "canceled import cb was called")
	})
	assert.Nil(t, err, "import return error")
	res := canceler()
	assert.Equal(t, true, res, "canceler result")
}