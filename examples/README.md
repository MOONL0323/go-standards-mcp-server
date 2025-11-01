# Examples

This directory contains example Go code files for testing the MCP server.

## Files

### sample.go

A simple Go program with intentional code quality issues for demonstration:
- Unused variables
- Unchecked errors
- Unused functions
- Style issues

Use this file to test the analyzer:

```bash
# Using the CLI (when implemented)
./bin/mcp-server analyze --file ./examples/sample.go --standard standard

# Using MCP tool via Cursor or other MCP client
{
  "file_path": "./examples/sample.go",
  "standard": "standard",
  "format": "markdown"
}
```

## Expected Results

When analyzing `sample.go` with standard mode, you should see issues related to:
- Code style (using println instead of fmt.Println)
- Unused variables (x)
- Dead code (unused function)
- Potential improvements

## Adding Your Own Examples

Feel free to add more example files to test different scenarios:
- Complex business logic
- Web servers
- Database operations
- Concurrent code
- Error handling patterns

Each example should demonstrate specific code quality issues or patterns.
