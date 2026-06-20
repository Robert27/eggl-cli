//go:build !linux && !darwin && !windows

package kill

func platformFinder() Finder {
	return unsupportedFinder{}
}
