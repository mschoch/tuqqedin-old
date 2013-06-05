package main

// #cgo LDFLAGS: -lantlr3c
// #include "UNQL2013Util.h"
import (
	"C"
)

import (
	"log"
)

//export GoWalkTree
func GoWalkTree(tree *C.struct_ANTLR3_BASE_TREE_struct) {

	s := C.myPrintTree(tree)
	gs := C.GoString(s)

	log.Printf("Gone walking %s", gs)

	if tree.children == nil || C.myVectorSize(tree.children) == 0 {
		return
	}

	childrenSize := C.myVectorSize(tree.children)
	for i := C.ANTLR3_UINT32(0); i < childrenSize; i++ {
		child := C.myVectorGet(tree.children, i)
		GoWalkTree(child)
	}
}

func parseUNQL(input string) {
	_ = C.parse(C.CString(input))
	log.Printf("after all")

	//x.tree.toStringTree(x.tree)

	//log.Printf("%+v", x)
	//log.Printf("%T", x.start)
	//log.Printf("%+v", x.start)
	//log.Printf("%T", x.tree)
	//log.Printf("%+v", x.tree)

	//log.Printf("%T", x.tree.super)
	//log.Printf("Here its:%+v", x.tree.super)
 
	//var y C.pANTLR3_COMMON_TREE = (C.pANTLR3_COMMON_TREE)(x.tree.super)

	//log.Printf("%T", y)
	//log.Printf("Now its: %+v", y)
	//log.Printf("Now its really type: %T", *y)
	//log.Printf("Now its really: %+v", *y)

	//log.Printf("Now its really type: %T", *(y.token))
	//log.Printf("Now its really: %+v", *(y.token))

}

func main() {
	parseUNQL("SELECT * WHERE abv > 7")
}
