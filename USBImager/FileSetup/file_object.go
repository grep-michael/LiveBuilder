package filesetup

import (
	"fmt"
	"os"
	"syscall"
)

type FileType int

const (
	TypeRegularFile FileType = iota
	TypeBlockDevice
	TypeDirectory
	TypeSymlink
	TypeCharDevice
	TypeFIFO
	TypeSocket
	TypeUnknown
)
const ALLOWEDTYPES = (1 << TypeRegularFile) | (1 << TypeBlockDevice)

func (ft FileType) String() string {
	switch ft {
	case TypeRegularFile:
		return "regular file"
	case TypeDirectory:
		return "directory"
	case TypeSymlink:
		return "symbolic link"
	case TypeBlockDevice:
		return "block device"
	case TypeCharDevice:
		return "character device"
	case TypeFIFO:
		return "FIFO/named pipe"
	case TypeSocket:
		return "socket"
	default:
		return "unknown"
	}
}

type fileObject struct {
	path     string
	fileType FileType
	size     int64
}

func (fileobj *fileObject) isNotAllowedFile() bool {
	return (1<<fileobj.fileType)&ALLOWEDTYPES == 0
}

func newFileObject(path string, should_make bool) *fileObject {
	fileinfo, err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
		if should_make {
			f, _ := os.Create(path)
			f.Close()
			return newFileObject(path, false)
		} else {
			return nil
		}

	}
	stat, ok := fileinfo.Sys().(*syscall.Stat_t)
	if !ok {
		fmt.Println("failed to cast to *syscall.Stat_t")
		return nil
	}
	return &fileObject{
		path:     path,
		fileType: getFileTypeFromMode(stat.Mode),
		size:     stat.Size,
	}

}

func getFileTypeFromMode(mode uint32) FileType {
	const (
		S_IFMT   = 0170000 // bit mask for the file type bit field
		S_IFSOCK = 0140000 // socket
		S_IFLNK  = 0120000 // symbolic link
		S_IFREG  = 0100000 // regular file
		S_IFBLK  = 0060000 // block device
		S_IFDIR  = 0040000 // directory
		S_IFCHR  = 0020000 // character device
		S_IFIFO  = 0010000 // FIFO
	)

	fileType := mode & S_IFMT

	switch fileType {
	case S_IFREG:
		return TypeRegularFile
	case S_IFDIR:
		return TypeDirectory
	case S_IFLNK:
		return TypeSymlink
	case S_IFBLK:
		return TypeBlockDevice
	case S_IFCHR:
		return TypeCharDevice
	case S_IFIFO:
		return TypeFIFO
	case S_IFSOCK:
		return TypeSocket
	default:
		return TypeUnknown
	}
}
