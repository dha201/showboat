# Shell Output and Pop

*2026-02-08T14:03:47Z*

This demo shows two new showboat features: `build run` now prints shell output and reflects the exit code, and the new `pop` command removes the most recent entry from a document.

## Build run output

Previously `build run` was silent — the output was captured into the document but not shown. Now it prints the shell output to stdout and exits with the same exit code as the shell command.

Let's create a scratch document and run some commands against it:

```bash
rm -f /tmp/demo.md
/tmp/showboat init /tmp/demo.md "My Demo"
/tmp/showboat build /tmp/demo.md run bash "echo Hello from the inner document"

```

```output
Hello from the inner document
```

## Exit code reflection

When a command fails, `build run` exits with the same exit code as the shell. The output is still captured in the document:

```bash
/tmp/showboat build /tmp/demo.md run bash "echo about to fail && exit 1"
echo "showboat exit code: $?"

```

```output
about to fail
showboat exit code: 1
```

## The pop command

The agent sees the non-zero exit code and can decide to remove the failed entry using `pop`. Here is the document before and after popping:

```bash
echo "=== Before pop: $(grep -c '```' /tmp/demo.md) fenced blocks ==="
grep '```' /tmp/demo.md
echo ""
/tmp/showboat pop /tmp/demo.md
echo "=== After pop: $(grep -c '```' /tmp/demo.md) fenced blocks ==="
grep '```' /tmp/demo.md

```

````output
=== Before pop: 8 fenced blocks ===
```bash
```
```output
```
```bash
```
```output
```

=== After pop: 4 fenced blocks ===
```bash
```
```output
```
````

The `pop` removed both the code block (`echo about to fail && exit 1`) and its output block in one operation. The document now only contains the successful entry.
