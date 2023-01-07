# Go-Rel

This is a small library to assist in creating relational data in Go. 
It is suspected to offer some potential performance advantages to using standard Go maps, 
since data more easily occupies contiguous memory. 
The downside is, that only integer keys are supported. 

## Installation

To use Go-Rel in your project, import it

    import "github.com/DeltaBoyBZ/go-rel/rel"
    
and run

    $ go mod tidy

## Basic Usage
    
Go-Rel centres around the creation of a key value. 
Keys are used to identify unique values of an array. 
Some key values remain unused in some arrays, and this is managed by Go-Rel.

### Creating a Key with Fields  

    var ID rel.Key
    var a rel.Field[int]
    var b rel.Field[float32]
    ID.AddField(&a)
    ID.AddField(&b)
    
### Generating Key Values and Inserting Data

    key1 := ID.Gen()
    key2 := ID.Gen()
    a.Assign(key1, 69)
    b.Assign(key1, 3.14)
    b.Assign(key2, 2.72)
    
In this example, we have created two keys `key1` and `key2`. 
While field `b` has values associated with each of these keys, `a` has a value associated with `key1` only. 

### Deleting Values or Entire Keys

We can choose to 'delete' data in an array, in the sense that we mark they key as being used by that array. 

    b.Del(key1)
    
We may also delete a key entirely, meaning that key is free to be regenerated. 

    ID.Del(key2)
    
### Iterating Over All Relevant Array Values

If we want to iterate over all the used values of a key for a particular array,
we used the `GetUsed` method:

    sum := 0
    for _, elem := range a.GetUsed() {
        sum += a.Vals[elem] 
    }
    
### Desperate Allocation

The main disadvantage to implementating data structures at they are in Go-Rel, 
is that we end up with large pieces of virtually useless memory. 
Go-Rel allows us to make use of this memory,
by providing a means of *desperate allocation*. 
We can allocate data to occupy the unused parts of an array. 

    var foo Desperate[int, float32] // a desperate variable of type float32, to occupy an int array  
    foo.Alloc(&a)                   // allocates foo to occupy array a
    *foo.Get() = 3.14               // sets the desperately allocate variable to 3.14

Here the `Get` method gives the pointer to the data we're actually interested in. 
If there was not enough space in the array for our variable, 
then the variable is instead allocated outside the array, on the heap. 

