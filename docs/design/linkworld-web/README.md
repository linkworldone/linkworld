# LinkWorld Web Baseline Documentation

This directory contains comprehensive baseline documentation for the Link World React web frontend project, created as input for the design/plan phase of the UI redesign pipeline.

## Files

### 1. `components.md` (49 lines)
**Component Inventory** — All 13 UI/shared/layout/wallet components with:
- File paths
- Component purposes
- Key props
- Reusability ratings (高/中/needs-redo)
- Key findings on shadcn/base-ui usage, hardcoded colors, and redesign risks

### 2. `color-mapping.md` (205 lines)
**Color System Analysis** — Current theme system + new theme mapping:
- Current CSS variables (light & dark mode)
- Tailwind config extensions (semantic tokens, hardcoded accents)
- Hardcoded gradient definitions (3 in CSS, 15+ in components)
- Hardcoded color hotspots (6 scattered hex values)
- New theme color mapping (deep blue #0C2340 → bright blue #1E40AF → gold #B8895A)
- Proposed token remapping with 4-phase migration strategy
- Risk assessment table (3 scatter types, severity levels)

### 3. `utils.md` (229 lines)
**Utilities & Helpers Inventory** — 35+ exported utilities:
- Core utils: `cn()`, 6 format functions, 6 pendingSync functions
- 7 data hooks (useUser, useNotification, useTransactionFlow, useBilling, useDeposit, useOperator)
- 5 contract hooks (deposit, registry, service manager, traffic card, payment)
- 5 API services + 1 mock service
- Configuration files (contracts, API, wagmi, chains)

### 4. `project-scan.md` (231 lines)
**Project Structure Overview** — Full system baseline:
- Tech stack & versions (React 19, TS 5.6, Vite 6, Tailwind 3.4, wagmi 2.14)
- Directory structure with descriptions
- Routes & pages (9 routes, 1 landing + 8 app pages)
- Data flow architecture (React Query, wagmi, localStorage)
- Services & config files
- Component counts & LOC metrics
- Build/dev scripts
- Critical dependencies & design patterns
- **Repainting Risk Points** — Highest/medium/low risk categories with specific file locations
- Quality notes & redesign checklist (14 items)

## Key Findings Summary

### Color Scatter & Migration Burden
- **6 hardcoded hex color values** scattered across `index.css` and `tailwind.config.ts`
- **20+ component usages** of hardcoded gradient classes (`accent-gradient-bg`, `accent-gradient-text`, `.btn-primary`)
- Highest risk: Landing page (5 gradient usages), Dashboard (4), Header/TabBar/ConnectButton (2-3 each)
- Migration requires: CSS variable updates + gradient class replacements + component class refactoring

### Component Reusability Assessment
- **High reusability** (7 components): Button, Badge, Tabs, AmountDisplay, BottomSheet, AppLayout, TabBar, ConnectButton
- **Medium reusability** (4 components): EmptyState, GuardCard (router-coupled), Header (location-branched), RegisterSheet (business-specific)
- **Refactor needed**: StatusBadge (hardcoded statusConfig), GuardCard (tight navigation coupling), Header (complex useLocation logic)

### Redesign Risk Points (Top 5)
1. **Gradient class scatter** — 12+ files use `accent-gradient-bg/text`, all need batch replacement
2. **Header complexity** — useLocation-based 3-way branching logic is brittle
3. **GuardCard coupling** — Navigation deeply embedded, makes it hard to reuse
4. **StatusBadge config** — Hardcoded JS object instead of theme-driven, not maintainable
5. **Missing component docs** — No JSDoc comments, prop changes risk breaking changes across codebase

## Usage

These files are designed to feed directly into the **design** and **plan** phases:

1. **Design Phase**: Use `components.md` + `color-mapping.md` to:
   - Identify which components need visual updates
   - Plan new gradient definitions for deep-blue + gold theme
   - Prioritize redesign effort (high reusability components → medium → needs-redo)

2. **Plan Phase**: Use `project-scan.md` + `utils.md` to:
   - Understand implementation scope (14-item redesign checklist)
   - Plan file change order (CSS vars → tailwind → gradient classes → components)
   - Identify refactoring opportunities (Header, GuardCard, StatusBadge)
   - Estimate effort (color hotspot count, risk points)

3. **Implementation Phase**: Reference exact line numbers and file paths from `color-mapping.md` for surgical edits

## Tech Stack Context

- **React 19 + TypeScript 5.6**: Modern React with strict typing
- **Vite + Tailwind 3.4**: Fast builds, utility-first styling with CSS variable integration
- **wagmi 2.14 + RainbowKit**: Web3 wallet integration
- **React Router 7.1 + React Query 5.62**: Client-side routing + server state management
- **shadcn/base-ui**: Headless components with CVA variant system

## Not Included (Out of Scope)

- Individual page layout redesigns (Landing, Dashboard, etc.) — handled separately per page
- Design tokens/design system Figma file — to be created in design phase
- Component Storybook setup — to be added as quality improvement post-redesign
- Copy/UX text — focusing on visual structure only

---

**Created**: 2026-06-04 (Scan phase, pipeline baseline)  
**For**: Link World web frontend UI redesign (theme: deep blue #0C2340 → bright blue #1E40AF → gold #B8895A)
