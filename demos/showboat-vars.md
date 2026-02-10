# Persistent Variables with showboat var

*2026-02-10T05:09:08Z*

This demo shows how `showboat var` lets you pass state between cells. Code blocks executed via `showboat exec` can set and read persistent variables — useful for recording PIDs, ports, file paths, and other values that later cells need.

## Basic usage

Variables are set with `showboat var set KEY VALUE` and read with `showboat var get KEY`. The showboat binary is automatically available on PATH inside executed cells.

```bash
showboat var set GREETING "Hello from an earlier cell"
echo "Variable set."
```

```output
Variable set.
```

```bash
echo "The greeting is: $(showboat var get GREETING)"
```

```output
The greeting is: Hello from an earlier cell
```

## Starting a background HTTP server

A common use case: start a server in one cell, do work against it, then shut it down. Without `showboat var`, there is no way to pass the PID between cells since each runs in a fresh process.

```bash
mkdir -p /tmp/showboat-demo-site
echo "<h1>It works!</h1>" > /tmp/showboat-demo-site/index.html
python3 -m http.server 8642 --directory /tmp/showboat-demo-site > /dev/null 2>&1 &
showboat var set SERVER_PID $!
showboat var set SERVER_PORT 8642
sleep 0.3
echo "Server started on port 8642."
```

```output
Server started on port 8642.
```

Now a later cell can fetch from the server using the saved port, without needing to know the port number:

```bash
PORT=$(showboat var get SERVER_PORT)
curl -s http://localhost:$PORT/
```

```output
<h1>It works!</h1>
```

## Stopping the server

The PID was saved earlier. A later cell can use it to cleanly shut down the server:

```bash
PID=$(showboat var get SERVER_PID)
kill $PID 2>/dev/null && echo "Server stopped." || echo "Server already stopped."
```

```output
Server stopped.
```

## Listing and cleaning up variables

You can list all variables and delete ones you no longer need:

```bash
echo "All variables:"
showboat var list
echo ""
echo "Deleting SERVER_PID..."
showboat var del SERVER_PID
echo "Remaining variables:"
showboat var list
```

```output
All variables:
GREETING
SERVER_PID
SERVER_PORT

Deleting SERVER_PID...
Remaining variables:
GREETING
SERVER_PORT
```

## Cross-language support

Variables work across languages. Any language can call the `showboat` binary via subprocess:

```python3
import subprocess

# Read a variable set by a bash cell
result = subprocess.run(["showboat", "var", "get", "GREETING"], capture_output=True, text=True)
print(f"From Python: {result.stdout}")

# Set a new variable from Python
subprocess.run(["showboat", "var", "set", "PYTHON_VERSION", "3.x"])
print("Set PYTHON_VERSION from Python.")
```

```output
From Python: Hello from an earlier cell
Set PYTHON_VERSION from Python.
```

```bash
echo "Python set PYTHON_VERSION to: $(showboat var get PYTHON_VERSION)"
```

```output
Python set PYTHON_VERSION to: 3.x
```
