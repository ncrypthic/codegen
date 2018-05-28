contextgen
==========

Import `context.Context` package (if not previously imported) and adding `ctx context.Context` object as first argument in existing
functions (if it doesn't have `context.Context` as their first arguments).

e.g
```go
package somepackage

import "fmt"

type SomeFn func(string) error
```

will translate into:
```go
package somepackage

import (
        "context"
        "fmt"
)

type SomeFn func(context.Context, string) error
```

This command line will help refactoring existing application writen in go to implement [opentracing.io](opentracing.io)
