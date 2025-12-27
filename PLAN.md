# Implementation Plan

---

## Phase 1: Web Scraping (Day 2 - High Priority)
**Goal**: Extract real prices from product URLs

**Tasks**:
1. **Create `internal/scraper/` package**
   - `scraper.go`: Core `ScrapePrice(url)` function
   - Use `goquery` for HTML parsing
   - Implement multiple CSS selectors for price extraction
   - Add regex to parse currency formats ($19.99 → 19.99)
   - Check for "out of stock" indicators
   - Add proper User-Agent headers

2. **Create scraper service**
   - Error handling and retry logic
   - Domain extraction for rate limiting
   - Integration-ready interface

3. **Integrate with items service**
   - Modify `CreateItem` to call scraper
   - Save real prices to database
   - Create initial price history entry
   - Add `UpdateItemPrice` method for background updates

4. **Test with real URLs**
   - Amazon, eBay, Target product pages
   - Verify price extraction accuracy
   - Debug CSS selectors as needed

**Success Criteria**:
- [ ] Can add product URL and get actual price
- [ ] Prices stored in database with history
- [ ] Scraper handles various e-commerce sites

## Phase 2: Background Price Refresh (Day 3)
**Goal**: Automatic price updates every 30 minutes

**Tasks**:
1. **Add `RefreshAllPrices()` to items service**
   - Get items due for checking (last_checked_at)
   - Call scraper for each item
   - Compare and update prices
   - Save to price_history if changed
   - Update last_checked_at timestamp

2. **Add scheduler to main.go**
   - Use `time.Ticker` for 30-minute intervals
   - Run in goroutine to avoid blocking
   - Add proper error handling and logging

3. **Add manual trigger endpoint**
   - POST `/admin/refresh` endpoint
   - Protected by auth middleware
   - For testing and debugging

4. **Add error handling**
   - Skip failed items, don't crash
   - Log errors for debugging
   - Consider retry logic for transient failures

**Success Criteria**:
- [ ] Background job runs automatically
- [ ] Prices update in database
- [ ] Can manually trigger refresh
- [ ] System handles scraping failures gracefully

## Phase 3: Frontend Items UI (Day 5-6)
**Goal**: Complete user interface for items management

**Tasks**:
1. **Create items dashboard page**
   - Replace placeholder with real dashboard
   - Fetch items from `/v1/items` endpoint
   - Display in responsive table:
     - Name, URL, Current Price, Target Price
     - Last Checked, In Stock status
     - Price change indicators
   - Add delete functionality

2. **Create "Add Item" form**
   - Modal dialog for adding items
   - Fields: URL (required), Name, Target Price
   - URL validation
   - Loading state during scraping
   - Success/error feedback

3. **Add price history visualization**
   - Use Recharts for line charts
   - Fetch from `/v1/items/{id}/history`
   - Show price changes over time
   - Add time range selector (7d, 30d, 90d)

4. **Add notifications**
   - Toast notifications for actions
   - Price drop alerts
   - Error messages from API

**Success Criteria**:
- [ ] Can view all tracked items
- [ ] Can add new items with real prices
- [ ] Can delete items
- [ ] Can see price history charts
- [ ] Responsive and user-friendly UI

## Phase 4: Concurrency & Performance (Day 4)
**Goal**: Faster scraping with parallel processing

**Tasks**:
1. **Refactor `RefreshAllPrices` with worker pool**
   - Create jobs and results channels
   - Spawn 10 worker goroutines
   - Implement domain-based rate limiting
   - Collect and process all results

2. **Add configuration**
   - Make worker count configurable
   - Add domain delay settings
   - Environment variable support

3. **Test performance**
   - Benchmark with 50+ items
   - Compare sequential vs concurrent
   - Verify no race conditions

**Success Criteria**:
- [ ] Can scrape 50 items in seconds
- [ ] Workers run concurrently
- [ ] Rate limiting prevents blocks
- [ ] No data corruption

## Phase 5: Polish & Deploy (Day 7)
**Goal**: Production-ready application

**Tasks**:
1. **Backend improvements**
   - Add comprehensive input validation
   - Better error messages and logging
   - Handle edge cases (duplicates, not found)
   - Add API documentation

2. **Frontend improvements**
   - Form validation with clear errors
   - Loading states and skeleton screens
   - Mobile responsive design
   - Accessibility improvements

3. **Testing**
   - End-to-end test of all flows
   - Test with various e-commerce sites
   - Performance testing
   - Error scenario testing

4. **Deployment**
   - Backend: Railway/Fly.io/VPS
   - Database: Supabase/Railway Postgres
   - Frontend: Vercel
   - Set environment variables
   - Configure CI/CD

**Success Criteria**:
- [ ] App deployed and working
- [ ] All features tested
- [ ] Production-ready error handling
- [ ] Documentation complete

## Week 1: Core Functionality
- **Day 1**: Web Scraping Implementation
- **Day 2**: Background Jobs & Scheduler
- **Day 3**: Basic Frontend UI

## Week 2: Polish & Launch
- **Day 4**: Concurrency & Performance
- **Day 5**: Frontend Polish & UX
- **Day 6**: Testing & Bug Fixes
- **Day 7**: Deployment & Documentation

## Key Technical Decisions

### Web Scraping Approach
- Use `goquery` for HTML parsing
- Multiple CSS selectors for different sites
- Regex fallback for price extraction
- User-Agent rotation to avoid blocks

### Concurrency Pattern
- Worker pool with 10 workers
- Domain-based rate limiting
- Channels for job distribution
- Graceful error handling

### Frontend Architecture
- React Query for data fetching
- Recharts for visualization
- ShadCN for UI components
- Supabase Auth for authentication

### Deployment Strategy
- Monorepo with separate backend/frontend
- Environment variables for configuration
- CI/CD for automated deployments
- Monitoring and logging setup
