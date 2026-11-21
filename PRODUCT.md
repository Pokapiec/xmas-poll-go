# Product

<!-- impeccable:product-schema 1 -->

## Platform

web

## Users

- **Organizer**: one person in a family or friend group who runs a Secret Santa-style gift exchange. They create a party, add everyone's name, shuffle who buys a gift for whom, then send each person their own private link (usually over a messaging app).
- **Participants**: family members and friends of all ages and technical comfort. They open their personal link, usually on a phone, write a wishlist, and come back to see who they are buying for and what that person wants.

Everyone is on a phone most of the time, the organizer included.

## Product Purpose

Make a family gift exchange painless: no accounts, no sign-up, just links. Success means every participant has written a wishlist and knows who they are buying for, and nobody draws themselves.

## Positioning

A private, link-based tool for one family's Christmas. No accounts, no emails, no ads. Each person's link is the only key they need.

## Operating Context

- The organizer's page URL is effectively their admin key; each person's link is their private key. Links get pasted into WhatsApp, Messenger and the like.
- A participant only sees their draw after adding at least 3 gifts to their own wishlist (deliberate nudge).
- Shuffling is disabled once everyone has a receiver; adding a new person allows a reshuffle (which reassigns everyone).
- Used in the weeks before Christmas, opened a handful of times per person.

## Capabilities and Constraints

- Go server-rendered HTML (`html/template`), htmx `hx-boost` for smooth form posts, Tailwind v4 + daisyUI 5 built with the standalone CLI. No Node.js, no client framework.
- Works without JavaScript (plain forms with redirects).
- Bilingual: Polish (default) and English on every page, with a language switch. Polish copy uses proper diacritics.
- No authentication; secrecy relies on unguessable UUID links.

## Brand Commitments

Name: "Xmas Gift Poll". Tone: warm and a little festive, for family and friends, while staying a tidy utility.

## Evidence on Hand

No testimonials, logos, imagery or usage data. Do not invent any.

## Product Principles

1. Phone-first: every action is reachable with a thumb and readable at arm's length.
2. One obvious next step per screen: the organizer and participants should always know what to do next.
3. Links are the product: make sharing a person's link effortless and make it clear whose link is whose.
4. Friendly to non-technical relatives: plain language, no jargon, forgiving forms.

## Accessibility & Inclusion

Participants span all ages: comfortable tap targets, good contrast in both light and dark modes, readable type sizes.
