<script lang="ts">
  import Setup from "@/components/Setup.svelte";
  import Explore from "@/components/Explore.svelte";
  import { initLogs } from "@/lib/logs.svelte";
  import { router } from "@/lib/router.svelte";

  initLogs();
</script>

<main>
  {#if router.current === "explore"}
    <Explore />
  {:else}
    <Setup />
  {/if}
</main>

<style>
</style>
