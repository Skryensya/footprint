package hooks

func Script(fpPath string, source string) string {
	// Run fp record with the source environment variable
	// Redirect stdout to /dev/null (suppress normal output)
	// Errors are now logged internally by fp record via the logger
	return "#!/bin/sh\n" +
		"FP_SOURCE=" + source + " " +
		fpPath + " record >/dev/null 2>&1 || true\n"
}
