gitversion
==========

Include git commit hash in golang file

usage
=====
gitversion -i pathToRepository -o pathToVersionFile -p packageName

-i  defaults to .<br>
-o  defaults to version.go<br>
-p  defaults to version<br>
-tf adds human readable timestamp format (see [time](http://golang.org/pkg/time/#pkg-constants) package)<br>
-v  adds version string (e.g., "1.0-Final")<br>
-s  uses short git commit hash

template
========
```go
package %s

const (
	GIT_COMMIT_HASH = "%s"

	// Unix time (seconds since January 1, 1970 UTC)
	GENERATED = %d

	// human readable timestamp
	GENERATED_FMT = "%s"  // only if "-tf" option specified

	VERSION = "%s"        // only if "-v" option specified
)
```
