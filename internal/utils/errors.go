package utils

import "errors"

var (
	// Client errors
	ErrValidation     = errors.New("validation failed") // input tidak valid
	ErrUnauthorized   = errors.New("unauthorized")      // tidak ada login / token invalid
	ErrForbidden      = errors.New("forbidden")         // tidak punya akses
	ErrNotFound       = errors.New("data not found")    // resource tidak ditemukan
	ErrConflict       = errors.New("conflict")          // sudah ada (duplicate)
	ErrTooManyRequest = errors.New("too many requests") // rate limit / throttle

	// Server errors
	ErrInternal    = errors.New("internal server error") // kesalahan server
	ErrUnavailable = errors.New("service unavailable")   // service down / maintenance
	ErrTimeout     = errors.New("request timeout")       // koneksi lama / gagal
)
