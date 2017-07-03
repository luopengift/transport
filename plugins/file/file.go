package file

import (
	"github.com/luopengift/golibs/logger"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
)

type File struct {
	name  string
	model int
	inode uint64
	fd    *os.File
	seek  int64
}

func NewFile(name string, model int) *File {
	file := &File{
		name:  name,
		model: model,
	}
	file.Open()
	return file
}

// 文件名称
func (self *File) Name() string {
	return self.name
}

func (self *File) Dir() string {
	return filepath.Dir(self.name)
}

// 文件名
func (self *File) BaseName() string {
	return filepath.Base(self.name)
}

func (self *File) Open() (err error) {
	self.fd, err = os.OpenFile(self.name, self.model, 0660)
	if err != nil {
		logger.Error("<file %s can not open>:%v", self.name, err)
		return
	}

	self.inode, err = self.Inode()
	if err != nil {
		logger.Error("< %s can not get inode>:%v", self.name, err)
		return
	}
	return
}

func (self *File) Close() error {
	return self.fd.Close()
}

func (self *File) Fd() *os.File {
	return self.fd
}

// os.SEEK_CUR int = 1 // seek relative to the current offset
// os.SEEK_SET int = 0 // seek relative to the origin of the file
// os.SEEK_END int = 2 // seek relative to the end
func (self *File) Seek(offset int64) (err error) {
	self.seek, err = self.fd.Seek(offset, os.SEEK_SET)
	return
}

func (self *File) ReadOneByte(offset int64) ([]byte, error) {
	buf := make([]byte, 1)
	_, err := self.fd.ReadAt(buf, offset)
	return buf, err

}

// 根据offset值,往前计算该行的起始偏移量
func (self *File) TrancateOffsetByLF(offset int64) (int64, error) {
	for ; offset >= 0; offset-- {
		buf, err := self.ReadOneByte(offset)
		if err != nil {
			return 0, err
		}
		if string(buf) == "\n" {
			return offset + 1, nil //pos为"\n"的位置,需要加1才是行首的位置
		}
	}
	return 0, nil
}

// 根据offset值,往后计算该行的起始偏移量
func (self *File) CeilingOffsetByLF(offset int64) (int64, error) {
	for ; ; offset++ {
		buf, err := self.ReadOneByte(offset)
		if err != nil {
			return 0, err
		}
		if string(buf) == "\n" {
			return offset + 1, nil //pos为"\n"的位置,需要加1才是行首的位置
		}
	}
	return 0, nil
}

func (self *File) Offset() int64 {
	return self.seek
}

func (self *File) Size() int64 {
	if stat, err := self.fd.Stat(); err != nil {
		return 0
	} else {
		return stat.Size()
	}
}

func (self *File) ReadAll() (file []byte, err error) {
	file, err = ioutil.ReadAll(self.fd)
	return
}

func (self *File) Inode() (uint64, error) {
	if stat, err := self.fd.Stat(); err != nil {
		return 0, err
	} else {
		inode := stat.Sys().(*syscall.Stat_t).Ino
		return inode, nil
	}
}
