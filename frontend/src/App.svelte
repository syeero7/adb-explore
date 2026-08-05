<script lang="ts">
  import Setup from "@/components/Setup.svelte";
  import Explorer from "@/components/Explorer.svelte";
  import { initLogs } from "@/lib/logs.svelte";
  import { router } from "@/lib/router.svelte";

  initLogs();
</script>

<main>
  {#if router.current === "explore"}
    <Explorer />
  {:else}
    <Setup />
  {/if}
</main>

<style>
  main {
    padding: 0.5em;
    max-width: fit-content;
    margin-inline: auto;
  }
</style>
