.PHONY: dev

dev: ## Start development servers with hot reload (Air for Go, Next.js for frontend)
	@trap 'kill 0' INT; cd apps/web && pnpm run dev & cd apps/api && air & wait
