# ints
Ints for golang int128

# GUID

`ints.NextUUID()` Gets next sequentional GUID  
`ints.RandUUID()` Gets next realy random GUID  

### Bench test
```
BenchmarkUInt128Rand-16         20646118                54.48 ns/op           16 B/op          0 allocs/op
BenchmarkUInt128Next-16         13945237                77.94 ns/op            0 B/op          0 allocs/op
```

