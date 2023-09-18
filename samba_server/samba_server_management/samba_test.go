package sambaservermanagement_test

import (
	"testing"

	sambaservermanagement "github.com/PranoSA/samba_share_backend/samba_server/samba_server_management"
)

type SambaMountTest struct {
	Name       string
	Dev        string
	Mount_Path string
	Correct    bool
}

var SambaMountTests []SambaMountTest = []SambaMountTest{
	{
		Name:       "Correct Path",
		Dev:        "/dev/sda12",
		Mount_Path: "/mount/samba_server/1",
		Correct:    true,
	}, {
		Name:       "Wrong Mount Location",
		Dev:        "/dev/sda12",
		Mount_Path: "/mount/samba_server/21",
		Correct:    false,
	}, {
		Name:       "Nonexistent Disk",
		Dev:        "/dev/sda18",
		Mount_Path: "/mount/samba_server/1",
		Correct:    false,
	},
}

func TestGetMount(t *testing.T) {

	t.Log("Starting Shit ")
	for _, v := range SambaMountTests {
		t.Logf("Test %v \n", v.Name)
		t.Run(v.Name, func(t *testing.T) {
			result, err := sambaservermanagement.EnsureMount(v.Dev, v.Mount_Path)

			if result != v.Correct {
				t.Errorf("Expected %v, got %v, Error : %v", v.Correct, result, err)
			}
		})
	}
}
