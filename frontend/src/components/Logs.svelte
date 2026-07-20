<script lang="ts">
  import { useLogs } from "@/lib/logs.svelte";

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
  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 -960 960 960">
    <path
      d="M120-160v-80h480v80zm378.5-338.5Q440-557 440-640t58.5-141.5T640-840t141.5 58.5T840-640t-58.5 141.5T640-440t-141.5-58.5M120-480v-80h252q7 22 16 42t22 38zm0 160v-80h376q23 14 49 23.5t55 13.5v43zm500-200h40v-160h-40zm34-206q6-6 6-14t-6-14-14-6-14 6-6 14 6 14 14 6 14-6" /></svg>
</button>

<dialog bind:this={dialog}>
  <div>
    <header>
      <span>Logs</span>
      <button title="close" onclick={closeDialog}>
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 -960 960 960">
          <path
            d="m336-280-56-56 144-144-144-143 56-56 144 144 143-144 56 56-144 143 144 144-56 56-143-144z" /></svg>
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
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 -960 960 960">
          <path
            d="m376-300 104-104 104 104 56-56-104-104 104-104-56-56-104 104-104-104-56 56 104 104-104 104zm-96 180q-33 0-56.5-23.5T200-200v-520h-40v-80h200v-40h240v40h200v80h-40v520q0 33-23.5 56.5T680-120zm400-600H280v520h400zm-400 0v520z" /></svg>
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
