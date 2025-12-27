import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';
import { createServerClient } from '@supabase/ssr';

/**
 * Middleware to protect auth routes
 *
 * - Redirects authenticated users away from root (login) and callback routes
 * - Allows unauthenticated users to access login and callback routes
 * - Preserves the intended destination for redirect after login
 */
export async function proxy(request: NextRequest) {
  const { pathname } = request.nextUrl;

  // Create Supabase client for middleware using the existing supabase client setup
  const supabase = createServerClient(
    process.env.NEXT_PUBLIC_SUPABASE_URL!,
    process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY!,
    {
      cookies: {
        getAll() {
          return request.cookies.getAll();
        },
        setAll(cookiesToSet) {
          cookiesToSet.forEach(({ name, value, options }) => {
            request.cookies.set({
              name,
              value,
              ...options,
            });
          });
        },
      },
    }
  );

  // Get the user's session
  const {
    data: { session },
  } = await supabase.auth.getSession();

  // Check if user is trying to access auth routes (root page is now login)
  const isAuthRoute = pathname === '/' || pathname === '/callback';

  // If user is authenticated and trying to access auth routes, redirect to /app
  if (session && isAuthRoute) {
    const redirectUrl = new URL('/app', request.url);
    return NextResponse.redirect(redirectUrl);
  }

  // If user is not authenticated and not on auth routes, let them proceed
  // (Additional logic for protected routes would go here)

  return NextResponse.next();
}

// Apply middleware to all routes
export const config = {
  matcher: ['/((?!api|_next/static|_next/image|favicon.ico).*)'],
};
