package sambaservermanagement

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

func getMount(mountdir string) (string, error) {
	dir, err := exec.Command("sh", "-c", fmt.Sprintf("mount | grep -w %s", mountdir)).Output()
	if err != nil {
		return "", err
	}

	num_lines := strings.Split(string(dir[:]), "\n")

	if len(num_lines) != 1 {
		return "", errors.New("Mounted At Less or More than 1 places")
	}

	return "", nil

}

func EnsureMount(dev string, mount_path string) (bool, error) {
	dir, err := exec.Command("sh", "-c", fmt.Sprintf("mount | grep -w %s", mount_path)).Output()
	if err != nil {
		return false, err //errors.New("Failed TO Execute Command")
	}

	fmt.Println(string(dir))

	num_lines := strings.Split(string(dir[:]), "\n")

	num_lines = num_lines[0:1]

	if len(num_lines) != 1 {
		return false, errors.New("More than Expected Mount Points ")
	}

	fields := strings.Fields(num_lines[0])

	if len(fields) != 6 {
		return false, errors.New("Output Different Than Expected, expected 6 fields ")
	}

	device := fields[0]

	if dev == device {
		return true, nil
	}
	deviceString := fmt.Sprintf("Requested%v:got %v", dev, device)
	return false, errors.New(deviceString)
}

/*func getMount(mountdir string) (string, error){
res, err := exec.Command(“sh”, “-c”, fmt.Sprintf(“mount | grep -w %s”, mountdir)).Output()
if err != nil{
return “”, err
}
lines := strings.Split(string(res[:]), “\n”)
if len(lines) != 1{
return “”, fmt.Errorf(“bad mount output”)
}
fields := strings.Fields(lines[0])
if len(fields) != 6{
return “”, fmt.Errorf(“bad mount line formating”)
}
device=fields[0]
return device, nill
}
*/
/**
 *
 *  Mount_Path Will Be SOmething like /mount/samba_server/spareid/shareid
 *
 * Name Will Be Something Like Owner + shareid
 *
 * New Samba User Will Be Something Like Owner+Shareid,
 * The Linux Group Name Will All Be ShareID
 *
 * Set First Password
 *
 * First Find If Mount Path is Valid Mount Path With Device Name
 *
 *
 *
 */
func AddSambaShare(mount_path string, name string, admin_password string) {
	exec.Command("sh", "-c", "")

}

func CheckMountPoint(mount_path string, device string) error {

	return nil
}
