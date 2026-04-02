# Ansizalizer

A terminal UI application for converting images to ANSI art. Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea).

<!-- TODO: Add hero screenshot showing the full app with an image rendered -->
![Hero Screenshot](https://placeholder.com/hero.png)

## Features

- Keyboard-navigable TUI with contextual help tooltips
- File browser for selecting `.png`, `.jpg`, `.jpeg`, and `.gif` images
- Real-time preview of rendered ANSI art
- Animated GIF support with per-frame rendering
- Alpha transparency with partial-block rendering for smooth edges
- Export to `.ansi` text files or copy directly to clipboard
- Batch export entire directories (with subdirectory support)
- Browse [Lospec.com](https://lospec.com) for color palettes
- Generate palettes by sampling image colors
- Save and load settings presets as JSON
- 6 built-in color themes (light/dark, paletted, transparent)
- User preferences (help visibility, splash screen, theme/settings restore)
- Animated splash screen on startup

## Installation

### From source

Requires [Go 1.21+](https://go.dev/dl/).

```bash
git clone https://github.com/Zebbeni/ansizalizer.git
cd ansizalizer
go build
```

**Windows:**
```bash
start ansizalizer.exe
```

**Mac/Linux:**
```bash
./ansizalizer
```

### Debug mode

Run with `-debug` to enable logging to `console.log`:

```bash
./ansizalizer -debug
```

## Screenshots

### Image Browser and Preview

<!-- TODO: Screenshot showing the file browser with an image selected and ANSI preview in the viewer -->
![Browse and Preview](https://placeholder.com/browse.png)

### Character Modes

Render with Unicode block characters, ASCII characters, or custom character sets.

<!-- TODO: Side-by-side screenshots showing the same image in Unicode, ASCII, and Custom modes -->
![Character Modes](https://placeholder.com/characters.png)

### Color Palettes

Use true color (24-bit RGB) or limited palettes from Lospec.com. Generate palettes by sampling colors from any image.

<!-- TODO: Screenshot showing palette browser or an image rendered with a limited palette -->
![Color Palettes](https://placeholder.com/palettes.png)

### Alpha Transparency

Transparent pixels render as empty space. Partial transparency at edges uses the best-fit Unicode block character for smooth outlines.

<!-- TODO: Screenshot showing a transparent PNG rendered with alpha, demonstrating the partial block characters at edges -->
![Alpha Transparency](https://placeholder.com/alpha.png)

### Animated GIFs

Preview animated GIFs frame-by-frame with adjustable delay.

<!-- TODO: Screenshot or GIF recording showing an animated GIF rendering in the viewer -->
![Animated GIF](https://placeholder.com/animation.png)

### Themes

Switch between 6 built-in themes. Paletted themes derive colors from your selected palette.

<!-- TODO: Screenshot showing the theme selector dropdown, or side-by-side of light vs dark theme -->
![Themes](https://placeholder.com/themes.png)

### Batch Export

Export all images in a directory as `.ansi` files with a single operation.

<!-- TODO: Screenshot of the batch export modal with source/destination panels -->
![Batch Export](https://placeholder.com/batch.png)

## Render Settings

| Setting | Options | Description |
|---------|---------|-------------|
| **Characters** | Unicode (Block), ASCII, Custom | Character set for rendering |
| **Unicode Blocks** | Full, Half, Quarter, Shade (Light/Med/Heavy) | Block character granularity |
| **Colors** | True Color, Limited Palette | 24-bit RGB or palette-matched |
| **Palette Source** | Lospec, From File, Generated | Where to get limited palettes |
| **Size** | Width, Height, Fit/Fill/Stretch | Output dimensions in characters |
| **Adjust** | Brightness, Contrast | Image adjustments (-100 to 100) |
| **Text Style** | Bold, Italic, Underline, Strikethrough | ANSI text attributes |
| **Advanced** | Sampling, Dithering, Serpentine | Resampling and dither options |
| **Animation** | Delay (ms) | GIF frame delay (10-2000ms) |
| **Alpha** | Render Threshold, Transparency, Trim | Alpha channel handling |

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| Arrow keys | Navigate menus and settings |
| Enter | Select / expand menu |
| Esc | Back / exit menu |
| Tab | Next item |
| `+` / `-` | Expand / collapse submenu |
| Ctrl+C | Copy rendered output to clipboard |
| Ctrl+S | Save rendered output to file |
| Ctrl+H | Toggle help tooltips |
| Ctrl+Q | Quit |

## Rendering Engine

Ansizalizer uses [ansipx](https://github.com/Zebbeni/ansipx) as its rendering engine -- a standalone Go library for converting images to ANSI art. You can use ansipx independently in your own projects.

## FAQ / Troubleshooting

**Q: The UI isn't rendering correctly**

Check your terminal's font settings. Make sure your font supports Unicode block characters. Monospace fonts like *DejaVu Sans Mono* or *JetBrains Mono* work well.

**Q: My images look squashed / stretched**

Adjust the Char Size Ratio under Settings > Size. Different terminal fonts have different width-to-height ratios. Typical values are 0.45-0.50.

**Q: My exported .ansi files are larger than the original image**

ANSI escape codes are verbose -- each colored character requires multiple bytes. For use in applications, consider compressing the `.ansi` files.

**Q: How do I use a custom color palette?**

Save a `.hex` file (one hex color per line, e.g. `#ff0000`) to the palettes directory, then select it from Settings > Colors > From File.
