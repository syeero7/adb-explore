import type { main } from "@wails/go/models";
import { List } from "@wails/go/main/App";

type SortBy = `${"name" | "size" | "lastMod"}:${"asc" | "desc"}`;

type Dir = {
  current: string;
  query: string;
  sortBy: SortBy;
};

const STORAGE_DIR = "/storage/";

export const directory = $state<Dir>({
  current: STORAGE_DIR,
  query: "",
  sortBy: "name:asc",
});

let result: Omit<main.DirEntries, "entries" | "convertValues"> | undefined;

export function useIsStorageDir() {
  const value = $derived(
    directory.current === STORAGE_DIR ||
      (result && result.parent === "/") ||
      directory.current === STORAGE_DIR.slice(0, -1)
  );

  return {
    get value() {
      return value;
    },
  };
}

export async function getEntries(dir: Dir) {
  const { entries, ...rest } = await List(dir.current, dir.query, dir.sortBy);
  result = rest;
  return entries;
}

export function toParentDir() {
  directory.current = result!.parent;
}

export function toStorageDir() {
  directory.current = STORAGE_DIR;
}

export async function sortBy() {}
