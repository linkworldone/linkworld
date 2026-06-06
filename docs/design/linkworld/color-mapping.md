# LinkWorld Color System Mapping

## CSS Variables (Root :root) from index.css

| Semantic Name | CSS Variable | Current Value | OKLCH Parse | Purpose |
|---------------|--------------|---------------|-------------|---------|
| **Background** | `--background` | `oklch(0.07 0.02 270)` | L:7%, C:0.02, H:270 | Primary dark background |
| **Foreground** | `--foreground` | `oklch(0.88 0 0)` | L:88%, C:0, H:0 | Primary text color (near white) |
| **Card** | `--card` | `oklch(0.10 0.02 270)` | L:10%, C:0.02, H:270 | Card/panel backgrounds |
| **Card Foreground** | `--card-foreground` | `oklch(0.88 0 0)` | L:88%, C:0, H:0 | Text on cards |
| **Popover** | `--popover` | `oklch(0.10 0.02 270)` | L:10%, C:0.02, H:270 | Popover/dropdown backgrounds |
| **Popover Foreground** | `--popover-foreground` | `oklch(0.88 0 0)` | L:88%, C:0, H:0 | Text on popovers |
| **Primary (MAIN THEME)** | `--primary` | `oklch(0.62 0.21 255)` | L:62%, C:0.21, H:255 (blue hue) | Button primary, active states, links |
| **Primary Foreground** | `--primary-foreground` | `oklch(1 0 0)` | L:100%, C:0, H:0 | Text on primary buttons (white) |
| **Secondary** | `--secondary` | `oklch(0.15 0.02 270)` | L:15%, C:0.02, H:270 | Secondary button/surface |
| **Secondary Foreground** | `--secondary-foreground` | `oklch(0.88 0 0)` | L:88%, C:0, H:0 | Text on secondary |
| **Muted** | `--muted` | `oklch(0.15 0.02 270)` | L:15%, C:0.02, H:270 | Disabled/muted backgrounds |
| **Muted Foreground** | `--muted-foreground` | `oklch(0.53 0 0)` | L:53%, C:0, H:0 | Disabled/muted text (gray) |
| **Accent** | `--accent` | `oklch(0.15 0.02 270)` | L:15%, C:0.02, H:270 | Accent elements |
| **Accent Foreground** | `--accent-foreground` | `oklch(0.88 0 0)` | L:88%, C:0, H:0 | Text on accent |
| **Destructive** | `--destructive` | `oklch(0.577 0.245 27.325)` | L:57.7%, C:0.245, H:27.325 (red-orange) | Delete/error states |
| **Border** | `--border` | `oklch(0.18 0.02 270)` | L:18%, C:0.02, H:270 | Border colors |
| **Input** | `--input` | `oklch(0.18 0.02 270)` | L:18%, C:0.02, H:270 | Form input backgrounds |
| **Ring** | `--ring` | `oklch(0.62 0.21 255)` | L:62%, C:0.21, H:255 (blue) | Focus ring color (matches primary) |
| **Chart 1** | `--chart-1` | `oklch(0.62 0.21 255)` | Blue | Chart/visualization |
| **Chart 2** | `--chart-2` | `oklch(0.65 0.18 290)` | Purple-blue | Chart variant |
| **Chart 3** | `--chart-3` | `oklch(0.70 0.15 195)` | Cyan | Chart variant |
| **Chart 4** | `--chart-4` | `oklch(0.75 0.18 145)` | Green | Chart variant |
| **Chart 5** | `--chart-5` | `oklch(0.70 0.18 55)` | Yellow-orange | Chart variant |
| **Radius** | `--radius` | `0.75rem` | 12px | Standard border radius |
| **Sidebar** | `--sidebar` | `oklch(0.10 0.02 270)` | L:10%, C:0.02, H:270 | Sidebar background |
| **Sidebar Foreground** | `--sidebar-foreground` | `oklch(0.88 0 0)` | L:88%, C:0, H:0 | Sidebar text |
| **Sidebar Primary** | `--sidebar-primary` | `oklch(0.62 0.21 255)` | Blue | Sidebar active/primary |
| **Sidebar Primary Foreground** | `--sidebar-primary-foreground` | `oklch(1 0 0)` | White | Sidebar primary text |
| **Sidebar Accent** | `--sidebar-accent` | `oklch(0.15 0.02 270)` | L:15%, C:0.02, H:270 | Sidebar accent |
| **Sidebar Accent Foreground** | `--sidebar-accent-foreground` | `oklch(0.88 0 0)` | L:88%, C:0, H:0 | Sidebar accent text |
| **Sidebar Border** | `--sidebar-border` | `oklch(0.18 0.02 270)` | L:18%, C:0.02, H:270 | Sidebar border |
| **Sidebar Ring** | `--sidebar-ring` | `oklch(0.62 0.21 255)` | Blue | Sidebar focus ring |

---

## Tailwind Config Color Extensions (`tailwind.config.ts`)

| Tailwind Class | Hex Color | Purpose |
|----------------|-----------|---------|
| **surface** DEFAULT | `#0a0a14` | Main dark background |
| **surface** card | `#0f0f1a` | Card/component backgrounds |
| **surface** secondary | `#1a1a2e` | Secondary surface |
| **surface** gradient-from | `#1a1a3e` | Gradient start (darker blue-gray) |
| **surface** gradient-to | `#0f1a2e` | Gradient end (darker blue) |
| **brand** blue | `#3b82f6` | Primary brand blue |
| **brand** purple | `#8b5cf6` | Brand purple accent |
| **brand** cyan | `#06b6d4` | Brand cyan |
| **status** success | `#22c55e` | Positive state (green) |
| **status** warning | `#f59e0b` | Warning state (amber) |
| **status** danger | `#ef4444` | Error/danger state (red) |
| **text** primary | `#e0e0e0` | Primary text (light gray) |
| **text** secondary | `#888888` | Secondary text (medium gray) |
| **text** muted | `#555555` | Muted text (dark gray) |
| **border** DEFAULT | `#1a1a2e` | Border color |
| **maxWidth** mobile | `430px` | Mobile viewport max-width |

---

## Current Theme Color (Primary & Background) - CURRENT STATE

### Primary Button/Active Color
- **Current**: `oklch(0.62 0.21 255)` = bright blue (~#3b82f6 in sRGB)
- **Used In**: Button default variant, active tab states, active icon color in TabBar, focus rings
- **Visual**: Bright blue, good contrast on dark backgrounds

### Background Surfaces
- **Root Background** (`--background`): `oklch(0.07 0.02 270)` = near-black with slight blue tint (#0a0a14)
- **Card Backgrounds**: `oklch(0.10 0.02 270)` = very dark blue-gray (#0f0f1a)
- **Gradient Surfaces**: `#1a1a3e` → `#0f1a2e` (gradient from darker blue to deeper blue)

### Summary
- **Overall Theme**: Dark mode with blue-dominant color scheme
- **Primary accent**: Bright blue (#3b82f6)
- **Secondary accent**: Purple (#8b5cf6)
- **Status colors**: Green (success), Amber (warning), Red (danger)

---

## Hardcoded Colors in Source Code

| File | Line | Hex/Color | Usage |
|------|------|-----------|-------|
| `/src/components/wallet/ConnectButton.tsx` | 12 | `bg-brand-blue` | Connect button background |
| `/src/components/layout/TabBar.tsx` | 34 | `text-brand-blue` | Active tab icon color |
| `/src/pages/Dashboard.tsx` | 64 | `text-brand-blue` | "Data Used" label color |
| `/src/pages/Dashboard.tsx` | 70 | `text-brand-purple` | "Calls" label color |
| `/src/pages/Dashboard.tsx` | 76 | `text-status-warning` | "Est. Bill" label color |
| `/src/pages/Landing.tsx` | 27 | `from-brand-blue to-brand-purple` | Logo gradient |
| `/src/pages/Landing.tsx` | 33 | `from-brand-blue to-brand-purple` | Hero circle gradient |
| `/src/pages/Landing.tsx` | 52 | `text-brand-blue` | "50+ Countries" number |
| `/src/pages/Landing.tsx` | 56 | `text-status-success` | "2.5% Platform Fee" number |
| `/src/pages/Landing.tsx` | 60 | `text-brand-purple` | "0 KYC" number |
| `/src/pages/Billing.tsx` | 45 | `bg-brand-blue` | Filter button active state |
| `/src/components/layout/Header.tsx` | 27 | `from-brand-blue to-brand-purple` | Avatar gradient |

---

## Replacement Plan for Theme Recolor (Deep Blue Gradient)

**Target Primary Gradient**: `linear-gradient(135deg, #0C2340 0%, #1E40AF 50%)`
- `#0C2340`: Deep navy/midnight blue
- `#1E40AF`: Strong blue
- Direction: 135° diagonal

### CSS Variables to Update in `index.css`
1. `--primary`: Change from `oklch(0.62 0.21 255)` to blue value matching `#1E40AF`
2. `--background`: Possibly adjust to match gradient tone
3. `--card`: Consider slight adjustment for harmony
4. Consider adding new gradient variables for gradient applications

### Tailwind Config Updates in `tailwind.config.ts`
1. `brand.blue`: Update from `#3b82f6` to `#1E40AF` (or keep if used elsewhere)
2. Add new `gradient` utilities if needed:
   - `bg-gradient-primary`: `linear-gradient(135deg, #0C2340 0%, #1E40AF 50%)`
3. Consider `primary-dark`: `#0C2340` and `primary-main`: `#1E40AF`

### Files Affected (Using Colors)
- `/src/pages/Landing.tsx`: Hero gradient (logo, circle)
- `/src/components/layout/Header.tsx`: Avatar gradient
- `/src/pages/Dashboard.tsx`: Data/Calls/Bill labels
- `/src/pages/Billing.tsx`: Filter active button
- `/src/components/layout/TabBar.tsx`: Active tab color

---

## Summary Statistics
- **Total CSS Variables**: 28 defined in `:root`
- **Total Tailwind Extensions**: 15 color groups + 1 size extension
- **Hardcoded Colors**: 12 instances across 5 files
- **Status Colors**: 3 (success, warning, danger)
- **Primary Color Hue**: 255° (blue range)
