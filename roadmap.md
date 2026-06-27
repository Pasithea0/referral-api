# Roadmap — Referral API

## Phase 1 — MVP (Complete)
- [x] Referral code generation (10-char, no ambiguous chars)
- [x] Campaign management (create, list)
- [x] Email + Discord collection on campaign pages
- [x] Duplicate detection (same email/discord returns existing code)
- [x] Return-visit cookie on campaign pages
- [x] Webhook for signup tracking (called from theintrodb-web)
- [x] Referrer dashboard (lookup code stats)
- [x] Admin dashboard (campaigns, stats, recent redemptions)
- [x] Admin password protection (env var)

## Phase 2 — Security & Hardening
- [ ] API key authentication for webhook endpoint
- [ ] Rate limiting on code generation and webhook
- [ ] Input sanitization and validation
- [ ] Request logging and monitoring
- [ ] CORS tightening per-environment

## Phase 3 — Scale & Features
- [ ] Pagination for redemptions and admin stats
- [ ] Email notifications to referrers on new redemption
- [ ] Campaign analytics (conversion rates, click tracking)
- [ ] CSV export of redemptions
- [ ] Referrer leaderboard (public page)
- [ ] Multiple base URLs per campaign (A/B testing)
- [ ] Webhook retry with dead-letter queue
