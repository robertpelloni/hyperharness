import React, { useState, useEffect, useRef, useMemo, useCallback, lazy, Suspense } from 'react';
import { EventsOn } from '../wailsjs/runtime/runtime';
import * as GoApp from '../wailsjs/go/main/App';
	if (!window.maestro) {
		return (
			<div
				className="flex flex-col items-center justify-center min-h-screen p-8 text-center"
				style={{
					backgroundColor: '#1a1a1a',
					color: '#f5f5f5',
					fontFamily: 'system-ui, -apple-system, sans-serif',
				}}
			>
				<div className="text-6xl mb-6">🚫</div>
				<h1 className="text-3xl font-bold mb-4 text-white">Electron Required</h1>
				<p className="text-lg opacity-80 max-w-md mb-8">
					Maestro must be run inside the Electron desktop application to access system APIs and AI
					agents.
				</p>
				<div
					className="p-4 rounded-lg text-sm text-left font-mono mb-8"
					style={{ backgroundColor: '#2d2d2d', border: '1px solid #444' }}
				>
					<div className="text-blue-400"># To run in development mode:</div>
					<div className="mt-1 text-green-400">npm run dev:win</div>
				</div>
				<p className="text-sm opacity-60">
					If you are running in Electron and still see this, check if the preload script built
					correctly.
				</p>
			</div>
		);
	}

	return (
		<InlineWizardProvider>
			<InputProvider>
				<MaestroConsoleInner />
			</InputProvider>
		</InlineWizardProvider>
	);
}
