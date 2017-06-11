package handlers

import (
	"bufio"
	"bytes"
	"github.com/kwk-super-snippets/cli/src/out"
	"github.com/rjarmstrong/kwk-types"
	"github.com/rjarmstrong/kwk-types/errs"
	"os"
)

func resolveCreateArgs(args []string) (string, *types.Alias, error) {
	if len(args) == 1 {
		return resolveOneCreateArg(args[0])
	}
	return resolveTwoCreateArgs(args[0], args[1])
}

func resolveOneCreateArg(str string) (string, *types.Alias, error) {
	if types.IsDefNotPouchedSnippetURI(str) {
		out.Debug("Assuming the only arg is the content.")
		return str, &types.Alias{}, nil
	}
	a, err := types.ParseAlias(str)
	if err != nil {
		return "", nil, err
	}
	return "", a, nil
}

func resolveTwoCreateArgs(a string, b string) (string, *types.Alias, error) {
	aIsNotSUri := types.IsDefNotPouchedSnippetURI(a)
	bIsNotSUri := types.IsDefNotPouchedSnippetURI(b)

	out.Debug("ARG1:%v (snip=%t) ARG2:%v (snip=%t)", a, aIsNotSUri, b, bIsNotSUri)

	if aIsNotSUri && bIsNotSUri {
		return "", nil, errs.New(errs.CodeInvalidArgument,
			"Please specify a pouch to create a snippet."+
				"\n   e.g."+
				"\n  `kwk new <pouch_name>/<snip_name>[.<ext>]  <snippet>`"+
				"\n  `kwk new <snippet> <pouch_name>/<snip_name>[.<ext>]`"+
				"\n  `<cmd> | kwk new <pouch_name>/<snip_name>[.<ext>]`",
		)
	}
	if !aIsNotSUri && !bIsNotSUri {
		return "", nil, errs.New(errs.CodeInvalidArgument,
			"It looks like both arguments could be either be a path or kwk URIs, "+
				"please add an extension or fully quality the kwk URI. e.g. kwk.co/richard/dill/name.path")
	}
	if !aIsNotSUri {
		alias, err := types.ParseAlias(a)
		if err != nil {
			return "", nil, err
		}
		return b, alias, nil

	}
	alias, err := types.ParseAlias(b)
	if err != nil {
		return "", nil, err
	}
	return a, alias, nil
}

func stdInAsString() string {
	scanner := bufio.NewScanner(os.Stdin)
	in := bytes.Buffer{}
	for scanner.Scan() {
		in.WriteString(scanner.Text() + "\n")
	}
	return in.String()
}
