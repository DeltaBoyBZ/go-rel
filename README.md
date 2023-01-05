# Go-Rel

This is a small library to assist in creating relational data in Go. 
It is suspected to offer some potential performance advantages to using standard Go maps, 
since data more easily occupies contiguous memory. 
The downside is, that only integer keys are supported. 

## Installation

To use Go-Rel in your project, download this repository and include the following in your `go.mod` file:

    require "delta/rel" v0.0.0
    replace "delta/rel" v0.0.0 => "<location of the `rel` subdirectory of this repository on your disk>"
    
## Basic Usage

To  use Go-Rel, we must of course import it:

    import "delta/rel"
    
Go-Rel centres around the creation of a key value. 
Keys are used to identify unique values of an array. 
Some key values remain unused in some arrays, and this is managed by Go-Rel.

### Creating a Key with Fields  

    var ID rel.Key
    var a rel.Field[int]
    var b rel.Field[string]
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
    
