package usbimager

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
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

type SystemFileInfo struct {
	Dev        uint64 // Device ID
	Inode      uint64 // Inode number
	Mode       uint32 // File mode
	Nlink      uint64 // Number of hard links
	UID        uint32 // User ID
	GID        uint32 // Group ID
	GroupName  string
	UserName   string
	Rdev       uint64    // Device ID if special file
	Size       int64     // File size
	BlockSize  int64     // Block size
	Blocks     int64     // Number of blocks
	AccessTime time.Time // Last access time
	ModifyTime time.Time // Last modify time
	ChangeTime time.Time // Last change time
	CreateTime time.Time // Creation time (Windows only)
	Attributes uint32    // File attributes (Windows only)
}

func NewFileObject(path string, create_new bool) (FileObject, error) {
	fileinfo, err := NewSystemFileInfoFromPath(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return FileObject{}, err
		}
		if !create_new {
			return FileObject{}, err
		}
		log.Printf("Making new File %s\n", path)
		_, err := os.Create(path)
		if err != nil {
			log.Printf("Failed create file %v\n", err)
			return FileObject{}, err
		}
		return NewFileObject(path, false)
	}
	T := getFileTypeFromMode(fileinfo.Mode)
	return FileObject{
		Path: path,
		Type: T,
		Info: fileinfo,
	}, nil
}

func NewSystemFileInfoFromPath(file string) (*SystemFileInfo, error) {

	fileinfo, err := os.Stat(file)
	if err != nil {
		return nil, err
	}
	stat, ok := fileinfo.Sys().(*syscall.Stat_t)
	if !ok {
		return nil, fmt.Errorf("failed to cast to *syscall.Stat_t")
	}
	info := &SystemFileInfo{}
	info.Dev = stat.Dev
	info.Inode = stat.Ino
	info.Mode = stat.Mode
	info.Nlink = stat.Nlink
	info.UID = stat.Uid
	info.GID = stat.Gid
	info.GroupName = getGroupNameForGID(strconv.FormatUint(uint64(stat.Gid), 10))
	info.UserName = getUsernameForUID((strconv.FormatUint(uint64(stat.Uid), 10)))
	info.Rdev = stat.Rdev
	info.Size = stat.Size
	info.BlockSize = stat.Blksize
	info.Blocks = stat.Blocks
	info.AccessTime = time.Unix(stat.Atim.Sec, stat.Atim.Nsec)
	info.ModifyTime = time.Unix(stat.Mtim.Sec, stat.Mtim.Nsec)
	info.ChangeTime = time.Unix(stat.Ctim.Sec, stat.Ctim.Nsec)
	return info, nil

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

func getGroupNameForGID(gid string) string {
	var outbuf bytes.Buffer
	cmd := exec.Command("sh", "-c", fmt.Sprintf("getent group %s | cut -d: -f1", gid))
	cmd.Stdout = &outbuf
	//cmd.Stderr = &outbuf
	if err := cmd.Run(); err != nil {
		return gid
	}

	name := strings.TrimSpace(outbuf.String())
	return name
}

func getUsernameForUID(uid string) string {
	var outbuf bytes.Buffer
	cmd := exec.Command("sh", "-c", fmt.Sprintf("id -nu %s", uid))
	cmd.Stdout = &outbuf
	//cmd.Stderr = &outbuf
	if err := cmd.Run(); err != nil {
		return uid
	}

	name := strings.TrimSpace(outbuf.String())
	return name
}

func calculateNeededSizeForISO(iso FileObject) int64 {
	bytes := iso.Info.Size
	const GB = 1024 * 1024 * 1024
	if bytes <= 0 {
		return 0
	}
	gbNeeded := (bytes + GB - 1) / GB //ceiling division
	return gbNeeded * GB
}

func run(cmd string, args ...string) (string, string, error) {
	var out bytes.Buffer
	var err bytes.Buffer
	c := exec.Command(cmd, args...)
	c.Stdout = &out
	c.Stderr = &err
	if err := c.Run(); err != nil {
		return "", "", err
	}
	return out.String(), err.String(), nil
}
