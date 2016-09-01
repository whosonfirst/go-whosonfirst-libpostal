# walk

This is a fork of [MichaelTJones](https://github.com/MichaelTJones)' original [walk](https://github.com/MichaelTJones/walk) package.

## The original version

_This is what [MichaelTJones](https://github.com/MichaelTJones) wrote in his original [README.md](https://github.com/MichaelTJones/walk/blob/master/README.md) file_:

Fast parallel version of golang filepath.Walk()

Performs traversals in parallel so set GOMAXPROCS appropriately. Vaues of 8 to 16 seem to work best on my 
4-CPU plus 4 SMT pseudo-CPU MacBookPro. The result is about 4x-6x the traversal rate of the standard Walk().
The two are not identical since we are walking the file system in a tumult of asynchronous walkFunc calls by
a number of goroutines. So, take note of the following:

1. This walk honors all of the walkFunc error semantics but as multiple user-supplied walkFuncs may simultaneously encounter a traversal error or generate one to stop traversal, only the FIRST of these will be returned as the Walk() result. 

2. Further, since there may be a few files in flight at the instant of  error discovery, a few more walkFunc calls may happen after the first error-generating call has signaled its desire to stop. In general this is a non-issue but it could matter so pay attention when designing your walkFunc. (For example, if you accumulate results then you need to have your own means to know to stop accumulating once you signal an error.)

3. Because the walkFunc is called concurrently in multiple goroutines, it needs to be careful about what it does with external data to avoid collisions. Results may be printed using fmt, but generally the best plan is to send results over a channel or accumulate counts using a locked mutex.

These issues are illustrated/handled in the simple traversal programs supplied with walk. There is also a test file that is just the tests from filepath in the Go language's standard library. Walk passes these tests when run in single process mode, and passes most of them in concurrent mode (GOMAXPROCS > 1). The problem is not a real problem, but one of the test expecting a specific number of errors to be found based on presumed sequential traversals.

## The changes

### Set number of walkers from runtime.GOMAXPROCS 

This package incorporates [avleen](https://github.com/avleen)'s [fork](https://github.com/avleen/walk) of the walk package to [set the number of walkers from runtime.GOMAXPROCS ](https://github.com/MichaelTJones/walk/compare/master...avleen:master).

### walk.WalkWithNFSKludge

This introduces a new package function called `WalkWithNFSKludge` that will trap and ignore `readdirent: errno 523` errors which can occur when traversing NFS mounts. You should use this function with caution and your eyes wide open.

There is _nothing_ magic happening here. It is a leap of faith that the error in question, which is raised by the operating system, is not really a big deal for the purposes of your application and shouldn't yield a fatal error by the `walk` package.

File under: 🙈
