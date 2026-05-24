# SPEC.md — Inventaris App (Visual / Frontend)

This document is the visual specification for the first two pages of the inventory app:
the Login page and the Dashboard page. Stack: plain HTML, CSS, and JS (no framework).

---

## Page 1 — Login

**URL:** `/login`

### Layout
- Page background: light gray (`#f7f6f3`)
- Content centered both horizontally and vertically
- Single white card in the center, max-width 300px, border-radius 12px, thin border

### Card contents (top to bottom)
1. **Logo/icon** — small box with a package icon, centered
2. **Title** — "Inventaris" (font-weight 500, 15px)
3. **Subtitle** — "Masuk ke akun Anda" (12px, gray)
4. **Username field** — label "Username" above, input with user icon on the left
5. **Password field** — label "Password" above, input with lock icon on the left, value shown as `••••••••`
6. **"Masuk" button** — full width, black background (`#1a1a18`), white text, border-radius 8px

### Behavior
- Click "Masuk" → redirect to `/dashboard`
- No "Forgot password" feature for now

---

## Page 2 — Dashboard

**URL:** `/dashboard`

### Main layout
Two columns: **sidebar** on the left + **main content** on the right.
Page background: `#f7f6f3`

---

### Sidebar (left)

- Width: 190px
- Background: white (`#fff`)
- Right border: 0.5px solid `#e0ded8`

**Top section — Logo**
- Small package icon + "Inventaris" text
- Bottom border separates it from the navigation

**Navigation menu**
Section label: "MENU" (small, gray, uppercase)

| Menu item | Icon (Tabler) |
|-----------|---------------|
| Dashboard | `ti-layout-dashboard` |
| Produk | `ti-box` |
| Kategori | `ti-tag` |
| Stok | `ti-arrows-exchange` |
| Supplier | `ti-truck` |

- Active item: background `#f7f6f3`, black text, border-left 2px solid `#1a1a18`
- Inactive item: gray text `#888`

**No** Settings or Logout items in the sidebar.

---

### Main content (right)

**Topbar**
- Position: top-right of the content area
- Contains only a **circular avatar** with user initials (e.g. "LA")
- Avatar background: `#e8e6e0`, border-radius 50%
- Click avatar → show **dropdown** (see below)

**Avatar dropdown**
- Position: absolute, appears below the avatar, right-aligned
- Contents:
  - User name: "Lion"
  - Role: "Administrator" (small, gray)
  - Divider line
  - "Keluar" button with `ti-logout` icon, red text (`#c0392b`)
- Click outside dropdown → dropdown closes

**Dashboard content area**
- Heading: "Dashboard" (18px, font-weight 500)
- Subheading: "Selamat datang kembali, Lion" (13px, gray)
- Below: placeholder text "Konten belum tersedia." (light gray)
- No widgets, cards, or data for now

---

## General notes

- Font: system sans-serif
- Primary colors: `#1a1a18` (black/text), `#f7f6f3` (background), `#e0ded8` (border)
- All borders: 0.5px solid
- Icons: Tabler Icons (outline only), loaded via CDN
- No dark mode — light mode only
- No JS framework — interactivity via vanilla JS
