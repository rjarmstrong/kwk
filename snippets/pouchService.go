package snippets

import "bitbucket.com/sharingmachine/kwkcli/models"

type PouchService interface {
	GetRoot (username string, all bool) (*Root, error)
	Create (pouch string) (string, error)
	Rename (pouch string, newPouch string) (string, error)
	MakePrivate (pouch string, private bool) (bool, error)
	Delete (pouch string) (bool, error)
}

type Pouches struct {

}

func (*Pouches) GetRoot(username string, all bool) (*Root, error) {
	panic("implement me")
}

func (*Pouches) Create(pouch string) (string, error) {
	panic("implement me")
}

func (*Pouches) Rename(pouch string, newPouch string) (string, error) {
	panic("implement me")
}

func (*Pouches) MakePrivate(pouch string, private bool) (bool, error) {
	panic("implement me")
}

func (*Pouches) Delete(pouch string) (bool, error) {
	panic("implement me")
}

type Root struct {
	Pouches []models.Pouch
	Snippets []*models.Snippet
}

func (rt *Root) IsPouch(name string) bool {
	return false
}