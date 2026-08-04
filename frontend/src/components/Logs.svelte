<script lang="ts">
  import { useLogs } from "@/lib/logs.svelte";
  import { svg, CLOSE, DELETE, VIEW_LOGS } from "@/lib/svg";

  let dialog: HTMLDialogElement | undefined;
  const logs = useLogs();

  const closeDialog = () => {
    if (!dialog) return;
    dialog.close();
  };

  const openDialog = () => {
    if (!dialog) return;
    dialog.showModal();
  };
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
        <span>Clear Logs</span>
      </button>
    </footer>
  </div>
</dialog>

<style>
  dialog {
    --shadow: 0 0 4px var(--danger);
    --btn-focus: var(--danger);

    & > div {
      display: flex;
      flex-direction: column;
      gap: 0.5em;
      min-width: 28em;
      min-height: 60vh;
    }
  }

  header {
    display: flex;
    justify-content: space-between;

    span {
      font-weight: 600;
      font-size: 1.25em;
    }
  }

  button {
    width: var(--btn-size);
    height: var(--btn-size);
    border-radius: var(--radius);
    box-shadow: var(--shadow);
  }

  section {
    flex: 1;
    overflow: scroll;
    margin-block: 0.5em;
    max-height: calc(60vh - 6em);

    p {
      font-size: 0.95em;
    }
  }

  footer {
    button {
      min-width: fit-content;
      margin-left: auto;
      display: flex;
      align-items: center;
      gap: 0.25em;
    }

    span {
      font-size: 0.75em;
    }
  }

  .err {
    color: var(--danger);
  }

  .warn {
    color: var(--warning);
  }

  .info {
    color: var(--info);
  }
</style>
