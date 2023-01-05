package rel

import "sync"

//The interace used by Key objects. To access all arrays using that Key. 
type Field interface {
    Push()
    Reserve(k int)
    Activate(k int)
    Delete(k int) 
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
func (key *Key) Generate () int {
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
        field.Delete(index)
    }
    key.free = append(key.free, index)
}

//Constain a list of values, as well as information about which are relevant, and the encompassing key. 
type Array [T any] struct {
    Vals []T
    used []bool
    key  *Key
}

//Assigns a value to be associated with a particular key in this array. 
func (a *Array[T]) Assign(k int, val T) {
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
