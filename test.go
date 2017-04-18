
package main

import (
	"strings"
	"fmt"
	"strconv"
)

func main()  {
	// bytes:= []byte{10,12,99,114,101,97,116,101,67,105,114,99,108,101,10,6,51,50,52,50,50,56,10,6,97,107,115,104,97,121,10,4,55,46,51,50}
	bytes:= []byte{38,240,239,242,254,125,54,16,115,72,237,63,178,162,236,200,98,212,245,7,164,227,169,213,237,93,28,217,223,115,198,196}
	str:= convert(bytes[:])
	fmt.Println(str)
}

func convert( b []byte ) string {
    s := make([]string,len(b))
    for i := range b {
        s[i] = strconv.Itoa(int(b[i]))
    }
    return strings.Join(s,",")
}
