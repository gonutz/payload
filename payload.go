package payload

import (
	"encoding/binary"
	"errors"
	"github.com/kardianos/osext"
	"io"
	"os"
)

func Read() ([]byte, error) {
	// find the path currently executed file
	path, err := osext.Executable()
	if err != nil {
		return nil, err
	}

	// The last 8 bytes in the file are a uint64 that gives us the original
	// exe's file size. This means that the data starts at that offset and ends
	// 8 bytes before the end of the file (the uint64 is not part of the
	// original data).

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// the end of the data is 8 bytes before the end of the file
	end, err := file.Seek(-8, os.SEEK_END)
	if err != nil {
		return nil, err
	}

	var originalSize uint64
	err = binary.Read(file, binary.LittleEndian, &originalSize)
	if err != nil {
		return nil, err
	}

	if originalSize > uint64(end) {
		return nil, errors.New("reading payload: invalid data size at file end")
	}

	// go to an offset of the original exe's size
	_, err = file.Seek(int64(originalSize), os.SEEK_SET)
	if err != nil {
		return nil, err
	}

	data := make([]byte, end-int64(originalSize))
	_, err = io.ReadFull(file, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
