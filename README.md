# ANSIZALIZER
A TUI to convert Images to ANSI strings using bubbletea

![Screenshot 2024-04-02 150412](https://github.com/Zebbeni/ansizalizer/assets/3377325/141c3662-7e70-4e82-ac0c-5db77adbf1c7)

## Features
- A keyboard-navigable Text-based UI
- File browser: Search .png and .jpeg image files and preview in real-time
- Export ANSI image strings to '.ansi' text files or copy directly to your Clipboard
- Save files individually or Batch Process All Images in a chosen directory
- Browse Lospec.com for cool color palettes

## Render Options
- Set output Width and Height of rendered text images (in characters)
- Choose character sets to use in output (ASCII, Unicode, or Custom)
- Render images with "true" colors or convert using Limited Color Palettes
- Generate new color palettes by sampling previewed image files
- Use Advanced settings to tweak pixel Sampling mode and Dithering options

![Screenshot 2024-04-02 155820](https://github.com/Zebbeni/ansizalizer/assets/3377325/24095f45-5c73-4654-a5e1-b491cda9dc66)

## To Run

(On Windows)
```bash
go install
go build
start ansizalizer.exe
```

![Screenshot 2024-04-02 155006](https://github.com/Zebbeni/ansizalizer/assets/3377325/d41df628-6c84-44e0-aa34-f7fcb72ed827)

## FAQ / Troubleshooting
**Q: The UI isn't rendering correctly**

Check your default console appearance settings. Make sure your chosen font, font size, and line height aren't the cause of the problem. 'DejaVu Sans Mono' works well for me on Windows.

**Q: My images look squashed / stretched**

Try adjusting the value of Char Size Ratio under Settings > Size. Depending on what font your console uses, your characters may have a width-to-height ratio different than 0.5.

**Q: My exported .ansi files take up more space than the original image**

The ANSI code that produces the text-rendered images isn't (currently) optimized for file size. If using this tool to batch process lots of text art for use in a game or application, I'd consider compressing the resulting text files and decompressing them as needed.
