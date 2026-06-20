//go:build linux

package kill

func platformFinder() Finder {
	return linuxFinder{}
}
