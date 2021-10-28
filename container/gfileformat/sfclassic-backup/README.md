# sfclassic

Use [siegfried](github.com/richardlehane/siegfried) for file format identification in go programs, without needing a
signature file.

This package writes the signature file into a go source file so you can just do "import
github.com/ross-spencer/sfclassic" and not worry about bundling signature files with your program.

You can configure the signature file used by your program by:

- replacing classic/classic.sig with one of your choice.
- Then just run `go generate`.

Future work: this package already wraps the siegfried Identify method to make it a bit more ergonomic. Further
simplifications could be made to the API in this package (e.g. to make it easier to get a MIME type).