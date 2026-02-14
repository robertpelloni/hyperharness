import { NextResponse } from 'next/server';
import { markResetRequested } from '@/lib/authStore';

export async function POST(req: Request) {
    try {
        const body = await req.json();
        const email = String(body?.email ?? '').trim();

        if (!email || !email.includes('@')) {
            return NextResponse.json({ ok: false, error: 'Please provide a valid email address.' }, { status: 400 });
        }

        await markResetRequested(email);
        return NextResponse.json({ ok: true, message: 'If this account exists, a reset link has been queued.' });
    } catch {
        return NextResponse.json({ ok: false, error: 'Invalid request payload.' }, { status: 400 });
    }
}
