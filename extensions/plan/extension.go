package plan

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"sync"

	"github.com/MichaelKinsy/PiG/extensions/sdk"
)

var readOnlyTools = []string{"read", "grep", "find", "ls"}

const planPrompt = `

# Plan mode
You are in plan mode. Only read-only tools are available: ` + "read, grep, find, ls" + `.
Do not modify files or run commands. Explore the code, then reply with a concrete,
numbered implementation plan: files to change, the changes, risks and open questions.
The user leaves plan mode with /plan before any change is made.`

func Extension() *sdk.Extension {
	e := sdk.New("plan")

	var (
		mu       sync.Mutex
		active   bool
		previous []string
	)

	enable := func(ctx sdk.Context) error {
		current, err := ctx.GetActiveTools()
		if err != nil {
			return fmt.Errorf("read active tools: %w", err)
		}
		previous = current
		ctx.SetActiveTools(readOnlyTools)
		active = true
		ctx.SetStatus("plan", "plan mode (read-only)")
		ctx.Notify("Plan mode on: read-only tools ("+strings.Join(readOnlyTools, ", ")+")", "info")
		return nil
	}

	disable := func(ctx sdk.Context) {
		restore := previous
		if len(restore) == 0 {
			restore = []string{"read", "bash", "edit", "write"}
		}
		ctx.SetActiveTools(restore)
		active = false
		previous = nil
		ctx.ClearStatus("plan")
		ctx.Notify("Plan mode off: tools restored", "info")
	}

	e.Command("plan", "Toggle read-only plan mode (/plan [on|off])", func(ctx sdk.Context, args string) error {
		mu.Lock()
		defer mu.Unlock()

		switch strings.TrimSpace(strings.ToLower(args)) {
		case "on":
			if !active {
				return enable(ctx)
			}
		case "off":
			if active {
				disable(ctx)
			}
		case "":
			if active {
				disable(ctx)
				return nil
			}
			return enable(ctx)
		default:
			ctx.Notify("usage: /plan [on|off]", "warn")
		}
		return nil
	})

	e.OnSessionStart(func(ctx sdk.Context, _ map[string]any) (any, error) {
		switch strings.ToLower(os.Getenv("PLAN_MODE_ENABLED")) {
		case "1", "true", "yes":
			mu.Lock()
			defer mu.Unlock()
			if !active {
				return nil, enable(ctx)
			}
		}
		return nil, nil
	})

	e.OnEvent("before_agent_start", func(ctx sdk.Context, data map[string]any) (any, error) {
		mu.Lock()
		defer mu.Unlock()
		if !active {
			return nil, nil
		}
		prompt, _ := data["systemPrompt"].(string)
		return map[string]any{"systemPrompt": prompt + planPrompt}, nil
	})

	e.OnEvent("tool_call", func(ctx sdk.Context, data map[string]any) (any, error) {
		mu.Lock()
		defer mu.Unlock()
		name, _ := data["toolName"].(string)
		if !active || slices.Contains(readOnlyTools, name) {
			return nil, nil
		}
		return map[string]any{"block": true, "reason": "Plan mode is on: " + name + " is not allowed. Use /plan to leave plan mode."}, nil
	})

	return e
}
