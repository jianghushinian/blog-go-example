package main

import (
	"io"
	"os"
)

var (
	version        = "dev"
	commit         = "none"
	builtGoVersion = "unknown"
	changeLogPath  = "CHANGELOG.md"
)

type ChangeLogSpec struct {
	Version        string
	Commit         string
	BuiltGoVersion string
	ChangeLog      string
}

func GetChangeLog() (ChangeLogSpec, error) {
	data, err := os.ReadFile(changeLogPath)
	if err != nil {
		return ChangeLogSpec{}, err
	}

	return ChangeLogSpec{
		Version:        version,
		Commit:         commit,
		BuiltGoVersion: builtGoVersion,
		ChangeLog:      string(data),
	}, nil
}

func GetChangeLogByIOReader(reader io.Reader) (ChangeLogSpec, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return ChangeLogSpec{}, err
	}

	return ChangeLogSpec{
		Version:        version,
		Commit:         commit,
		BuiltGoVersion: builtGoVersion,
		ChangeLog:      string(data),
	}, nil
}
