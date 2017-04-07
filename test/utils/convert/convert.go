//package main
//
//import (
//	"server/utils/convert"
//	"fmt"
//)
//
//func main(){
//	s := convert.StringToBytes("123")
//	str := convert.BytesToString(s)
//	fmt.Println(str)
//
//
//}

package main

import (

	"bytes"

	"encoding/binary"

	"fmt"

)



func main() {

	//b  := []byte{0x00, 0x00, 0x03, 0xe8}
	//
	//b_buf  :=  bytes.NewBuffer(b)
	//
	//
	//var x int32
	//binary.Read(b_buf, binary.BigEndian, &x)
	//
	//fmt.Println(x)
	//
	//fmt.Println(strings.Repeat("-", 100))
	//

	b_buf  :=  bytes.NewBuffer([]byte{})

	binary.Write(b_buf, binary.BigEndian, int32(10))

	fmt.Println(b_buf.Bytes())

}