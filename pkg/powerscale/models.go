package powerscale

type User struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

type ACL struct {
	Grantee    *AclUser `json:"grantee"`
	Permission string   `json:"permission"`
}

type AclUser struct {
	Type string `json:"type,omitempty"`
	// ID string `json:"id,omitempty"`
	Name string `json:"name"`
}

type Bucket struct {
	Description     string `json:"description"`
	Acl             []ACL  `json:"acl"`
	ObjectACLPolicy string `json:"object_acl_policy"`
	CreatePath      bool   `json:"create_path"`
	Owner           string `json:"owner"`
	Path            string `json:"path"`
	Name            string `json:"name"`
}

type PartialBucket struct {
	Acl []ACL `json:"acl"`
}

type BucketList struct {
	Buckets []*Bucket `json:"buckets"`
	Total   int       `json:"total"`
}

type Key struct {
	AccessID  string `json:"access_id"`
	SecretKey string `json:"secret_key"`
}

type Keys struct {
	Keys Key `json:"keys"`
}

type UserList struct {
	Users []*User `json:"users"`
	Total int     `json:"total"`
}
