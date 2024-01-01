//go:build (!linux && !freebsd) || android

package main

func getLinuxPrimaryClipboard() string { return "" }
func setLinuxPrimaryClipboard(string)  {}
