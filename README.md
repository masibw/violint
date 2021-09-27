# violint
lint tool for videos.

# PreRequirements
This tool uses following tools. You have to install them before.
- opencv
- tesseract

# Install

```
go install github.com/masibw/violint@latest
```

# How to use
Run the command.
```
violint [target_file]
```

A window will appear, with the areas that have been detected by OCR surrounded by a blue frame and the areas with low contrast surrounded by a red frame.

![result](images/result.jpg)