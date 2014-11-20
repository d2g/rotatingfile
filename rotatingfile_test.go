package rotatingfile_test

import (
	"crypto/rand"
	"os"
	"testing"

	"github.com/d2g/rotatingfile"
)

func TestRotatingFile(test *testing.T) {
	file, err := os.Create(os.TempDir() + "/file.tmp")
	if err != nil {
		test.Errorf("Error Creating Temp File %s\n", err.Error())
	}

	rotatingfile := rotatingfile.File{
		File:           file,
		MaxFileSize:    10485760,
		MaxBackupIndex: 10,
	}

	dummy := make([]byte, rotatingfile.MaxFileSize) // Cause a file rotation.

	_, err = rand.Read(dummy)
	if err != nil {
		test.Errorf("Error Creating Dummy Data %s\n", err.Error())
	}

	_, err = rotatingfile.Write(dummy)
	if err != nil {
		test.Errorf("Error Writing Data %s\n", err.Error())
	}

	// We should now have 2 files.
	err = rotatingfile.Close()
	if err != nil {
		test.Errorf("Error Closing File %s\n", err.Error())
	}

	// Tidy Up.
	err = os.Remove(os.TempDir() + "/file.tmp")
	if err != nil {
		test.Errorf("Error Removing original File %s\n", err.Error())
	}
	err = os.Remove(os.TempDir() + "/file.tmp.1")
	if err != nil {
		test.Errorf("Error Removing .1 File %s\n", err.Error())
	}
}

func TestRemovingAtMaxRotation(test *testing.T) {
	file, err := os.Create(os.TempDir() + "/file.tmp")
	if err != nil {
		test.Errorf("Error Creating Temp File %s\n", err.Error())
	}

	rotatingfile := rotatingfile.File{
		File:           file,
		MaxFileSize:    10485760,
		MaxBackupIndex: 2,
	}

	dummy := make([]byte, rotatingfile.MaxFileSize) // Cause a file rotation.

	_, err = rand.Read(dummy)
	if err != nil {
		test.Errorf("Error Creating Dummy Data %s\n", err.Error())
	}

	// Create Original an .1
	_, err = rotatingfile.Write(dummy)
	if err != nil {
		test.Errorf("Error Writing Data %s\n", err.Error())
	}

	//Create original, .1, .2
	_, err = rotatingfile.Write(dummy)
	if err != nil {
		test.Errorf("Error Writing Data %s\n", err.Error())
	}

	//Create original, .1, .2, -> Dump existing .2
	_, err = rotatingfile.Write(dummy)
	if err != nil {
		test.Errorf("Error Writing Data %s\n", err.Error())
	}

	// We should now have 3 files.
	err = rotatingfile.Close()
	if err != nil {
		test.Errorf("Error Closing File %s\n", err.Error())
	}

	// Tidy Up.
	err = os.Remove(os.TempDir() + "/file.tmp")
	if err != nil {
		test.Errorf("Error Removing original File %s\n", err.Error())
	}
	err = os.Remove(os.TempDir() + "/file.tmp.1")
	if err != nil {
		test.Errorf("Error Removing .1 File %s\n", err.Error())
	}
	err = os.Remove(os.TempDir() + "/file.tmp.2")
	if err != nil {
		test.Errorf("Error Removing .1 File %s\n", err.Error())
	}

	//Try and remove 3 which shouldn't be there.
	err = os.Remove(os.TempDir() + "/file.tmp.3")
	if !os.IsNotExist(err) {
		if err == nil {
			test.Errorf("Error .3 File shouldn't exists?\n")
		} else {
			test.Errorf("Error .3 File shouldn't exists? %s\n", err.Error())
		}
	}

}
