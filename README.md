# Fix for 0.7

## Invocation

The `fix-for-0.7` fixes elvishscript code written before 0.7. It can be invoked in one of two ways:

*   Without arguments, it reads stdin and writes stdout.

*   With filename arguments, it rewrites each given file.

It does not accept any flags.

## Scope of this tool

This tool fixes all breaking changes when upgrading from 0.6 to 0.7. For a comprehensive list, see the [notes](https://github.com/elves/elvish/releases/tag/0.7) for the 0.7 pre-release.

This tool does one more thing. Version 0.6 introduced a more readable syntax for assignment, allowing spaces around the equal sign: for instance, this program will fix `a=b` to `a = b`. Not all assignments are fixed:

*   The new syntax does not support doing multiple assignments at once, so
    `a=b c=d` will be left untouched.

*   The old syntax is still the only syntax for temporary assignments, so `a=b
    ls` will be left untouched.
