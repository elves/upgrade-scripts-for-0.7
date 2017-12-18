# Upgrade scripts for 0.7

## Invocation

This tool fixes Elvishscript code written before 0.7. It can be invoked in one of two ways:

*   Without arguments, it reads stdin and writes stdout.

*   With filename arguments, it rewrites each given file.

It does not accept any flags.

## Scope of this tool

This tool fixes the syntax part of all breaking changes when upgrading from 0.6 to 0.7. For a comprehensive list, see the [notes](https://github.com/elves/elvish/releases/tag/0.7) for the 0.7 pre-release.

Note that since control structures now introduce their own scopes, variables
defined within them are no longer accessible outside. For instance, this used
to work:

```sh
if true; then
  x = 1
else
  x = 2
fi
put $x $x
```

After automatic fixing, the above becomes the following, which does not work:

```elvish
if (true) {
  x = 1
} else {
  x = 2
}
put $x $x
```

To fix this, either declare `$x` beforehand:

```elvish
x = ''
if (true) {
  x = 1
} else {
  x = 2
}
put $x $x
```

Or rewrite the code like this:

```elvish
x = (
  if (true) {
    put 1
  } else {
    put 2
  })
put $x $x
```

### Assignment syntax

This tool does one more thing. Version 0.6 introduced a more readable syntax for assignment, allowing spaces around the equal sign: for instance, this program will fix `a=b` to `a = b`. The new syntax does not support doing multiple assignments at once, so `a=b c=d` will be split into two commands: `a = b; c = d`.

Note that the old syntax is not going away though: it is still the syntax for
temporary assignments, so `a=b ls` will be left untouched. Support for using
this syntax for non-temporary assignments will be dropped in a later version.
