package har

type Cache struct {
	BeforeRequest CacheObject `json:"beforeRequest,omitempty"`
	AfterRequest  CacheObject `json:"afterRequest,omitempty"`
	Comment       string      `json:"comment,omitempty"`
}

type CacheObject struct {
	Expires    string `json:"expires,omitempty"`
	LastAccess string `json:"lastAccess"`
	ETag       string `json:"eTag"`
	HitCount   int    `json:"hitCount"`
	Comment    string `json:"comment,omitempty"`
}
