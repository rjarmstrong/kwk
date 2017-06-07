package app

import (
	"github.com/kwk-super-snippets/cli/src/cli"
	"github.com/kwk-super-snippets/cli/src/out"
	"github.com/kwk-super-snippets/cli/src/runtime"
	"github.com/kwk-super-snippets/types"
	"github.com/kwk-super-snippets/types/vwrite"
)

func rootPrinter(prefs *out.Prefs, writer vwrite.Writer, user *types.User) cli.RootPrinter {
	return func(rr *types.RootResponse) error {
		return writer.EWrite(out.PrintRoot(prefs, &info, rr, user))
	}
}

func snippetPatcher(cxf cli.ContextFunc, sc types.SnippetsClient) runtime.SnippetPatcher {
	return func(req *types.PatchRequest) (*types.PatchResponse, error) {
		return sc.Patch(cxf(), req)
	}
}

func rootGetter(cxf cli.ContextFunc, sc types.SnippetsClient) runtime.RootGetter {
	return func(req *types.RootRequest) (*types.RootResponse, error) {
		return sc.GetRoot(cxf(), req)
	}
}

func useLogger(cxf cli.ContextFunc, sc types.SnippetsClient) runtime.UseLogger {
	return func(req *types.UseContext) (*types.LogUseResponse, error) {
		return sc.LogUse(cxf(), req)
	}
}

func snippetGetter(cxf cli.ContextFunc, sc types.SnippetsClient) runtime.SnippetGetter {
	return func(req *types.GetRequest) (*types.ListResponse, error) {
		return sc.Get(cxf(), req)
	}
}

func snippetMaker(cxf cli.ContextFunc, sc types.SnippetsClient) runtime.SnippetMaker {
	return func(req *types.CreateRequest) error {
		_, err := sc.Create(cxf(), req)
		return err
	}
}
