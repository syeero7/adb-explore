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
        <input required type="text" bind:value={dialog.value} autofocus />
      </label>

      <div>
        <button onclick={closeDialog} type="button">Cancel</button>
        <button type="submit">Submit</button>
      </div>
    </form>
  </div>
</dialog>

<style>
  form {
    display: grid;
    gap: 1em;

    & > div {
      display: flex;
      gap: 1em;
      justify-content: end;
    }

    --shadow: 0 0 4px var(--danger);
    --btn-focus: var(--danger);

    button {
      box-shadow: var(--shadow);
      padding: 0.4em 0.8em;
      border-radius: var(--radius);

      &[type="submit"] {
        --shadow: 0 0 4px var(--info);
        --btn-focus: var(--info);
      }
    }
  }

  label {
    display: grid;
    gap: 0.5em;
  }
</style>
