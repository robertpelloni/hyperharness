import sys
import re

with open("tools/goose_opencode_kimi_parity.go", "r") as f:
    content = f.read()

import_pattern = """import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)"""

import_replacement = """import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"github.com/robertpelloni/hyperharness/internal/extensions"
)"""

content = content.replace(import_pattern, import_replacement)

with open("tools/goose_opencode_kimi_parity.go", "w") as f:
    f.write(content)
print("Added extensions import to tools/goose_opencode_kimi_parity.go")
