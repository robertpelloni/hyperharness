
'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import { motion } from 'framer-motion';

export default function ForgotPasswordPage() {
    const [email, setEmail] = useState('');
    const [sent, setSent] = useState(false);

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        setSent(true);
        // TODO: Trigger password reset email
    };

    return (
        <div className="min-h-screen flex items-center justify-center bg-zinc-50 dark:bg-black relative overflow-hidden">
            {/* Background Ambience */}
            <div className="absolute inset-0 z-0">
                <div className="absolute top-[-20%] left-[20%] w-[50%] h-[50%] bg-blue-500/20 rounded-full blur-[120px]" />
            </div>

            <motion.div
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                className="relative z-10 w-full max-w-md p-8 bg-white/50 dark:bg-zinc-900/50 backdrop-blur-xl border border-zinc-200 dark:border-zinc-800 rounded-2xl shadow-2xl"
            >
                <div className="mb-8 text-center">
                    <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-100">
                        Reset Password
                    </h1>
                    <p className="text-zinc-500 dark:text-zinc-400 mt-2 text-sm">
                        Enter your email to receive a reset link
                    </p>
                </div>

                {!sent ? (
                    <form onSubmit={handleSubmit} className="space-y-4">
                        <input
                            type="email"
                            placeholder="Email address"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            className="w-full h-10 px-3 bg-white dark:bg-zinc-950 border border-zinc-200 dark:border-zinc-800 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500/50 transition-all text-sm"
                            required
                        />
                        <button type="submit" className="w-full h-10 bg-blue-600 hover:bg-blue-500 text-white font-medium rounded-lg transition-colors shadow-lg shadow-blue-500/20">
                            Send Reset Link
                        </button>
                    </form>
                ) : (
                    <div className="text-center space-y-4">
                        <div className="w-12 h-12 bg-green-100 dark:bg-green-900/30 text-green-600 dark:text-green-400 rounded-full flex items-center justify-center mx-auto text-xl">
                            ✓
                        </div>
                        <p className="text-zinc-600 dark:text-zinc-300">
                            Check your email for instructions to reset your password.
                        </p>
                    </div>
                )}

                <div className="mt-6 text-center text-xs">
                    <Link href="/login" className="text-zinc-400 hover:text-zinc-900 dark:hover:text-zinc-100 transition-colors">
                        ← Back to Login
                    </Link>
                </div>
            </motion.div>
        </div>
    );
}
