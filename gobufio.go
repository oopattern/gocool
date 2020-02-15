package main

import (
	"bufio"
	"fmt"
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

func main() {
	// TestWriteBufio()
	// TestBatchWrite()
	// TestResetWrite()
	TestPeekRead()
}