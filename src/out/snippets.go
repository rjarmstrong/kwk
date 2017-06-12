package out

import (
	"fmt"
	"github.com/rjarmstrong/kwk-types"
	"github.com/rjarmstrong/kwk-types/vwrite"
	"github.com/rjarmstrong/kwk/src/style"
	"io"
	"strings"
)

func SnippetDescriptionUpdated(uri string, desc string) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		Success(w, "%s description updated:\n%s\n\n\n", uri, desc)
	})
}

func SnippetClonedAs(list *types.ListResponse, prefs *Prefs, newName string) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		printPouchSnippets(w, prefs, list)
		Success(w, "Cloned as %s\n\n", newName)
	})
}

func SnippetCreated(s *types.Snippet) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		Success(w, "%s %s created \n\n", snippetIcon(s), s.Alias.URI())
	})
}

func SnippetCat(s *types.Snippet) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		fmt.Fprintln(w, s.Content)
	})
}

func SnippetEdited(s *types.Snippet) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		Success(w, "Successfully updated %s %s\n\n", snippetIcon(s), s.Alias.VersionURI())
	})
}

func SnippetNoChanges(s *types.Snippet) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		Success(w, "No changes to %s %s\n\n", snippetIcon(s), s.Alias.VersionURI())
	})
}

func SnippetEditing(s *types.Snippet) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		fmt.Fprintf(w, "\n%sEditing:  %s %s...\n\n", style.Margin, snippetIcon(s), s.Alias.VersionURI())
		fmt.Fprintf(w, "%sCTRL+C to CANCEL | Any key to COMMIT\n", style.Margin)
	})
}

func SnippetEditNewPrompt(uri string) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		Info(w, "%s doesn't exist - would you like create it? [y/n] \n", uri)
	})
}

func SnippetList(prefs *Prefs, list *types.ListResponse) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		printPouchSnippets(w, prefs, list)
	})
}

func PrintRoot(prefs *Prefs, cli *types.AppInfo, rr *types.RootResponse, u *types.User) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		printRoot(w, prefs, cli, rr, u)
	})
}

func Tagged(uri string, tags []string) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		Success(w, "%s tagged with: %s\n", uri, strings.Join(tags, ", "))
	})
}

func UnTagged(uri string, tags []string) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		Success(w, "Tags: %s removed from %s\n", strings.Join(tags, ", "), uri)
	})
}

func SnippetAmbiguous(callerUri string, uri string) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		Warn(w, "The uri %s called from app %s is ambiguous, please provide the username, pouch, name and extension.", uri, callerUri)
	})
}

func DidYouMean(uri string) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		Info(w, "Did you mean %s ? y/n ", uri)
	})
}

func RunAllSnippetsNotTrue(callerUri string, uri string) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		Warn(w, "You have not enabled 'runallsnippets' in prefs and the usersname for uri %s called from app %s is not yours. Either change the uri or enable runallsnippets.", uri, callerUri)
	})
}

func SnippetRenamed(list *types.ListResponse, prefs *Prefs, originalUri string, newUri string) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		printPouchSnippets(w, prefs, list)
		Success(w, "%s renamed to %s", originalUri, newUri)
	})
}

func SnippetPatched(uri string) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		Success(w, "%s patched", uri)
	})
}

func SnippetCheckDelete(snipNames []*types.SnipName) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		printSnipNames(w, snipNames)
		Warn(w, "\nAre you sure you want to delete these snippets? [y/n] ")
	})
}

func SnippetsDeleted(snipNames []*types.SnipName) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		Success(w, "Snippets deleted: ")
		printSnipNames(w, snipNames)
		fmt.Fprintln(w, "")
	})
}

func SnippetsNotDeleted(snipNames []*types.SnipName) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		Info(w, "Snippets NOT deleted: ")
		printSnipNames(w, snipNames)
		fmt.Fprintln(w, "")
	})
}

func SnippetsMoved(list *types.ListResponse, prefs *Prefs, snipNames []*types.SnipName, pouch string) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		printPouchSnippets(w, prefs, list)
		printSnipNames(w, snipNames)
		Info(w, " moved to pouch %s\n", pouch)
	})
}

func SnippetView(prefs *Prefs, s *types.Snippet) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		printSnippetView(w, prefs, s)
	})
}
