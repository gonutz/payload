package payload

import (
	"encoding/binary"
	"errors"
	"io"
	"os"

	"github.com/kardianos/osext"
)

func Read() ([]byte, error) {
	annotate := func(msg string, err error) error {
		return errors.New("read payload: " + msg + ": " + err.Error())
	}

	// find the path currently executed file
	path, err := osext.Executable()
	if err != nil {
		return nil, annotate("unable to find executable name", err)
	}

	// The last 16 bytes in the file are the magic string "payload " folowed by
	// a uint64 that gives us the original exe's file size. This means that the
	// data starts at that offset and ends 16 bytes before the end of the file
	// (the 16 byte trailer is not part of the original data).

	file, err := os.Open(path)
	if err != nil {
		return nil, annotate("cannot open executable", err)
	}
	defer file.Close()

	// the end of the data is 16 bytes before the end of the file, due to the
	// trailer
	end, err := file.Seek(-16, os.SEEK_END)
	if err != nil {
		return nil, annotate("unable to seek to executable's end", err)
	}

	var magic [8]byte
	_, err = io.ReadFull(file, magic[:])
	if err != nil {
		return nil, annotate("unable to read magic string", err)
	}

	if string(magic[:]) != "payload " {
		return nil, errors.New("read payload: the executable does not contain a payload")
	}

	var originalSize uint64
	err = binary.Read(file, binary.LittleEndian, &originalSize)
	if err != nil {
		return nil, annotate("unable to read payload size", err)
	}

	if originalSize > uint64(end) {
		return nil, errors.New("read payload: invalid data size at file end")
	}

	// go to the original exe's end, at this point the payload data starts
	_, err = file.Seek(int64(originalSize), os.SEEK_SET)
	if err != nil {
		return nil, annotate("unable to seek to payload start", err)
	}

	data := make([]byte, end-int64(originalSize))
	_, err = io.ReadFull(file, data)
	if err != nil {
		return nil, annotate("unable to read payload data", err)
	}

	return data, nil
}
