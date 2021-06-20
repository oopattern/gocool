package util

import (
	"bufio"
	"github.com/oopattern/gocool/log"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// 文件行的处理函数接口
type Parser interface {
	// line表示行内容, tp表示处理行的方式
	ParseLine(line string, tp string) error
}

// 解析文件, 按行解析
func HandleFile(file string, tp string, pFunc Parser) error {
	f, err := os.OpenFile(file, os.O_RDONLY, 0666)
	if err != nil {
		log.Error("open file[%s] err[%v]", file, err)
		return err
	}

	defer f.Close()
	buf := bufio.NewReader(f)
	// 死循环怎么办
	for {
		line, _, err := buf.ReadLine()
		if io.EOF == err {
			log.Info("file end")
			break
		}
		if err != nil {
			log.Error("file[%s] read line err[%v]", file, err)
			continue
		}

		strline := string(line)
		strline = strings.TrimSpace(strline)

		// 执行行处理行数, 通过接口方式
		if err := pFunc.ParseLine(strline, tp); err != nil {
			log.Error("parse line handle error")
			continue
		}
	}

	return nil
}

// 将数据按行记录到文件中, 文件名称为绝对路径
func AppendFile(file string, line string) {
	// 先检查目录是否存在, 如果目录不存在, 生成上层目录
	var _ = os.MkdirAll(filepath.Dir(file), os.ModePerm)

	f, err := os.OpenFile(file, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		log.Error("open file[%s] err[%v]", file, err)
		return
	}

	defer f.Close()

	w := bufio.NewWriter(f)
	s := line + "\n"
	if n, err := w.WriteString(s); err != nil {
		log.Error("file[%s] write[%s] n[%d] len[%d] error", file, s, n, len(s))
	}
	if err := w.Flush(); err != nil {
		log.Error("file[s] flush error[%v]", file, err)
	}
}
