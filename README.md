# ints
Ints for golang int128

# Twitter snowflake ID
`ints.NextID()` Get next ID  
`ints.DefaultSnowflakeGenerator.ServerID = 5` Set Server ID (0-31 when uses DataCenterID or 0-1023 when not)  
`ints.DefaultSnowflakeGenerator.DataCenterID = 7` Set data center ID (0-31)  

## SnowflakeGenerator
Generator for Twitter snowflake ID.  
Format: 1 - 0 | 41 - time in milliseconds | 5 - data center | 5 server | 12 sequence  
You can configure the following settings:  
1. `DataCenterID` - data center id; acceptable values 0-31  
1. `ServerID` - server id; acceptable values 0-31  
1. `TimeShift` - time shift from UNIXTIME in milliseconds; use `DefaultTimeShift` for compatible with Twitter snowflake ID  


# GUID

`ints.NextUUID()` Gets next sequentional GUID  
`ints.RandUUID()` Gets next realy random GUID  

### Bench test
```
BenchmarkUInt128Rand-16         20646118                54.48 ns/op           16 B/op          0 allocs/op
BenchmarkUInt128Next-16         13945237                77.94 ns/op            0 B/op          0 allocs/op
```

