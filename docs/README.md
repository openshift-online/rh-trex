# TRex Documentation

Welcome to TRex - Red Hat's **T**rusted **R**est **EX**ample for rapid API development.

## Quick Navigation

### 🚀 I'm New Here
**[Getting Started](getting-started/)** - Choose your path and get up and running quickly

### 🏗️ I Want to Build Something
- **[Template Cloning](template-cloning/)** - Create a new microservice from TRex template
- **[Entity Development](entity-development/)** - Add business entities to existing projects

### 🔧 I Need to Deploy/Operate
**[Operations](operations/)** - Local development, deployment, and maintenance

### 📚 I Need Reference Information
**[Reference](reference/)** - API specs, configuration, commands, and technical details

### 🐛 Something's Not Working
**[Troubleshooting](troubleshooting/)** - Common problems and solutions

### 🛠️ I Want to Contribute to TRex
**[Framework Development](framework-development/)** - Architecture, contributing, and extending TRex

## Documentation Philosophy

This documentation is organized by **user goals** rather than technical categories. Each section has:

- **README.md** - Overview and navigation within that section
- **Guides** - Step-by-step instructions for specific tasks
- **References** - Detailed technical information
- **Examples** - Real-world usage patterns

## Quick Start Options

**Option 1: Clone TRex Template → New Microservice**
```bash
go run ./scripts/clone/main.go --name my-service --destination ~/projects/my-service
```

**Option 2: Generate Entity → Add to Existing Project**
```bash
go run ./scripts/generate/main.go --kind Product
```

**Option 3: Run TRex Locally → Explore and Learn**
```bash
make db/setup && make run
```

Choose the path that matches your goal and dive into the relevant section above.