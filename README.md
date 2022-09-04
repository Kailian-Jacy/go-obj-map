# go-obj-map
A map indexed with object but not string.
**A toy project that cannot be used for industry.**

A map index could be a normal struct like:
```
type someStruct struct {
	// ... `objmap:"ignore"`
}
```
Struct name and type are considered for comparison of equality.

By default, struct instance of "different struct type" and "same struct type with different values" are both considered different.

```
objk = someStruct{...}

d := objmap.New()
d.Set(objk, v)

if v1, ok := d.Get(objk); ok {
	fmt.Println(v)
	fmt.Println(v1.(someStruct))
}
```
You may add tags to configure the equality judgement for different objects. For now, only "ignoring" some fields are permitted.