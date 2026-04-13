import sys

with open("internal/mcp/mcp.go", "r") as f:
    content = f.read()

bad1 = """		if _, err := server.stdin.Write(append(reqData, '
')); err != nil {"""
good1 = """		if _, err := server.stdin.Write(append(reqData, '\\n')); err != nil {"""

bad2 = """		line, err := server.stdout.ReadString('
')"""
good2 = """		line, err := server.stdout.ReadString('\\n')"""

bad3 = """	server.stdin.Write(append(initData, '
'))"""
good3 = """	server.stdin.Write(append(initData, '\\n'))"""

bad4 = """	line, err := server.stdout.ReadString('
')"""
good4 = """	line, err := server.stdout.ReadString('\\n')"""

bad5 = """	line, err = server.stdout.ReadString('
')"""
good5 = """	line, err = server.stdout.ReadString('\\n')"""

content = content.replace(bad1, good1).replace(bad2, good2).replace(bad3, good3).replace(bad4, good4).replace(bad5, good5)

with open("internal/mcp/mcp.go", "w") as f:
    f.write(content)
