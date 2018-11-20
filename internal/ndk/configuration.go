// android app

// +build android

package app

import (
	"unsafe"
)

/*
#include <android/configuration.h>
#if __ANDROID_API__ < 13
#define ACONFIGURATION_DENSITY_TV											   0
#define ACONFIGURATION_DENSITY_XHIGH                                           0
#define ACONFIGURATION_DENSITY_XXHIGH                                          0
#define ACONFIGURATION_DENSITY_XXXHIGH                                         0
#define ACONFIGURATION_LAYOUTDIR                                               0
#define ACONFIGURATION_LAYOUTDIR_ANY                                           0
#define ACONFIGURATION_LAYOUTDIR_LTR                                           0
#define ACONFIGURATION_LAYOUTDIR_RTL                                           0
#define ACONFIGURATION_SCREEN_HEIGHT_DP_ANY                                    0
#define ACONFIGURATION_SCREEN_WIDTH_DP_ANY                                     0
#define ACONFIGURATION_SMALLEST_SCREEN_SIZE                                    0
#define ACONFIGURATION_SMALLEST_SCREEN_WIDTH_DP_ANY                            0
#define ACONFIGURATION_UI_MODE_TYPE_APPLIANCE                                  0
#define ACONFIGURATION_UI_MODE_TYPE_TELEVISION                                 0
static int AConfiguration_getLayoutDirection									(AConfiguration* obj) { return 0; }
static int AConfiguration_getScreenHeightDp                                    (AConfiguration* obj) { return 0; }
static int AConfiguration_getScreenWidthDp                                     (AConfiguration* obj) { return 0; }
static int AConfiguration_getSmallestScreenWidthDp                             (AConfiguration* obj) { return 0; }
static int AConfiguration_setLayoutDirection                                   (AConfiguration* obj, int32_t v) { return 0; }
static int AConfiguration_setScreenHeightDp                                    (AConfiguration* obj, int32_t v) { return 0; }
static int AConfiguration_setScreenWidthDp                                     (AConfiguration* obj, int32_t v) { return 0; }
static int AConfiguration_setSmallestScreenWidthDp								(AConfiguration* obj, int32_t v) { return 0; }

#endif
*/
import "C"

const (
	CONFIGURATION_ORIENTATION_ANY    = C.ACONFIGURATION_ORIENTATION_ANY
	CONFIGURATION_ORIENTATION_PORT   = C.ACONFIGURATION_ORIENTATION_PORT
	CONFIGURATION_ORIENTATION_LAND   = C.ACONFIGURATION_ORIENTATION_LAND
	CONFIGURATION_ORIENTATION_SQUARE = C.ACONFIGURATION_ORIENTATION_SQUARE

	CONFIGURATION_TOUCHSCREEN_ANY     = C.ACONFIGURATION_TOUCHSCREEN_ANY
	CONFIGURATION_TOUCHSCREEN_NOTOUCH = C.ACONFIGURATION_TOUCHSCREEN_NOTOUCH
	CONFIGURATION_TOUCHSCREEN_STYLUS  = C.ACONFIGURATION_TOUCHSCREEN_STYLUS
	CONFIGURATION_TOUCHSCREEN_FINGER  = C.ACONFIGURATION_TOUCHSCREEN_FINGER

	CONFIGURATION_DENSITY_DEFAULT = C.ACONFIGURATION_DENSITY_DEFAULT
	CONFIGURATION_DENSITY_LOW     = C.ACONFIGURATION_DENSITY_LOW
	CONFIGURATION_DENSITY_MEDIUM  = C.ACONFIGURATION_DENSITY_MEDIUM
	CONFIGURATION_DENSITY_TV      = C.ACONFIGURATION_DENSITY_TV
	CONFIGURATION_DENSITY_HIGH    = C.ACONFIGURATION_DENSITY_HIGH
	CONFIGURATION_DENSITY_XHIGH   = C.ACONFIGURATION_DENSITY_XHIGH
	CONFIGURATION_DENSITY_XXHIGH  = C.ACONFIGURATION_DENSITY_XXHIGH
	CONFIGURATION_DENSITY_XXXHIGH = C.ACONFIGURATION_DENSITY_XXXHIGH
	CONFIGURATION_DENSITY_NONE    = C.ACONFIGURATION_DENSITY_NONE

	CONFIGURATION_KEYBOARD_ANY    = C.ACONFIGURATION_KEYBOARD_ANY
	CONFIGURATION_KEYBOARD_NOKEYS = C.ACONFIGURATION_KEYBOARD_NOKEYS
	CONFIGURATION_KEYBOARD_QWERTY = C.ACONFIGURATION_KEYBOARD_QWERTY
	CONFIGURATION_KEYBOARD_12KEY  = C.ACONFIGURATION_KEYBOARD_12KEY

	CONFIGURATION_NAVIGATION_ANY       = C.ACONFIGURATION_NAVIGATION_ANY
	CONFIGURATION_NAVIGATION_NONAV     = C.ACONFIGURATION_NAVIGATION_NONAV
	CONFIGURATION_NAVIGATION_DPAD      = C.ACONFIGURATION_NAVIGATION_DPAD
	CONFIGURATION_NAVIGATION_TRACKBALL = C.ACONFIGURATION_NAVIGATION_TRACKBALL
	CONFIGURATION_NAVIGATION_WHEEL     = C.ACONFIGURATION_NAVIGATION_WHEEL

	CONFIGURATION_KEYSHIDDEN_ANY  = C.ACONFIGURATION_KEYSHIDDEN_ANY
	CONFIGURATION_KEYSHIDDEN_NO   = C.ACONFIGURATION_KEYSHIDDEN_NO
	CONFIGURATION_KEYSHIDDEN_YES  = C.ACONFIGURATION_KEYSHIDDEN_YES
	CONFIGURATION_KEYSHIDDEN_SOFT = C.ACONFIGURATION_KEYSHIDDEN_SOFT

	CONFIGURATION_NAVHIDDEN_ANY = C.ACONFIGURATION_NAVHIDDEN_ANY
	CONFIGURATION_NAVHIDDEN_NO  = C.ACONFIGURATION_NAVHIDDEN_NO
	CONFIGURATION_NAVHIDDEN_YES = C.ACONFIGURATION_NAVHIDDEN_YES

	CONFIGURATION_SCREENSIZE_ANY    = C.ACONFIGURATION_SCREENSIZE_ANY
	CONFIGURATION_SCREENSIZE_SMALL  = C.ACONFIGURATION_SCREENSIZE_SMALL
	CONFIGURATION_SCREENSIZE_NORMAL = C.ACONFIGURATION_SCREENSIZE_NORMAL
	CONFIGURATION_SCREENSIZE_LARGE  = C.ACONFIGURATION_SCREENSIZE_LARGE
	CONFIGURATION_SCREENSIZE_XLARGE = C.ACONFIGURATION_SCREENSIZE_XLARGE

	CONFIGURATION_SCREENLONG_ANY = C.ACONFIGURATION_SCREENLONG_ANY
	CONFIGURATION_SCREENLONG_NO  = C.ACONFIGURATION_SCREENLONG_NO
	CONFIGURATION_SCREENLONG_YES = C.ACONFIGURATION_SCREENLONG_YES

	CONFIGURATION_UI_MODE_TYPE_ANY        = C.ACONFIGURATION_UI_MODE_TYPE_ANY
	CONFIGURATION_UI_MODE_TYPE_NORMAL     = C.ACONFIGURATION_UI_MODE_TYPE_NORMAL
	CONFIGURATION_UI_MODE_TYPE_DESK       = C.ACONFIGURATION_UI_MODE_TYPE_DESK
	CONFIGURATION_UI_MODE_TYPE_CAR        = C.ACONFIGURATION_UI_MODE_TYPE_CAR
	CONFIGURATION_UI_MODE_TYPE_TELEVISION = C.ACONFIGURATION_UI_MODE_TYPE_TELEVISION
	CONFIGURATION_UI_MODE_TYPE_APPLIANCE  = C.ACONFIGURATION_UI_MODE_TYPE_APPLIANCE

	CONFIGURATION_UI_MODE_NIGHT_ANY = C.ACONFIGURATION_UI_MODE_NIGHT_ANY
	CONFIGURATION_UI_MODE_NIGHT_NO  = C.ACONFIGURATION_UI_MODE_NIGHT_NO
	CONFIGURATION_UI_MODE_NIGHT_YES = C.ACONFIGURATION_UI_MODE_NIGHT_YES

	CONFIGURATION_SCREEN_WIDTH_DP_ANY = C.ACONFIGURATION_SCREEN_WIDTH_DP_ANY

	CONFIGURATION_SCREEN_HEIGHT_DP_ANY = C.ACONFIGURATION_SCREEN_HEIGHT_DP_ANY

	CONFIGURATION_SMALLEST_SCREEN_WIDTH_DP_ANY = C.ACONFIGURATION_SMALLEST_SCREEN_WIDTH_DP_ANY

	CONFIGURATION_LAYOUTDIR_ANY = C.ACONFIGURATION_LAYOUTDIR_ANY
	CONFIGURATION_LAYOUTDIR_LTR = C.ACONFIGURATION_LAYOUTDIR_LTR
	CONFIGURATION_LAYOUTDIR_RTL = C.ACONFIGURATION_LAYOUTDIR_RTL

	CONFIGURATION_MCC                  = C.ACONFIGURATION_MCC
	CONFIGURATION_MNC                  = C.ACONFIGURATION_MNC
	CONFIGURATION_LOCALE               = C.ACONFIGURATION_LOCALE
	CONFIGURATION_TOUCHSCREEN          = C.ACONFIGURATION_TOUCHSCREEN
	CONFIGURATION_KEYBOARD             = C.ACONFIGURATION_KEYBOARD
	CONFIGURATION_KEYBOARD_HIDDEN      = C.ACONFIGURATION_KEYBOARD_HIDDEN
	CONFIGURATION_NAVIGATION           = C.ACONFIGURATION_NAVIGATION
	CONFIGURATION_ORIENTATION          = C.ACONFIGURATION_ORIENTATION
	CONFIGURATION_DENSITY              = C.ACONFIGURATION_DENSITY
	CONFIGURATION_SCREEN_SIZE          = C.ACONFIGURATION_SCREEN_SIZE
	CONFIGURATION_VERSION              = C.ACONFIGURATION_VERSION
	CONFIGURATION_SCREEN_LAYOUT        = C.ACONFIGURATION_SCREEN_LAYOUT
	CONFIGURATION_UI_MODE              = C.ACONFIGURATION_UI_MODE
	CONFIGURATION_SMALLEST_SCREEN_SIZE = C.ACONFIGURATION_SMALLEST_SCREEN_SIZE
	CONFIGURATION_LAYOUTDIR            = C.ACONFIGURATION_LAYOUTDIR
)

type Configuration C.AConfiguration

func (config *Configuration) cptr() *C.AConfiguration {
	return (*C.AConfiguration)(config)
}

/**
 * Create a new AConfiguration, initialized with no values set.
 */
func NewConfiguration() *Configuration {
	return (*Configuration)(C.AConfiguration_new())
}

/**
 * Free an AConfiguration that was previously created with
 * AConfiguration_new().
 */
func (config *Configuration) Delete() {
	C.AConfiguration_delete(config.cptr())
}

/**
 * Create and return a new AConfiguration based on the current configuration in
 * use in the given AssetManager.
 */
func fromAssetManager(am *AssetManager) *Configuration {
	config := NewConfiguration()
	C.AConfiguration_fromAssetManager(config.cptr(), am.cptr())
	return config
}

/**
 * Copy the contents of 'src' to 'dest'.
 */
func (config *Configuration) Copy(dest *Configuration) {
	C.AConfiguration_copy(dest.cptr(), config.cptr())
}

/**
 * Return the current MCC set in the configuration.  0 if not set.
 */
//int32_t AConfiguration_getMcc(AConfiguration* config);
func (config *Configuration) GetMcc() int {
	return int(C.AConfiguration_getMcc(config.cptr()))
}

/**
 * Set the current MCC in the configuration.  0 to clear.
 */
//void AConfiguration_setMcc(AConfiguration* config, int32_t mcc);
func (config *Configuration) SetMcc(mcc int) {
	C.AConfiguration_setMcc(config.cptr(), C.int32_t(mcc))
}

/**
 * Return the current MNC set in the configuration.  0 if not set.
 */
//int32_t AConfiguration_getMnc(AConfiguration* config);
func (config *Configuration) GetMnc() int {
	return int(C.AConfiguration_getMnc(config.cptr()))
}

/**
 * Set the current MNC in the configuration.  0 to clear.
 */
//void AConfiguration_setMnc(AConfiguration* config, int32_t mnc);
func (config *Configuration) SetMnc(mnc int) {
	C.AConfiguration_setMnc(config.cptr(), C.int32_t(mnc))
}

/**
 * Return the current language code set in the configuration.  The output will
 * be filled with an array of two characters.  They are not 0-terminated.  If
 * a language is not set, they will be 0.
 */
//void AConfiguration_getLanguage(AConfiguration* config, char* outLanguage);
func (config *Configuration) GetLanguage() []byte {
	var chars [2]byte
	C.AConfiguration_getLanguage(config.cptr(), (*C.char)(unsafe.Pointer(&chars[0])))
	return chars[:]
}

/**
 * Set the current language code in the configuration, from the first two
 * characters in the string.
 */
//void AConfiguration_setLanguage(AConfiguration* config, const char* language);
func (config *Configuration) SetLanguage(lang string) {
	var chars [2]byte
	chars[0] = lang[0]
	chars[1] = lang[1]
	C.AConfiguration_setLanguage(config.cptr(), (*C.char)(unsafe.Pointer(&chars[0])))
}

/**
 * Return the current country code set in the configuration.  The output will
 * be filled with an array of two characters.  They are not 0-terminated.  If
 * a country is not set, they will be 0.
 */
//void AConfiguration_getCountry(AConfiguration* config, char* outCountry);
func (config *Configuration) GetCountry() string {
	var chars [2]byte
	C.AConfiguration_getCountry(config.cptr(), (*C.char)(unsafe.Pointer(&chars[0])))
	if chars[0] == 0 {
		return ""
	}
	return string(chars[:])
}

/**
 * Set the current country code in the configuration, from the first two
 * characters in the string.
 */
//void AConfiguration_setCountry(AConfiguration* config, const char* country);
func (config *Configuration) SetCountry(lang []byte) {
	var chars [2]byte
	chars[0] = lang[0]
	chars[1] = lang[1]
	C.AConfiguration_setCountry(config.cptr(), (*C.char)(unsafe.Pointer(&chars[0])))
}

/**
 * Return the current ACONFIGURATION_ORIENTATION_* set in the configuration.
 */
//int32_t AConfiguration_getOrientation(AConfiguration* config);
func (config *Configuration) GetOrientation() int {
	return int(C.AConfiguration_getOrientation(config.cptr()))
}

/**
 * Set the current orientation in the configuration.
 */
//void AConfiguration_setOrientation(AConfiguration* config, int32_t orientation);
func (config *Configuration) SetOrientation(orientation int) {
	C.AConfiguration_setOrientation(config.cptr(), C.int32_t(orientation))
}

/**
 * Return the current ACONFIGURATION_TOUCHSCREEN_* set in the configuration.
 */
//int32_t AConfiguration_getTouchscreen(AConfiguration* config);
func (config *Configuration) GetTouchscreen() int {
	return int(C.AConfiguration_getTouchscreen(config.cptr()))
}

/**
 * Set the current touchscreen in the configuration.
 */
//void AConfiguration_setTouchscreen(AConfiguration* config, int32_t touchscreen);
func (config *Configuration) SetTouchscreen(touchscreen int) {
	C.AConfiguration_setTouchscreen(config.cptr(), C.int32_t(touchscreen))
}

/**
 * Return the current ACONFIGURATION_DENSITY_* set in the configuration.
 */
//int32_t AConfiguration_getDensity(AConfiguration* config);
func (config *Configuration) GetDensity() int {
	return int(C.AConfiguration_getDensity(config.cptr()))
}

/**
 * Set the current density in the configuration.
 */
//void AConfiguration_setDensity(AConfiguration* config, int32_t density);
func (config *Configuration) SetDensity(density int) {
	C.AConfiguration_setDensity(config.cptr(), C.int32_t(density))
}

/**
 * Return the current ACONFIGURATION_KEYBOARD_* set in the configuration.
 */
//int32_t AConfiguration_getKeyboard(AConfiguration* config);
func (config *Configuration) GetKeyboard() int {
	return int(C.AConfiguration_getKeyboard(config.cptr()))
}

/**
 * Set the current keyboard in the configuration.
 */
//void AConfiguration_setKeyboard(AConfiguration* config, int32_t keyboard);
func (config *Configuration) SetKeyboard(keyboard int) {
	C.AConfiguration_setKeyboard(config.cptr(), C.int32_t(keyboard))
}

/**
 * Return the current ACONFIGURATION_NAVIGATION_* set in the configuration.
 */
//int32_t AConfiguration_getNavigation(AConfiguration* config);
func (config *Configuration) GetNavigation() int {
	return int(C.AConfiguration_getNavigation(config.cptr()))
}

/**
 * Set the current navigation in the configuration.
 */
//void AConfiguration_setNavigation(AConfiguration* config, int32_t navigation);
func (config *Configuration) SetNavigation(navigation int) {
	C.AConfiguration_setNavigation(config.cptr(), C.int32_t(navigation))
}

/**
 * Return the current ACONFIGURATION_KEYSHIDDEN_* set in the configuration.
 */
//int32_t AConfiguration_getKeysHidden(AConfiguration* config);
func (config *Configuration) GetKeysHidden() int {
	return int(C.AConfiguration_getKeysHidden(config.cptr()))
}

/**
 * Set the current keys hidden in the configuration.
 */
//void AConfiguration_setKeysHidden(AConfiguration* config, int32_t keysHidden);
func (config *Configuration) SetKeysHidden(keysHidden int) {
	C.AConfiguration_setKeysHidden(config.cptr(), C.int32_t(keysHidden))
}

/**
 * Return the current ACONFIGURATION_NAVHIDDEN_* set in the configuration.
 */
//int32_t AConfiguration_getNavHidden(AConfiguration* config);
func (config *Configuration) GetNavHidden() int {
	return int(C.AConfiguration_getNavHidden(config.cptr()))
}

/**
 * Set the current nav hidden in the configuration.
 */
//void AConfiguration_setNavHidden(AConfiguration* config, int32_t navHidden);
func (config *Configuration) SetNavHidden(navHidden int) {
	C.AConfiguration_setNavHidden(config.cptr(), C.int32_t(navHidden))
}

/**
 * Return the current SDK (API) version set in the configuration.
 */
//int32_t AConfiguration_getSdkVersion(AConfiguration* config);
func (config *Configuration) GetSdkVersion() int {
	return int(C.AConfiguration_getSdkVersion(config.cptr()))
}

/**
 * Set the current SDK version in the configuration.
 */
//void AConfiguration_setSdkVersion(AConfiguration* config, int32_t sdkVersion);
func (config *Configuration) SetSdkVersion(sdkVersion int) {
	C.AConfiguration_setSdkVersion(config.cptr(), C.int32_t(sdkVersion))
}

/**
 * Return the current ACONFIGURATION_SCREENSIZE_* set in the configuration.
 */
//int32_t AConfiguration_getScreenSize(AConfiguration* config);
func (config *Configuration) GetScreenSize() int {
	return int(C.AConfiguration_getScreenSize(config.cptr()))
}

/**
 * Set the current screen size in the configuration.
 */
//void AConfiguration_setScreenSize(AConfiguration* config, int32_t screenSize);
func (config *Configuration) SetScreenSize(screenSize int) {
	C.AConfiguration_setScreenSize(config.cptr(), C.int32_t(screenSize))
}

/**
 * Return the current ACONFIGURATION_SCREENLONG_* set in the configuration.
 */
//int32_t AConfiguration_getScreenLong(AConfiguration* config);
func (config *Configuration) GetScreenLong() int {
	return int(C.AConfiguration_getScreenLong(config.cptr()))
}

/**
 * Set the current screen long in the configuration.
 */
//void AConfiguration_setScreenLong(AConfiguration* config, int32_t screenLong);
func (config *Configuration) SetScreenLong(screenLong int) {
	C.AConfiguration_setScreenLong(config.cptr(), C.int32_t(screenLong))
}

/**
 * Return the current ACONFIGURATION_UI_MODE_TYPE_* set in the configuration.
 */
//int32_t AConfiguration_getUiModeType(AConfiguration* config);
func (config *Configuration) GetUiModeType() int {
	return int(C.AConfiguration_getUiModeType(config.cptr()))
}

/**
 * Set the current UI mode type in the configuration.
 */
//void AConfiguration_setUiModeType(AConfiguration* config, int32_t uiModeType);
func (config *Configuration) SetUiModeType(uiModeType int) {
	C.AConfiguration_setUiModeType(config.cptr(), C.int32_t(uiModeType))
}

/**
 * Return the current ACONFIGURATION_UI_MODE_NIGHT_* set in the configuration.
 */
//int32_t AConfiguration_getUiModeNight(AConfiguration* config);
func (config *Configuration) GetUiModeNight() int {
	return int(C.AConfiguration_getUiModeNight(config.cptr()))
}

/**
 * Set the current UI mode night in the configuration.
 */
//void AConfiguration_setUiModeNight(AConfiguration* config, int32_t uiModeNight);
func (config *Configuration) SetUiModeNight(uiModeNight int) {
	C.AConfiguration_setUiModeNight(config.cptr(), C.int32_t(uiModeNight))
}

/**
 * Return the current configuration screen width in dp units, or
 * ACONFIGURATION_SCREEN_WIDTH_DP_ANY if not set.
 */
//int32_t AConfiguration_getScreenWidthDp(AConfiguration* config);
func (config *Configuration) GetScreenWidthDp() int {
	return int(C.AConfiguration_getScreenWidthDp(config.cptr()))
}

/**
 * Set the configuration's current screen width in dp units.
 */
//void AConfiguration_setScreenWidthDp(AConfiguration* config, int32_t value);
func (config *Configuration) SetScreenWidthDp(value int) {
	C.AConfiguration_setScreenWidthDp(config.cptr(), C.int32_t(value))
}

/**
 * Return the current configuration screen height in dp units, or
 * ACONFIGURATION_SCREEN_HEIGHT_DP_ANY if not set.
 */
//int32_t AConfiguration_getScreenHeightDp(AConfiguration* config);
func (config *Configuration) GetScreenHeightDp() int {
	return int(C.AConfiguration_getScreenHeightDp(config.cptr()))
}

/**
 * Set the configuration's current screen width in dp units.
 */
//void AConfiguration_setScreenHeightDp(AConfiguration* config, int32_t value);
func (config *Configuration) SetScreenHeightDp(value int) {
	C.AConfiguration_setScreenHeightDp(config.cptr(), C.int32_t(value))
}

/**
 * Return the configuration's smallest screen width in dp units, or
 * ACONFIGURATION_SMALLEST_SCREEN_WIDTH_DP_ANY if not set.
 */
//int32_t AConfiguration_getSmallestScreenWidthDp(AConfiguration* config);
func (config *Configuration) GetSmallestScreenWidthDp() int {
	return int(C.AConfiguration_getSmallestScreenWidthDp(config.cptr()))
}

/**
 * Set the configuration's smallest screen width in dp units.
 */
//void AConfiguration_setSmallestScreenWidthDp(AConfiguration* config, int32_t value);
func (config *Configuration) SetSmallestScreenWidthDp(value int) {
	C.AConfiguration_setSmallestScreenWidthDp(config.cptr(), C.int32_t(value))
}

/**
 * Return the configuration's layout direction, or
 * ACONFIGURATION_LAYOUTDIR_ANY if not set.
 */
//int32_t AConfiguration_getLayoutDirection(AConfiguration* config);
func (config *Configuration) GetLayoutDirection() int {
	return int(C.AConfiguration_getLayoutDirection(config.cptr()))
}

/**
 * Set the configuration's layout direction.
 */
//void AConfiguration_setLayoutDirection(AConfiguration* config, int32_t value);
func (config *Configuration) SetLayoutDirection(value int) {
	C.AConfiguration_setLayoutDirection(config.cptr(), C.int32_t(value))
}

/**
 * Perform a diff between two configurations.  Returns a bit mask of
 * ACONFIGURATION_* constants, each bit set meaning that configuration element
 * is different between them.
 */
//int32_t AConfiguration_diff(AConfiguration* config1, AConfiguration* config2);
func (config *Configuration) Diff(config2 *Configuration) int {
	return int(C.AConfiguration_diff(config.cptr(), config2.cptr()))
}

/**
 * Determine whether 'base' is a valid configuration for use within the
 * environment 'requested'.  Returns 0 if there are any values in 'base'
 * that conflict with 'requested'.  Returns 1 if it does not conflict.
 */
//int32_t AConfiguration_match(AConfiguration* base, AConfiguration* requested);
func (config *Configuration) Match(requested *Configuration) bool {
	return 0 != int(C.AConfiguration_diff(config.cptr(), requested.cptr()))
}

/**
 * Determine whether the configuration in 'test' is better than the existing
 * configuration in 'base'.  If 'requested' is non-NULL, this decision is based
 * on the overall configuration given there.  If it is NULL, this decision is
 * simply based on which configuration is more specific.  Returns non-0 if
 * 'test' is better than 'base'.
 *
 * This assumes you have already filtered the configurations with
 * AConfiguration_match().
 */
//int32_t AConfiguration_isBetterThan(AConfiguration* base, AConfiguration* test,
//        AConfiguration* requested);
func (config *Configuration) IsBetterThan(test, requested *Configuration) bool {
	return 0 != int(C.AConfiguration_isBetterThan(config.cptr(), test.cptr(), requested.cptr()))
}
