# Makefile Options

This project includes a \`Makefile\` to simplify common tasks such as building, running, testing, and managing the application. Below is a detailed explanation of the available \`Makefile\` targets:

---

### **1. \`make build\`**
- **Description:** Compiles the Go source code into a binary.
- **Output:** The compiled binary is placed in the \`build/\` directory with the name \`proc-monitor\`.
- **Usage:**
  \`\`\`bash
  make build
  \`\`\`

---

### **2. \`make run\`**
- **Description:** Runs the compiled binary with default flags (\`-interval=1 -html=true -port=8080\`).
- **Prerequisite:** Requires the binary to be built first using \`make build\`.
- **Usage:**
  \`\`\`bash
  make run
  \`\`\`

---

### **3. \`make test\`**
- **Description:** Executes all Go tests in the project.
- **Details:** Runs tests with verbose output and generates a coverage profile (\`coverage.out\`).
- **Usage:**
  \`\`\`bash
  make test
  \`\`\`

---

### **4. \`make coverage\`**
- **Description:** Generates an HTML coverage report from the \`coverage.out\` file.
- **Details:** Opens a browser-friendly \`coverage.html\` file to visualize test coverage.
- **Prerequisite:** Requires \`make test\` to be run first to generate the \`coverage.out\` file.
- **Usage:**
  \`\`\`bash
  make coverage
  \`\`\`

---

### **5. \`make clean\`**
- **Description:** Removes all build artifacts, including the \`build/\` directory, \`coverage.out\`, and \`coverage.html\`.
- **Usage:**
  \`\`\`bash
  make clean
  \`\`\`

---

### **6. \`make install\`**
- **Description:** Installs the compiled binary to \`/usr/local/bin\` so it can be run globally.
- **Prerequisite:** Requires the binary to be built first using \`make build\`.
- **Usage:**
  \`\`\`bash
  make install
  \`\`\`

---

### **7. \`make uninstall\`**
- **Description:** Removes the binary from \`/usr/local/bin\`.
- **Usage:**
  \`\`\`bash
  make uninstall
  \`\`\`

---

### **8. \`make all\`**
- **Description:** Default target that builds the binary.
- **Details:** Equivalent to running \`make build\`.
- **Usage:**
  \`\`\`bash
  make all
  \`\`\`

---

### **Notes**
- Ensure you have Go installed and properly configured on your system before using these commands.
- Modify the \`Makefile\` variables (e.g., \`BIN_NAME\`, \`TEST_FLAGS\`) if needed to suit your project's requirements.
