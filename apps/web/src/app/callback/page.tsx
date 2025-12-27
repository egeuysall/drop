'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { createClient } from '@/lib/supabase/client';
import { Card, CardHeader, CardTitle } from '@/components/ui/card';

export default function AuthCallbackPage() {
  const router = useRouter();
  const supabase = createClient();

  useEffect(() => {
    const handleAuthCallback = async () => {
      try {
        const { data, error } = await supabase.auth.getSession();

        if (error) {
          throw error;
        }

        if (data.session) {
          router.push('/app');
        } else {
          router.push('/');
        }
      } catch (error) {
        console.error('Auth callback error:', error);
        router.push('/?error=auth_failed');
      }
    };

    handleAuthCallback();
  }, [router, supabase]);

  return (
    <main className="flex items-center justify-center">
      <Card className="w-full max-w-72">
        <CardHeader className="text-center">
          <CardTitle>Completing authentication...</CardTitle>
        </CardHeader>
      </Card>
    </main>
  );
}
