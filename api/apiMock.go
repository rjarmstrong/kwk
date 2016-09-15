package api


type ApiMock struct {
  PrintProfileCalled bool
}

func (a *ApiMock) PrintProfile() {
	a.PrintProfileCalled = true
}


