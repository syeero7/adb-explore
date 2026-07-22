<script lang="ts">
  import { useLogs } from "@/lib/logs.svelte";
  import { svg, CLOSE, DELETE, VIEW_LOGS } from "@/lib/svg";

  let dialog: HTMLDialogElement | undefined;
  const logs = useLogs();

  function closeDialog() {
    if (!dialog) return;
    dialog.close();
  }

  function openDialog() {
    if (!dialog) return;
    dialog.showModal();
  }
</script>

<button title="view logs" onclick={openDialog}>
  {@render svg({ d: VIEW_LOGS })}
</button>

<dialog bind:this={dialog}>
  <div>
    <header>
      <span>Logs</span>
      <button title="close" onclick={closeDialog}>
        {@render svg({ d: CLOSE })}
      </button>
    </header>

    <section>
      {#each logs.data as msg}
        {const logType = msg.startsWith("e: ") ? "err" : msg.startsWith("w: ") ? "warn" : "info"}
        <p class={logType}>{msg.slice(3)}</p>
      {/each}
    </section>

    <footer>
      <button title="clear" onclick={() => logs.clear()}>
        {@render svg({ d: DELETE })}
        <span>Clear logs</span>
      </button>
    </footer>
  </div>
</dialog>

<style>
  button {
    height: 2em;
    width: 2em;
  }

  .err {
    color: red;
  }

  .warn {
    color: orange;
  }

  .info {
    color: lightblue;
  }
</style>
