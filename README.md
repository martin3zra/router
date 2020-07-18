# Router

Wrapper in top of [Gorilla Tool Kit](https://www.gorillatoolkit.org/) inspired in [Laravel Framework](https://laravel.com/docs/7.x/routing#basic-routing).

* Easy to write and read
* Route prefix
* Route group
* Middleware group

## Install

You might as well just copy the [router.go](https://github.com/martin3zra/router/blob/master/router.go) file into your own project (and the [router_test.go](https://github.com/martin3zra/router/blob/master/router_test.go) while you're at it for future generations) rather than adding a dependency.

But it is maintained as a Go module which you can get with:

```bash
go get github.com/martin3zra/router
```


## Usage
Import it:

```shell script
import (
	"github.com/martin3zra/router"
)
```
The package will show you a clean way and guide you through the necessary steps to set up your routes.

![router](screenshot.png?raw=true)