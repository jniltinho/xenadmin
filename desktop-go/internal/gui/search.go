package gui

import (
	"errors"
	"os"
	"strings"

	"xcpng-center/internal/app"
	"xcpng-center/internal/messages"
	"xcpng-center/internal/ops"
	"xcpng-center/internal/xapi"
)

// The Search tab's commands (XenAdmin.Desktop SearchView, SearchEditor; SearchRules): New and Edit Search, Save to a
// pool, Delete, Export and Import of .xensearch files.

// SearchEditor is the New/Edit Search dialog: the form and the pools a search can be saved to.
type SearchEditor struct {
	app.SearchForm
	Pools []SearchPool `json:"pools"`
}

// SearchPool is a connected server's pool (Helpers.GetName).
type SearchPool struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

// SavedSearch is a search saved in a pool and the server it was read from.
type SavedSearch struct {
	XML    string `json:"xml"`
	Server string `json:"server"`
}

// EditSearch fills the editor from a search's XML ("" for New Search).
func (w *Workspace) EditSearch(text string) (SearchEditor, error) {
	f, err := app.EditSearchForm(text)
	if err != nil {
		return SearchEditor{}, err
	}
	out := SearchEditor{SearchForm: f, Pools: []SearchPool{}}
	for _, k := range w.connectedKeys() {
		_ = w.with(k, func(c *xapi.Connection) error {
			out.Pools = append(out.Pools, SearchPool{k, app.PoolName(c)})
			return nil
		})
	}
	return out, nil
}

// BuildSearch is the editor's search as XML (Apply and Save); the error is the validation text.
func (w *Workspace) BuildSearch(e app.SearchEdit) (string, error) { return app.BuildSearch(e) }

// OverviewSearch is the overview as a search (Edit Search and Export of the overview start from it).
func (w *Workspace) OverviewSearch() string { return app.OverviewSearch() }

// SearchDefinition is a saved search's XML, from the first connected pool that has it.
func (w *Workspace) SearchDefinition(uuid string) (SavedSearch, error) {
	for _, k := range w.connectedKeys() {
		var text string
		_ = w.with(k, func(c *xapi.Connection) error { var err error; text, err = app.PoolSearch(c, uuid); return err })
		if text != "" {
			return SavedSearch{text, k}, nil
		}
	}
	return SavedSearch{}, errors.New("this search is not saved in a connected pool any more")
}

// SaveSearch is SearchRules.SaveAction: the search into this server's pool (it joins Saved Searches); its uuid.
func (w *Workspace) SaveSearch(key, text string) (string, error) {
	uuid := app.SearchUUID(text) // "" when the text is not a search: ops.SaveSearch refuses it
	return uuid, w.with(key, func(c *xapi.Connection) error { return ops.SaveSearch(c, uuid, text) })
}

// DeleteSearch is SearchRules.DeleteAction (Search.Delete): the search removed from every connected pool.
func (w *Workspace) DeleteSearch(uuid string) error {
	var errs []error
	for _, k := range w.connectedKeys() {
		if err := w.with(k, func(c *xapi.Connection) error { return ops.DeleteSearch(c, uuid) }); err != nil && !errors.Is(err, errNotConnected) {
			errs = append(errs, errors.New(k+": "+err.Error()))
		}
	}
	return errors.Join(errs...)
}

// searchFilter is the file dialogs' filter (SearchPage: Messages.XENSEARCH_SAVED_SEARCH).
func searchFilter() [][2]string {
	text, _, _ := strings.Cut(messages.Format(messages.Get("XENSEARCH_SAVED_SEARCH"), "xensearch"), "|")
	return [][2]string{{text, "*.xensearch"}}
}

// ExportSearch is SearchView.Export (Search.Save(filename)): the search to a .xensearch file; the file, "" when
// cancelled.
func (w *Workspace) ExportSearch(text string) (string, error) {
	if w.Files == nil {
		return "", errNoWindow
	}
	content, err := app.SearchFile(text)
	if err != nil {
		return "", err
	}
	name := app.SearchName(text)
	if name == "" {
		name = "Overview"
	}
	file, err := w.Files(true, "Export Search", name+".xensearch", searchFilter())
	if err != nil || file == "" {
		return "", err
	}
	// shortcut: no DefaultExt added to a name typed without it (the dialog confirmed overwriting that very file)
	return file, os.WriteFile(file, []byte(content), 0o600)
}

// ImportedSearch is the first search of an imported file, with a new identity.
type ImportedSearch struct {
	XML  string `json:"xml"`
	Name string `json:"name"`
}

// ImportSearch is SearchView.Import: a .xensearch file's search, to show (Edit Search can then save it); empty when
// cancelled.
func (w *Workspace) ImportSearch() (ImportedSearch, error) {
	if w.Files == nil {
		return ImportedSearch{}, errNoWindow
	}
	file, err := w.Files(false, "Import Search", "", searchFilter())
	if err != nil || file == "" {
		return ImportedSearch{}, err
	}
	data, err := os.ReadFile(file)
	if err != nil {
		return ImportedSearch{}, err
	}
	list, err := app.ImportSearches(data)
	if err != nil {
		return ImportedSearch{}, errors.New("The file is not a search exported by XCP-ng Center.")
	}
	return ImportedSearch{list[0], app.SearchName(list[0])}, nil
}

// RunSearch runs a search that is shown but not saved (Apply, an imported search) on every connected server.
func (w *Workspace) RunSearch(text string) ([]*app.ViewNode, error) { return w.runView("", text) }
