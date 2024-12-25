package bean

const (
	GoogleIDEmpty = "00000000-0000-0000-0000-000000000000"

	//platform
	PlatformMobileIos      = "mobile-ios"
	PlatformMobileAndroid  = "mobile-android"
	PlatformDesktopWindows = "desktop-windows"
	PlatformDesktopMacOs   = "desktop-macos"

	//channel
	ChannelApple     = "apple"
	ChannelGoogle    = "google"
	ChannelUniversal = "universal"
)

func CheckPlatform(platform string) bool {
	return platform == PlatformMobileIos || platform == PlatformMobileAndroid || platform == PlatformDesktopWindows || platform == PlatformDesktopMacOs
}

func CheckChannel(channel string) bool {
	return channel == ChannelApple || channel == ChannelGoogle || channel == ChannelUniversal
}
