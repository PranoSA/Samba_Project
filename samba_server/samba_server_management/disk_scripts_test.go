package sambaservermanagement_test

import (
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestDiskScripts(t *testing.T) {

	device := "/dev/sda3"

	res, err := exec.Command("sh", "-c", "df -h").Output()
	if err != nil {
		t.Errorf("%v", err)
	}
	entries := strings.Split(string(res), "\n")

	var correctEntry string

	for _, e := range entries[1:] {
		fields := strings.Fields(e)
		if len(fields) != 6 {
			break
		}
		if fields[0] == device {
			correctEntry = fields[3]
		}
	}

	if correctEntry == "" {
		t.Error("Failed To Find Disk")
	}

	re, err := regexp.Compile("^[0-9]+")

	if err != nil {
		t.Error("Failed to Compile Regex")
	}

	disk := re.FindAllString(correctEntry, -1)
	if len(disk) != 1 {
		t.Error("Improper Disk Format ")
	}

	er, err := regexp.Compile("^[0-9]+(M|G|T)$")

	units := er.FindAllStringSubmatch(correctEntry, -1)
	if len(units) != 1 {
		t.Errorf("Not Correct Unit ")
	}

	size, err := strconv.Atoi(disk[0])

	if err != nil {
		t.Error("")
	}

	if units[0][1] == "G" {
		size = size * 1000
	}
	if units[0][1] == "T" {
		size = size * 1_000_000
	}
	if units[0][1] == "K" {
		size = size / 1000
	}
	//return size

}
