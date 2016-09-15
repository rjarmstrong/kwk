package system


type ISystem interface {
	Upgrade()
	GetVersion() string
	ChangeDirectory(username string)
}

type SystemMock struct {
  UpgradeCalled bool
  VersionCalled bool
  ChangeDirectoryCalledWith string
}

func (s *SystemMock) Upgrade(){
  s.UpgradeCalled = true
}

func (s *SystemMock) GetVersion() string {
    s.VersionCalled = true
   return "0.0.1"
}


// get and check username exists
// Save to settings
// Print confirmation
//fmt.Println(gui.Colour(gui.LightBlue, "Switched to kwk.co/" + args[0] + "/"))
func (s *SystemMock) ChangeDirectory(username string) {
	s.ChangeDirectoryCalledWith = username
}