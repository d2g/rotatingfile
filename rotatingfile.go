package rotatingfile

import (
	"log"
	"os"
	"strconv"
)

type File struct {
	*os.File

	MaxFileSize    uint64 // Max number of bytes before rotation. 0 = Unlimited.
	MaxBackupIndex uint   // Max number of rotations stored 0 = Unlimited.

	cBytesWritten uint64 // Current number of bytes written.
}

func (t *File) Write(p []byte) (n int, err error) {
	n, err = t.File.Write(p)
	if err != nil {
		return
	}

	//Update the current number of bytes written.
	t.cBytesWritten += uint64(n)

	if t.MaxFileSize > 0 && t.cBytesWritten >= t.MaxFileSize {
		// Our file should be rotated.

		// Delete the File at the maximum index.
		filename := t.File.Name()

		err = os.Remove(filename + "." + strconv.FormatUint(uint64(t.MaxBackupIndex), 10))
		if err != nil && !os.IsNotExist(err) {
			return
		}

		// Cycle renaming the other files.
		for i := (t.MaxBackupIndex - 1); i > 0; i-- {
			err = os.Rename(filename+"."+strconv.FormatUint(uint64(i), 10), filename+"."+strconv.FormatUint(uint64(i)+1, 10))
			if err != nil && !os.IsNotExist(err) {
				return
			}
		}

		// Rename the current open file.
		err = t.File.Close()
		if err != nil {
			log.Printf("Warning: Error Closing \"%s\" %s\n", t.File.Name(), err.Error())
		}

		err = os.Rename(filename, filename+".1")
		if err != nil {
			return
		}

		t.File, err = os.Create(t.File.Name())
		if err != nil {
			return
		}

		t.cBytesWritten = 0
	}
	return
}
