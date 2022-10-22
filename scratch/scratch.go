package scratch

import "fmt"

type TestStruct struct {
	someFloat float64
	someInt   int
	someUint8 byte
}

type ITestStruct interface {
	GetNewFloat() float64
}

func (testStruct *TestStruct) GetNewFloat() float64 {
	testStruct.someFloat += 3.14
	return testStruct.someFloat
}

func Run() {
	//do stuff
	var testStruct TestStruct
	doWork(&testStruct)
}

func doWork(testStruct *TestStruct) {
	test := testStruct.GetNewFloat()
	fmt.Printf("test:%f\n", test)

}
