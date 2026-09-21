# Att använda hjälp

Denna anvisning förklarar hur man använder detta inbäddade hjälp system i CDF Console verktyget.

## For End Users

### Accessing Help

You can access help documentation in two ways:

1. **Direct Navigation** - Navigate to the Help page from the main menu
2. **Context Help Buttons** - Click the help icon (?) next to features throughout the application

### Browsing Documentation

- Use the folder tree on the left sidebar to browse topics
- Click on any document to view its content
- Click links within documents to navigate between related topics

## For Developers

### Adding Context Help to Pages

Import and use the `ContextHelp` component to add help buttons:

```vue
<template>
  <div>
    <h2>
      Device Management
      <ContextHelp docPath="devices/overview.md" />
    </h2>
  </div>
</template>

<script setup>
import ContextHelp from '@/components/ContextHelp.vue'
</script>
```

### ContextHelp Component Props

- **docPath** (required) - Path to the help document relative to docs folder (e.g., "devices/overview.md")
- **tooltip** (optional) - Custom tooltip text (default: "View help documentation")
- **icon** (optional) - MDI icon name (default: "mdi-help-circle-outline")
- **size** (optional) - Button size: "x-small", "small", "default", "large", "x-large" (default: "small")
- **variant** (optional) - Button variant: "text", "outlined", "flat", "elevated" (default: "text")
- **color** (optional) - Button color (default: "grey")

### Example with Custom Styling

```vue
<ContextHelp
  docPath="troubleshooting/common-issues.md"
  tooltip="Get help with troubleshooting"
  icon="mdi-help"
  size="default"
  variant="outlined"
  color="primary"
/>
```

### Adding New Documentation

1. Create markdown files in `server/help/docs/` folder
2. Organize files into topic folders (e.g., devices/, networking/, troubleshooting/)
3. Use relative links to reference other documents: `[Link Text](../other-folder/document.md)`
4. Add images to `server/help/docs/images/` folder
5. Rebuild the Go backend to embed new documentation

### Adding Images

Images are embedded in the Go binary and served via the API:

1. **Add image file** to `server/help/docs/images/` folder (PNG, JPG, GIF, SVG supported)
2. **Reference in markdown** using relative path:
   ```markdown
   ![Alt text description](images/screenshot.png)
   ```
3. **Rebuild backend** to embed the image:
   ```bash
   cd server && go build
   ```

**Example:**
```markdown
## Dashboard View

![Dashboard Overview](images/dashboard-overview.png)

The dashboard shows all devices in your fleet.
```

The frontend automatically converts `images/filename.png` to the API endpoint `/api/help/images/filename.png`.

### Markdown Link Format

For inter-document links, use relative paths:

- Same folder: `[Link](document.md)`
- Parent folder: `[Link](../folder/document.md)`
- Root level: `[Link](/README.md)`

External links work normally:
- `[External](https://example.com)`
