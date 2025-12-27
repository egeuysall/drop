'use client';

import { Github } from 'lucide-react';
import { useGitHubAuth, useUserSession } from '@/lib/api/hooks';
import { Button } from '../ui/button';
import { toast } from 'sonner';

export function GitHubAuthButton() {
  const { loginWithGitHub } = useGitHubAuth();

  const handleGitHubLogin = async () => {
    try {
      await loginWithGitHub();
    } catch (error) {
      toast.error('Login Failed', {
        description: error instanceof Error ? error.message : 'Failed to authenticate with GitHub',
      });
    }
  };

  return (
    <Button
      onClick={handleGitHubLogin}
      variant="default"
      size="default"
      className="flex w-full items-center justify-center gap-2"
      aria-label="Sign in with GitHub"
    >
      <Github className="h-4 w-4" />
      Sign in with GitHub
    </Button>
  );
}

export function AuthStatus() {
  const { data: user } = useUserSession();

  // If no user, show the login button
  if (!user) {
    return <GitHubAuthButton />;
  }

  // If user is logged in, middleware will redirect them away from auth routes
  // This fallback should never be reached due to middleware protection
  return <GitHubAuthButton />;
}
