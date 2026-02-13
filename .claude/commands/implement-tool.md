---
description: Implement a new DualSPHysics tool following the OpenCode Tool interface
argument-hint: "<tool-name>"
---

Implement a new DualSPHysics tool named `$ARGUMENTS` following the OpenCode BaseTool interface.

## TDD Approach (Red → Green)

### Step 1: Write Test First (RED)

Create `internal/llm/tools/$ARGUMENTS_test.go`:

```go
package tools

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func Test<Name>Tool_Info(t *testing.T) {
    tool := &<name>Tool{}
    info := tool.Info()
    assert.Equal(t, "$ARGUMENTS", info.Name)
    assert.NotEmpty(t, info.Description)
    assert.NotNil(t, info.Parameters)
}

func Test<Name>Tool_Run_Success(t *testing.T) {
    // Test with valid input
}

func Test<Name>Tool_Run_InvalidInput(t *testing.T) {
    // Test error handling
}
```

### Step 2: Implement Tool (GREEN)

Create `internal/llm/tools/$ARGUMENTS.go`:

```go
package tools

import (
    "context"
    "encoding/json"
    "fmt"
    "os/exec"
)

type <name>Tool struct{}

func (t *<name>Tool) Info() ToolInfo {
    return ToolInfo{
        Name:        "$ARGUMENTS",
        Description: "<tool description>",
        Parameters: map[string]any{
            "type": "object",
            "properties": map[string]any{
                // Define input schema
            },
            "required": []string{},
        },
    }
}

func (t *<name>Tool) Run(ctx context.Context, call ToolCall) (ToolResponse, error) {
    var input struct {
        // Parse input fields
    }
    if err := json.Unmarshal([]byte(call.Input), &input); err != nil {
        return ToolResponse{Content: err.Error(), IsError: true}, nil
    }

    // Execute via Docker
    cmd := exec.CommandContext(ctx, "docker", "compose", "run", "--rm", "dsph",
        // DualSPHysics command and arguments
    )
    output, err := cmd.CombinedOutput()
    if err != nil {
        return ToolResponse{Content: string(output), IsError: true}, nil
    }

    return ToolResponse{Content: string(output)}, nil
}
```

### Step 3: Register Tool

Add to `internal/llm/tools/tools.go` (or the tool registration function).

## DualSPHysics Integration Rules

1. **All execution via Docker**: `docker compose run --rm dsph <binary> <args>`
2. **GenCase paths**: Do NOT include `.xml` extension (auto-appended)
3. **XML format**: Attribute-only values (`<gravity x="0" y="0" z="-9.81" />`)
4. **Output paths**: Use `/data/` prefix (mapped to `./simulations/` via docker-compose)
5. **Error handling**: Parse DualSPHysics error output and translate to user-friendly Korean messages

## Available DualSPHysics Binaries

| Binary | Purpose | Example |
|--------|---------|---------|
| GenCase | XML → particles | `GenCase /cases/Tank_Def -save:/data/out` |
| DualSPHysics5.4_linux64 | GPU solver | `DualSPHysics5.4_linux64 /data/out/Tank_Def /data/out -gpu` |
| PartVTK | VTK export | `PartVTK -dirin /data/out -savevtk /data/vtk/part` |
| MeasureTool | Measurements | `MeasureTool -dirin /data/out -pointsfile /cases/probe_points.txt` |
| IsoSurface | Mesh generation | `IsoSurface -dirin /data/out -saveiso /data/iso/surface` |
