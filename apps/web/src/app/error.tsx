'use client';

import { useEffect } from 'react';
import { Button } from '@/components/ui/button';

export default function GlobalError({
	error,
	reset,
}: {
	error: Error & { digest?: string };
	reset: () => void;
}) {
	useEffect(() => {
		console.error('Global error occurred:', error);
	}, [error]);

	return (
		<main className='h-full flex items-center justify-center'>
			<div className='flex items-center justify-center flex-col gap-4 p-4'>
				<h2>Something went wrong!</h2>
				<Button onClick={() => reset()}>
					Try again
				</Button>
			</div>
		</main>
	);
}
