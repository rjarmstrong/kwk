package settings

type SettingsMock struct {
	ChangeDirectoryCalledWith string
}

// get and check username exists
// Save to settings
// Print confirmation
//fmt.Println(gui.Colour(gui.LightBlue, "Switched to kwk.co/" + args[0] + "/"))
func (s *SettingsMock) ChangeDirectory(username string) {
	s.ChangeDirectoryCalledWith = username
}

func (s *SettingsMock) Delete(key string) error {
  return nil
}

func (s *SettingsMock) Get(key string, input interface{}) error {
	return nil
}

func (s *SettingsMock) Upsert(dir string, data interface{}) {
}
