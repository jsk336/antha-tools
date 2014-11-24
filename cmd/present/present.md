# present
--
Present displays slide presentations and articles. It runs a web server that
presents slide and article files from the current directory.

It may be run as a stand-alone command or an App Engine app. Instructions for
deployment to App Engine are in the README of the antha-tools repository.

Usage of present:

    -base="": base path for slide template and static resources
    -http="127.0.0.1:3999": host:port to listen on

Input files are named foo.extension, where "extension" defines the format of the
generated output. The supported formats are:

    .slide        // HTML5 slide presentation
    .article      // article format, such as a blog post

The present file format is documented by the present package:
http://godoc.org/antha-tools/present
