//go:build windows

package kill

func platformFinder() Finder {
	return windowsFinder{}
}
