<script lang="ts">
  import { tick } from "svelte";
  import Logs from "./Logs.svelte";
  import { directory, isStorageDir, refresh, toParentDir, toStorageDir } from "@/lib/fs.svelte";
  import { svg, RELOAD, UP_ARROW, STORAGE, SEARCH, CLOSE } from "@/lib/svg";

  let isSearching = $state(false);
  let timeout: number | undefined;
  // svelte-ignore <non_reactive_update>
  let input: HTMLInputElement | undefined;

  const openSearch = async () => {
    isSearching = true;
    await tick();
    input?.focus();
  };

  const refreshDir = () => {
    clearTimeout(timeout);
    refresh();
  };

  const closeSearch = () => {
    clearTimeout(timeout);
    isSearching = false;
    directory.query = "";
  };

  const search = (e: Event) => {
    clearTimeout(timeout);
    const query = (e.target as HTMLInputElement).value.trim();
    timeout = setTimeout(() => {
      directory.query = query;
    }, 300);
  };
</script>

<header>
  <button title="refresh current directory" onclick={refreshDir}>
    {@render svg({ d: RELOAD })}
  </button>

  <button
    title="go to parent directory"
    onclick={toParentDir}
    disabled={isStorageDir(directory.current)}>
    {@render svg({ d: UP_ARROW })}
  </button>

  <button
    title="go to storage directory"
    onclick={toStorageDir}
    disabled={isStorageDir(directory.current)}>
    {@render svg({ d: STORAGE })}
  </button>

  {#if isSearching}
    <input type="text" bind:this={input} oninput={search} aria-label="search query" />
    <button title="close" onclick={closeSearch}>
      {@render svg({ d: CLOSE })}
    </button>
  {:else}
    <input
      type="text"
      bind:value={directory.current}
      readonly
      aria-label="current directory path" />
    <button title="search entry" onclick={openSearch}>
      {@render svg({ d: SEARCH })}
    </button>
  {/if}

  <Logs />
</header>

<style>
  button {
    height: 2.5em;
    width: 2.5em;
  }
</style>
