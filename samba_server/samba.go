package main

import (
	"exec"
	"os/exec"
)

func getMount(mountdir string) (string, error){
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

func CheckMountPoint(mount_path string, device string) error  {

	return nil 
}
