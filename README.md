# Stitch Language
This repository houses a proof-of-concept language for building Flow-based
programs with a minor emphasis on data collection.
The language will potentially make its way into the [Gaufre](https://github.com/nirosys/gaufre)
runtime which is used as a common graph execution runtime.

## Goals for the stitch language
The goals for the language itself are fairly simple:

  * Provide a syntax that is fairly intuitive to read, but not too verbose.
  * Provide constructs that allow code re-use, such as functions, macros, etc.
  * Provide constructs to allow automated schema discovery.

Some extra things that would be nice to have:

  * Provide in-code documentation capabilities (eg. documenting the above mentioned schema discovery)

## Building / Experimentation
The current state of stitch centers around a REPL that is implemented in this repository.
Building the REPL is simple and requires only the use of `make`:

```
$ make
Building: v0.1.3-04fee9f-dirty
$
```

This will build a versioned binary in the root of the repo named `stitch`.
You can run the binary with no arguments in order to explore the REPL:

```
./stitch
Stitch
v0.1.3-04fee9f-dirty

Use '.quit' to quit, or '.help' for help
>>
```

Once in the REPL, any valid [stitch syntax](doc/language-notes.md) can be entered.

## Goals for stitch CLI
In order to help build working stitch definitions
I'd like the CLI to offer the following feature.

* Compile stitch into gaufre graph
* Produce `dot` of resulting graph
* Output the schema produced by the stitch
* Error on unknown variables, with suggestions for misspellings.
* Detect cycles.
* Detect unused subgraphs.


## REPL Goals

* Provide an interactive method for building and testing graphs.
  * Run graphs for testing.
  * Allow Watches
  * Allow Breakpoints
  * Output generated data
