package out

import (
	"bitbucket.com/sharingmachine/kwkcli/src/models"
	"bitbucket.com/sharingmachine/types"
	"bitbucket.com/sharingmachine/types/vwrite"
	"fmt"
	"io"
	"strings"
)

func SnippetDescriptionUpdated(uri string, desc string) vwrite.Handler {
	return Success(vwrite.HandlerFunc(func(w io.Writer) {
		fmt.Fprintf(w, "%s description updated:\n%s\n\n\n", uri, desc)
	}))
}

func SnippetClonedAs(newName string) vwrite.Handler {
	return Success(vwrite.HandlerFunc(func(w io.Writer) {
		fmt.Fprintf(w, "Cloned as %s\n\n", newName)
	}))
}

func SnippetCreated(s *types.Snippet) vwrite.Handler {
	return Success(vwrite.HandlerFunc(func(w io.Writer) {
		fmt.Fprintf(w, "%s %s created \n\n", snippetIcon(s), s.String())
	}))
}

func SnippetCat(s *types.Snippet) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		fmt.Fprint(w, s.Snip)
	})
}

func SnippetAmbiguousCat(snippets []*types.Snippet) vwrite.Handler {
	return Info(vwrite.HandlerFunc(func(w io.Writer) {
		fmt.Fprintln(w, "That snippet name is ambiguous, please choose one of the following: ")
	}))
}

func SnippetEdited(s *types.Snippet) vwrite.Handler {
	return Success(vwrite.HandlerFunc(func(w io.Writer) {
		fmt.Fprintf(w, "Successfully updated %s %s\n\n", snippetIcon(s), s.String())
	}))
}

func SnippetEditing(s *types.Snippet) vwrite.Handler {
	return Success(vwrite.HandlerFunc(func(w io.Writer) {
		fmt.Fprintf(w, "Editing: %s %s... \nHit ENTER to upload changes. CTRL+C to cancel.\n", snippetIcon(s), s.String())
	}))
}

func SnippetEditNewPrompt(uri string) vwrite.Handler {
	return Info(vwrite.HandlerFunc(func(w io.Writer) {
		fmt.Fprintf(w, "%s doesn't exist - would you like create it? [y/n] \n", uri)
	}))
}

func SnippetList(list *models.ListView) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		printPouchSnippets(w, list)
	})
}

func PrintRoot(list *models.ListView) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		printRoot(w, list)
	})
}

func Tagged(uri string, tags []string) vwrite.Handler {
	return Success(vwrite.HandlerFunc(func(w io.Writer) {
		fmt.Fprintf(w, "%s tagged with: %s\n", uri, strings.Join(tags, ", "))
	}))
}

func UnTag(uri string, tags []string) vwrite.Handler {
	return Success(vwrite.HandlerFunc(func(w io.Writer) {
		fmt.Fprintf(w, "Tags: %s removed from %s\n", strings.Join(tags, ", "), uri)
	}))
}

func SnippetRenamed(originalUri string, newUri string) vwrite.Handler {
	return Success(vwrite.HandlerFunc(func(w io.Writer) {
		fmt.Fprintf(w, "%s renamed to %s", originalUri, newUri)
	}))
}

func SnippetPatched(uri string) vwrite.Handler {
	return Success(vwrite.HandlerFunc(func(w io.Writer) {
		fmt.Fprintf(w, "%s patched", uri)
	}))
}

func SnippetCheckDelete(snipNames []*types.SnipName) vwrite.Handler {
	return Warn(vwrite.HandlerFunc(func(w io.Writer) {
		printSnipNames(w, snipNames)
		fmt.Fprint(w, "\nAre you sure you want to delete these snippets? [y/n] ")
	}))
}

func SnippetsDeleted(snipNames []*types.SnipName) vwrite.Handler {
	return Success(vwrite.HandlerFunc(func(w io.Writer) {
		fmt.Fprint(w, "Snippets deleted: ")
		printSnipNames(w, snipNames)
		fmt.Fprintln(w, "")
	}))
}

func SnippetsNotDeleted(snipNames []*types.SnipName) vwrite.Handler {
	return Info(vwrite.HandlerFunc(func(w io.Writer) {
		fmt.Fprint(w, "Snippets NOT deleted: ")
		printSnipNames(w, snipNames)
		fmt.Fprintln(w, "")
	}))
}

func SnippetsMoved(snipNames []*types.SnipName, pouch string) vwrite.Handler {
	return Info(vwrite.HandlerFunc(func(w io.Writer) {
		printSnipNames(w, snipNames)
		fmt.Fprintf(w, " moved to pouch %s\n", pouch)
	}))
}

func SnippetPouchCreatePrompt() vwrite.Handler {
	return Info(vwrite.HandlerFunc(func(w io.Writer) {
		fmt.Fprint(w, "Would you like to create the snippet in a new pouch? [y/n] ")
	}))
}

func SnippetView(s *types.Snippet) vwrite.Handler {
	return Info(vwrite.HandlerFunc(func(w io.Writer) {
		printSnippetView(w, s)
	}))
}
