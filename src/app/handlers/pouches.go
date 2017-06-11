package handlers

import (
	"github.com/rjarmstrong/kwk/src/out"
	"github.com/rjarmstrong/kwk-types"
)

func (sc *Snippets) CreatePouch(name string) error {
	res, err := sc.client.CreatePouch(sc.cxf(), &types.CreatePouchRequest{Name: name})
	if err != nil {
		return err
	}
	sc.rootPrinter(res.Root)
	return sc.EWrite(out.PouchCreated(name))
}

func (sc *Snippets) deletePouch(pouch string) error {
	res := sc.Dialog.Modal(out.PouchCheckDelete(pouch), false)
	if !res.Ok {
		return sc.EWrite(out.PouchNotDeleted(pouch))
	}
	dres, err := sc.client.DeletePouch(sc.cxf(), &types.DeletePouchRequest{Name: pouch})
	if err != nil {
		return err
	}
	sc.rootPrinter(dres.Root)
	return sc.EWrite(out.PouchDeleted(pouch))
}

func (sc *Snippets) Lock(pouch string) error {
	_, err := sc.client.MakePouchPrivate(sc.cxf(), &types.MakePrivateRequest{Name: pouch, MakePrivate: true})
	if err != nil {
		return err
	}
	return sc.EWrite(out.PouchLocked(pouch))
}

func (sc *Snippets) UnLock(pouch string) error {
	res := sc.Dialog.Modal(out.PouchCheckUnLock(pouch), false)
	if res.Ok {
		_, err := sc.client.MakePouchPrivate(sc.cxf(), &types.MakePrivateRequest{Name: pouch, MakePrivate: false})
		if err != nil {
			return err
		}
		return sc.EWrite(out.PouchUnLocked(pouch))
	}
	return sc.EWrite(out.PouchNotUnLocked(pouch))
}
