package har

type Cache struct {
	BeforeRequest CacheObject `json:"beforeRequest,omitzero"`
	AfterRequest  CacheObject `json:"afterRequest,omitzero"`
	Comment       string      `json:"comment,omitzero"`
}

type CacheObject struct {
	Expires    string `json:"expires,omitzero"`
	LastAccess string `json:"lastAccess"`
	ETag       string `json:"eTag"`
	HitCount   int    `json:"hitCount"`
	Comment    string `json:"comment,omitzero"`
}
