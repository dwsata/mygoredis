package caches

type Options struct {
	MaxEntrySize int64  //unit GB,write protect
	MaxGCCount   int    //auto expire
	GCDuration   int64  //minutes,auto expire interval
	DumpFile     string //file path
	DumpDuration int64  //minutes
}

func DefaultOptions() *Options {
	return &Options{
		MaxEntrySize: int64(4), // 4gb
		MaxGCCount:   1000,     // 1000
		GCDuration:   60,       //1hours
		DumpFile:     "cache-server.dump",
		DumpDuration: 30,
	}
}
