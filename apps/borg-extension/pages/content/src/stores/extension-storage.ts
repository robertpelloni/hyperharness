import type { StateStorage } from 'zustand/middleware';

interface ExtensionStorageArea {
  get(key: string): Promise<unknown>;
  set(items: Record<string, string>): Promise<void>;
  remove(key: string): Promise<void>;
}

interface BrowserStorageArea {
  get(key: string): Promise<Record<string, unknown>>;
  set(items: Record<string, string>): Promise<void>;
  remove(key: string): Promise<void>;
}

function hasLocalStorage(): boolean {
  try {
    return typeof localStorage !== 'undefined';
  } catch {
    return false;
  }
}

function getLocalStorageValue(key: string): string | null {
  if (!hasLocalStorage()) {
    return null;
  }

  try {
    return localStorage.getItem(key);
  } catch {
    return null;
  }
}

function setLocalStorageValue(key: string, value: string): void {
  if (!hasLocalStorage()) {
    return;
  }

  try {
    localStorage.setItem(key, value);
  } catch {
    // Ignore storage quota/private mode errors and let extension storage be the source of truth.
  }
}

function removeLocalStorageValue(key: string): void {
  if (!hasLocalStorage()) {
    return;
  }

  try {
    localStorage.removeItem(key);
  } catch {
    // Ignore cleanup failures.
  }
}

function createChromeStorageArea(): ExtensionStorageArea | null {
  if (typeof chrome === 'undefined' || !chrome.storage?.local) {
    return null;
  }

  return {
    async get(key: string) {
      const maybePromise = chrome.storage.local.get(key);
      if (maybePromise && typeof (maybePromise as Promise<Record<string, unknown>>).then === 'function') {
        const result = await (maybePromise as Promise<Record<string, unknown>>);
        return result?.[key];
      }

      return await new Promise<unknown>((resolve, reject) => {
        chrome.storage.local.get(key, result => {
          const error = chrome.runtime?.lastError;
          if (error) {
            reject(new Error(error.message));
            return;
          }

          resolve(result?.[key]);
        });
      });
    },
    async set(items: Record<string, string>) {
      const maybePromise = chrome.storage.local.set(items);
      if (maybePromise && typeof (maybePromise as Promise<void>).then === 'function') {
        await maybePromise;
        return;
      }

      await new Promise<void>((resolve, reject) => {
        chrome.storage.local.set(items, () => {
          const error = chrome.runtime?.lastError;
          if (error) {
            reject(new Error(error.message));
            return;
          }

          resolve();
        });
      });
    },
    async remove(key: string) {
      const maybePromise = chrome.storage.local.remove(key);
      if (maybePromise && typeof (maybePromise as Promise<void>).then === 'function') {
        await maybePromise;
        return;
      }

      await new Promise<void>((resolve, reject) => {
        chrome.storage.local.remove(key, () => {
          const error = chrome.runtime?.lastError;
          if (error) {
            reject(new Error(error.message));
            return;
          }

          resolve();
        });
      });
    },
  };
}

function createBrowserStorageArea(): ExtensionStorageArea | null {
  const extensionBrowser = (globalThis as typeof globalThis & {
    browser?: { storage?: { local?: BrowserStorageArea } };
  }).browser;

  const storageArea = extensionBrowser?.storage?.local;

  if (!storageArea) {
    return null;
  }

  return {
    async get(key: string) {
      const result = await storageArea.get(key);
      return result?.[key];
    },
    async set(items: Record<string, string>) {
      await storageArea.set(items);
    },
    async remove(key: string) {
      await storageArea.remove(key);
    },
  };
}

function getExtensionStorageArea(): ExtensionStorageArea | null {
  return createBrowserStorageArea() ?? createChromeStorageArea();
}

/**
 * Persist state in extension-scoped storage when available.
 *
 * Content-script `localStorage` is page-origin scoped, which fragments saved data by site.
 * This adapter promotes the source of truth to extension storage while still falling back to
 * localStorage in tests or non-extension environments.
 */
export function createExtensionStateStorage(): StateStorage {
  return {
    async getItem(name: string) {
      const extensionStorage = getExtensionStorageArea();

      if (!extensionStorage) {
        return getLocalStorageValue(name);
      }

      const extensionValue = await extensionStorage.get(name);
      if (typeof extensionValue === 'string') {
        return extensionValue;
      }

      const legacyValue = getLocalStorageValue(name);
      if (legacyValue !== null) {
        await extensionStorage.set({ [name]: legacyValue });
        return legacyValue;
      }

      return null;
    },
    async setItem(name: string, value: string) {
      const extensionStorage = getExtensionStorageArea();
      if (extensionStorage) {
        await extensionStorage.set({ [name]: value });
      }

      setLocalStorageValue(name, value);
    },
    async removeItem(name: string) {
      const extensionStorage = getExtensionStorageArea();
      if (extensionStorage) {
        await extensionStorage.remove(name);
      }

      removeLocalStorageValue(name);
    },
  };
}
