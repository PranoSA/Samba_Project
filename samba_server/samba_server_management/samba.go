package sambaservermanagement

import (
	"errors"
	"fmt"
	"log"
	"os"
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

/**
 * Space_mount_path acts as the moutn point all samba shares in the space will mount
 * shareid is used to mount the folder at space_mount_path/shareid
 * owner and password are used by smbpasswd to create samba share
 * spaceid is also used to identify Linux Users Groups In THe Space
 * And shareid is Used to identify users as within the share as well
 */
func CreateSambaShare(space_mount_path string, shareid string, owner string, password string, spaceid string) {

	os.Mkdir(space_mount_path+"/"+shareid, 0770)

	_, err := exec.Command("sh", "./create_samba.sh", space_mount_path, shareid, owner, password, spaceid).Output()

	if err == nil {
		log.Fatalf("Failed TO ALlocat Samba : %v ", err)
	}

}

func AddUserToShareId(user string, password string, shareid string, spaceid string) {

	sambaname := strings.Replace(user, "@", "-", -1)

	_, err := exec.Command("sh", "./add_user_samba.sh", sambaname, password, shareid, spaceid).Output()

	if err != nil {

	}

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
