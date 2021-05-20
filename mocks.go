package main

//go:generate mockgen -package connlib_mock -destination ./internal/connlib_mock/mock_connection_source.go github.com/kbence/conndetect/internal/connlib ConnectionSource
//go:generate mockgen -package utils_mock -destination ./internal/utils_mock/mocks.go github.com/kbence/conndetect/internal/utils Printer,Time
