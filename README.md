# Daily Deals Tracker

A simple full-stack app that tracks product prices and shows when they drop.

## Stack

- Backend: Go
- Database: PostgreSQL
- Frontend: Next.js + TypeScript
- Auth: JWT
- Scraping: goquery

## Features

- User signup & login
- Track product URLs
- Store current price + price history
- Background price refresh
- Target price alerts
- Web dashboard

## Database Tables

- users
- tracked_items
- price_history

## Project Structure

- cmd/server – app entry point
- internal/auth – authentication logic
- internal/items – item tracking logic
- internal/scraper – price scraping
- db/migrations – SQL migrations
- web – Next.js frontend

## Development Plan

- Day 1: Auth, DB schema, item tracking (fake prices)
- Day 2: Real price scraping
- Day 3: Background price refresh
- Day 4: Concurrent scraping
- Day 5: Frontend setup
- Day 6: UI polish + charts
- Day 7: Validation, testing, deployment
