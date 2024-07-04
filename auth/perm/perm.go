package perm

import (
	apiAuth "github.com/wayne011872/api-toolkit/auth"
)

const (
	Guest  = apiAuth.ApiPerm("guest")
	Admin  = apiAuth.ApiPerm("admin")
	Owner  = apiAuth.ApiPerm("owner")
	Editor = apiAuth.ApiPerm("editor")
)
