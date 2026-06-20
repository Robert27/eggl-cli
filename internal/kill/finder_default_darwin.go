//go:build darwin

package kill

func platformFinder() Finder {
	return darwinFinder{}
}
