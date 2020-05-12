Codegen
=======

Codegen is collection of code generation tools for golang

List
----

- [contextgen](contextgen)

  Command line to modify existing code to import `context.Context`
  package and append it as first argument to every functions in the file.

- [opentracing](opentracing)

  Command line to modify existing functions with first argument `context.Context`
  to add `span := opentracing.StartSpanFromContext` and `defer span.Close()`

License
-------
MIT
