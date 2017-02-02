package snippets

type PouchMock struct {

}

func (*PouchMock) GetRoot(username string, all bool) (*Root, error) {
	panic("implement me")
}

func (*PouchMock) Create(pouch string) (string, error) {
	panic("implement me")
}

func (*PouchMock) Rename(pouch string, newPouch string) (string, error) {
	panic("implement me")
}

func (*PouchMock) MakePrivate(pouch string, private bool) (bool, error) {
	panic("implement me")
}

func (*PouchMock) Delete(pouch string) (bool, error) {
	panic("implement me")
}

