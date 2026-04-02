# Ansizalizer

A versatile terminal UI application for converting images to ANSI art. Built with [Ansipx](https://github.com/Zebbeni/ansipx/tree/main) and [Bubble Tea](https://github.com/charmbracelet/bubbletea).
<img width="1005" alt="image" src="https://github.com/user-attachments/assets/089f847a-67a0-42d5-8548-b2a38b93f195" />
<img width="500" alt="image" src="https://github.com/user-attachments/assets/cc9c2b0a-6aad-4a5d-ab34-c3cd5ac8bdc7" />
<img width="500" alt="image" src="https://github.com/user-attachments/assets/49e6b118-3ae8-4dba-93d5-15549fdf71e7" />

## Features

- Keyboard-navigable text-based UI with contextual help tooltips
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
<img width="700" alt="image" src="https://github.com/user-attachments/assets/a1484601-104a-4c04-b68e-4ccaa93d4186" />

### Character Modes

Render with Unicode block characters, ASCII characters, or custom character sets.

#### Unicode blocks (ex. ▀▄)
<img width="800" alt="unicode blocks" src="https://github.com/user-attachments/assets/a1e22d7d-2963-483e-a9f6-907e069ae100" />

#### Ascii characters (ex. 0-9)
<img width="800" alt="ascii numbers" src="https://github.com/user-attachments/assets/cac00dbf-2a97-44b3-acc8-bec1f6b7400f" />

#### Custom text (ex. "isle of the dead " | repeat)
<img width="800" alt="custom text" src="https://github.com/user-attachments/assets/0a1b77c2-a74f-4498-ae20-15c00fd7e01b" />

### Color Palettes

Use true color (24-bit RGB) or limited palettes from Lospec.com. Generate palettes by sampling colors from any image.
#### Lospec (ex. Mona Lisa | Ammo8, Iridescent Crystal)
<img height="500" alt="image" src="https://github.com/user-attachments/assets/86d82973-16f0-41e6-b738-a166d18756b3" />
<img height="500" alt="image" src="https://github.com/user-attachments/assets/8ce8acee-233d-4649-b354-a778e38733b2" />


### Alpha Transparency

Transparent pixels render as empty space. Partial transparency at edges uses the best-fit Unicode block character for smooth outlines.

<img width="600" height="970" alt="image" src="https://github.com/user-attachments/assets/d3c331f5-cc69-4959-a928-35d9e70049b1" />

### Animated GIFs

Preview animated GIFs frame-by-frame with adjustable delay.
![Animated GIF](https://github.com/user-attachments/assets/fb364211-739e-428d-a58c-40a4482d987d)

### Themes

Switch between 6 built-in themes. Paletted themes derive colors from your selected palette.
#### Dark on Light (paletted)
<img width="500" alt="image" src="https://github.com/user-attachments/assets/a689fd90-4370-405e-8c10-4834d8051b3f" />

#### Light on Dark (paletted)
<img width="500" alt="image" src="https://github.com/user-attachments/assets/508888fe-897c-4d5b-afbe-36467055063f" />


### Batch Export

Export all images in a directory as `.ansi` files with a single operation.
<img width="1915" height="1043" alt="image" src="https://github.com/user-attachments/assets/531b02c9-36b3-439f-993e-e9d8523b5c7a" />


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
