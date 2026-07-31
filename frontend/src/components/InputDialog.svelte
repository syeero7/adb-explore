<script lang="ts">
  import { svg, CLOSE } from "@/lib/svg";

  export type Dialog = {
    isOpen: boolean;
    value: string;
    title: string;
    oldValue?: string;
    onSubmit: (e: SubmitEvent) => Promise<void>;
  };

  type InputDialogProps = { dialog: Dialog };

  let { dialog = $bindable() }: InputDialogProps = $props();

  function closeDialog() {
    dialog.isOpen = false;
  }
</script>

{dialog.isOpen}
<dialog open={dialog.isOpen}>
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
