package handlers

import (
	"github.com/rjarmstrong/kwk-types"
	"github.com/rjarmstrong/kwk/src/out"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPouches_Create(t *testing.T) {
	err := snippets.CreatePouch("pouch1")
	assert.Nil(t, err)
	funcName := writer.PopCalled("EWrite")
	assert.Equal(t, "PouchCreated", funcName)
}

func TestPouches_Lock(t *testing.T) {
	err := snippets.Lock("pouch1")
	assert.Nil(t, err)
	req := snippetClient.PopCalled("MakePouchPrivate").(*types.MakePrivateRequest)
	assert.Equal(t, true, req.MakePrivate)
	funcName := writer.PopCalled("EWrite")
	assert.Equal(t, "PouchLocked", funcName)
}

func TestPouches_UnLock(t *testing.T) {
	t.Log("UNLOCKED")
	err := snippets.UnLock("pouch1")
	assert.Nil(t, err)
	req := snippetClient.PopCalled("MakePouchPrivate").(*types.MakePrivateRequest)
	assert.Equal(t, false, req.MakePrivate)
	funcName := writer.PopCalled("EWrite")
	assert.Equal(t, "PouchUnLocked", funcName)

	t.Log("UNLOCK CANCELLED")
	dlg.returnsFor["Modal"] = response{val: &out.DialogResponse{Ok: false}}
	err = snippets.UnLock("pouch1")
	assert.Nil(t, err)
	funcName = writer.PopCalled("EWrite")
	assert.Equal(t, "PouchNotUnLocked", funcName)
	dlg.returnsFor["Modal"] = response{val: &out.DialogResponse{Ok: true}}
}
