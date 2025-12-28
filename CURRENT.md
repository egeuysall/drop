# Drop Plan

Following your most recent progress, here is the detailed breakdown of the project’s core systems:

## 1. Automated Scraping & Amazon Optimization (Recent Focus)

The project now features a specialized **Scraper module** designed to handle the complexities of e-commerce sites.

- **Split-Price Extraction:** You implemented logic specifically for Amazon to handle "split prices," where the application extracts the `.a-price-whole` and `.a-price-fraction` elements separately and combines them into a single string for accuracy.
- **Currency Parsing:** A robust `parsePrice` utility cleans currency strings (e.g., "$1,145.99") by removing symbols and thousands separators before converting them into `float64` values using `strconv.ParseFloat`.
- **Stock Detection:** The system interprets a missing price as an "out of stock" state, allowing the service layer to update the item's status without losing the last known price.

## 2. The Background Scheduler

You implemented a **PriceRefreshScheduler** to automate the tracking process without user intervention.

- **Ticker System:** The scheduler utilizes a `time.NewTicker` to run a `refreshAllPrices` cycle every 30 minutes.
- **Background Goroutine:** By running the loop in a **goroutine**, the scheduler performs non-blocking updates, allowing the main API to remain responsive while prices are being refreshed in the background.
- **Orchestration:** The scheduler calls `GetItemsDueForCheck` from the service layer to find items not checked in over 6 hours and then iterates through them to trigger a `RefreshPrice` call for each.

## 3. Business Logic & Orchestration (Service Layer)

The **Items Service** acts as the central coordinator, injecting both the **Repository** for data and the **Scraper** for real-time updates.

- **Initial Tracking:** When a user adds a new URL via `CreateItem`, the service validates the URL, checks for duplicates, and immediately triggers the scraper to populate the initial `current_price`.
- **Price Refreshing:** The `RefreshPrice` method handles the high-level logic of updating an existing item, including catching "out of stock" errors from the scraper and updating the repository accordingly.

## 4. Data Persistence & Performance

Your data layer is designed for integrity and high performance.

- **sqlc & Type Safety:** You use **sqlc** to generate type-safe Go code from pure SQL, ensuring that database interactions are efficient.
- **Automatic History:** A PostgreSQL **trigger** (`trigger_add_price_to_history`) is used to automatically log a new entry in the `price_history` table whenever a `current_price` is updated, keeping the Go code lean and the data audit-trail perfect.
- **Security:** Data privacy is enforced via **Row Level Security (RLS)** in Supabase, ensuring users only interact with items tied to their `user_id`.

---

**Analogy for Understanding:**
Think of your project as a **High-Tech Garden**. The **Scraper** is an automated sensor that checks the height of your plants (prices). The **Repository** is the soil where all information is stored. The **Service Layer** is the Head Gardener who decides what needs watering and when. Finally, the **Scheduler** is like a programmable timer that wakes the Gardener up every few hours to make sure the sensors are checked, ensuring the garden stays healthy without you ever having to step outside.
