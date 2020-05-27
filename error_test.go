package gerror

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	er := New(-1001, "error cur", 'a', []byte("hello world!"), "-23234343")
	fmt.Println(er.Error())
	fmt.Println(er.String())
	assert.Equal(t, er.GetCode(), -1001)

	l1 := len(er.ErrorBytes())
	l2 := len(er.String())
	assert.Equal(t, l1, len(er.Error()))
	assert.Equal(t, l2, len(er.Error()))
}

func TestPanicFunc(t *testing.T) {
	defer test_recovery()

	Panic(-100001, "painc can be modify by recover!")
}

func TestPanicFuncNoErrno(t *testing.T) {
	defer test_recovery()
	Panic("no error number .painc can be modify by recover!")
}

func TestPanicFuncPanicErrno(t *testing.T) {
	defer test_recovery()
	New(CONST_ERRNO_PANIC, "some error")
}

func TestErrorFunc(t *testing.T) {
	// error return error string.
	str := Error(-1002, "2323", "what error it is?")
	fmt.Println(str)

	fmt.Println("FORMAT FMT.PRINTF", Errorf("[%s]", "dfdf"))

	er := New("error length need to longer than 50.", "I am in", "who want join to ous!", "enough to 50", 232)
	assert.Equal(t, true, 50 < er.Size())
	er_byte := New([]byte("error length need to longer than 50." + "I am in" + "who want join to ous!" + "enough to 50"))
	assert.Equal(t, true, 50 < er_byte.Size())

}

func TestAcceptStruct(t *testing.T) {
	a := &strings_stu{
		name: "aixgle",
	}
	er := New(a)

	fmt.Println("testing struct implement String()strin,", "err:", er.Size(), (er.Error()), "name:", len(a.String()), a.String())
	assert.True(t, er.Error() == " "+a.String())
}

func TestEmpty(t *testing.T) {
	New()
	test_r(3234, 'a', "want a good work! hard !", []byte("now you cand do better"))
}

type strings_stu struct {
	name string
}

func (ss *strings_stu) String() string {
	return ss.name
}

var _ ToString = &strings_stu{}

func test_recovery() {
	if err := recover(); err != nil {
		fmt.Println(err)
	}
}
