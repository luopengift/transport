package file

import (
	"bufio"
	"github.com/luopengift/golibs/logger"
	"io"
	"os"
	"strings"
	"time"
)

type Tail struct {
	*File
	cname    string //config name
	line     chan *string
	reader   *bufio.Reader
	interval int64

	//EOF
	//@ true: stop
	//@ false: wait
	endstop bool
}

func NewTail(cname string) *Tail {
	name := HandlerRule(cname)
	file := NewFile(name, os.O_RDONLY)
	return &Tail{
		file,
		cname,
		make(chan *string),
		bufio.NewReader(file.fd),
		1000, //ms
		false,
	}
}

func (self *Tail) EndStop(b bool) {
	self.endstop = b
}

func (self *Tail) ReOpen() error {
	if err := self.Close(); err != nil {
		logger.Error("<file %v close fail:%v>", self.name, err)
	}
	self.name = HandlerRule(self.cname)
	err := self.Open()
	if err != nil {
		return err
	}
	self.reader = bufio.NewReader(self.fd)
	return nil
}

func (self *Tail) Stop() {
	self.Close()
	close(self.line)
}

func (self *Tail) ReadLine() {
	go func() {

		offset, err := self.TrancateOffsetByLF(self.seek)
		if err != nil {
			logger.Error("<Trancate offset:%d,Error:%+v>", self.seek, err)
		}
		err = self.Seek(offset)
		if err != nil {
			logger.Error("<seek offset[%d] error:%+v>", self.seek, err)
		}

		for {
			line, err := self.reader.ReadString('\n')
			switch {
			case err == io.EOF:
				if self.endstop {
					logger.Warn("<file %s is END:%+v>", self.name, err)
					self.Stop()
					return
				}
				time.Sleep(time.Duration(self.interval) * time.Millisecond)
				if self.name == self.cname {
					if inode, err := self.Inode(); err != nil { //检测是否需要重新打开新的文件
						continue
					} else {
						if inode != self.inode {
							self.ReOpen()
						}
					}
				} else {
					if self.name == HandlerRule(self.cname) { //检测是否需要按时间轮转新文件
						continue
					} else {
						self.ReOpen()
					}
				}

			case err != nil && err != io.EOF:
				time.Sleep(time.Duration(self.interval) * time.Millisecond)
				logger.Error("<Read file error:%v,%v>", line, err)
				self.ReOpen()
				continue
			default:
				msg := strings.TrimRight(line, "\n")
				self.line <- &msg
				self.seek += int64(len(line))
			}
		}
	}()
}

func (self *Tail) NextLine() chan *string {
	return self.line
}
