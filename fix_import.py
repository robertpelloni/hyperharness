import sys

with open("tools/claude_code_parity.go", "r") as f:
    content = f.read()

import_pattern = """import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"context"
	"github.com/robertpelloni/hyperharness/foundation/pi"
	"github.com/robertpelloni/hyperharness/internal/subagents"
	"github.com/robertpelloni/hyperharness/internal/skills"
	"github.com/robertpelloni/hyperharness/config"
)"""

old_import = """import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"context"
	"github.com/robertpelloni/hyperharness/config"
	"github.com/robertpelloni/hyperharness/foundation/pi"
	"github.com/robertpelloni/hyperharness/internal/skills"
	"github.com/robertpelloni/hyperharness/internal/subagents"
)"""

if "github.com/robertpelloni/hyperharness/config" not in content:
    content = content.replace('"context"\n\t"github.com/robertpelloni/hyperharness/foundation/pi"', '"context"\n\t"github.com/robertpelloni/hyperharness/config"\n\t"github.com/robertpelloni/hyperharness/foundation/pi"')

with open("tools/claude_code_parity.go", "w") as f:
    f.write(content)
