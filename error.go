package gerror

import (
	"bytes"
	//"encoding/binary"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"sync"
	"unsafe"
)

/**
 * implement error interface
 *
 *	tyep error interface {
 *		Error()string
 *  }
 * sopport error stack.
 * errno int32 indicated.
 */

const (
	// .
	POSITION_CODE_SIZE = 0
	// panic error
	CONST_ERRNO_PANIC = -100000
)

type ToString interface {
	String() string
}

type gerror struct {
	code int
	mem  *bytes.Buffer //[]byte
	//pointer int
}

var err_pool sync.Pool
var _ io.Writer = &gerror{}

func init() {
	err_pool.New = func() interface{} {
		return &gerror{
			//pointer: POSITION_CODE_SIZE,
			//mem:     make([]byte, 20),
			//mem: new(bytes.Buffer),
		}
	}
}

func New(args ...interface{}) *gerror {
	er := &gerror{
		//pointer: POSITION_CODE_SIZE,
		//mem: new(bytes.Buffer),
	}
	er.W(args...)

	if er.GetCode() == CONST_ERRNO_PANIC {
		panic(er.Error())
	}
	return er
}

func Panic(args ...interface{}) {
	er := New(args...)
	if er.GetCode() == 0 {
		er.SetCode(CONST_ERRNO_PANIC)
	}
	panic(er.Error())
}

func Error(args ...interface{}) string {
	er := err_pool.Get().(*gerror)
	defer err_pool.Put(er)

	er.W(args...)
	str := er.Error()
	er.Reset()
	return str
}

// just use ftm.Srpintf
func Errorf(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

func (err *gerror) Error() string {
	//return ""
	return BytesToString(err.mem.Bytes())
}

func (err *gerror) Reset() {
	err.code = 0
	err.mem.Reset()
}

func (err *gerror) ErrorBytes() []byte {
	return err.mem.Bytes()
}

func (err *gerror) String() string {
	return err.mem.String()
}

func (err *gerror) GetCode() int {
	return err.code
}

func (err *gerror) SetCode(code int) {
	err.code = code

	str := strconv.Itoa(code)

	if err.mem == nil {
		err.mem = bytes.NewBufferString(str)
		err.mem.Reset()
	}
	err.mem.WriteString(str)
	err.mem.WriteByte(':')

}

func (err *gerror) Write(bs []byte) (n int, errr error) {
	if err.mem == nil {
		err.mem = bytes.NewBuffer(bs)
		err.mem.Reset()
	}
	err.mem.WriteByte(' ')
	err.mem.Write(bs)
	n = err.mem.Len()
	return
}

func (err *gerror) WriteString(e string) {
	if err.mem == nil {
		err.mem = bytes.NewBufferString(e)
		err.mem.Reset()
	}
	err.mem.WriteByte(' ')
	err.mem.WriteString(e)
}

func (err *gerror) W(args ...interface{}) {
	if len(args) == 0 {
		return
	}

	code, ok := args[0].(int)
	begin := 0
	if ok {
		err.SetCode(code)
		begin++
	}
	for i := begin; i < len(args); i++ {
		if b, ok := args[i].([]byte); ok {
			err.Write(b)
			continue
		}

		if b, ok := args[i].(byte); ok {
			err.Write([]byte{b})
			continue
		}

		if b, ok := args[i].(string); ok {
			err.WriteString(b)
			continue
		}

		if b, ok := args[i].(ToString); ok {
			err.WriteString(b.String())
			continue
		}
		if b, ok := args[i].(int); ok {
			err.WriteString(strconv.Itoa(b))
			continue
		}
	}
}

func (err *gerror) Size() int {
	//return err.pointer
	return err.mem.Len()
}

func test_r(args ...interface{}) string {
	if len(args) == 0 {
		return ""
	}
	var str = ""
	for _, v := range args {
		if b, ok := v.([]byte); ok {
			str += string(b)
			continue
		}

		if b, ok := v.(byte); ok {
			str += string([]byte{b})
			continue
		}

		if b, ok := v.(string); ok {
			str += b
			continue
		}

		if b, ok := v.(ToString); ok {
			str += (b.String())
			continue
		}
		if b, ok := v.(int); ok {
			str += (strconv.Itoa(b))
			continue
		}

	}
	return str
}

// ======================utile ===========================
//func Int32ToBytes(i int) []byte {
//	var buf = make([]byte, 4)
//	binary.BigEndian.PutUint32(buf, uint32(i))
//	return buf
//}
//
//func BytesToInt32(buf []byte) int {
//	return int(binary.BigEndian.Uint32(buf))
//}

// StringToBytes converts string to byte slice without a memory allocation.
func StringToBytes(s string) (b []byte) {
	sh := *(*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh.Data, bh.Len, bh.Cap = sh.Data, sh.Len, sh.Len
	return b
}

// BytesToString converts byte slice to string without a memory allocation.
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
