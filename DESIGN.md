---
name: Xmas Gift Poll
description: A family gift draw built as a printed advent calendar; each person's draw is a door that opens after three wishes.
colors:
  calendar-red: "#9e1b32"
  door-ink: "#7a1426"
  pine: "#1e4d3a"
  gold-foil: "#d9a441"
  gold-foil-ink: "#3a2606"
  paper: "#ffffff"
  snow: "#f7f6f2"
  frost-line: "#e6e3da"
  ink: "#1c1c1c"
  holly-green: "#1e6b45"
  amber: "#b7791f"
  error-red: "#b3261e"
  night-ground: "#101914"
  night-paper: "#18241e"
  night-line: "#2a3a31"
  night-ink: "#f1efe8"
  night-pine: "#173d2e"
  night-gold: "#e2b259"
  night-mint: "#6fb894"
  night-error: "#f2726b"
typography:
  display:
    fontFamily: "Bricolage Grotesque, ui-sans-serif, system-ui, sans-serif"
    fontSize: "1.875rem"
    fontWeight: 800
    lineHeight: 1.25
  headline:
    fontFamily: "Bricolage Grotesque, ui-sans-serif, system-ui, sans-serif"
    fontSize: "1.5rem"
    fontWeight: 800
    lineHeight: 1.333
  title:
    fontFamily: "Bricolage Grotesque, ui-sans-serif, system-ui, sans-serif"
    fontSize: "1.25rem"
    fontWeight: 700
    lineHeight: 1.4
  body:
    fontFamily: "ui-sans-serif, system-ui, sans-serif"
    fontSize: "1rem"
    fontWeight: 400
    lineHeight: 1.5
  body-sm:
    fontFamily: "ui-sans-serif, system-ui, sans-serif"
    fontSize: "0.875rem"
    fontWeight: 400
    lineHeight: 1.43
  label:
    fontFamily: "ui-sans-serif, system-ui, sans-serif"
    fontSize: "0.875rem"
    fontWeight: 600
    lineHeight: 1.43
  numeral:
    fontFamily: "Bricolage Grotesque, ui-sans-serif, system-ui, sans-serif"
    fontSize: "2.25rem"
    fontWeight: 800
    lineHeight: 1
    fontFeature: "\"tnum\""
rounded:
  tile: "8px"
  field: "12px"
  box: "16px"
  pill: "9999px"
spacing:
  list-gap: "8px"
  gutter: "16px"
  card: "20px"
  stack: "20px"
  card-roomy: "24px"
  band-tuck: "32px"
components:
  header-band:
    backgroundColor: "{colors.pine}"
    textColor: "{colors.paper}"
  panel:
    backgroundColor: "{colors.paper}"
    textColor: "{colors.ink}"
    rounded: "{rounded.box}"
    padding: "20px"
  door:
    backgroundColor: "{colors.calendar-red}"
    textColor: "{colors.paper}"
    rounded: "{rounded.box}"
    padding: "20px"
  door-open:
    backgroundColor: "{colors.paper}"
    textColor: "{colors.ink}"
    rounded: "{rounded.box}"
    padding: "20px"
  door-number:
    backgroundColor: "{colors.gold-foil}"
    textColor: "{colors.door-ink}"
    typography: "{typography.numeral}"
    rounded: "{rounded.tile}"
    size: "72px"
  door-number-list:
    backgroundColor: "{colors.gold-foil}"
    textColor: "{colors.door-ink}"
    rounded: "{rounded.tile}"
    size: "40px"
  button-primary:
    backgroundColor: "{colors.calendar-red}"
    textColor: "{colors.paper}"
    rounded: "{rounded.field}"
    height: "48px"
    padding: "0 20px"
  button-accent:
    backgroundColor: "{colors.gold-foil}"
    textColor: "{colors.gold-foil-ink}"
    rounded: "{rounded.field}"
    height: "48px"
    padding: "0 20px"
  button-soft-pine:
    textColor: "{colors.pine}"
    rounded: "{rounded.field}"
    height: "40px"
    padding: "0 16px"
  input:
    backgroundColor: "{colors.paper}"
    textColor: "{colors.ink}"
    rounded: "{rounded.field}"
    height: "48px"
    padding: "0 12px"
  wish-row:
    backgroundColor: "{colors.snow}"
    textColor: "{colors.ink}"
    rounded: "{rounded.field}"
    padding: "12px 16px"
  language-switch-active:
    backgroundColor: "{colors.paper}"
    textColor: "{colors.pine}"
    rounded: "{rounded.pill}"
    padding: "6px 12px"
---

# Design System: Xmas Gift Poll

## Overview

**Creative North Star: "The Printed Advent Calendar"**

Every page is a sheet of a printed advent calendar. A pine-green band runs across the top, white paper cards are tucked into it, and the thing you are waiting for is a calendar-red door with a gold-foil number on it. The participant's draw stays behind a closed door until they have written three wishes; then the door swings open once and the inside is plain paper with the name of the person they buy for. The organizer sees the same vocabulary: each person is a numbered gold tile, and the draw itself is a red door with a gold button inside it.

The world is warm but tidy. It is a utility for relatives of every age, so density is low, type is large, and there is one obvious action per card. Festivity comes from materials (red, pine, gold foil, perforated door edges, a chunky grotesque numeral), never from decoration. The user's chosen direction explicitly refuses two things: the generic white SaaS card stack, and the snowflake or clip-art Christmas page.

Dark mode is the same calendar at night: the ground and paper turn pine-black, while calendar red and gold foil stay put.

**Key Characteristics:**
- Pine header band with the first card tucked 32px up into it.
- Calendar red is reserved for doors and the primary action.
- Gold-foil tiles with Bricolage Grotesque numerals mark every numbered or iconic item.
- Single phone-width column (max 576px) on every screen size.
- Exactly one authored motion: the door swing.
- Lucide line icons inlined as SVG, sized to the text.

## Colors

A printed-calendar palette: one saturated red, one deep pine, one gold foil, on warm snow-white paper with near-black ink.

### Primary
- **Calendar Red** (`calendar-red`): the closed door surface, the open door's frame, the primary button (`Add gift`, `Create the draw`, `Add`), in-text links, and the input caret. White text on it measures 7.9:1.
- **Door Ink** (`door-ink`): the deep red of numerals and icons printed on gold tiles, and of the favicon's gift glyph. Never used as a surface.

### Secondary
- **Pine** (`pine`): the header band on every page (also the `theme-color`), the soft link-copy buttons, the organizer's wish pips, and the active language pill's text.

### Tertiary
- **Gold Foil** (`gold-foil`): door-number tiles, the accent button that lives inside a red door (`Draw now`, `Go to the start page`), wish pips on the red door, the bullet dots in the receiver's wishlist, positive status text on the pine band, the text selection tint and the keyboard focus ring.
- **Gold Foil Ink** (`gold-foil-ink`): text on gold buttons (6.4:1).

### Neutral
- **Paper** (`paper`): cards (panels and opened doors), inputs.
- **Snow** (`snow`): the page ground below the band, and the inset rows that hold each wish.
- **Frost Line** (`frost-line`): card borders, list dividers, and the rule above the add-wish form.
- **Ink** (`ink`): all body text on paper and snow; secondary text uses ink at 70 to 75% opacity.
- **Holly Green, Amber, Error Red** (`holly-green`, `amber`, `error-red`): daisyUI status roles only (success toast and the copied state, the "needs a new draw" badge, form errors and error toasts).

### Night (dark theme `advent-night`, selected by `prefers-color-scheme: dark`)
- **Night Ground** (`night-ground`) replaces snow; **Night Paper** (`night-paper`) replaces paper; **Night Line** (`night-line`) replaces frost line; **Night Ink** (`night-ink`) replaces ink.
- **Night Pine** (`night-pine`) is the band. **Night Gold** (`night-gold`) is the accent button and pips; the door-number tiles keep the printed `gold-foil`.
- **Night Mint** (`night-mint`) carries info/success and replaces pine for the organizer's wish pips, because pine vanishes on night paper. **Night Error** (`night-error`) is the lifted error red.
- Calendar red is identical in both themes.

### Named Rules
**The Red Means Door Rule.** Calendar red covers a large surface only when it is a door (closed card or open frame). Elsewhere it appears as one primary button per form, links and the caret. No red section backgrounds, red headers or red decorative fills.

**The Foil Sits on Red or Pine Rule.** Gold is foil: it appears as tiles, pips, dots, buttons inside a red door, and status text on the pine band. It is never body text on paper or snow.

**The Same Calendar at Night Rule.** Dark mode swaps the ground, paper, lines and ink to pine-black values and keeps calendar red and gold foil. It is a lighting change, not a second palette.

## Typography

**Display Font:** Bricolage Grotesque, weights 600 to 800, self-hosted (OFL, latin and latin-ext subsets so Polish diacritics render), with ui-sans-serif and system-ui fallback.
**Body Font:** the system sans stack (ui-sans-serif, system-ui, sans-serif).

**Character:** a chunky, slightly quirky grotesque for everything that is printed on the calendar (headings, the brand, the numerals), against a neutral system sans that keeps forms and wish text plain and legible for every relative.

### Hierarchy
- **Display** (800, 1.875rem, line-height 1.25): the h1 greeting or party name in the pine band. The home page steps it up to 2.25rem from 640px. Balanced wrapping, breaks long names.
- **Headline** (800, 1.5rem): the receiver's name inside the opened door; the payoff of the whole product.
- **Title** (700, 1.25rem, line-height 1.4): card headings (`Your list`, `People`, the door's own heading). `How it works` uses the same weight one step down (1.125rem).
- **Body** (400, 1rem, line-height 1.5): wishes, explanations, inputs. Wish text preserves line breaks and wraps long URLs anywhere.
- **Body small** (400, 0.875rem): supporting lines under titles and inside doors, at 70 to 85% opacity of their surface's text color.
- **Label** (600, 0.875rem): form labels and person names' meta; field help drops to 0.75rem.
- **Numeral** (Bricolage 800, tabular figures): door numbers, sized to the tile: 2.25rem in the 72px door tile, 1.125rem in 40px list tiles, 3rem in the 80px door leaf.

### Named Rules
**The Printed Numeral Rule.** Anything numbered or emblematic sits in a gold tile with a Bricolage 800 numeral or a Lucide icon in door ink. Counts in running text use tabular figures.

**The Grotesque Is for Print Rule.** Bricolage Grotesque is used only for h1 to h3, the brand wordmark and door numerals. Body copy, labels, buttons and inputs stay in the system sans.

## Layout

One centered column, max 576px wide with a 16px side gutter, on phone and desktop alike; the desktop screen is the phone layout with more snow around it. The header band holds a 56px bar (brand left, PL/EN pill right) and the hero text, and ends with 48px of bottom padding; the main column pulls up 32px so the first card sits half on the band. Cards stack with a 20px gap. Card padding is 20px; doors grow to 24px from 640px. Lists inside cards use an 8px gap between wish rows, or full-bleed rows separated by frost-line dividers for people. Primary buttons and the main inputs are full width and 48px tall. The body respects the bottom safe-area inset (at least 32px).

**The Tucked Card Rule.** The first card on every page overlaps the pine band by 32px. The band is never a standalone header with content floating below it.

## Elevation & Depth

Mostly flat, with depth reserved for doors. daisyUI depth and noise are off, so controls have no bevel or shadow. Paper panels carry a barely visible ambient shadow plus a frost-line border. The closed red door is the most lifted object on any page, with a shadow tinted from its own red. The opened door drops back to a soft neutral shadow and lets its 3px red frame do the work. Toasts use a standard large shadow.

### Shadow Vocabulary
- **Paper rest** (`box-shadow: 0 1px 2px rgb(0 0 0 / 0.04), 0 4px 16px -8px rgb(0 0 0 / 0.12)`): panels.
- **Door lift** (`box-shadow: 0 2px 4px rgb(0 0 0 / 0.08), 0 12px 24px -12px color-mix(in oklab, var(--color-primary) 70%, black)`): the closed red door.
- **Open door** (`box-shadow: 0 2px 4px rgb(0 0 0 / 0.06), 0 12px 24px -14px rgb(0 0 0 / 0.25)`): the opened door.

### Named Rules
**The Doors Cast, Paper Rests Rule.** Only the closed door casts a coloured shadow. Every other surface is flat or sits on the paper-rest shadow; no new elevation levels.

## Shapes

Soft, printed-card corners on a three-step scale: 16px for cards and doors, 12px for fields, buttons and inset wish rows, 8px for gold tiles. Pills (fully rounded) are used only for the language switch and wish pips. A closed door carries a perforation: a 2px dashed white line at 30% opacity, inset 6px from the edge, following the door's corner (12px). An opened door replaces it with a solid 3px calendar-red frame, as if the door leaf has been torn away and only the frame is left.

**The Perforation Rule.** Dashed inner edges belong only to closed doors. Opened doors use the solid red frame. No other surface gets a dashed or double border.

## Components

### Buttons
Tactile, full-width and thumb-sized; daisyUI buttons with depth turned off.
- **Shape:** gently rounded (12px).
- **Primary:** calendar red with white text; large size is 48px tall with 1.125rem text and 20px side padding, used full width for the one main action on a card. Inline variant (40px) sits next to an input (`Add`).
- **Accent:** gold foil with gold-foil ink, used only inside a red door. Disabled on a door it becomes white at 15% with white text at 60%.
- **Soft pine:** pine text on an 8% pine tint, for per-person `Copy link` and `Share`; collapses to an icon-only square under 640px.
- **On the band:** white at 15% with white text (hover 25%), for the organizer's own copy link.
- **Ghost square:** icon-only, ink at 70%, for delete (hover turns error red) and rename.
- **Hover / Focus:** daisyUI's 0.2s colour transition; keyboard focus everywhere is a 2px gold outline offset 2px.
- **Copied state:** the copy button briefly switches to the success colour and its label to `Copied!` for 2 seconds.

### Chips
- **Language switch:** a pill track of white at 10% on the band; the active language is a white pill with pine text, the other is white at 80%.
- **Badge:** daisyUI soft warning badge for "needs a new draw".

### Cards / Containers
- **Panel:** paper, 16px corners, 1px frost-line border, paper-rest shadow, 20px padding. Holds lists and forms.
- **Door (closed):** calendar red with white text, 16px corners, perforated inner edge, door-lift shadow, 20px padding (24px from 640px). Contents sit above the perforation (`position: relative`).
- **Door (open):** paper with a 3px calendar-red frame, open-door shadow, 20px padding.
- **Wish row:** snow inset rows with 12px corners inside a panel or open door; the receiver's wishes lead with a small gold dot.
- **Band note:** white at 10% inset row on the pine band, used for the organizer's "save this link" note with a gold lock icon.

### Inputs / Fields
- **Style:** daisyUI input and textarea on paper, 12px corners, 1px border at 20% ink. The main inputs are 48px (`input-lg`); the wish textarea grows with its content from 80px.
- **Focus:** the border darkens to ink; keyboard focus adds the gold ring. The caret is calendar red.
- **Error:** error-red border, `aria-invalid`, and an error line below with a circle-alert icon. On the band, the rename error is gold.

### Navigation
The pine band is the only navigation: brand wordmark (Bricolage 700) with a 28px gold gift tile on the left, the PL/EN pill switch on the right. No menus.

### Door Number (signature)
A gold-foil square tile with 8px corners holding a Bricolage 800 numeral or a Lucide icon in door ink. Sizes in use: 28px (brand), 40px (steps and people list), 48 to 56px (draw doors, receiver), 72px (the participant's closed door), 80px (the swinging door leaf).

### Wish Pips (signature)
A row of three small circles (10px, 12px on the door) with a 2px border in the current colour: filled for each written wish, hollow for each still owed. Gold on the red door, pine (night mint in dark mode) in the people list. Always paired with the "n of 3 gifts" text.

### Door Swing (signature motion)
When the third wish unlocks the draw, the closed door is laid over the opened one and swings away once: `rotateY(-112deg)` around its left edge over 1.1s with `cubic-bezier(0.16, 1, 0.3, 1)` after a 0.35s pause, in 1200px perspective, then fades and hides. Under reduced motion the leaf is not rendered. Info toasts fade up and away after 4s (0.3s ease-in), or simply hide under reduced motion.

**The One Swing Rule.** The door swing is the only authored animation, and it plays once per unlock. New screens do not add entrance animations, looping effects or decorative motion.

## Do's and Don'ts

### Do:
- **Do** tuck the first card 32px into the pine band on every page.
- **Do** mark numbered or emblematic items with a gold-foil door-number tile in door ink.
- **Do** keep one full-width 48px primary button per card, and use the gold accent button only inside a red door.
- **Do** inline Lucide icons through the `icon` template function at 1.25em, stroked in `currentColor`.
- **Do** pair every visual progress mark (pips) with its written count.
- **Do** design every surface for both `advent` and `advent-night`; check that pine-coloured marks switch to night mint in dark mode.

### Don't:
- **Don't** add snowflakes, clip-art, ornaments or other Christmas decoration; the festive feel comes from red, pine, gold and the door.
- **Don't** fall back to a plain stack of white cards on a white page without the pine band and the door.
- **Don't** use calendar red for anything larger than a button unless it is a door.
- **Don't** set gold foil as text on paper or snow.
- **Don't** set body copy, labels or buttons in Bricolage Grotesque.
- **Don't** add a second animation or replay the door swing.
- **Don't** use emoji or text glyphs as icons.
