package vm

import (
	"encoding/binary"

	"github.com/guilhermeg2k/dusklang/dusk"
)

func iLoadConst(stack *Stack, function *Function) {
	offset, err := function.readBytes(8)
	if err != nil {
		handleError(err)
	}
	offsetValue, _ := binary.Uvarint(offset)
	push(stack, function.Consts[offsetValue])
}

func iStore(stack *Stack, function *Function) {
	offset, err := function.readBytes(8)
	if err != nil {
		handleError(err)
	}
	offsetValue, _ := binary.Uvarint(offset)
	bytes := pop(stack, 8)
	store(function.Storage, offsetValue, bytes)
}

func iLoad(stack *Stack, function *Function) {
	offset, err := function.readBytes(8)
	if err != nil {
		handleError(err)
	}
	offsetValue, _ := binary.Uvarint(offset)
	bytes := load(function.Storage, offsetValue)
	push(stack, bytes)
}

func iAdd(stack *Stack) {
	right := pop(stack, 8)
	left := pop(stack, 8)
	res := dusk.IAdd(left, right)
	push(stack, res)
}

func iSub(stack *Stack) {
	right := pop(stack, 8)
	left := pop(stack, 8)
	res := dusk.ISub(left, right)
	push(stack, res)
}

func iMult(stack *Stack) {
	right := pop(stack, 8)
	left := pop(stack, 8)
	res := dusk.IMult(left, right)
	push(stack, res)
}

func iDiv(stack *Stack) {
	right := pop(stack, 8)
	left := pop(stack, 8)
	res := dusk.IDiv(left, right)
	push(stack, res)
}

func iMod(stack *Stack) {
	right := pop(stack, 8)
	left := pop(stack, 8)
	res := dusk.IMod(left, right)
	push(stack, res)
}

func fLoadConst(stack *Stack, function *Function) {
	offset, err := function.readBytes(8)
	if err != nil {
		handleError(err)
	}
	offsetValue, _ := binary.Uvarint(offset)
	push(stack, function.Consts[offsetValue])
}

func fStore(stack *Stack, function *Function) {
	offset, err := function.readBytes(8)
	if err != nil {
		handleError(err)
	}
	offsetValue, _ := binary.Uvarint(offset)
	bytes := pop(stack, 8)
	store(function.Storage, offsetValue, bytes)
}

func fLoad(stack *Stack, function *Function) {
	offset, err := function.readBytes(8)
	if err != nil {
		handleError(err)
	}
	offsetValue, _ := binary.Uvarint(offset)
	bytes := load(function.Storage, offsetValue)
	push(stack, bytes)
}

func bLoadConst(stack *Stack, function *Function) {
	offset, err := function.readBytes(8)
	if err != nil {
		handleError(err)
	}
	offsetValue, _ := binary.Uvarint(offset)
	push(stack, function.Consts[offsetValue])
}

func bStore(stack *Stack, function *Function) {
	offset, err := function.readBytes(8)
	if err != nil {
		handleError(err)
	}
	offsetValue, _ := binary.Uvarint(offset)
	bytes := pop(stack, 1)
	store(function.Storage, offsetValue, bytes)
}

func bLoad(stack *Stack, function *Function) {
	offset, err := function.readBytes(8)
	if err != nil {
		handleError(err)
	}
	offsetValue, _ := binary.Uvarint(offset)
	bytes := load(function.Storage, offsetValue)
	push(stack, bytes)
}

func fAdd(stack *Stack) {
	right := pop(stack, 8)
	left := pop(stack, 8)
	res := dusk.FAdd(left, right)
	push(stack, res)
}

func fSub(stack *Stack) {
	right := pop(stack, 8)
	left := pop(stack, 8)
	res := dusk.FSub(left, right)
	push(stack, res)
}

func fMult(stack *Stack) {
	right := pop(stack, 8)
	left := pop(stack, 8)
	res := dusk.FMult(left, right)
	push(stack, res)
}

func fDiv(stack *Stack) {
	right := pop(stack, 8)
	left := pop(stack, 8)
	res := dusk.FDiv(left, right)
	push(stack, res)
}

func iCmpEquals(stack *Stack) {
	right := pop(stack, 8)
	left := pop(stack, 8)
	res := dusk.ICmpEquals(left, right)
	push(stack, res)
}

func iCmpLessEquals(stack *Stack) {
	right := pop(stack, 8)
	left := pop(stack, 8)
	res := dusk.ICmpLessEquals(left, right)
	push(stack, res)
}

func iCmpGreaterEquals(stack *Stack) {
	right := pop(stack, 8)
	left := pop(stack, 8)
	res := dusk.ICmpGreaterEquals(left, right)
	push(stack, res)
}

func iCmpLessThen(stack *Stack) {
	right := pop(stack, 8)
	left := pop(stack, 8)
	res := dusk.ICmpLessThen(left, right)
	push(stack, res)
}

func iCmpGreaterThen(stack *Stack) {
	right := pop(stack, 8)
	left := pop(stack, 8)
	res := dusk.ICmpGreaterThen(left, right)
	push(stack, res)
}

func fCmpEquals(stack *Stack) {
	right := pop(stack, 8)
	left := pop(stack, 8)
	res := dusk.FCmpEquals(left, right)
	push(stack, res)
}

func fCmpLessEquals(stack *Stack) {
	right := pop(stack, 8)
	left := pop(stack, 8)
	res := dusk.FCmpLessEquals(left, right)
	push(stack, res)
}

func fCmpGreaterEquals(stack *Stack) {
	right := pop(stack, 8)
	left := pop(stack, 8)
	res := dusk.FCmpGreaterEquals(left, right)
	push(stack, res)
}

func fCmpLessThen(stack *Stack) {
	right := pop(stack, 8)
	left := pop(stack, 8)
	res := dusk.FCmpLessThen(left, right)
	push(stack, res)
}

func fCmpGreaterThen(stack *Stack) {
	right := pop(stack, 8)
	left := pop(stack, 8)
	res := dusk.FCmpGreaterThen(left, right)
	push(stack, res)
}

func jumpIfElse(stack *Stack, function *Function) {
	offset, err := function.readBytes(8)
	if err != nil {
		handleError(err)
	}
	offsetValue, _ := binary.Uvarint(offset)
	value := pop(stack, 1)
	if value[0] == 0 {
		function.CurrentOffset = function.Labels[offsetValue]
	}
}

func jumpIfTrue(stack *Stack, function *Function) {
	offset, err := function.readBytes(8)
	if err != nil {
		handleError(err)
	}
	offsetValue, _ := binary.Uvarint(offset)
	value := pop(stack, 1)
	if value[0] == 1 {
		function.CurrentOffset = function.Labels[offsetValue]
	}
}

//Temporary
func Print(stack *Stack) {
	value := pop(stack, 8)
	dusk.Print(value)
}
