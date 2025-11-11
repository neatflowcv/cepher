package repository

import "errors"

var (
	ErrClusterAlreadyExists = errors.New("cluster already exists")
	ErrClusterNotFound      = errors.New("cluster not found")
)
