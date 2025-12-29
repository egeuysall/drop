'use client';

import { useUserSession, useSignOut } from '@/lib/api/hooks';
import { Button } from '@/components/ui/button';
import { LogOut, ShoppingCart } from 'lucide-react';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import Link from 'next/link';

export default function AppPage() {
  const { data: user, isLoading } = useUserSession();
  const { mutate: signOut } = useSignOut();

  if (isLoading) {
    return (
      <main className="flex min-h-screen items-center justify-center p-4">
        <Card className="w-full max-w-72">
          <CardHeader className="text-center">
            <CardTitle>Loading...</CardTitle>
          </CardHeader>
          <CardDescription className="text-center text-neutral-700 dark:text-neutral-300">
            Please wait while we prepare your app
          </CardDescription>
        </Card>
      </main>
    );
  }

  if (!user) {
    return (
      <main className="flex min-h-screen items-center justify-center p-4">
        <Card className="w-full max-w-72">
          <CardHeader className="text-center">
            <CardTitle>Not authenticated</CardTitle>
            <CardDescription className="text-center">
              Please sign in to access this page
            </CardDescription>
          </CardHeader>
        </Card>
      </main>
    );
  }

  return (
    <main className="flex min-h-screen items-center justify-center p-4">
      <Card className="w-full max-w-72">
        <CardHeader className="text-center">
          <CardTitle>Welcome to Drop</CardTitle>
          <CardDescription className="text-center text-neutral-700 dark:text-neutral-300">
            You are signed in as {user.email}
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <Button
            onClick={() => signOut()}
            variant="destructive"
            className="flex w-full items-center justify-center gap-2"
          >
            <LogOut className="h-4 w-4" />
            <span>Sign Out</span>
          </Button>

          <Button
            asChild
            className="flex w-full items-center justify-center gap-2"
          >
            <Link href="/app/items" className='no-underline'>
              <ShoppingCart className="h-4 w-4" />
              <span>View Your Items</span>
            </Link>
          </Button>
        </CardContent>
      </Card>
    </main>
  );
}
