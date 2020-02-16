package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

type Writer int
type Wrt2 int

func (*Writer) Write(p []byte) (n int, err error)  {
	fmt.Printf("write: %q\n", p)
	return len(p), nil
	// return len(p), errors.New("boom!")
}

func (*Wrt2) Write(p []byte) (n int, err error)  {
	fmt.Printf("wrt2: %q\n", p)
	return len(p), nil
}

func TestWriteBufio()  {
	fmt.Println("Unbuffer IO")
	w := new(Writer)
	w.Write([]byte{'a'})
	w.Write([]byte{'b'})
	w.Write([]byte{'c'})
	w.Write([]byte{'d'})
	fmt.Println("Buffer IO")
	bw := bufio.NewWriterSize(w, 3)
	bw.Write([]byte{'a'})
	bw.Write([]byte{'b'})
	bw.Write([]byte{'c'})
	bw.Write([]byte{'d'})
	err := bw.Flush()
	if err != nil {
		fmt.Println(err)
	}
}

func TestBatchWrite()  {
	w := new(Writer)
	bw := bufio.NewWriterSize(w, 3)
	bw.Write([]byte("abcd"))
	err := bw.Flush()
	if err != nil {
		fmt.Println(err)
	}
}

func TestResetWrite()  {
	w := new(Writer)
	w2 := new(Wrt2)
	bw := bufio.NewWriterSize(w, 2)
	bw.Write([]byte("ab"))
	bw.Write([]byte("cd"))
	// bw.Flush()
	bw.Reset(w2)
	bw.Write([]byte("ef"))
	bw.Flush()
}

func TestPeekRead()  {
	s := strings.NewReader(strings.Repeat("a", 20))
	r := bufio.NewReaderSize(s, 16)
	b, err := r.Peek(3)
	if err != nil {
	 	fmt.Println(err)
	}
	fmt.Printf("%q\n", b)
	b, err = r.Peek(17)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%q\n", b)
	sn := strings.NewReader("aaa")
	r.Reset(sn)
	b, err = r.Peek(10)
	if err != nil {
		fmt.Println(err)
	}
}

func TestRead()  {
	s := strings.NewReader(strings.Repeat("a", 16) + strings.Repeat("b", 16))
	r := bufio.NewReaderSize(s, 16)
	b, _ := r.Peek(3)
	fmt.Printf("%q\n", b)
	r.Read(make([]byte, 16))
	r.Read(make([]byte, 15))
	// b共用bufio底层的buffer,
	// 第一次Read会将string的前16个字符a读到bufio底层buffer中, 此时打印b会显示aaa
	// 第二次Read会将string的后16个字符b读到bufio底层的buffer中, 此时打印b会显示bbb
	fmt.Printf("%q\n", b)
}

type R struct {
	n int
}

func (r *R) Read(b []byte) (n int, err error) {
	// fmt.Println("Read")
	// copy(b, "abcdefghijklmnop")
	// return 16, nil
	fmt.Printf("read#: %d\n", r.n)
	if r.n >= 10 {
		return 0, io.EOF
	}
	copy(b, "abcd")
	r.n += 1
	return 4, nil
}

func TestDiscardRead()  {
	r := new(R)
	br := bufio.NewReaderSize(r, 16)
	buf := make([]byte, 4)
	// 第一次Read会将调用io.Read读入16个字符到bufio底层buffer中,
	// 然后从bufio底层buffer取出4个字符, bufio底层buffer还剩12个字符未取出
	br.Read(buf)
	fmt.Printf("%q\n", buf)
	// 要丢弃13个字符超过bufio底层buffer的12个字符, 剩余要丢弃的1个字符, 触发第2次调用io.Read
	br.Discard(13)
	br.Read(buf)
	fmt.Printf("%q\n", buf)
}

func TestWriteTo()  {
	br := bufio.NewReaderSize(new(R), 16)
	// 这个功能要再研究一下...
	n, err := br.WriteTo(ioutil.Discard)
	if err != nil {
		panic(err)
	}
	fmt.Printf("writen bytes:%d\n", n)
}

func main() {
	// TestWriteBufio()
	// TestBatchWrite()
	// TestResetWrite()
	// TestPeekRead()
	// TestRead()
	// TestDiscardRead()
	TestWriteTo()
}