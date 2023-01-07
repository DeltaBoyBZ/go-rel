package rel

import "sync"
import "unsafe"

//The interace used by Key objects. To access all arrays using that Key. 
type Field interface {
    Push()
    Reserve(k int)
    Activate(k int)
    Del(k int) 
}

type GenericDesperate interface {
    Realloc()
    GetIndex() int
    GetSize() int
}

//Elements of arrays are indexed by keys. Keys can be generated and deleted, amounting to setting them to used or unsed. 
type Key struct {
    used []bool
    free []int 
    fields []Field
    genMutex sync.Mutex
}

//Initialise a key so it's ready for use.
func (key *Key) Init () {
    key.used = make([]bool, 0,  32)
    key.free = make([]int, 0, 32)
    key.fields = make([]Field, 0, 32)
}

//Add an array to the key, so that it may be managed alongside said key. 
func (key *Key) AddField (field Field) {
    key.fields = append(key.fields, field)
}

//If there are some free keys (from deletion) then reactivate an old key. Otherwise add a new higher value key. 
func (key *Key) Gen () int {
    key.genMutex.Lock()
    if len(key.free) == 0 {
        key.used = append(key.used, true)
        index := len(key.used) - 1
        for _, field := range key.fields {
            field.Push()
        }
        key.genMutex.Unlock()
        return index
    }
    index := key.free[0]
    key.used[index] = true
    key.free = key.free[1:]
    key.genMutex.Unlock()
    return index
}

//Indicate a key value as unused. 
func (key *Key) Del (index int) {
    key.used[index] = false
    for _ , free := range key.free {
        if free == index { return }
    }
    for _, field := range key.fields {
        field.Push()
        field.Del(index)
    }
    key.free = append(key.free, index)
}

//Constain a list of values, as well as information about which are relevant, and the encompassing key. 
type Array [T any] struct {
    Vals []T
    used []bool
    key  *Key
    desp []GenericDesperate
}

//Assigns a value to be associated with a particular key in this array. 
func (a *Array[T]) Assign(k int, val T) {
    for _, desp := range a.desp {
        if k >= desp.GetIndex() && k < desp.GetIndex() + desp.GetSize() {
            desp.Realloc()
        }
    }
    a.used[k] = true
    a.Vals[k] = val 
}

//Expands an array. 
func (a *Array[T]) Reserve(k int) {
    val_appendix := make([]T, k)  
    used_appendix := make([]bool, k)
    a.Vals = append(a.Vals, val_appendix...)
    a.used = append(a.used, used_appendix...)
    for i := 0 ; i < k ; i++ {
        a.used[len(a.used) - k - 1] = false
    }
}

//Adds an element to an array. 
func (a *Array[T]) Push() {
    var dummy T 
    a.Vals = append(a.Vals, dummy)
    a.used = append(a.used, false) 
}

//Activates an array element. 
func (a *Array[T]) Activate(k int) {
    a.used[k] = true
}

//Deletes an array element, in the sense that it is marked as unused.
func (a *Array[T]) Del(k int) {
    a.used[k] = false
}

//Returns a list of active array indices.
func (a *Array[T]) GetUsed () []int {
    used := make([]int, 0, 32)
    for index, elem := range a.used {
        if elem { 
            used = append(used, index) 
        }
    }
    return used
}

func (desp *Desparate[AllocType, ElemType]) Alloc (a *Array[ElemType]) bool {
    return desp.AllocWithOffset(a, -1)
}

func (desp *Desparate[AllocType, ElemType]) AllocWithOffset (a *Array[ElemType], offset int) bool {
    var dummyAllocType AllocType
    var dummyElemType ElemType 
    allocSize := unsafe.Sizeof(dummyAllocType)
    desp.Size = int(allocSize) 
    elemSize := unsafe.Sizeof(dummyElemType)
    available := 0
    availableStart := len(a.Vals)
    x := len(a.Vals) - 1
    if offset >= 0 { x = offset - 1 } 
    for i := x ; i >= 0 ; i += -1 {
        if a.used[i] { 
            available = 0
        } else {
            available += int(elemSize)
            availableStart = i 
        }
        if available >= int(allocSize) {
            desp.Arr = a
            desp.StartIndex = availableStart 
            a.desp = append(a.desp, desp)
            return true
        }
    }
    return false
}

type Desparate [AllocType any, ArrayElemType any] struct {
    Arr *Array[ArrayElemType]
    StartIndex int
    Fallback *AllocType
    Size int
}

func (desp *Desparate[AllocType, ElemType]) Get () *AllocType {
    if desp.Arr != nil {
        arb_pointer := unsafe.Pointer(&desp.Arr.Vals[desp.StartIndex])
        return (*AllocType)(arb_pointer)
    }
    return desp.Fallback
}

func (desp *Desparate[AllocType, ElemType]) Realloc() {
    var dummy AllocType
    dummy = *desp.Get()
    //first try to reallocate within array
    if !desp.AllocWithOffset(desp.Arr, desp.StartIndex) {
        desp.Fallback = new(AllocType)
        desp.Arr = nil
    }
    *desp.Get() = dummy
}

func (desp *Desparate[AllocType, ElemType]) GetIndex () int {
    return desp.StartIndex
}

func (desp *Desparate[AllocType, ElemType]) GetSize () int {
    return desp.Size
}




