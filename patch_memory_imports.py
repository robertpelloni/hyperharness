import sys

with open("tools/pi_exact_parity.go", "r") as f:
    content = f.read()

import_old = """	"context"
	"encoding/json"
	"path/filepath"
	"fmt"
	"strings"
	"sync"
"""
import_new = """	"context"
	"encoding/json"
	"path/filepath"
	"fmt"
	"strings"
	"sync"
	"os"
	"time"
"""

content = content.replace(import_old, import_new)

with open("tools/pi_exact_parity.go", "w") as f:
    f.write(content)
