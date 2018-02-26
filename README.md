# go-versions

[![](https://godoc.org/github.com/apparentlymart/go-versions/versions?status.svg)](https://godoc.org/github.com/apparentlymart/go-versions/versions)

`versions` is a library for wrangling versions, lists of versions, and sets
of versions. Its idea of "version" is that defined by
[semantic versioning](https://semver.org/).

There are _many_ Go libraries out there for dealing with versions in general
and semantic versioning in particular, but many of them don't meet all of
the following requirements that this library seeks to meet:

* Version string and constraint string parsing _with good, user-oriented error
  messages in case of syntax problems_.
* Built-in mechanisms for filtering and sorting lists of candidate versions
  based on constraints.
* Ergonomic API for the calling application.

Whether _this_ library meets those requirements is of course subjective, but
these are certainly its goals.

To whet your appetite, here's an example program that solves the common
problem of taking a list of available versions and a version constraint and
then returning the newest available version that meets the constraint.

```go
package main

import "fmt"
import "os"
import "github.com/apparentlymart/go-versions/versions"

func main() {
	// In a real program, the version list would probably come
	// from some registry API, but we'll hard-code it for
	// example here.
	available := versions.List{
		versions.MustParseVersion("0.8.0"),
		versions.MustParseVersion("1.0.1"),
		versions.MustParseVersion("0.9.1"),
		versions.MustParseVersion("2.0.0-beta.1"),
		versions.MustParseVersion("2.1.0"),
		versions.MustParseVersion("1.0.0"),
		versions.MustParseVersion("0.9.0"),
		versions.MustParseVersion("1.1.0"),
		versions.MustParseVersion("2.0.0"),
	}

	allowed, err := versions.MeetingConstraintsString("^1.0.0")
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid version constraint: %s", err)
		os.Exit(1)
	}

	candidates := available.Filter(allowed)
	chosen := candidates.Newest()
	fmt.Printf("Would install v%s\n", chosen)
	// => Would install v1.1.0
}
```

## Version Sets and Version Lists

Many version libraries stop at just parsing and representing exact versions,
but most applications that need to process versions need also to represent
version constraints, ordered lists of versions, etc.

This library has a simple representation of versions as defined in the
semver spec, but its main focus is on _version sets_ and _version lists_,
which is reflected in the plural package name `versions`.

A version set is primarily used to represent _permitted_ versions, and
version sets are usually created from user-supplied constraint strings
that specify concisely which versions are members of the set:

```go
allowed, err := versions.MeetingConstraintsString("^1.0.0")
// (handle error)
fmt.Println(allowed.Has(MustParseVersion("1.0.0"))) // => true
fmt.Println(allowed.Has(MustParseVersion("0.0.1"))) // => false
fmt.Println(allowed.Has(MustParseVersion("2.0.0"))) // => false
```

Version sets can also be created and composed using the `versions`
package API, with the following predefined sets and set functions:

| Expression | Returns |
| ---------- | ------- |
| `versions.All` | Set of all possible versions. |
| `versions.None` | Set containing no versions at all. |
| `versions.Released` | Set of all "released" versions (not betas, alphas, etc). |
| `versions.Prerelease` | The opposite of `versions.Released`. |
| `versions.InitialDevelopment` | Contains all versions less than `1.0.0`, defined by semver as initial development releases where semver promises do not necessarily apply. |
| `versions.AtLeast(v)` | Set of versions greater than or equal to `v`. |
| `versions.AtMost(v)` | Set of versions less than or equal to `v`. |
| `versions.NewerThan(v)` | Set of versions greater than `v`. |
| `versions.OlderThan(v)` | Set of versions less than `v`. |
| `versions.Only(v)` | Set containing only the given version `v`. |
| `versions.Selection(vs...)` | Set containing only the given versions `vs`. |
| `versions.Intersection(sets...)` | Set containing the versions that all of the given sets have in common. |
| `versions.Union(sets...)` | Set containing all of the versions from all of the given sets. |
| `set1.Subtract(set2)` | Set containing the versions from `set1` that are not in `set2`. |

```go
v1 := versions.MustParseVersion("1.0.0")
fmt.Println(versions.All.Has(v1))                // => true
fmt.Println(versions.Releasaed.Has(v1))          // => true
fmt.Println(versions.None.Has(v1))               // => false
fmt.Println(versions.AtLeast(v1).Has(v1))        // => true
fmt.Println(versions.NewerThan(v1).Has(v1))      // => false
fmt.Println(versions.InitialDevelopment.Has(v1)) // => false
```

Whereas version sets contain an unordered collection of possibly-infinite
versions, version _lists_ are finite and ordered. A version list is in fact
just a named type around `[]Version` which adds some additional helper
methods for common operations with versions:

| Statement | Effect |
| --------- | ------ |
| `list = list.Filter(set)` | Removes from the list any members not in the given set, modifying the backing array in-place, and returns the new slice. |
| `list.Sort()` | Sorts in-place the list in increasing order by version, so the newest versions are at the end of the list. |
| `v = list.Newest()` | Returns the newest version in the list. (Strictly, _one of_ the newest versions, if the same version appears multiple times with different build metadata) |
| `v = list.NewestList()` | Returns a `List` of all of the versions that are newest in the list. May return more than one if there are multiple versions differing only in build metadata. |
| `v = list.NewestInSet(set)` | Like `Newest`, but considers only versions that are in the given set, without modifying the list. |

## Version values

The representation of versions themselves is, by comparison, very simple.
As defined by the semver spec, versions have major, minor, and patch segments
that are numeric, and also have more free-form strings representing
prerelease versions and build metadata.

Versions are usually passed as values and so can be compared for exact
equality using the standard `==` operator. However, most operations are
instead concerned primarily with the notion of _priority_ defined by the
semantic version spec, which is implemented in the following methods:

| Expression | Returns |
| ---------- | ------- |
| `v1.Same(v2)` | True if `v1` and `v2` are identical aside from their "metadata" |
| `v1.LessThan(v2)` | True if `v1` has a lower semver priority than `v2`. |
| `v1.GreaterThan(v2)` | True if `v1` has a higher semver priority than `v2`. |

These comparison functions are the basis of the `List.Sort` method. Note that
it is possible for two non-equal versions to be neither less than nor
greater than each other if they have differing `Metadata`.

When considering set membership, the entire version value is considered
including metadata. Applications that do not have any need for metadata
may choose to strip it out using `v.Comparable()`, which returns a new
version that is identical to the receiver except that its metadata is
empty.

### Unspecified Versions

The special version value `versions.Unspecified` is the zero value of
`Version` and represents the absense of a version. Its representation is
the same as for the version string `0.0.0`, and so that string is not a
valid version number according to this package.

The only version set that contains `versions.Unspecified` is `versions.All`.
This is true even of the set returned by `versions.Only(versions.Unspecified)`,
which is a useless expression.

## Text and JSON serialization

The `versions.Version` type implements `encoding.TextMarshaler` and
`encoding.TextUnmarshaler`, using the same syntax expected by
`versions.ParseVersion` and `versions.MustParseVersion`. This allows
version values to be included in structs used with encoding packages that
make use of these interfaces, including `encoding/json`:

```go
type Package struct {
    Name    string           `json:"name"`
    Version versions.Version `json:"version"`
}
```

The `versions.Set` type also supports `encoding.TextUnmarshaler`, so
it can be used for _unmarshalling_ of constrants into sets via the
canonical constraint syntax. Sets cannot be marshalled because the
set model implemented by this package contains features that cannot be
expressed in the constraint language.

```go
type Requirement struct {
    PackageName string       `json:"packageName"`
    Versions    versions.Set `json:"versions"`
}
```

In practice the asymmetry of version set marshalling is not usually a problem
because constraint sets are more often written by humans than by machines. In
future the `constraints` package may get support for serializing its own
constraint model, should a compelling use-case emerge. If you have one, please
open a GitHub issue to discuss it!

## Finite vs. Infinite Version Sets

Most version sets contain an infinite number of versions that lay within
some bounds, such as the set returned by `versions.AtLeast(...)`.
Some version sets contain only a finite number of versions, though. For
example, `versions.Only(...)` returns a set containing only one version.

The `set.IsFinite()` method allows a calling application to recognize if
a particular set is finite. Some set operations in this package guarantee a
finite set when certain conditions are met, avoiding the need to check this
method; see
[the package godoc](https://godoc.org/github.com/apparentlymart/go-versions/versions)
for full details.

A finite set can be converted into a list using `set.List()`. (This method
will panic if used on an infinite set.)

## Constraint String Parsers

Earlier examples showed the function `versions.MeetingConstraintsString`, which
is a straightforward way to take a version string provided by the user and
obtain a version set containing all of the versions it selects.

The constraint syntax is implemented by the sub-package
[`constraints`](https://godoc.org/github.com/apparentlymart/go-versions/versions/constraints),
which contains a model for representing constraint specifications and some
parser functions. Applications with more specific needs may wish to call
directly into the functions in this package, for example to parse constraints
using a `rubygems`-like syntax rather than the "npm-like" syntax this package
uses by default.

Full details of this package's _canonical_ constraint syntax (the "npm-like"
one) are in the documentation for
[`constraints.Parse`](https://godoc.org/github.com/apparentlymart/go-versions/versions/constraints#Parse).

The "ruby-like" parsers use the same basic structure but use alternative
operators inspired by the `rubygems` constraint syntax, including the
"pessimistic" operator `~>`.

Neither constraint syntax is 100% compatible with the system it takes
inspiration from, but the goal is to be familiar enough to allow for a good
user experience for users that have worked in these other systems.

The constraint string parsers are designed to produce helpful error messages
that are suitable to return directly an English-speaking end-user that has
authored an invalid constraint string. For example:

| Invalid String | Error Message |
| -------------- | ------------- |
| `1.0.0.0` | too many numbered portions; only three are allowed (major, minor, patch) |
| `=>1.1.1` | invalid constraint operator `=>`; did you mean `>=`? |
| `1.0.0, 2.0.0` | commas are not needed to separate version selections; separate with spaces instead |

## Requested Versions

In addition to the usual idea of a set either containing or not containing
a version, a version set has an additional concept of a version being
_requested_. The requested versions of a set form a subset of that set,
inferred from any exact version selections (`versions.Only(...)` and
`versions.Selection(...)`) made in the construction of that set.

In most cases this distinction is unimportant, but it is particularly
interesting when dealing with pre-release versions, since these should
generally be considered only if explicitly requested.

Constraints processed using `versions.MeetingConstraints(...)` and
`versions.MeetingConstraintsString(...)` will automatically exclude all
unreleased versions that are not explicitly requested:

```go
beta1 := versions.MustParseVersion("2.0-beta.1")
allowed := versions.MustMakeSet(versions.MeetingConstraintsString(">=1.0"))
fmt.Println(allowed.Has(beta1)) // false
```

Version sets constructed manually using the constructor functions do not have
this characteristic, and will return pre-release versions unless they are
specifically excluded from the set:

```go
beta1 := versions.MustParseVersion("2.0-beta.1")
beta2 := versions.MustParseVersion("2.0-beta.2")
min := versions.MustParseVersion("1.0.0")

allowed := versions.AtLeast(min)
fmt.Println(allowed.Has(beta1)) // => true
fmt.Println(allowed.Has(beta2)) // => true

// Construct a new version set containing only _released_ versions that meet
// our constraint.
onlyReleased = allowed.Intersection(versions.Released)
fmt.Println(onlyReleased.Has(beta1)) // => false
fmt.Println(onlyReleased.Has(beta2)) // => false
```

The _requested versions set_ of a set can be used to obtain any versions
that are requested exactly by a set, in order to implement the pre-release
version selection behavior done automatically by `versions.MeetingConstraints`:

```go
beta1 := versions.MustParseVersion("2.0-beta.1")
beta2 := versions.MustParseVersion("2.0-beta.2")
min := versions.MustParseVersion("1.0.0")

allowed := versions.Union(
    versions.AtLeast(min), // allow any version >1.0.0
    versions.Only(beta1),  // also allow beta1
)
fmt.Println(allowed.Has(beta1)) // => true
fmt.Println(allowed.Has(beta2)) // => true

// Exclude pre-release versions
onlyReleased = allowed.Intersection(versions.Released)
fmt.Println(onlyReleased.Has(beta1)) // => false
fmt.Println(onlyReleased.Has(beta2)) // => false

// Now re-allow the explicitly-requested version, beta1
allowed = Union(
    allowed.AllRequested(), // set containing only beta1
    onlyReleased,           // set containing released versions >=1.0.0
)
fmt.Println(allowed.Has(beta1)) // => true, because it was requested
fmt.Println(allowed.Has(beta2)) // => false, because it was not requested
```

Because excluding pre-releases unless explicitly requested is usually
desirable, a helper method is provided to automatically implement the above
for any arbitrary set:

```go
beta1 := versions.MustParseVersion("2.0-beta.1")
beta2 := versions.MustParseVersion("2.0-beta.2")
min := versions.MustParseVersion("1.0.0")

allowed := versions.Union(
    versions.AtLeast(min), // allow any version >1.0.0
    versions.Only(beta1),  // also allow beta1
)
fmt.Println(allowed.Has(beta1)) // => true
fmt.Println(allowed.Has(beta2)) // => true

allowed = allowed.WithoutUnrequestedPrereleases()
fmt.Println(allowed.Has(beta1)) // => true, because it was requested
fmt.Println(allowed.Has(beta2)) // => false, because it was not requested
```

The set of requested versions for a set is always a finite set, by definition.
It can therefore be converted to a version list with `set.List()` if required.

Requested versions are subject to the same set operations as normal set
members, due to the rule that all requested versions must also be set members.
For example,
`versions.Only(versions.MustParseVersion("1.0-beta.1")).Subtract(versions.Prerelease)` does not request `1.0-beta.1`, because that
member was removed by the `Subtract` operation.

Most callers will just pass in constraint strings authored by the user and
thus not need to worry about requested version sets. However, the functionality
is available to directly interact with this concept for the benefit of
applications that wish to implement different rules for pre-release versions.

