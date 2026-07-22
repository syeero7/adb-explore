<script lang="ts">
  import { directory, toParentDir, toStorageDir } from "@/lib/fs.svelte";
  import { svg, RELOAD, UP_ARROW, STORAGE, SEARCH, CLOSE } from "@/lib/svg";

  let isSearching = $state(false);

  function openSearch() {
    isSearching = true;
  }

  function closeSearch() {
    isSearching = false;
    directory.query = "";
  }
</script>

<header>
  <button title="refresh current directory">
    {@render svg({ d: RELOAD })}
  </button>

  <button title="go to parent directory" onclick={toParentDir}>
    {@render svg({ d: UP_ARROW })}
  </button>

  <button title="go to storage directory" onclick={toStorageDir}>
    {@render svg({ d: STORAGE })}
  </button>

  {#if isSearching}
    <input type="text" bind:value={directory.query} aria-label="search query" />
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
</header>

<style>
  button {
    height: 2.5em;
    width: 2.5em;
  }
</style>
