package time

import (
	"time"
)

func UnixTime() int32 {
	return int32(time.Now().Unix())
}

func UnixNanoTime() int64 {
	return time.Now().UnixNano()
}

//func main() {
//	//时间戳
//	t := time.Now().Unix()
//	fmt.Println(int32(t))
//
//	//时间戳到具体显示的转化
//	fmt.Println(time.Unix(t, 0).String())
//
//	//带纳秒的时间戳
//	t = time.Now().UnixNano()
//	fmt.Println(t)
//	fmt.Println("------------------")
//
//	//基本格式化的时间表示
//	fmt.Println(time.Now().String())
//
//	fmt.Println(time.Now().Format("2006year 01month 02day"))
//
//}

func main() {
	UnixTime()
}
