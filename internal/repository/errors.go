package repository

import "errors"

var (
	// ErrNotFound 记录不存在
	ErrNotFound = errors.New("record not found")

	// ErrDuplicate 记录重复
	ErrDuplicate = errors.New("record already exists")
)
