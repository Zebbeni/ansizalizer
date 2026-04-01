package testutil

import "github.com/Zebbeni/ansizalizer/style"

// Setup initializes global state needed by tests (theme, styles).
// Call this from TestMain or at the start of tests that exercise View/Update.
func Setup() {
	style.SetTheme(style.LightOnDark)
}
