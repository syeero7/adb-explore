import { List } from "@wails/go/main/App";

export type InfoTitle = "name" | "size" | "date modified";
type SortBy = `${InfoTitle}:${"asc" | "desc"}`;

type Dir = {
  current: string;
  query: string;
  sortBy: SortBy;
  refresh: number;
};

const STORAGE_DIR = "/storage/";

export const directory = $state<Dir>({
  current: STORAGE_DIR,
  query: "",
  sortBy: "name:asc",
  refresh: 0,
});

let parentDir = $state("/");

export function isStorageDir(dir: string) {
  return dir === STORAGE_DIR || dir === STORAGE_DIR.slice(0, -1);
}

export async function getEntries(dir: Dir) {
  const { entries, parent } = await List(dir.current, dir.query, dir.sortBy, dir.refresh);
  parentDir = parent;
  return entries;
}

export function toParentDir() {
  directory.current = parentDir;
  directory.refresh = 0;
}

export function toStorageDir() {
  directory.current = STORAGE_DIR;
  directory.refresh = 0;
}

export function toDir(path: string) {
  directory.current = path;
  directory.refresh = 0;
}

export function sortBy(title: InfoTitle, sortBy: SortBy) {
  const isActive = sortBy.startsWith(title);
  const isAsc = sortBy.endsWith("asc");
  const handler = () => {
    if (!isActive) {
      directory.sortBy = `${title}:asc`;
      return;
    }

    directory.sortBy = `${title}:${isAsc ? "desc" : "asc"}`;
  };

  return { isActive, isAsc, handler };
}
