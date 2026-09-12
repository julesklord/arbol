with open("cmd/arbol/sysinfo.go", "r") as f:
    text = f.read()

count = text.count("func getMemory() string {")
print(f"func getMemory() string {{ count: {count}")
