# Fix for 0.7

## Invocation

The `fix-for-0.7` fixes elvishscript code written before 0.7. It can be
invoked in one of two ways:

*   Without arguments, it reads stdin and writes stdout.

*   With filename arguments, it rewrites each given file.

It does not accept any flags.

## Changes in 0.7

Version 0.7 introduces several major backward-incompatible changes:

1.  Control structures used to mimic POSIX shell, but now use curly braces.

2.  There used to be a special "boolean return value" for functions like `==`
    or `eq` to indicate whether it has suceeded or failed. Now they simply
    output a boolean value `$true` or `$false` 
    ([#319](https://github.com/elves/elvish/issues/319)).

3.  The `if` and `while` control structures now take values instead of
    pipelines (also see [#319](https://github.com/elves/elvish/issues/319)).
    Together with 1 and 2, this code

    ```
    if == $x 1; then
        echo 1
    fi
    ```

    should now be written as

    ```
    if (== $x 1) {
        echo 1
    }
    ```

4.  The `for` control structure has been changed to operate on one iterable
    value, instead of a bunch of values. The `in` keyword is no longer needed.
    For example, this code

    ```
    for x in lorem ipsum; do
        echo $x
    done
    ```

    should now be written as

    ```
    for x [lorem ipsum] {
        echo $x
    }
    ```
