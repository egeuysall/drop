import { AuthStatus } from '@/components/auth/github-auth';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';

/**
 * Simple Login Page with GitHub OAuth
 *
 * Features:
 * - Minimal design focused on GitHub authentication
 * - Responsive layout using Card component
 * - Clean, modern UI with proper spacing
 */
export default function LoginPage() {
  return (
    <main className="flex min-h-screen items-center justify-center p-4">
      <Card className="w-full max-w-72">
        <CardHeader className="text-center">
          <CardTitle>Welcome to Drop</CardTitle>
          <CardDescription className="text-center">
            Sign in to continue to your account
          </CardDescription>
        </CardHeader>

        <CardContent>
          <AuthStatus />
        </CardContent>
      </Card>
    </main>
  );
}
