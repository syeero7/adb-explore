<script lang="ts">
  import { svg, CLOSE } from "@/lib/svg";

  export type Dialog = {
    value: string;
    title: string;
    oldValue?: string;
    node: HTMLDialogElement | undefined;
    onSubmit: (e: SubmitEvent) => Promise<void>;
  };

  type InputDialogProps = { dialog: Dialog };

  let { dialog = $bindable() }: InputDialogProps = $props();

  const closeDialog = () => {
    dialog.node?.close();
  };
</script>

<dialog bind:this={dialog.node}>
  <div>
    <header>
      <span>{dialog.title}</span>
      <button title="close" onclick={closeDialog}>
        {@render svg({ d: CLOSE })}
      </button>
    </header>

    <form onsubmit={dialog.onSubmit}>
      <label>
        <span>Name</span>
        <input required type="text" bind:value={dialog.value} />
      </label>

      <button onclick={closeDialog} type="button">Cancel</button>
      <button type="submit">Submit</button>
    </form>
  </div>
</dialog>

<style>
</style>
